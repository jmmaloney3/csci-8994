package pgg;

import java.util.Map;
import java.util.Set;

import sim.engine.Steppable;

/**
 * A class used to do bookkeeping for a public goods game.
 */
public class PGGame {
    // constants
    static final public boolean VERBOSE = false;
    
    // payout parameters
    // -- cost to participate in game
    final private double COST;
    // -- factor used to multiply contributions
    final private double R;
    // -- non-participant payout
    final private double SIGMA;
    // -- punishment cost imposed on defectors by each punisher
    final private double BETA;
    // -- cost to punish a defector
    final private double GAMMA;
    
    // instance variables
    private int roundNum;
    private int gameNum;
    
    // map to hold actions taken by the participants
    final private Map<Steppable,PGGAction> agentActions = new java.util.HashMap<Steppable, PGGAction>();
    
    // map to hold the count of the number of agents that have taken each action
    final private Map<PGGAction, Integer> actionCounts =
        new java.util.EnumMap<PGGAction, Integer>(PGGAction.class);
    
    /**
     * Package constructor used to create the singleton instance.
     */
    /*package*/ PGGame(final double cost, final double r, final double sigma, final double beta, final double gamma) {
        super();
        this.COST = cost;
        this.R = r;
        this.SIGMA = sigma;
        this.BETA = beta;
        this.GAMMA = gamma;
    }
    
    /**
     * Create a new public goods game with the specified set of agents.
     */
    public void newGame(final Steppable[] agents, final int roundNum, final int gameNum) {
        // clear the previous game results
        this.clear();

        // copy participants into results map
        for (Steppable agent : agents) {
            agentActions.put(agent, null);
        }
        
        // set round and game number
        this.roundNum = roundNum;
        this.gameNum = gameNum;
    }
    
    /**
     * Return the list of agents selected to participate in this game.
     */
    public Set<Steppable> getParticipants() {
        return agentActions.keySet();
    }
    
    /**
     * Take an action for an agent.  This method is used by game participants
     * to register their action with the simulation engine.
     */
    protected void takeAction(final Steppable agent, final PGGAction action) {
        if (PGGSim.VERBOSE) { System.out.println(agent + " => " + action); }
        
        // set the action taken by this agent
        this.agentActions.put(agent, action);
        
        // update action counts
        actionCounts.put(action, actionCounts.get(action) + 1);
    }
    
    /**
     * Get a map identifying the actions taken by each agent in the game.
     *
     * The returned map is unmodifiable.
     */
    public Map<Steppable, PGGAction> getAgentActions() {
        return java.util.Collections.unmodifiableMap(agentActions);
    }
    
    /**
     * Calculate and return the payouts for each action.
     */
    protected Map<PGGAction, Double> getActionPayouts() {
        return getActionPayouts(this.actionCounts);
    }
    
    /**
     * Calculate and return the payouts based on the provided action counts map.
     *
     * This method allows the payout calculation code to be used to estimate
     * payouts.
     */
    public Map<PGGAction, Double> getActionPayouts(final Map<PGGAction, Integer> actionCounts) {
        Map<PGGAction, Double> payouts = new java.util.EnumMap<PGGAction, Double>(PGGAction.class);
            
        // -- find number of cooperators, defectors, noon-participants & punishers
        int cooperators = 0;
        if (actionCounts.containsKey(PGGAction.COOPERATE)) {
            cooperators = actionCounts.get(PGGAction.COOPERATE);
        }
        int defectors   = 0;
        if (actionCounts.containsKey(PGGAction.DEFECT)) {
            defectors = actionCounts.get(PGGAction.DEFECT);
        }
        int abstainers  = 0;
        if (actionCounts.containsKey(PGGAction.ABSTAIN)) {
            abstainers = actionCounts.get(PGGAction.ABSTAIN);
        }
        int punishers   = 0;
        if (actionCounts.containsKey(PGGAction.PUNISH)) {
            punishers = actionCounts.get(PGGAction.PUNISH);
        }
        
        // -- calculate total number of agents that contribute
        int contributors = cooperators + punishers;

        // -- calculate total participants in the game
        int total = cooperators + defectors + punishers;
        
        // -- calculate base payout
        double basePayout = (contributors == 1 && total == 1) ? SIGMA : (contributors*COST*R)/((double)total);
        if (PGGSim.VERBOSE || PGGame.VERBOSE) { System.out.println("base payout: " + basePayout); }

        // -- calculate payout for defectors
        if (defectors > 0) {
            double payout = basePayout - BETA*punishers;
            if (PGGSim.VERBOSE || PGGame.VERBOSE) { System.out.println(PGGAction.DEFECT + " payout: " + payout); }
            payouts.put(PGGAction.DEFECT, payout);
        }
        
        // -- calculate payout for cooperators
        if (cooperators > 0) {
            double payout = basePayout - COST;
            if (PGGSim.VERBOSE || PGGame.VERBOSE) { System.out.println(PGGAction.COOPERATE + " payout: " + payout); }
            payouts.put(PGGAction.COOPERATE, payout);
        }

        // -- calculate payout for punisher
        if (punishers > 0) {
            double payout = basePayout - COST - GAMMA*defectors;
            if (PGGSim.VERBOSE || PGGame.VERBOSE) { System.out.println(PGGAction.PUNISH + " payout: " + payout); }
            payouts.put(PGGAction.PUNISH, payout);
        }
        
        // -- calculate payouts for abstainers
        if (abstainers > 0) {
            if (PGGSim.VERBOSE || PGGame.VERBOSE) { System.out.println(PGGAction.ABSTAIN + " payout: " + SIGMA); }
            payouts.put(PGGAction.ABSTAIN, SIGMA);
        }
        
        // return results
        return payouts;
    }
    
    /**
     * Get a list of the agents that should be notified of the game results.
     */
    public Set<HistoryTrackingAgent> getAgentsToNotify() {
        Set<HistoryTrackingAgent> agentsToNotify = new java.util.HashSet<HistoryTrackingAgent>();
        Steppable agent;
        for (Map.Entry<Steppable, PGGAction> entry : agentActions.entrySet()) {
            agent = entry.getKey();
            if ( (agent instanceof HistoryTrackingAgent) &&
                 (entry.getValue() != PGGAction.ABSTAIN) ) {
                agentsToNotify.add((HistoryTrackingAgent)agent);
            }
        }
        
        return agentsToNotify;
    }
    
    /**
     * Clear the game to prepare for a new game.
     */
    private void clear() {
        // reset action counts back to zero
        for (PGGAction action : PGGAction.values()) {
            actionCounts.put(action, 0);
        }
        
        // clear out actions taken by previous participants
        agentActions.clear();
    }
}