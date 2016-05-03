package pgg;

import java.util.Collection;
import java.util.Map;
import java.util.Set;

import sim.engine.SimState;
import sim.engine.Steppable;

/**
 * An agent that attempts to extend the binary image score ("good" or "bad")
 * approach to the n-person public goods game.
 * 
 * See:
 *  Nowak, et. al., "Evolution of indirect reciprocity by image scoring,"
 *      Nature, vol. 393, pp. 573-577, 1998.
 *
 *  Ohtsuki, et. al., "The leading eight: Social norms that maintain cooperation
 *      by indirect reciprocity," Journal of Theoretical Biology, vol. 239,
 *      pp. 435-444, 2006.
 *
 * The model described in (Ohtsuki 2006) is extended to public goods games by
 * classifying a group as "good" if all players in the group have a "good"
 * reputation and as "bad" if one or more players in the group have a "bad"
 * reputation.
 */
public class ImageScoreAgent implements Steppable, AllSeeingAgent {

    // Reputation Dynamics (aka, reputation assessment rule)
    // -- d(i, j, X) defines the reputation of an agent given that
    //    the agent's current reputation is i, the reputation of the
    //    group is j and the action taken by the agent is X
    private ImageScore[][][] d = new ImageScore[ImageScore.length][ImageScore.length][PGGAction.length];

    // Strategy
    // -- p(i, j) defines the action to be taken when this agent's
    //    current reputation is i and the reputation of the group
    //    is j
    private PGGAction[][] p = new PGGAction[ImageScore.length][ImageScore.length];

    // Image Scores
    // -- map to hold agent reputations
    final private Map<Steppable, ImageScore> iScores = new java.util.HashMap<Steppable, ImageScore>();
    
    // Default Image Score
    // -- assigned to agents by default in the absence of any other information
    private ImageScore defaultIScore = ImageScore.GOOD;
    
    // Bad Agent Threshold
    // -- the threshold used to determine whether a group is bad or good
    // -- if the percentage of "bad" agents in a group exceeds this threshold
    // -- then the group is condsidered "bad"
    private double badThresh = 0.0D;
    
    // Abstain or Defect
    // -- this flag determines whether the agent abstains or defects when
    // -- asked to participate in a public goods game with a "bad" group
    private boolean punishWithAbstain = true;
    
    /**
     * Choose the action to take in the current game.
     */
    public void step(SimState state) {
        if (state instanceof PGGSim) {
            PGGAction action = this.selectAction(((PGGSim)state).GAME);
            if (PGGSim.VERBOSE) { System.out.println("selected action: " + action); }
            ((PGGSim)state).takeAction(this, action);
        }
    }
    
    /**
     * Update this agent's reputation dynamics (aka, reputation assessment rule).
     */
    /*package*/ void updateDynamics(ImageScore rep, ImageScore grpRep, PGGAction act, ImageScore newRep) {
        d[rep.ordinal()][grpRep.ordinal()][act.ordinal()] = newRep;
    }
    
    /**
     * Apply the reputation dynamics (aka, reputation assessment rule) to the
     * specified situation and return the reputation that should be assigned
     * to the agent.
     */
    /*package*/ ImageScore applyDynamics(ImageScore rep, ImageScore grpRep, PGGAction act) {
        return d[rep.ordinal()][grpRep.ordinal()][act.ordinal()];
    }

    /**
     * Update this agent's strategy.
     *
     * If punishWithAbstain is true, a defect action is switched to abstain.
     */
    /*package*/ void updateStrategy(ImageScore rep, ImageScore grpRep, PGGAction act) {
        if (punishWithAbstain && PGGAction.DEFECT.equals(act)) {
            act = PGGAction.ABSTAIN;
        }
        p[rep.ordinal()][grpRep.ordinal()] = act;
    }
    
