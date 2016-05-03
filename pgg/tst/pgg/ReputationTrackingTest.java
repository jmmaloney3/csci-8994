package pgg;

import java.util.Map;

import org.junit.Test;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;

import static pgg.TestUtils.createGame;
import static pgg.TestUtils.assertDoubleEquals;

import sim.engine.Steppable;

public class ReputationTrackingTest {
    
    public LocalReputationAgent createAgent() {
        return new LocalReputationAgent();
    }
    
    /**
     * Test the following methods in LocalReputationAgent:
     *   updateHistory
     *   selectAction
     *   clearHistory
     */
    @Test
    public void testUpdateHistory() {
        PGGame game = createGame();
        LocalReputationAgent agent = createAgent();
        
        // create set of agents
        Steppable grAgent = new GlobalReputationAgent();
        Steppable lrAgent = new LocalReputationAgent();
        Steppable defector = new Defector();
        Steppable cooperator = new Cooperator();
        Steppable abstainer = new NonParticipant();
        Steppable punisher = new Punisher();
        Steppable[] agents = { grAgent, lrAgent, defector, cooperator, abstainer, punisher };
        
        // start new game
        game.newGame(agents, 1, 2);
        
        // agents take some actions
        game.takeAction(grAgent, PGGAction.COOPERATE);
        game.takeAction(lrAgent, PGGAction.COOPERATE);
        game.takeAction(defector, PGGAction.DEFECT);
        game.takeAction(cooperator, PGGAction.COOPERATE);
        game.takeAction(abstainer, PGGAction.ABSTAIN);
        game.takeAction(punisher, PGGAction.PUNISH);
        
        // update agent history
        agent.updateHistory(game);
        
        // test the reputation scores
        Map<Steppable, Integer> repScores = agent.getRepScores();
        assertEquals((int)repScores.get(grAgent), 1);
        assertEquals((int)repScores.get(lrAgent), 1);
        assertEquals((int)repScores.get(defector), -1);
        assertEquals((int)repScores.get(cooperator), 1);
        assertFalse(repScores.containsKey(abstainer));
        assertEquals((int)repScores.get(punisher), 1);
        
        // test select action
        // -- estimated COOPERATE payout: 5*R/6 - COST = 2.5 - 1.0 = 1.5
        // -- estimated DEFECT payout:    4*R/6        = 2.0
        PGGAction action = agent.selectAction(game);
        assertEquals(PGGAction.DEFECT, action);
        
        // agents take some more actions
        game.takeAction(grAgent, PGGAction.ABSTAIN);   // no change
        game.takeAction(lrAgent, PGGAction.DEFECT);    // decrease
        game.takeAction(defector, PGGAction.ABSTAIN);  // no change
        game.takeAction(cooperator, PGGAction.DEFECT); // decrease
        game.takeAction(abstainer, PGGAction.DEFECT);  // decrease
        game.takeAction(punisher, PGGAction.DEFECT);   // decrease

        // update agent history again
        agent.updateHistory(game);
        
        // test the reputation scores
        repScores = agent.getRepScores();
        assertEquals((int)repScores.get(grAgent), 1);
        assertEquals((int)repScores.get(lrAgent), 0);
        assertEquals((int)repScores.get(defector), -1);
        assertEquals((int)repScores.get(cooperator), 0);
        assertEquals((int)repScores.get(abstainer), -1);
        assertEquals((int)repScores.get(punisher), 0);
        
        // test select action
        // -- estimated COOPERATE payout: 2*R/4 - COST = 1.5 - 1.0 = 0.5
        // -- estimated DEFECT payout:    1*R/4        = 0.75
        // -- estimated ABSTAIN payout:   SIGMA        = 1.0
        action = agent.selectAction(game);
        assertEquals(PGGAction.ABSTAIN, action);

        // agents take some more actions
        game.takeAction(grAgent, PGGAction.COOPERATE); // increase
        game.takeAction(lrAgent, PGGAction.PUNISH);    // increase
        game.takeAction(defector, PGGAction.COOPERATE);// increase
        game.takeAction(cooperator, PGGAction.PUNISH); // increase
        game.takeAction(abstainer, PGGAction.PUNISH);  // increase
        game.takeAction(punisher, PGGAction.COOPERATE);// increase
        
        // update agent history again
        agent.updateHistory(game);
        
        // test the reputation scores
        repScores = agent.getRepScores();
        assertEquals((int)repScores.get(grAgent), 2);
        assertEquals((int)repScores.get(lrAgent), 1);
        assertEquals((int)repScores.get(defector), 0);
        assertEquals((int)repScores.get(cooperator), 1);
        assertEquals((int)repScores.get(abstainer), 0);
        assertEquals((int)repScores.get(punisher), 1);
        
        // test select action
        // -- estimated COOPERATE payout: 5*R/5 - COST = 3.0 - 1.0 = 2.0
        // -- estimated DEFECT payout:    4*R/5        = 2.4
        // -- estimated ABSTAIN payout:   SIGMA        = 1.0
        action = agent.selectAction(game);
        assertEquals(PGGAction.DEFECT, action);

        // clear the history
        agent.clearHistory();
        
        // test the reputation scores
        repScores = agent.getRepScores();
        assertFalse(repScores.containsKey(grAgent));
        assertFalse(repScores.containsKey(lrAgent));
        assertFalse(repScores.containsKey(defector));
        assertFalse(repScores.containsKey(cooperator));
        assertFalse(repScores.containsKey(abstainer));
        assertFalse(repScores.containsKey(punisher));
    }
    
    /**
     * Test the following methods in LocalReputationAgent:
     *   estimatePayouts
     */
    @Test
    public void testEstmatePayouts() {
        PGGame game = createGame();
        LocalReputationAgent agent = createAgent();
        
        int numDefectors = 1;
        int numCooperators = 3;

        Map<PGGAction, Integer> actionCounts = new java.util.EnumMap<PGGAction, Integer>(PGGAction.class);
        actionCounts.put(PGGAction.DEFECT, numDefectors);
        actionCounts.put(PGGAction.COOPERATE, numCooperators);
        
        // estimate payouts
        // -- COOPERATE payout: 4*R/5 - COST
        double expCPayout = TestUtils.getExpectedPayout(PGGAction.COOPERATE, numCooperators+1, numDefectors, 0);
        double cPayout = agent.estimatePayout(game, PGGAction.COOPERATE, actionCounts);
        assertDoubleEquals(expCPayout, cPayout);
        
        // -- DEFECT payout:    3*R/5
        double expDPayout = TestUtils.getExpectedPayout(PGGAction.DEFECT, numCooperators, numDefectors+1, 0);
        double dPayout = agent.estimatePayout(game, PGGAction.DEFECT, actionCounts);
        assertDoubleEquals(expDPayout, dPayout);
    }
}