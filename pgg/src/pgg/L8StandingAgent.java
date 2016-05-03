package pgg;

/**
 * A "leading eight" agent that whose reputation dynamics are consistent with
 * the "Standing" reputation assessment rule.  The three remaining reputation
 * dynamics properties are set as follows:
 *
 *   updateDynamics(ImageScore.GOOD, ImageScore.BAD, PGGAction.COOPERATE, ImageScore.GOOD);
 *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.COOPERATE, ImageScore.GOOD);
 *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.DEFECT,    ImageScore.BAD);
 *
 * See: (need reference for "Standing" assessment rule)
 */
public class L8StandingAgent extends L8Agent {
    /**
     * Create a leading eight whose unspecified reputation dynamics properties are set
     * to be consistent with the "Standing" reputation assessment rule:
     *
     *   updateDynamics(ImageScore.GOOD, ImageScore.BAD, PGGAction.COOPERATE, ImageScore.GOOD);
     *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.COOPERATE, ImageScore.GOOD);
     *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.DEFECT,    ImageScore.BAD);
     */
    public L8StandingAgent() {
        super();

        updateDynamics(ImageScore.GOOD, ImageScore.BAD, PGGAction.COOPERATE, ImageScore.GOOD);
        updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.COOPERATE, ImageScore.GOOD);
        updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.DEFECT,    ImageScore.BAD);
    }
}
