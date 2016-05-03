package pgg;

import java.util.Map;
import java.util.Set;

import org.junit.Test;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;
import static org.junit.Assert.assertFalse;

import static pgg.TestUtils.assertDoubleEquals;
import static pgg.TestUtils.assertPayoutsCorrect;
import static pgg.TestUtils.getExpectedPayout;
import static pgg.TestUtils.createGame;

import sim.engine.Steppable;

public class PGGameTest {
    
    public Map<PGGAction, Integer> createActionCounts() {
        return new java.util.EnumMap<PGGAction, Integer>(PGGAction.class);
    }
    
    @Test
    public void emptyPayoutsTest() {
        PGGame game = createGame();
        Map<PGGAction, Integer> actionCounts = createActionCounts();
        // try with empty actionCounts map
        Map<PGGAction, Double> payouts = game.getActionPayouts(actionCounts);
        // validate results
        assertTrue(payouts.isEmpty());
    }
    
    // SINGLE STRATEGY PAYOUT TESTS
    
    /**
     * Utility method for testing payout calculations.
     */
    public void runSingleStrategyPayoutTest(PGGAction action, int numAgents, double expectedPayout) {
        PGGame game = createGame();
        Map<PGGAction, Integer> actionCounts = createActionCounts();
        // populate action counts
        actionCounts.put(action, numAgents);
        // calculate payouts
        Map<PGGAction, Double> payouts = game.getActionPayouts(actionCounts);
        
        // test size
        assertEquals(payouts.size(), 1);
        
        // test payout
        // System.out.println("Test: " + payouts.get(action) + " == " + expectedPayout);
        assertDoubleEquals(payouts.get(action), expectedPayout);
    }
    
    @Test
    public void defectorsOnlyPayoutsTest() {
        int numAgents = 5;
        // expected payout is: zero
        double expectedPayout = getExpectedPayout(PGGAction.DEFECT, 0, numAgents, 0);
        // run test
        runSingleStrategyPayoutTest(PGGAction.DEFECT, numAgents, expectedPayout);
    }
    
    @Test
    public void cooperatorsOnlyPayoutsTest() {
        int numAgents = 5;
        // expected payout is: (COST*numAgents*R)/numAgents - COST
        double expectedPayout = getExpectedPayout(PGGAction.COOPERATE, numAgents, 0, 0);
        // run test
        runSingleStrategyPayoutTest(PGGAction.COOPERATE, numAgents, expectedPayout);
    }
    
    @Test
    public void punishersOnlyPayoutsTest() {
        int numAgents = 5;
        // expected payout is: (COST*numAgents*R)/numAgents - COST
        double expectedPayout = getExpectedPayout(PGGAction.PUNISH, 0, 0, numAgents);
        // run test
        runSingleStrategyPayoutTest(PGGAction.PUNISH, numAgents, expectedPayout);
    }

    @Test
    public void abstainersOnlyPayoutsTest() {
        int numAgents = 5;
        // expected payout is: SIGMA
        double expectedPayout = getExpectedPayout(PGGAction.ABSTAIN, 0, 0, 0);
        // run test
        runSingleStrategyPayoutTest(PGGAction.ABSTAIN, numAgents, expectedPayout);
    }
    
    // MULTIPLE STRATEGY PAYOUT TESTS
    
    @Test
    public void allTypesPayoutTest() {
        PGGame game = createGame();
        Map<PGGAction, Integer> actionCounts = createActionCounts();
        // populate action counts
        int numCooperators = 7;
        actionCounts.put(PGGAction.COOPERATE, numCooperators);
        int numDefectors = 5;
        actionCounts.put(PGGAction.DEFECT, numDefectors);
        int numAbstainers = 10;
        actionCounts.put(PGGAction.ABSTAIN, numAbstainers);
        int numPunishers = 4;
        actionCounts.put(PGGAction.PUNISH, numPunishers);
        
        // calculate payouts
        Map<PGGAction, Double> payouts = game.getActionPayouts(actionCounts);
        
        // test size
        assertEquals(payouts.size(), 4);
        
        // test payouts
        assertPayoutsCorrect(payouts, numCooperators, numDefectors, numPunishers);
    }
    
    // TEST GAME START & FINISH
    @Test public void startFinishTest() {
        PGGame game = createGame();
        Steppable grAgent = new GlobalReputationAgent();
        Steppable lrAgent = new LocalReputationAgent();
        Steppable defector = new Defector();
        Steppable cooperator = new Cooperator();
        Steppable abstainer = new NonParticipant();
        Steppable punisher = new Punisher();
        Steppable[] agents = { grAgent, lrAgent, defector, cooperator, abstainer, punisher };
        
        // start new game
        game.newGame(agents, 1, 2);
        
        // test participants
        Set participants = game.getParticipants();
        assertEquals(participants.size(), 6);
        for (Steppable agent : agents) {
            assertTrue(participants.contains(agent));
        }
        
        // take some actions
        game.takeAction(grAgent, PGGAction.COOPERATE);
        game.takeAction(lrAgent, PGGAction.COOPERATE);
        game.takeAction(defector, PGGAction.DEFECT);
        game.takeAction(cooperator, PGGAction.COOPERATE);
        game.takeAction(abstainer, PGGAction.ABSTAIN);
        game.takeAction(punisher, PGGAction.PUNISH);
        
        // test actions list
        Map<Steppable, PGGAction> agentActions = game.getAgentActions();
        assertEquals(agentActions.get(grAgent), PGGAction.COOPERATE);
        assertEquals(agentActions.get(lrAgent), PGGAction.COOPERATE);
        assertEquals(agentActions.get(defector), PGGAction.DEFECT);
        assertEquals(agentActions.get(cooperator), PGGAction.COOPERATE);
        assertEquals(agentActions.get(abstainer), PGGAction.ABSTAIN);
        assertEquals(agentActions.get(punisher), PGGAction.PUNISH);
        
        // test payouts
        Map<PGGAction, Double> payouts = game.getActionPayouts();
        int numCooperators = 3;
        int numDefectors = 1;
        int numPunishers = 1;
        assertPayoutsCorrect(payouts, numCooperators, numDefectors, numPunishers);
        
        // test get list of agents to notify
        Set agentsToNotify = game.getAgentsToNotify();
        assertEquals("notify 2 agents", agentsToNotify.size(), 2);
        assertTrue("notify grAgent", agentsToNotify.contains(grAgent));
        assertTrue("notify lrAgent", agentsToNotify.contains(lrAgent));
        
        // test list of agents to notify when history tracking agents abstain
        game.takeAction(grAgent, PGGAction.ABSTAIN);
        game.takeAction(lrAgent, PGGAction.ABSTAIN);
        assertEquals(game.getAgentActions().get(grAgent), PGGAction.ABSTAIN);
        assertEquals(game.getAgentActions().get(lrAgent), PGGAction.ABSTAIN);
        agentsToNotify = game.getAgentsToNotify();
        assertEquals("notify zero agents", agentsToNotify.size(), 0);
        assertFalse(agentsToNotify.contains(grAgent));
        assertFalse(agentsToNotify.contains(lrAgent));
        
        // create new game without participants
        agents = new Steppable[0];
        game.newGame(agents, 3, 4);
        assertEquals(game.getParticipants().size(), 0);
        assertEquals(game.getAgentActions().size(), 0);
        assertEquals(game.getActionPayouts().size(), 0);
        assertEquals(game.getAgentsToNotify().size(), 0);
    }
}