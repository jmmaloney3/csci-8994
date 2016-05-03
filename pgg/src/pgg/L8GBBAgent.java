package pgg;

/**
 * A "leading eight" agent that whose three remaining reputation dynamics
 * properties are set as follows:
 *
 *   updateDynamics(ImageScore.GOOD, ImageScore.BAD, PGGAction.COOPERATE, ImageScore.GOOD);
 *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.COOPERATE, ImageScore.BAD);
 *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.DEFECT,    ImageScore.BAD);
 */
public class L8GBBAgent extends L8Agent {
    /**
     * Create a leading eight whose unspecified reputation dynamics properties are set
     * as follows:
     *
     *   updateDynamics(ImageScore.GOOD, ImageScore.BAD, PGGAction.COOPERATE, ImageScore.GOOD);
     *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.COOPERATE, ImageScore.BAD);
     *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.DEFECT,    ImageScore.BAD);
     */
    public L8GBBAgent() {
        super();

        updateDynamics(ImageScore.GOOD, ImageScore.BAD, PGGAction.COOPERATE, ImageScore.GOOD);
        updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.COOPERATE, ImageScore.BAD);
        updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.DEFECT,    ImageScore.BAD);
    }
}
