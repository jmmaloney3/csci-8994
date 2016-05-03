package pgg;

import java.util.Map;

import sim.engine.Steppable;

public class TestUtils {
    
    // payout parameters
    // -- cost to participate in game
    final static public double COST = 1.0;
    // -- factor used to multiply contributions
    final static public double R = 3.0;
    // -- non-participant payout
    final static public double SIGMA = 1.0;
    // -- punishment cost imposed on defectors by each punisher
    final static public double BETA = 1.0;
    // -- cost to punish a defector
    final static public double GAMMA = 0.3;
    
    /**
     * Create a new simulation that can be used for test cases.
     */
    static public PGGSim createSim(Map<Class<? extends Steppable>, Double> stratProps,
                                   int numAgents, int pggSize) {
        return new PGGSim(1L, stratProps, numAgents, pggSize, 100, 10, 0.249, null, true, false);
        
    }
    
    /**
     * Create a new empty game that can be used for test cases.
     */
    static public PGGame createGame() {
        return new PGGame(COST, R, SIGMA, BETA, GAMMA);
    }
    
    /**
     * Compare that two doubles are within PGGSim.TOL of each other.
     */
    static public void assertDoubleEquals(double expected, double actual) {
        org.junit.Assert.assertEquals(expected, actual, PGGSim.TOL);
    }
    
    /**
     * Calculate the expected payout for the specified actiongiven the
     * specified game participants.
     */
    static public double getExpectedPayout(PGGAction action, int numCooperators, int numDefectors, int numPunishers) {
        // expected payouts
        // -- calculate total contributors
        int numContributors = numCooperators + numPunishers;
        // -- calculate total participants
        int total = numContributors + numDefectors;
        // -- base payout: (cooperators*R)/(cooperators + defectors)
        double basePayout = (numContributors*R)/((double)total);
        
        // calculate payout for the specified action
        if (PGGAction.COOPERATE.equals(action)) {
            return basePayout - COST;
        }
        else if (PGGAction.DEFECT.equals(action)) {
            return basePayout - BETA*numPunishers;
        }
        else if (PGGAction.ABSTAIN.equals(action)) {
            return SIGMA;
        }
        else if (PGGAction.PUNISH.equals(action)) {
            return basePayout - COST - GAMMA*numDefectors;
        }
        else {
            throw new RuntimeException("Unknown action: " + action);
        }
    }
    
    /**
     * Utility method for calculating and testing payouts.
     */
    static public void assertPayoutsCorrect(Map<PGGAction, Double> payouts, int numCooperators, int numDefectors,
                                     int numPunishers) {
        // expected payouts
        // -- calculate total contributors
        int numContributors = numCooperators + numPunishers;
        // -- calculate total participants
        int total = numContributors + numDefectors;
        // -- base payout: (cooperators*R)/(cooperators + defectors)
        double basePayout = (numContributors*R)/((double)total);
        // -- cooperator payout: basePayout - COST
        // -- defector payout: basePayout - BETA*numPunishers
        // -- abstainer payout: SIGMA
        // -- punisher payout: basePayout - COST - GAMMA*numDefector4
        
        // test payouts
        assertDoubleEquals(payouts.get(PGGAction.COOPERATE), basePayout - COST);
        assertDoubleEquals(payouts.get(PGGAction.DEFECT), basePayout - BETA*numPunishers);
        assertDoubleEquals(payouts.get(PGGAction.ABSTAIN), SIGMA);
        assertDoubleEquals(payouts.get(PGGAction.PUNISH), basePayout - COST - GAMMA*numDefectors);
    }
    
}