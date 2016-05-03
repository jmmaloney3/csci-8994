package pgg;

import java.util.List;

import org.junit.Test;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;
import static org.junit.Assert.assertNull;

import static pgg.TestUtils.createGame;

import sim.engine.Steppable;

/**
 * Test cases for the pgg.ImageScoreAgent class.
 */
public class ImageScoreAgentTest {

    /**
     * Test the reputation dynamics.  Test the following methods:
     *   updateDynamics
     *   applyDynamics
     */
    @Test
    public void testDynamics() {
        ImageScoreAgent agent = new ImageScoreAgent();
        
        // configure the reputation dynamics
        agent.updateDynamics(ImageScore.GOOD, ImageScore.GOOD, PGGAction.COOPERATE, ImageScore.BAD);
        agent.updateDynamics(ImageScore.BAD, ImageScore.BAD, PGGAction.DEFECT, ImageScore.GOOD);
        
        // test application of the reputation dynamics
        ImageScore rep = agent.applyDynamics(ImageScore.GOOD, ImageScore.GOOD, PGGAction.COOPERATE);
        assertEquals(ImageScore.BAD, rep);
        rep = agent.applyDynamics(ImageScore.BAD, ImageScore.BAD, PGGAction.DEFECT);
        assertEquals(ImageScore.GOOD, rep);
        rep = agent.applyDynamics(ImageScore.GOOD, ImageScore.BAD, PGGAction.DEFECT);
        assertNull(rep);
    }
    
    /**
     * Test the strategy.  Test the following methods:
     *   updateStrategy
     *   applyStrategy
     */
    @Test
    public void testStrategy() {
        ImageScoreAgent agent = new ImageScoreAgent();

        // configure strategy
        agent.updateStrategy(ImageScore.GOOD, ImageScore.GOOD, PGGAction.DEFECT);
        agent.updateStrategy(ImageScore.BAD, ImageScore.BAD, PGGAction.COOPERATE);
        
        // test application of strategy
        PGGAction act = agent.applyStrategy(ImageScore.GOOD, ImageScore.GOOD);
        assertEquals(agent.getDefectAction(), act);
        act = agent.applyStrategy(ImageScore.BAD, ImageScore.BAD);
        assertEquals(PGGAction.COOPERATE, act);
        act = agent.applyStrategy(ImageScore.GOOD, ImageScore.BAD);
        assertNull(act);
    }
    
    /**
     * Test updating and retrieving reputations.  Test the following methods:
     *   updateReputation
     *   getReputation
     *   setDefaultReputation
     *   getGroupReputation
     *   setBadThreshold
     */
    @Test
    public void testGetReputation() {
        // create a set of agents
        ImageScoreAgent agent1 = new ImageScoreAgent();
        Steppable agent2 = new Defector();
        Steppable agent3 = new Cooperator();
        Steppable agent4 = new NonParticipant();
        Steppable agent5 = new Punisher();
        Steppable[] agents = new Steppable[]{ agent1, agent2, agent3, agent4, agent5 };
        
        // test get reputation without any history
        // -- default reputation is "good" by default
        ImageScore rep = agent1.getReputation(agent1);
        assertEquals(ImageScore.GOOD, rep);
        rep = agent1.getReputation(agent3);
        assertEquals(ImageScore.GOOD, rep);
        // -- set default reputation to "bad"
        agent1.setDefaultReputation(ImageScore.BAD);
        rep = agent1.getReputation(agent1);
        assertEquals(ImageScore.BAD, rep);
        rep = agent1.getReputation(agent3);
        assertEquals(ImageScore.BAD, rep);
        // -- set default reputation back to "good"
        agent1.setDefaultReputation(ImageScore.GOOD);
        
        // set reputation for each agent
        agent1.updateReputation(agent1, ImageScore.GOOD);
        agent1.updateReputation(agent2, ImageScore.GOOD);
        agent1.updateReputation(agent3, ImageScore.GOOD);
        agent1.updateReputation(agent4, ImageScore.BAD);
        agent1.updateReputation(agent5, ImageScore.BAD);
        
        // test get reputation with history
        rep = agent1.getReputation(agent1);
        assertEquals(ImageScore.GOOD, rep);
        rep = agent1.getReputation(agent2);
        assertEquals(ImageScore.GOOD, rep);
        rep = agent1.getReputation(agent3);
        assertEquals(ImageScore.GOOD, rep);
        rep = agent1.getReputation(agent4);
        assertEquals(ImageScore.BAD, rep);
        rep = agent1.getReputation(agent5);
        assertEquals(ImageScore.BAD, rep);
        
        // test get group reputation
        List<Steppable> group = java.util.Arrays.asList(agents);
        // set "bad" threshold to 50%
        agent1.setBadThreshold(0.5);
        // group is 40% "bad" - so will have a "good" reputation
        rep = agent1.getGroupReputation(group);
        assertEquals(ImageScore.GOOD, rep);
        // update reputation of agent3 so that group becomes 60% "bad"
        agent1.updateReputation(agent3, ImageScore.BAD);
        rep = agent1.getReputation(agent3);
        assertEquals(ImageScore.BAD, rep);  // make sure agent3's rep was updated
        // group is 60% "bad" - so will have a "bad" reputation
        rep = agent1.getGroupReputation(group);
        assertEquals(ImageScore.BAD, rep);
        // change threshold to 60%
        agent1.setBadThreshold(0.6);
        // now 60% "bad" is "good" - so will have a "good" reputation
        rep = agent1.getGroupReputation(group);
        assertEquals(ImageScore.GOOD, rep);
    }
    
