package pgg;

/**
 * The actions that agents can take in a public goods game.
 */
public enum PGGAction {

    COOPERATE, DEFECT, ABSTAIN, PUNISH;

    public static final int length = PGGAction.values().length;

}

