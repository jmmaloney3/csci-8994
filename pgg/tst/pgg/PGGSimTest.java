package pgg;

import java.util.Collection;
import java.util.Map;
import java.util.Set;

import org.junit.Test;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;

import static pgg.TestUtils.createSim;
import static pgg.TestUtils.createGame;
import static pgg.TestUtils.assertDoubleEquals;
import static pgg.TestUtils.getExpectedPayout;

import sim.engine.Steppable;

public class PGGSimTest {
    
    /**
     * Utility method to take actions and update strategy statistics
     * given the specified set of actions.
     */
    public void updateStrategyStats(PGGSim sim, Map<Steppable, PGGAction> agentActions) {
        
        // get the game participants
        Steppable[] agents = agentActions.keySet().toArray(new Steppable[0]);
        
        // start new game
        sim.GAME.newGame(agents, 1, 1);
        
        // take actions
        for (Map.Entry<Steppable, PGGAction> entry : agentActions.entrySet()) {
            sim.GAME.takeAction(entry.getKey(), entry.getValue());
        }
        
        // update strategy statistics
        sim.updateStrategyStats(sim.GAME);
    }
    
    /**
     * Utility method to take actions for agents and validate the
     * strategy statistics are updated appropriately.
     */
    public void assertStrategyStatsCorrect(PGGSim sim, Map<Steppable, PGGAction> agentActions,
                                          Map<Class<? extends Steppable>, Double> expAvgPayouts,
                                          Map<Class<? extends Steppable>, Integer> expTotalGames) {
        
        // update strategy statistics given the set of specified actions
        updateStrategyStats(sim, agentActions);

        Class<? extends Steppable> strategy;
        // validate average payouts
        for (Map.Entry<Class<? extends Steppable>, Double> entry : sim.getStrategyAvgPayouts().entrySet()) {
            strategy = entry.getKey();
            assertDoubleEquals(entry.getValue(), expAvgPayouts.get(strategy));
        }

        // validate games played
        for (Map.Entry<Class<? extends Steppable>, Integer> entry : sim.getStrategyTotalGames().entrySet()) {
            strategy = entry.getKey();
            assertDoubleEquals(entry.getValue(), expTotalGames.get(strategy));
        }
    }
                                          

