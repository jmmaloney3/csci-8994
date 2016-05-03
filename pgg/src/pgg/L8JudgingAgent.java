package pgg;

/**
 * A "leading eight" agent that whose reputation dynamics are consistent with
 * the "Judging" reputation assessment rule.  The three remaining reputation
 * dynamics properties are set as follows:
 *
 *   updateDynamics(ImageScore.GOOD, ImageScore.BAD, PGGAction.COOPERATE, ImageScore.BAD);
 *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.COOPERATE, ImageScore.BAD);
 *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.DEFECT,    ImageScore.BAD);
 *
 * See: (need reference for "Judging" assessment rule)
 */
public class L8JudgingAgent extends L8Agent {
    /**
     * Create a leading eight whose unspecified reputation dynamics properties are set
     * to be consistent with the "Standing" reputation assessment rule:
     *
     *   updateDynamics(ImageScore.GOOD, ImageScore.BAD, PGGAction.COOPERATE, ImageScore.BAD);
     *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.COOPERATE, ImageScore.BAD);
     *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.DEFECT,    ImageScore.BAD);
     */
    public L8JudgingAgent() {
        super();

        updateDynamics(ImageScore.GOOD, ImageScore.BAD, PGGAction.COOPERATE, ImageScore.BAD);
        updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.COOPERATE, ImageScore.BAD);
        updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.DEFECT,    ImageScore.BAD);
        
        // set the threshold to 20%
        setBadThreshold(0.0);
    }
}
