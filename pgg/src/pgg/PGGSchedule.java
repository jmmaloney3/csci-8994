package pgg;

import sim.engine.SimState;
import sim.engine.Schedule;

/**
 * A variation of the sim.engine.Schedule class that supports public goods
 * games.  On each step, the schedule first calls a method on the simulation
 * engine to set up the next game.  After all the schedule agents have been
 * stepped for that time step, the schedule calls a method to finish the
 * game.
 */
public class PGGSchedule extends Schedule {
    
    /**
     * A variation of the Schedule.step method that does the following:
     * - call state.newGame()
     * - call super.step(state)
     * - call state.finishGame()
     */
    public boolean step(SimState state) {
        if (state instanceof PGGSim) {
            ((PGGSim)state).newGame();
        }
        boolean rval = super.step(state);
        if (state instanceof PGGSim) {
            ((PGGSim)state).finishGame();
        }
        return rval;
    }
}