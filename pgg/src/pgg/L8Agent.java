package pgg;

/**
 * An agent that follows a version of one of the leading eight combinations
 * of reputation dynamics and strategies that has been extended to
 * the n-person public goods game.
 *
 * See:
 *  Ohtsuki, et. al., "The leading eight: Social norms that maintain cooperation
 *      by indirect reciprocity," Journal of Theoretical Biology, vol. 239,
 *      pp. 435-444, 2006.
 */
public abstract class L8Agent extends ImageScoreAgent {
    /**
     * Create an agent that is preconfigured with the strategy and reputation
     * dynamics properties that are shared by all leading eight (d, p) pairs.
     *
     * Subclasses must fill in the following properties:
     *
     *   updateDynamics(ImageScore.GOOD, ImageScore.BAD, PGGAction.COOPERATE, <rep>);
     *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.COOPERATE, <rep>);
     *   updateDynamics(ImageScore.BAD,  ImageScore.BAD, PGGAction.DEFECT,    <rep>);
     *
     * The settings for the last two dynamics properties determines the setting
     * for the following strategy property:
     *
     *   updateStrategy(ImageScore.BAD, ImageScore.BAD, <action>);
     */
    public L8Agent() {
        super();

        // (1) maintenance of cooperation
        updateStrategy(ImageScore.GOOD, ImageScore.GOOD, PGGAction.COOPERATE);
        updateDynamics(ImageScore.GOOD, ImageScore.GOOD, PGGAction.COOPERATE, ImageScore.GOOD);
        
        // (2) identification of defectors
        updateDynamics(ImageScore.GOOD, ImageScore.GOOD, PGGAction.DEFECT, ImageScore.BAD);
        updateDynamics(ImageScore.BAD,  ImageScore.GOOD, PGGAction.DEFECT, ImageScore.BAD);
        
        // (3) punishment and justification of punishment
        updateStrategy(ImageScore.GOOD, ImageScore.BAD, PGGAction.DEFECT);
        updateDynamics(ImageScore.GOOD, ImageScore.BAD, PGGAction.DEFECT, ImageScore.GOOD);
        
        // (4) apology and forgiveness
        updateStrategy(ImageScore.BAD, ImageScore.GOOD, PGGAction.COOPERATE);
        updateDynamics(ImageScore.BAD, ImageScore.GOOD, PGGAction.COOPERATE, ImageScore.GOOD);
    }
    
    /**
     * Update this agent's reputation dynamics (aka, reputation assessment rule).
     *
     * For a "leading eight" agent, the setting of p[BAD][BAD} is determined by the
     * settings for d[BAD][BAD][COOPERATE] and d[BAD][BAD][DEFECT].  Ensure that the
     * value for p[BAD][BAD] is set appropriate after the dynamics are updated.
     */
    public void updateDynamics(ImageScore rep, ImageScore grpRep, PGGAction act, ImageScore newRep) {
        super.updateDynamics(rep, grpRep, act, newRep);
        
        // set value of p[BAD][BAD] appropriately based on the dynamics
        if (ImageScore.BAD.equals(rep) && ImageScore.BAD.equals(grpRep)) {
            if (ImageScore.GOOD.equals(applyDynamics(rep, grpRep, PGGAction.COOPERATE)) &&
                ImageScore.BAD.equals(applyDynamics(rep, grpRep, PGGAction.DEFECT)) ) {
                updateStrategy(rep, grpRep, PGGAction.COOPERATE);
            }
            else {
                updateStrategy(rep, grpRep, PGGAction.DEFECT);
            }
        }
    }
}