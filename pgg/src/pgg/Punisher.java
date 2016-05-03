package pgg;

public class Punisher extends SimpleAgent {
    public Punisher() {
        // contribute to the public goods game & punish defectors
        super(PGGAction.PUNISH);
    }
}