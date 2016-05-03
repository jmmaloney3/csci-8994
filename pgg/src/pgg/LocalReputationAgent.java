package pgg;

import java.util.Map;
import java.util.Set;

import sim.engine.SimState;
import sim.engine.Steppable;

public class LocalReputationAgent implements Steppable, HistoryTrackingAgent {
    // constants
    private final static int DEFAULT_REP_SCORE = 0;
    
    // map to hold agent reputations
    final private Map<Steppable, Integer> repScores = new java.util.HashMap<Steppable, Integer>();

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
     * Get a map of the agents reputation scores.
     *
     * The returned map is unmodifiable.
     */
    public Map<Steppable, Integer> getRepScores() {
        return java.util.Collections.unmodifiableMap(repScores);
    }
    
    /**
     * Update the agent's history based on the results of the specified game.
     */
    public void updateHistory(PGGame game) {
        // get the actions taken by the game participants
        Map<Steppable, PGGAction> agentActions = game.getAgentActions();
        
        // update agents' reputations (image score) based on the action taken
        // -- defecting lowers reputation by one
        // -- cooperating/punishing raises reputation by one
        // -- abstaining has no impact on reputation
        Steppable agent;
        PGGAction action;
        for (Map.Entry<Steppable, PGGAction> entry : agentActions.entrySet()) {
            agent = entry.getKey();
            action = entry.getValue();
            if (PGGAction.DEFECT.equals(action)) {
                if (PGGSim.VERBOSE) { System.out.println("decrease reputation score: " + action); }
                this.decreaseReputationScore(agent);
            }
            else if (PGGAction.COOPERATE.equals(action) ||
                     PGGAction.PUNISH.equals(action)) {
                if (PGGSim.VERBOSE) { System.out.println("increase reputation score: " + action); }
                this.increaseReputationScore(agent);
            }
            else {
                if (PGGSim.VERBOSE) { System.out.println("leave reputation unchanged: " + action); }
            }
        }
    }
    
    /**
     * Choose the action to take based on the reputations of the agents in
     * the game.
     */
    /*package*/ PGGAction selectAction(PGGame game) {
        // get game participants
        Set<Steppable> participants = game.getParticipants();

        // estimate the number of cooperators and defectors
        Map<PGGAction, Integer> actionCounts = new java.util.EnumMap<PGGAction, Integer>(PGGAction.class);
        // -- initialize action counts to zero
        for (PGGAction action : PGGAction.values()) {
            actionCounts.put(action, 0);
        }
        // -- estimate actions for participants based on repScore
        int repScore;
        int count;
        for (Steppable agent : participants) {
            // don't count self
            if (this.equals(agent)) {
                break;
            }
            // estimate agents action based on repScore
            repScore = DEFAULT_REP_SCORE;
            if (this.repScores.containsKey(agent)) {
                repScore = this.repScores.get(agent);
            }
            if (repScore > 0) {
                count = actionCounts.get(PGGAction.COOPERATE);
                actionCounts.put(PGGAction.COOPERATE, ++count);
            }
            else if (repScore < 0) {
                count = actionCounts.get(PGGAction.DEFECT);
                actionCounts.put(PGGAction.DEFECT, ++count);
            }
        }
        
        // select the action that produces the maximum payout
        double maxPayout = Double.NEGATIVE_INFINITY;
        PGGAction maxAction = null;
        double po;
        // -- estimate payout for each action and update maxAction
        for (PGGAction action : PGGAction.values()) {
            po = this.estimatePayout(game, action, actionCounts);
            if (PGGSim.VERBOSE) { System.out.println("Estimated payout for action " + action + ": " + po); }
            // update maxPayout and maxAction if appropriate
            if (po > maxPayout) {
                maxPayout = po;
                maxAction = action;
            }
        }
        
        // choose abstain if no action found
        if (maxAction == null) {
            maxAction = PGGAction.ABSTAIN;
        }
        
        // return the actin with the max payout
        return maxAction;
    }
    
    /**
     * Increase the specified agent's reputation score.
     */
    private void increaseReputationScore(final Steppable agent) {
        int repScore = 0;
        if (this.repScores.containsKey(agent)) {
            repScore = this.repScores.get(agent);
        }
        this.repScores.put(agent, ++repScore);
    }
    
    /**
     * Increase the specified agent's reputation score.
     */
    private void decreaseReputationScore(final Steppable agent) {
        int repScore = 0;
        if (this.repScores.containsKey(agent)) {
            repScore = this.repScores.get(agent);
        }
        this.repScores.put(agent, --repScore);
    }
    
    /**
     * Estimate payout that will be received if the specified action is selected.
     */
    /*package*/ double estimatePayout(PGGame game, PGGAction action, Map<PGGAction, Integer> actionCounts) {
        // -- estimate the payout if specified action is selected
        // ---- add one more agent using specified action (self)
        int count = actionCounts.get(action);
        actionCounts.put(action, ++count);
        // ---- use game logic to calculate payouts
        Map<PGGAction, Double> payouts = game.getActionPayouts(actionCounts);
        // ---- get the payout for the specified action
        double actionPayout = payouts.get(action);
        // ---- remove self from set of agets using the specified action
        count = actionCounts.get(action);
        actionCounts.put(action, --count);
        
        // return the estimate payout for the specified action
        return actionPayout;
    }

    /**
     * Clear the agent's game history.
     */
    public void clearHistory() {
        this.repScores.clear();
    }
}