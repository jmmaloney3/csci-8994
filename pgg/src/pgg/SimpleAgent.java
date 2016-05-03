package pgg;

import sim.engine.SimState;
import sim.engine.Steppable;

/**
 * A simple agent that plays the same action in every game.
 *
 * The simulator relies on uniqu classes existing for each strategy.
 * Therefore, a subclass of this class should be created for a
 * strategy.
 */
public abstract class SimpleAgent implements Steppable {
    protected PGGAction action;
    
    public SimpleAgent(PGGAction act) {
        this.action = act;
    }
    
    public void step(SimState state) {
        if (state instanceof PGGSim) {
            ((PGGSim)state).takeAction(this, this.action);
        }
    }
}