    @Test
    public void testUpdateStrategyStats() {
        
        // set up sim for testing
        int numAgents = 4;
        int pggSize = 4;
        Map<Class<? extends Steppable>, Double> stratProps
            = new java.util.LinkedHashMap<Class<? extends Steppable>, Double>();
        stratProps.put(Defector.class, 0.25);
        stratProps.put(Cooperator.class, 0.25);
        stratProps.put(NonParticipant.class, 0.25);
        stratProps.put(Punisher.class, 0.25);
        PGGSim sim = createSim(stratProps, numAgents, pggSize);

        // create set of agents
        Steppable defector = new Defector();
        Steppable cooperator = new Cooperator();
        Steppable abstainer = new NonParticipant();
        Steppable punisher = new Punisher();
        Steppable[] agents = { defector, cooperator, abstainer, punisher };

        int numD, numC, numP;
        
        // TEST CASE 1
        
        // create actions for agents to take
        Map<Steppable, PGGAction> agentActions = new java.util.HashMap<Steppable, PGGAction>();
        agentActions.put(defector, PGGAction.DEFECT);
        agentActions.put(cooperator, PGGAction.COOPERATE);
        agentActions.put(abstainer, PGGAction.ABSTAIN);
        agentActions.put(punisher, PGGAction.PUNISH);
        
        // count up number of actions of each type
        numD = 0;
        numC = 0;
        numP = 0;
        for (PGGAction action : agentActions.values()) {
            if (PGGAction.DEFECT.equals(action)) { numD++; }
            else if (PGGAction.COOPERATE.equals(action)) { numC++; }
            else if (PGGAction.PUNISH.equals(action)) { numP++; }
        }

        // get the expected average payouts and total games
        Map<Class<? extends Steppable>, Double> expAvgPayouts
            = new java.util.HashMap<Class<? extends Steppable>, Double>();
        Map<Class<? extends Steppable>, Integer> expTotalGames
            = new java.util.HashMap<Class<? extends Steppable>, Integer>();
        Class<? extends Steppable> strategy;
        double expPayout;
        for (Map.Entry<Steppable, PGGAction> entry : agentActions.entrySet()) {
            strategy = entry.getKey().getClass();
            expPayout = getExpectedPayout(entry.getValue(), numC, numD, numP);
            //System.out.println(entry.getValue() + " payout: " + expPayout);
            // for first game - average payout equals the payout
            expAvgPayouts.put(strategy, expPayout);
            // for first game - total games equals one
            expTotalGames.put(strategy, 1);
        }
        
        // run test and verify results
        assertStrategyStatsCorrect(sim, agentActions, expAvgPayouts, expTotalGames);
        
        // TEST CASE 2
        
        // create actions for agents to take
        agentActions.clear();
        agentActions.put(defector, PGGAction.COOPERATE);
        agentActions.put(cooperator, PGGAction.COOPERATE);
        agentActions.put(abstainer, PGGAction.PUNISH);
        agentActions.put(punisher, PGGAction.PUNISH);

        // count up number of actions of each type
        numD = 0;
        numC = 0;
        numP = 0;
        for (PGGAction action : agentActions.values()) {
            if (PGGAction.DEFECT.equals(action)) { numD++; }
            else if (PGGAction.COOPERATE.equals(action)) { numC++; }
            else if (PGGAction.PUNISH.equals(action)) { numP++; }
        }
        
        // get the expected average payouts and total games
        for (Map.Entry<Steppable, PGGAction> entry : agentActions.entrySet()) {
            strategy = entry.getKey().getClass();
            expPayout = getExpectedPayout(entry.getValue(), numC, numD, numP);
            //System.out.println(entry.getValue() + " payout: " + expPayout);
            // for second game - calculate average of first and second payouts
            expAvgPayouts.put(strategy, (expAvgPayouts.get(strategy) + expPayout)/2.0D);
            // for second game - total games equals two
            expTotalGames.put(strategy, 2);
        }

        // run test and verify results
        assertStrategyStatsCorrect(sim, agentActions, expAvgPayouts, expTotalGames);
        
        // TEST CASE 3

        // create actions for agents to take
        agentActions.clear();
        agentActions.put(defector, PGGAction.DEFECT);
        agentActions.put(cooperator, PGGAction.DEFECT);
        agentActions.put(abstainer, PGGAction.ABSTAIN);
        agentActions.put(punisher, PGGAction.ABSTAIN);
        
        // count up number of actions of each type
        numD = 0;
        numC = 0;
        numP = 0;
        for (PGGAction action : agentActions.values()) {
            if (PGGAction.DEFECT.equals(action)) { numD++; }
            else if (PGGAction.COOPERATE.equals(action)) { numC++; }
            else if (PGGAction.PUNISH.equals(action)) { numP++; }
        }
        
        // get the expected average payouts and total games
        for (Map.Entry<Steppable, PGGAction> entry : agentActions.entrySet()) {
            strategy = entry.getKey().getClass();
            expPayout = getExpectedPayout(entry.getValue(), numC, numD, numP);
            //System.out.println(entry.getValue() + " payout: " + expPayout);
            // for third game - calculate average of first, second and third payouts
            expAvgPayouts.put(strategy, (expAvgPayouts.get(strategy)*2 + expPayout)/3.0D);
            // for second game - total games equals three
            expTotalGames.put(strategy, 3);
        }
        
        // run test and verify results
        assertStrategyStatsCorrect(sim, agentActions, expAvgPayouts, expTotalGames);
        
        // TEST CASE 4
        
        // create actions for agents to take
        agentActions.clear();
        agentActions.put(defector, PGGAction.PUNISH);
        agentActions.put(cooperator, PGGAction.ABSTAIN);
        agentActions.put(abstainer, PGGAction.COOPERATE);
        agentActions.put(punisher, PGGAction.DEFECT);
        
        // count up number of actions of each type
        numD = 0;
        numC = 0;
        numP = 0;
        for (PGGAction action : agentActions.values()) {
            if (PGGAction.DEFECT.equals(action)) { numD++; }
            else if (PGGAction.COOPERATE.equals(action)) { numC++; }
            else if (PGGAction.PUNISH.equals(action)) { numP++; }
        }
        
        // get the expected average payouts and total games
        for (Map.Entry<Steppable, PGGAction> entry : agentActions.entrySet()) {
            strategy = entry.getKey().getClass();
            expPayout = getExpectedPayout(entry.getValue(), numC, numD, numP);
            //System.out.println(entry.getValue() + " payout: " + expPayout);
            // for third game - calculate average of first, second, third and fourth payouts
            expAvgPayouts.put(strategy, (expAvgPayouts.get(strategy)*3 + expPayout)/4.0D);
            // for second game - total games equals four
            expTotalGames.put(strategy, 4);
        }
        
        // run test and verify results
        assertStrategyStatsCorrect(sim, agentActions, expAvgPayouts, expTotalGames);
    }
    
