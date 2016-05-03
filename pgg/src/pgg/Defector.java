package pgg;

public class Defector extends SimpleAgent {
    public Defector() {
        // do not contribute to the public goods game
        super(PGGAction.DEFECT);
    }
}