    /**
     * Test the history tracking interface.  Test the following methods:
     *   selectAction
     *   updateHistory
     *   clearHistory
     */
    @Test
    public void testHistoryTracking() {
        PGGame game = createGame();
        
        // create a set of agents
        ImageScoreAgent agent1 = new ImageScoreAgent();
        Steppable agent2 = new Defector();
        Steppable agent3 = new Cooperator();
        Steppable agent4 = new NonParticipant();
        Steppable agent5 = new Punisher();
        Steppable[] agents = new Steppable[]{ agent1, agent2, agent3, agent4, agent5 };
        
        // start new game
        game.newGame(agents, 1, 2);
        
        // select an action without any strategy defined
        PGGAction act = agent1.selectAction(game);
        assertNull(act);
        
        // configure agent1's strategy
        agent1.updateStrategy(ImageScore.GOOD, ImageScore.GOOD, PGGAction.COOPERATE);
        agent1.updateStrategy(ImageScore.GOOD, ImageScore.BAD,  PGGAction.DEFECT);
        agent1.updateStrategy(ImageScore.BAD,  ImageScore.GOOD, PGGAction.COOPERATE);
        agent1.updateStrategy(ImageScore.BAD,  ImageScore.BAD,  PGGAction.DEFECT);
        
        // select an action without any history
        // -- with default reputation set to "good" (the default setting)
        act = agent1.selectAction(game);
        // "good" agent and "good" group => cooperate
        assertEquals(PGGAction.COOPERATE, act);
        // -- with default reputation set to "bad"
        agent1.setDefaultReputation(ImageScore.BAD);
        act = agent1.selectAction(game);
        // "bad" agent and "bad" group => defect
        assertEquals(agent1.getDefectAction(), act);
        // -- set default reputation back to "good"
        agent1.setDefaultReputation(ImageScore.GOOD);
        
        // agents take some actions - without any dynamics configured
        game.takeAction(agent1, PGGAction.COOPERATE);
        game.takeAction(agent2, PGGAction.COOPERATE);
        game.takeAction(agent3, PGGAction.DEFECT);
        game.takeAction(agent4, PGGAction.COOPERATE);
        game.takeAction(agent5, PGGAction.DEFECT);
        
        // update agent history - without dynamics configured
        agent1.updateHistory(game);
        
        // check agent's reputations
        // -- will be equal to default reputation (GOOD) because no reputation dynamics are defined
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent1));
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent2));
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent3));
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent4));
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent5));
        
        // set default reputation to "bad"
        agent1.setDefaultReputation(ImageScore.BAD);
        
        // update agent history - without dynamics configured
        agent1.updateHistory(game);
        
        // check agent's reputations
        // -- will be equal to default reputation (BAD) because no reputation dynamics are defined
        assertEquals(ImageScore.BAD, agent1.getReputation(agent1));
        assertEquals(ImageScore.BAD, agent1.getReputation(agent2));
        assertEquals(ImageScore.BAD, agent1.getReputation(agent3));
        assertEquals(ImageScore.BAD, agent1.getReputation(agent4));
        assertEquals(ImageScore.BAD, agent1.getReputation(agent5));
        
        // set default reputation back to "good"
        agent1.setDefaultReputation(ImageScore.GOOD);

        // configure reputation dynamics
        agent1.updateDynamics(ImageScore.GOOD, ImageScore.GOOD, PGGAction.COOPERATE, ImageScore.GOOD);
        agent1.updateDynamics(ImageScore.GOOD, ImageScore.GOOD, PGGAction.DEFECT,    ImageScore.BAD);
        agent1.updateDynamics(ImageScore.GOOD, ImageScore.BAD,  PGGAction.COOPERATE, ImageScore.BAD);
        agent1.updateDynamics(ImageScore.GOOD, ImageScore.BAD,  PGGAction.DEFECT,    ImageScore.GOOD);
        agent1.updateDynamics(ImageScore.BAD,  ImageScore.GOOD, PGGAction.COOPERATE, ImageScore.GOOD);
        agent1.updateDynamics(ImageScore.BAD,  ImageScore.GOOD, PGGAction.DEFECT,    ImageScore.BAD);
        agent1.updateDynamics(ImageScore.BAD,  ImageScore.BAD,  PGGAction.COOPERATE, ImageScore.BAD);
        agent1.updateDynamics(ImageScore.BAD,  ImageScore.BAD,  PGGAction.DEFECT,    ImageScore.BAD);
 
        // update agent history - with dynamics configured
        agent1.clearHistory();
        agent1.updateHistory(game);

        // check agent's reputations
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent1));
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent2));
        assertEquals(ImageScore.BAD,  agent1.getReputation(agent3)); // defector
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent4));
        assertEquals(ImageScore.BAD,  agent1.getReputation(agent5)); // defector

        // update agent history - with dynamics configured
        agent1.setBadThreshold(0.5);
        agent1.updateHistory(game);
        
        // check agent's reputations
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent1));
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent2));
        assertEquals(ImageScore.BAD,  agent1.getReputation(agent3)); // defector
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent4));
        assertEquals(ImageScore.BAD,  agent1.getReputation(agent5)); // defector
        
        // update agent history - with dynamics configured
        agent1.setBadThreshold(0.0);
        game.takeAction(agent1, PGGAction.DEFECT); // defect against "bad" group
        agent1.updateHistory(game);
        
        // check agent's reputations
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent1)); // punisher
        assertEquals(ImageScore.BAD,  agent1.getReputation(agent2)); // collaborator
        assertEquals(ImageScore.BAD,  agent1.getReputation(agent3)); // defector
        assertEquals(ImageScore.BAD,  agent1.getReputation(agent4)); // collaborator
        assertEquals(ImageScore.BAD,  agent1.getReputation(agent5)); // defector
        
        // test clear history
        agent1.clearHistory();
        agent1.updateHistory(game);
        
        // check agent's reputations
        assertEquals(ImageScore.BAD,  agent1.getReputation(agent1)); // defector
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent2));
        assertEquals(ImageScore.BAD,  agent1.getReputation(agent3)); // defector
        assertEquals(ImageScore.GOOD, agent1.getReputation(agent4));
        assertEquals(ImageScore.BAD,  agent1.getReputation(agent5)); // defector
        
        // PUNISH is not supported
        game.takeAction(agent2, PGGAction.PUNISH);
        boolean exceptionThrown = false;
        try {
            agent1.updateHistory(game);
        }
        catch (RuntimeException ex) {
            String expectedMsg = "Unsupported action: " + PGGAction.PUNISH;
            if (expectedMsg.equals(ex.getMessage())) {
                exceptionThrown = true;
            }
        }
        assertTrue(exceptionThrown);

        // ABSTAIN is not supported (unless it is used to punish)
        game.takeAction(agent2, PGGAction.ABSTAIN);
        exceptionThrown = false;
        System.out.println("Punish with abstain: " + agent1.isPunishWithAbstain());
        try {
            agent1.updateHistory(game);
        }
        catch (RuntimeException ex) {
            String expectedMsg = "Unsupported action: " + PGGAction.ABSTAIN;
            if (expectedMsg.equals(ex.getMessage())) {
                exceptionThrown = true;
            }
        }
        assertEquals(exceptionThrown, !agent1.isPunishWithAbstain());
    }
}