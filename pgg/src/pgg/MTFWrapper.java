package pgg;

import ec.util.MersenneTwisterFast;

/**
 * A wrapper for ec.util.MersenneTwisterFast that allows it to replace
 * uses of java.util.Random.
 */
public class MTFWrapper extends java.util.Random {
    // private
    MersenneTwisterFast mtf;
    
    public MTFWrapper() {
        super();
        mtf = new MersenneTwisterFast();
    }
    
    public MTFWrapper(int seed) {
        super(seed);
        mtf = new MersenneTwisterFast(seed);
    }
    
    public MTFWrapper(MersenneTwisterFast mtf) {
        super();
        this.mtf = mtf;
    }
    
    public boolean nextBoolean() { return mtf.nextBoolean(); }
    
    public void nextBytes(byte[] bytes) { mtf.nextBytes(bytes); }
    
    public double nextDouble() { return mtf.nextDouble(); }
    
    public float nextFloat() { return mtf.nextFloat(); }
    
    public double nextGaussian() { return mtf.nextGaussian(); }
    
    public int nextInt() { return mtf.nextInt(); }
    
    public long nextLong() { return mtf.nextLong(); }
    
    public void setSeed(long seed) {
        // set seed gets called during java.util.Random constructor
        // before mtf has a value
        if (mtf == null) {
            super.setSeed(seed);
        }
        else {
            mtf.setSeed(seed);
        }
    }
}