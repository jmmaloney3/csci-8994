package pgg;

public class NonParticipant extends SimpleAgent {
    public NonParticipant() {
        // abstain from the public goods game
        super(PGGAction.ABSTAIN);
    }
}