    /**
     * Apply the agent's strategy to the specified situation and return the
     * action that should be taken.
     */
    /*package*/ PGGAction applyStrategy(ImageScore rep, ImageScore grpRep) {
        return p[rep.ordinal()][grpRep.ordinal()];
    }

    /**
     * Choose the action to take based on the reputations of the agents in
     * the game.
     */
    /*package*/ PGGAction selectAction(PGGame game) {
        ImageScore groupRep = getGroupReputation(game.getParticipants());
        ImageScore myRep = getReputation(this);
        return applyStrategy(myRep, groupRep);
    }
    
    /**
     * Update the reputation for the specified agent.  Note that the specified agent
     * could be this agent.
     */
    /*package*/ void updateReputation(Steppable agent, ImageScore rep) {
        iScores.put(agent, rep);
    }

    /**
     * Set the default reputation that is used wen no history is available for
     * an agent.
     */
    /*package*/ void setDefaultReputation(ImageScore rep) {
        defaultIScore = rep;
    }
    
    /**
     * Return whether the agent is punishing by abstaining or by defecting.
     */
    /*package*/ boolean isPunishWithAbstain() {
        return punishWithAbstain;
    }
    
    /**
     * Get the action used to defect against "bad" groups.
     */
    /*package*/ PGGAction getDefectAction() {
        return punishWithAbstain ? PGGAction.ABSTAIN : PGGAction.DEFECT;
    }

    /**
     * Get the current reputation for the specified agent.  Note that the specified
     * agent could be this agent.
     */
    /*package*/ ImageScore getReputation(Steppable agent) {
        if (iScores.containsKey(agent)) {
            ImageScore rep = iScores.get(agent);
            if (rep == null) { rep = defaultIScore; }
            return rep;
        }
        else {
            return defaultIScore;
        }
    }
    
    /**
     * Set the threshold used to determine whether or not a group is "bad".
     * The threshold indicates the maximum percentage of agents that can be
     * "bad" in order for the group to still be considered "good".
     */
    /*package*/ void setBadThreshold(double t) {
        this.badThresh = t;
    }

    /**
     * Calculate the reputation of the group of agents.
     */
    /*package*/ ImageScore getGroupReputation(Collection<Steppable> group) {
        if (badThresh < 0) { return ImageScore.BAD; }
        
        // count number of bad agents in group
        int totalAgents = group.size();
        int badAgents = 0;
        double percentBad;
        for (Steppable player : group) {
            if (ImageScore.BAD.equals(getReputation(player))) {
                if (getReputation(player) == null) {
                    System.out.println("Agent reputation is null.");
                }
                badAgents++;
                percentBad = ((double)badAgents)/((double)totalAgents);
                if (percentBad > badThresh) {
                    return ImageScore.BAD;
                }
            }
        }
        
        return ImageScore.GOOD;
    }
    
    /**
     * Update the agent's history based on the results of the specified game.
     */
    public void updateHistory(PGGame game) {
        ImageScore grpRep = getGroupReputation(game.getParticipants());
        ImageScore oldRep, newRep;
        Steppable agent;
        PGGAction action;
        for (Map.Entry<Steppable, PGGAction> entry : game.getAgentActions().entrySet()) {
            agent =  entry.getKey();
            oldRep = getReputation(agent);
            action = entry.getValue();
            // if punishWithAbstain then switch defect to abstain
            if (punishWithAbstain && PGGAction.ABSTAIN.equals(action)) {
                action = PGGAction.DEFECT;
            }
            if ((PGGAction.PUNISH.equals(action) || (PGGAction.ABSTAIN.equals(action)))) {
                throw new RuntimeException("Unsupported action: " + action);
            }
            newRep = applyDynamics(oldRep, grpRep, action);
            //if (newRep == null) { newRep = defaultIScore; }
            updateReputation(agent, newRep);
        }
    }

    /**
     * Clear the agent's game history.
     */
    public void clearHistory() {
        iScores.clear();
    }
}
