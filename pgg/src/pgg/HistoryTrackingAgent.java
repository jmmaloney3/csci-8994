package pgg;

/**
 * An agent that tracks the history of the results of games in which
 * it has participated.
 */
public interface HistoryTrackingAgent {
    
    /**
     * Update the agent's history based on the results of the specified game.
     */
    public void updateHistory(PGGame game);
    
    /**
     * Clear the agent's game history.
     */
    public void clearHistory();
}