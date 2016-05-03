package pgg;

/**
 * An agent that tracks the history of all games that have been played
 * regardless of whether the agent particpanted in the game or not.
 *
 * This is a tagging interface used to identify those agents that should
 * receive the results of all games.
 */
public interface AllSeeingAgent extends HistoryTrackingAgent {

}