    @Test
    public void testUpdateEvolutionThresholds() {
        // set up sim for testing
        int numAgents = 4;
        int pggSize = 4;
        Map<Class<? extends Steppable>, Double> stratProps
            = new java.util.LinkedHashMap<Class<? extends Steppable>, Double>();
        stratProps.put(Defector.class, 0.25);
        stratProps.put(Cooperator.class, 0.15);
        stratProps.put(NonParticipant.class, 0.55);
        stratProps.put(Punisher.class, 0.05);
        PGGSim sim = createSim(stratProps, numAgents, pggSize);

        // test the birth thresholds
        Map<Class<? extends Steppable>, Double> birthProbs = sim.getStrategyBirthProbs();
        assertDoubleEquals(birthProbs.get(Punisher.class), 0.05D);
        assertDoubleEquals(birthProbs.get(Cooperator.class), 0.2D);
        assertDoubleEquals(birthProbs.get(Defector.class), 0.45D);
        assertDoubleEquals(birthProbs.get(NonParticipant.class), 1.0D);
        
        // test the death thresholds
        Map<Class<? extends Steppable>, Double> deathProbs = sim.getStrategyDeathProbs();
        assertDoubleEquals(deathProbs.get(Punisher.class), 0.55D);
        assertDoubleEquals(deathProbs.get(Cooperator.class), 0.8D);
        assertDoubleEquals(deathProbs.get(Defector.class), 0.95D);
        assertDoubleEquals(deathProbs.get(NonParticipant.class), 1.0D);

        // set new fitness values
        stratProps = new java.util.LinkedHashMap<Class<? extends Steppable>, Double>();
        stratProps.put(Defector.class, 0.4);
        stratProps.put(Cooperator.class, 0.3);
        stratProps.put(NonParticipant.class, 0.2);
        stratProps.put(Punisher.class, 0.1);

        // update the birth and death thresholds based on the fitness measures
        sim.updateEvolutionThresholds(stratProps);
        
        // test the birth thresholds
        birthProbs = sim.getStrategyBirthProbs();
        assertDoubleEquals(birthProbs.get(Punisher.class), 0.1D);
        assertDoubleEquals(birthProbs.get(NonParticipant.class), 0.3D);
        assertDoubleEquals(birthProbs.get(Cooperator.class), 0.6D);
        assertDoubleEquals(birthProbs.get(Defector.class), 1.0D);
        
        // test the death thresholds
        deathProbs = sim.getStrategyDeathProbs();
        assertDoubleEquals(deathProbs.get(Punisher.class), 0.4D);
        assertDoubleEquals(deathProbs.get(NonParticipant.class), 0.7D);
        assertDoubleEquals(deathProbs.get(Cooperator.class), 0.9D);
        assertDoubleEquals(deathProbs.get(Defector.class), 1.0D);
    }
    
    @Test
    public void testCalculateStrategyFitness() {
        // set up sim for testing
        int numAgents = 4;
        int pggSize = 4;
        Map<Class<? extends Steppable>, Double> stratProps
        = new java.util.LinkedHashMap<Class<? extends Steppable>, Double>();
        stratProps.put(Defector.class, 0.25);
        stratProps.put(Cooperator.class, 0.15);
        stratProps.put(NonParticipant.class, 0.55);
        stratProps.put(Punisher.class, 0.05);
        PGGSim sim = createSim(stratProps, numAgents, pggSize);
        
        // create set of agents
        Steppable defector = new Defector();
        Steppable cooperator = new Cooperator();
        Steppable abstainer = new NonParticipant();
        Steppable punisher = new Punisher();
        Steppable[] agents = { defector, cooperator, abstainer, punisher };
        
        int numD, numC, numP;
        
        // create actions for agents to take
        Map<Steppable, PGGAction> agentActions = new java.util.HashMap<Steppable, PGGAction>();
        agentActions.put(defector, PGGAction.DEFECT);
        agentActions.put(cooperator, PGGAction.COOPERATE);
        agentActions.put(abstainer, PGGAction.ABSTAIN);
        agentActions.put(punisher, PGGAction.PUNISH);

        // update strategy statistics based on acions taken by agents
        updateStrategyStats(sim, agentActions);
        
        // payouts:
        // -- Defector:   2*3/3 - 1 = 1
        // -- Cooperator: 2*3/3 - 1 = 1
        // -- Punisher:   2*3/3 - 1 - 0.3 = 0.7
        // -- Abstainer:  1.0
        
        // fitness:
        // -- base fitness: (1 - 0.249) = 0.751
        // -- Defector:     0.751 + (0.249 * 1)   = 1
        // -- Cooperator:   0.751 + (0.249 * 1)   = 1
        // -- Punisher:     0.751 + (0.249 * 0.7) = 0.9253
        // -- Abstainer:    0.751 + (0.249 * 1)   = 1
        //
        // -- total fitness: 3.9253
        
        // normalized fitness:
        // -- Defector:     1/3.9253      = 0.25475...
        // -- Cooperator:   1/3.9253      = 0.25475...
        // -- Punisher:     0.9253/3.9253 = 0.23572...
        // -- Abstainer:    1/3.9253      = 0.25475...
        
        // calculate fitness values
        Map<Class<? extends Steppable>, Double> fitnessValues = sim.calculateStrategyFitness();
        
        assertDoubleEquals(0.25475, fitnessValues.get(Defector.class));
        assertDoubleEquals(0.25475, fitnessValues.get(Cooperator.class));
        assertDoubleEquals(0.23572, fitnessValues.get(Punisher.class));
        assertDoubleEquals(0.25475, fitnessValues.get(NonParticipant.class));
    }
}