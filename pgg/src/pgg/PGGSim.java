package pgg;

import java.io.File;
import java.io.PrintStream;
import java.util.Collection;
import java.util.List;
import java.util.Map;
import java.util.Set;

import sim.engine.SimState;
import sim.engine.Steppable;
import sim.util.Bag;

/**
 * Generic simulation model for public goods games.
 */
public class PGGSim extends SimState {
    // constants
    public static final double TOL = 0.001;
    public static final String SCOUNTSFN  = "scounts.csv";
    public static final String SFITNESSFN = "sfitness.csv";
    public static final String SPAYOUTSFN = "spayouts.csv";
    
    // public goods game
    final protected double cost = 1.0; // the cost to invest in public goods game
    final protected double r = 3.0; // factor by which contributions are increased
    final protected double beta = 1.0; // punishment cost
    final protected double gamma = 0.3; // cost to impose punishment
    final protected double sigma = 1; // non-participant payoff
    protected double       selectionStrength = 0.249; // selection strength used in fitness calc
    final PGGame GAME = new PGGame(cost, r, sigma, beta, gamma);
    
    // file for writing strategy counts and fitness after each round
    File outDir;
    PrintStream scStream;
    PrintStream sfStream;
    PrintStream spStream;
    
    // wraper for random number generator
    final public MTFWrapper randomWrapper;

    // verbose setting
    static boolean VERBOSE = false;
    
    // agent lists
    // -- current agent population
    protected Steppable[] agents;
    // -- sortable and shuffle-able list of agent IDs
    protected List<Integer> agentIDs = new java.util.ArrayList<Integer>();
    // -- all seeing agents (need to be provided results from all games)
    protected List<AllSeeingAgent> allSeeingAgents = new java.util.LinkedList<AllSeeingAgent>();
    // -- history tracking agents - includes al seeing agents (history needs to be cleared after each round)
    protected List<HistoryTrackingAgent> historyTrackingAgents = new java.util.LinkedList<HistoryTrackingAgent>();

    // simulation variables
    protected int numRounds = 10000;
    protected int roundsCompleted = 0;
    protected double mu = 0.001; // the mutation rate
    protected boolean resetPayoffsAfterRound = true;
    protected int numAgents = 100; // M
    // reproduction probability limits for the strategies (use LinkedHashMap to guarantee
    // order of keys) - this map lists the limits that are used to determine which
    // strategy has been selected for reproduction
    protected Map<Class<? extends Steppable>, Double>  birthProbs
        = new java.util.LinkedHashMap<Class<? extends Steppable>, Double>();
    // death probability limits for the strategies (use LinkedHashMap to guarantee
    // order of keys) - this map lists the limits that are used to determine which
    // strategy has been selected to be replaced
    protected Map<Class<? extends Steppable>, Double>  deathProbs
        = new java.util.LinkedHashMap<Class<? extends Steppable>, Double>();
    // strategy population counts (use LinkedHashMap to guarantee order of keys)
    // ensure the counts are output to file in same order each time
    protected Map<Class<? extends Steppable>, Integer>  stratCounts
        = new java.util.TreeMap<Class<? extends Steppable>, Integer>(new ClassComparator());
    // = new java.util.LinkedHashMap<Class<? extends Steppable>, Integer>();

    // round variables
    protected int gamesPerRound = 100;
    protected int roundGamesCompleted = 0; // number of games completed in the current round
    // the average payout per strategy over all games in the round
    protected Map<Class<? extends Steppable>, Double>  avgPayoffs;
    // the total games participated per strategy over all games in the round
    protected Map<Class<? extends Steppable>, Integer> totalGames;

    // game variables
    protected int pggSize = 5; // N
    
    /**
     * Construct a public goods game simulation.
     */
    public PGGSim(long seed, Map<Class<? extends Steppable>, Double> stratProps, int numAgents, int pggSize,
                  int gamesPerRound, int numRounds, double ss, File oDir, boolean resetPayoffs, boolean verbose) {
        super(seed, new PGGSchedule());
        // set up simulation parameters
        this.numAgents = numAgents;
        this.pggSize = pggSize;
        this.gamesPerRound = gamesPerRound;
        this.numRounds = numRounds;
        this.outDir = oDir;
        this.resetPayoffsAfterRound = resetPayoffs;
        this.selectionStrength = ss;
        PGGSim.VERBOSE = verbose;
        
        // initialize birth and death pribability thresholds
        updateEvolutionThresholds(stratProps);
        if (VERBOSE) {
            System.out.println("Initial birth probability thresholds: " + birthProbs);
            System.out.println("Initial death probability thresholds: " + deathProbs);
        }
        
        // initialize the average payoff and total games maps
        initializePayoffMaps();
        
        // create wrapper around random number generator
        randomWrapper = new MTFWrapper(this.random);
    }

    /**
     * Initialize and start the simulation.
     */
    public void start() {
        super.start();
        
        // initialize list of agent IDs
        agents = new Steppable[numAgents];
        for (int i = 0; i < numAgents; i++) {
            this.agentIDs.add(i);
        }
        
        // create the agents
        createAgents();
        
        // open files to record strategy counts and fitness at the end of each round
        if (outDir != null) {
            try {
                // strategy counts file
                File sCountsFile = new File(outDir, SCOUNTSFN);
                scStream = new PrintStream(new java.io.FileOutputStream(sCountsFile), true);
                
                // strategy fitness file
                File sFitnessFile = new File(outDir, SFITNESSFN);
                sfStream = new PrintStream(new java.io.FileOutputStream(sFitnessFile), true);

                // strategy average payouts file
                File sPayoutsFile = new File(outDir, SPAYOUTSFN);
                spStream = new PrintStream(new java.io.FileOutputStream(sPayoutsFile), true);
}
            catch (java.io.FileNotFoundException fnfe) {
                throw new RuntimeException(fnfe);
            }
        
            // write headers to strategy counts and fitness files
            writeHeaders(stratCounts, scStream);
            writeHeaders(stratCounts, sfStream);
            writeHeaders(stratCounts, spStream);
        }
    }
    
    /**
     * Quit the simulation.
     */
    public void finish() {
        scStream.close();
        sfStream.close();
        spStream.close();
        super.finish();
    }

    /**
     * Select the participants for a new game and add them to the schedule.
     */
    protected void newGame() {
        // check to see if simulation is complete
        if (roundsCompleted == numRounds) {
            // simulation is complete - don't start another game
            return;
        }
        
        if (VERBOSE) { System.out.println("New Game"); }
        
        // create a random sample of agents
        Steppable[] participants = this.selectParticipants();
        
        // initialize new game
        this.GAME.newGame(participants, roundsCompleted, roundGamesCompleted);
        
        // add the selected agents to the schedule
        double time = this.schedule.getTime();
        if (time < schedule.EPOCH) { time = schedule.EPOCH; }
        for (int i = 0; i < this.pggSize; i++) {
            this.schedule.scheduleOnce(time, 1, participants[i]);
        }
    }
    
    /**
     * Take an action for an agent.  This method is used by game participants
     * to regiter their action with the simulation engine.
     */
    protected void takeAction(final Steppable agent, final PGGAction action) {
        this.GAME.takeAction(agent, action);
    }

    /**
     * Finish the game by distributing the payouts.
     */
    protected void finishGame() {
        // -- update statistics for each strategy represented in the current game
        updateStrategyStats(this.GAME);
        
        // push results to history tracking agents
        Set<HistoryTrackingAgent> agentsToNotify = new java.util.HashSet<HistoryTrackingAgent>();
        agentsToNotify.addAll(this.allSeeingAgents);
        agentsToNotify.addAll(this.GAME.getAgentsToNotify());
        for (HistoryTrackingAgent agent : agentsToNotify) {
            agent.updateHistory(this.GAME);
        }
        
        // print out current totals for the strategies
        if (VERBOSE) {
            for (Map.Entry<Class<? extends Steppable>, Double> entry : avgPayoffs.entrySet()) {
                System.out.println(entry.getKey() + " avg payoff => " + entry.getValue());
            }
            for (Map.Entry<Class<? extends Steppable>, Integer> entry : totalGames.entrySet()) {
                System.out.println(entry.getKey() + " tot games  => " + entry.getValue());
            }
            System.out.println("Finish Game");
        }

        // update population if this is the last game of the round
        roundGamesCompleted += 1; // increment game number
        
        if (roundGamesCompleted > gamesPerRound) {
            finishRound();
        }
        
    }
    
    /**
     * Finish the round and set up the next round.
     */
    protected void finishRound() {
        if (VERBOSE) { System.out.println("Finish Round"); }
        
        // write strategy counts to file
        writeValues(stratCounts, scStream);
        
        // write strategy average payouts to file
        writeValues(avgPayoffs, spStream);
        
        // update birth and death probability thresholds
        Map<Class<? extends Steppable>, Double> fitnessValues = calculateStrategyFitness();
        updateEvolutionThresholds(fitnessValues);
        
        // write fitness value to file
        writeValues(fitnessValues, sfStream);

        // update agent population
        //for (int i=0; i<100; i++) {
            evolvePopulation();
        //}
        
        // reset payouts and games played count
        if (resetPayoffsAfterRound) {
            // reset strategy statistics
            initializePayoffMaps();

            // reset agent histories
            for (HistoryTrackingAgent agent : this.historyTrackingAgents) {
                agent.clearHistory();
            }
        }
        
        // reset number of games played in the current round
        roundGamesCompleted = 0;
        
        // increment number of rounds played
        roundsCompleted++;
    }
    
    /**
     * Select participants for a new game.
     */
    protected Steppable[] selectParticipants() {
        // sort the list of IDs to ensure we have the same starting point for each shuffle
        java.util.Collections.sort(agentIDs);
        
        // shuffle the list of agent IDs using the MersenneTwisterFast random number generator
        java.util.Collections.shuffle(agentIDs, randomWrapper);
        
        // create and return list of participants
        Steppable[] participants = new Steppable[pggSize];
        for (int i = 0; i < pggSize; i++) {
            participants[i] = this.agents[agentIDs.get(i)];
        }
        return participants;
    }

    /**
     * Update the statistics (average payout and total games) for each strategy based on
     * the results of the specified game.
     */
    protected void updateStrategyStats(PGGame game) {
        Map<Steppable, PGGAction> agentActions = game.getAgentActions();
        Map<PGGAction, Double> actionPayouts   = game.getActionPayouts();

        Class<? extends Steppable> strategy;
        double actionPayout;
        double avg;
        int games;
        for (Map.Entry<Steppable, PGGAction> entry : agentActions.entrySet()) {
            // get the agent's strategy
            strategy = entry.getKey().getClass();
            
            // get the agent's payout
            actionPayout = actionPayouts.get(entry.getValue());

            // update total games participated in for the strategy
            if (totalGames.containsKey(strategy)) { games = totalGames.get(strategy); }
            else { games = 0; }
            games++;
            totalGames.put(strategy, games);

            // update average payoff for the strategy
            if (avgPayoffs.containsKey(strategy)) { avg = avgPayoffs.get(strategy); }
            else { avg = 0.0; }
            avg = ((avg * (games - 1)) + actionPayout)/((double)games);
            avgPayoffs.put(strategy, avg);
            
        }
    }
    
    /**
     * Evolve the population.
     *
     * The population evolves as follows (based on Nowak 2004 p.647):
     * 
     *    "At each time step, an individual is chosen for reproduction proportional
     *     to its fitness.  One identical offspring is being produced that replaces
     *     another randomly chosen individual.  Thus N is strictly constant."
     *
     * The fitness of an individual is calculated as follows (based on Nowak 2004 &
     * Hauert 2007):
     *
     *    "We define each players' fitness as 1 - s + sP, the convex combination of
     *     the 'baseline fitness', which is normalized to 1 for all players, and the
     *     payoff P from the optional public goods game with punishment."
     *
     * Let B be the baseine fitnes and P be the payoff received.  Then the equation is:
     *
     *    (1 - s)B + sP
     *
     * Since the baseline fitness is 1 for all individuals, this simplifies to:
     *
     *    (1 - s)B + sP = (1 - s)*1 + sP = 1 - s + sP
     *
     * The payoff P for a strategy is the expected payoff that an agent following that
     * strategy will receive in one instance of a public goods gaem.
     */
    protected void evolvePopulation() {
        // reproduce
        // SHOULD THIS SELECT AN INDIVIDUAL RATHER THAN A STRATEGY?
        Class<? extends Steppable> newStrategy = selectRandomStrategy(birthProbs);
        Steppable newAgent = createAgentFromStrategy(newStrategy);
        
        // identify agent to replace
        Class<? extends Steppable> oldStrategy = selectRandomStrategy(deathProbs);
        int i = 0;
        for (; i < numAgents; i++) {
            if (oldStrategy.isInstance(agents[i])) {
                break;
            }
        }
        
        // if no agent is following oldStrategy then randomly select an agent
        if (i >= numAgents) {
            if (VERBOSE) { System.out.println("no agent is following old strategy: " + oldStrategy); }
            i = this.random.nextInt(numAgents);
            oldStrategy = agents[i].getClass();
        }
        
        // remove old agent from history tracking lists
        if (agents[i] instanceof HistoryTrackingAgent) {
            this.historyTrackingAgents.remove((HistoryTrackingAgent)agents[i]);
            if (agents[i] instanceof AllSeeingAgent) {
                this.allSeeingAgents.remove(agents[i]);
            }
        }
        
        // replace the selected agant
        agents[i] = newAgent;
        
        // add new agent to history tracking lists
        if (agents[i] instanceof HistoryTrackingAgent) {
            this.historyTrackingAgents.add((HistoryTrackingAgent)agents[i]);
            if (agents[i] instanceof AllSeeingAgent) {
                this.allSeeingAgents.add((AllSeeingAgent)agents[i]);
            }
        }
        
        // update the count for new strategy
        int c = 0;
        if (stratCounts.containsKey(newStrategy)) {
            c = stratCounts.get(newStrategy);
        }
        stratCounts.put(newStrategy, ++c);
        
        // update count for old strategy
        c = 0;
        if (stratCounts.containsKey(oldStrategy)) {
            c = stratCounts.get(oldStrategy);
        }
        stratCounts.put(oldStrategy, --c);
    }
    
    /**
     * Calculate fitness measures for the strategies.
     *
     * The fitness of a strategy is calculated as follows:
     *
     *    T[i]: the total rewards earned by individuals following the strategy i
     *    G[i]: the total number of games played by individuals following the strategy i
     *    P[i]: the average payout earned in PGGs by individuals following strategy i
     *
     * The expected payout that an agent following strategy i will receive in one
     * instance of a public goods game is just the average payout received by all
     * individuals that followed strategy i in the previous round:
     *
     *    P[i] = T[i]/G[i]
     *
     * The fitness of strategy i is:
     *
     *    F[i] = 1 - s + sP[i] = 1 - s + sT[i]/G[i]
     *
     * The probability that an individual following strategy i is selected to
     * to reproduce is the normalized fitness values:
     *
     *    F[i]/(sum F[j] for all strategies j)
     */
    protected Map<Class<? extends Steppable>, Double> calculateStrategyFitness() {
        // map to hold the current fitness for each strategy
        Map<Class<? extends Steppable>, Double>  fitnessValues
            = new java.util.HashMap<Class<? extends Steppable>, Double>();

        // calculate the total fitness for all strategies (used for normalization)
        double fitness, totalFitness = 0.0;
        for (Class<? extends Steppable> strategy : birthProbs.keySet()) {
            if (avgPayoffs.get(strategy) < TOL) {
                if (VERBOSE) { System.out.println(strategy + " payout is zero"); }
            }
            fitness = (1 - selectionStrength) + selectionStrength*avgPayoffs.get(strategy);
            
            // magnify performance difference between strategies
            //fitness = Math.pow(fitness, 10);
            
            if (fitness < TOL) {
                fitness = 0.0;
            }
            if (VERBOSE) { System.out.println(strategy + " fitness: " + fitness); }
            fitnessValues.put(strategy, fitness);
            totalFitness +=  fitness;
        }
        if (VERBOSE) { System.out.println("total fitness: " + totalFitness); }
        
        // if total fitness <= zero then fitness values are invalid - return null
        if (totalFitness < TOL) { return null; }
        
        // calculate normalized fitness - aka, reproduction probabilities - and set the
        // limits that will be used for randomly selecting a strategy to reproduce
        for (Map.Entry<Class<? extends Steppable>, Double> entry : fitnessValues.entrySet()) {
            entry.setValue(entry.getValue()/totalFitness);
            if (VERBOSE) { System.out.println(entry.getKey() + " normalized fitness: " + entry.getValue()); }
        }
        
        return fitnessValues;
    }
    
    /**
     * Calculate the probability thresholds used to determine which agents will
     * reproduce and which will be removed from the population.
     *
     * The stratProps parameter is either the proportions of the population to be
     * allocated to each strategy (for simulation start up) or the normalized strategy
     * fitness measures (for mid-simulation updates).  The values in stratProps are
     * assumed to be normalized (sum to 1.0).
     *
     * If stratProbs is null then the birth and death probability thresholds are left unchanged.
     */
    public void updateEvolutionThresholds(Map<Class<? extends Steppable>, Double> stratProps) {
        // if stratProps is null then don't update the thresholds
        if (stratProps == null) {
            System.out.println("birth/death probabilty thresholds will not be updated: " + stratProps);
            return;
        }
        
        // validate that at least one threshold is greater than zero
        boolean invalid = true;
        for (Double p : stratProps.values()) {
            if (p > TOL) {
                invalid = false;
                break;
            }
        }
        if (invalid) {
            System.out.println("invalid birth/death probabilty thresholds: " + stratProps);
            return;
        }

        // clear out old probability thresholds
        birthProbs.clear();
        deathProbs.clear();

        // create a list of the strategies with proportions/fitness measures
        List<Map.Entry<Class<? extends Steppable>, Double>> entries
            = new java.util.ArrayList<Map.Entry<Class<? extends Steppable>, Double>>();
        entries.addAll(stratProps.entrySet());
        int numStrategies = entries.size();
        
        // list ot hold strategies ordered by value ascending
        List<Class<? extends Steppable>> orderedStrats
            = new java.util.ArrayList<Class<? extends Steppable>>();
        
        // populate thresholds for birth probabilities
        double p;
        double limit = 0.0D;
        int i = 1;
        java.util.Collections.sort(entries, new SortMapEntriesByValue(true));
        for (Map.Entry<Class<? extends Steppable>, Double> entry : entries) {
            // preserve strategy order for later use
            orderedStrats.add(entry.getKey());

            // get probability
            p = entry.getValue();
            
            // don't add threshold for zero probabilty strategy
            //if (p < TOL) { i++; continue; }

            // align last threshold with 1.0D
            if (i == numStrategies) { limit = 1.0D; }
            else { limit += p; }
            
            // update birth probability threshold
            birthProbs.put(entry.getKey(), limit);
            i++;
        }
        
        // populate thresholds for death probabilities
        limit = 0.0D;
        i = 1;
        java.util.Collections.sort(entries, new SortMapEntriesByValue(false));
        for (Map.Entry<Class<? extends Steppable>, Double> entry : entries) {
            // get probability
            p = entry.getValue();
            
            // don't add threshold for zero probabilty strategy
            //if (p < TOL) { i++; continue; }
            
            // align last threshold with 1.0D
            if (i == numStrategies) { limit = 1.0D; }
            else { limit += p; }

            // add the values in reverse order - but keep strategies in same order (orderedStrats)
            deathProbs.put(orderedStrats.get(i-1), limit);
            i++;
        }
    }
    
    /**
     * Create a population of agents based on the proportions specified
     * in the strategy reproduction probability limits map.
     */
    protected void createAgents() {
        // initialize strat counts
        for (Class<? extends Steppable> strategy : birthProbs.keySet()) {
            stratCounts.put(strategy, 0);
        }
        
        // initialize the agents, agent IDs and srategy proportions
        agents = new Steppable[numAgents];
        Class<? extends Steppable> strategy;
        int count;
        for (int i = 0; i < numAgents; i++) {
            // create new agent from randomly selected strategy
            strategy = selectRandomStrategy(birthProbs);
            agents[i] = createAgentFromStrategy(strategy);
            
            // add agent to history tracking lists if necessary
            if (agents[i] instanceof HistoryTrackingAgent) {
                this.historyTrackingAgents.add((HistoryTrackingAgent)agents[i]);
                if (agents[i] instanceof AllSeeingAgent) {
                    this.allSeeingAgents.add((AllSeeingAgent)agents[i]);
                }
            }

            // update strategy counts
            count = stratCounts.get(strategy);
            count++;
            stratCounts.put(strategy, count);
        }
    }
    
    /**
     * Select a random strategy ignoring the mutation rate.
     */
    protected Class<? extends Steppable>
        selectRandomStrategyNoMutation(Map<Class<? extends Steppable>, Double> probThresholds) {
        double p = this.random.nextDouble();
        Class<? extends Steppable> strategy = null;
        for (Map.Entry<Class<? extends Steppable>, Double> entry : probThresholds.entrySet()) {
            strategy = entry.getKey();
            if (p <= entry.getValue()) {
                return strategy;
            }
        }
        
        // the limits are not set up correctly - return highest probabilty strategy
        System.out.println("invalid probability thresholds: " + probThresholds);
        return strategy;
    }
    
    /**
     * Select a random strategy based on the reproduction probability limits.  Incorporates
     * specified mutation rate as well.
     */
    protected Class<? extends Steppable> selectRandomStrategy(Map<Class<? extends Steppable>, Double> probThresholds) {
        Class<? extends Steppable> strategy = null;
        if (this.random.nextDouble() > mu) { // no mutation
            strategy = selectRandomStrategyNoMutation(probThresholds);
        }
        else { // mutation occured
            if (VERBOSE) { System.out.println("=========== MUTATION OCCURED ==============="); }
            int numStrategies = probThresholds.keySet().size();
            int k = this.random.nextInt(numStrategies);
            if (VERBOSE) { System.out.println("  Mutation selected: " + k); }
            int i = 0;
            for (Class<? extends Steppable> s : probThresholds.keySet()) {
                if (VERBOSE) { System.out.println("  Is " + i + " the mutation?"); }
                if (i == k) {
                    strategy = s;
                    if (VERBOSE) { System.out.println("Mutated to: " + s); }
                    break;
                }
                i++;
            }
        }
        if (VERBOSE) { System.out.println("Selected: " + strategy); }
        return strategy;
    }
    
    /**
     * Create a new agent that follows the specified strategy.
     */
    protected Steppable createAgentFromStrategy(Class<? extends Steppable> strategy) {
        // create new agent
        try {
            java.lang.reflect.Constructor<? extends Steppable> ctor = strategy.getConstructor();
            return ctor.newInstance();
        }
        catch (Exception e) {
            throw new RuntimeException(e);
        }
    }
    
    /*
     * Write the keys in the map to file
     */
    protected void writeHeaders(Map<Class<? extends Steppable>, ?> map, java.io.PrintStream ps) {
        // create new line for file
        String line = "round";
        // add column headers to line
        for (Class<? extends Steppable> strategy : map.keySet()) {
            // append strategy name to current line
            line = String.format("%s, %s", line, strategy.getSimpleName());
        }
        // write line to file
        ps.println(line);
    }

    /**
     * Write the values in the map to file.
     */
    protected void writeValues(Map<Class<? extends Steppable>, ?> map, java.io.PrintStream ps) {
        // write the round number to file
        ps.print(roundsCompleted);
        if (map != null) {
            // write the values to file
            for (Object value : map.values()) {
                ps.print(", ");
                ps.print(value);
            }
        }
        else {
            // just output zero for each strategy
            for (int i=0; i<stratCounts.size(); i++) {
                ps.print(", ");
                ps.print(0);
            }
        }
        // write new line
        ps.println();
    }
    
    /**
     * Get the map of strategy average payouts.
     *
     * The returned map is unmodifiable.
     */
    public Map<Class<? extends Steppable>, Double> getStrategyAvgPayouts() {
        return java.util.Collections.unmodifiableMap(avgPayoffs);
    }
    
    /**
     * Get the map of strategy total games played.
     *
     * The returned map is unmodifiable.
     */
    public Map<Class<? extends Steppable>, Integer> getStrategyTotalGames() {
        return java.util.Collections.unmodifiableMap(totalGames);
    }
    
    /**
     * Get the map of reproduction probability thresholds.
     *
     * The returned map is unmodifiable.
     */
    public Map<Class<? extends Steppable>, Double> getStrategyBirthProbs() {
        return java.util.Collections.unmodifiableMap(birthProbs);
    }

    /**
     * Get the map of death probability thresholds.
     *
     * The returned map is unmodifiable.
     */
    public Map<Class<? extends Steppable>, Double> getStrategyDeathProbs() {
        return java.util.Collections.unmodifiableMap(deathProbs);
    }
    
    /**
     * Initialize the average payoff and total games maps.
     */
    protected void initializePayoffMaps() {
        avgPayoffs = new java.util.HashMap<Class<? extends Steppable>, Double>();
        totalGames = new java.util.HashMap<Class<? extends Steppable>, Integer>();
        for (Class<? extends Steppable> strategy : birthProbs.keySet()) {
            avgPayoffs.put(strategy, 0.0);
            totalGames.put(strategy, 0);
        }
    }
    
    /**
     * Comparator class used by stratCounts TreeMap
     */
    static class ClassComparator implements java.util.Comparator<Class<? extends Steppable>> {
        public int compare(Class<? extends Steppable> o1, Class<? extends Steppable> o2) {
            if (o1 == null) {
                if (o2 == null) {
                    return 0;
                }
                else {
                    // null is less than everything else
                    return -1;
                }
            }
            else if (o2 == null) {
                // already know that o1 is not null
                // return 1 since null is less than everything else
                return 1;
            }
            else {
                // neither argument is null - compare their names
                return o1.getName().compareTo(o2.getName());
            }
        }
        public boolean equals(Object obj) {
            if (obj instanceof ClassComparator) {
                return true;
            }
            else {
                return false;
            }
        }
    }
    
    /**
     * Comparator to sort a collection of Map.Entry object by value.
     */
    static class SortMapEntriesByValue implements java.util.Comparator<Map.Entry<Class<? extends Steppable>, Double>> {
        private boolean asc = true;
        public SortMapEntriesByValue() {
            this.asc = true;
        }
        public SortMapEntriesByValue(boolean asc) {
            this.asc = asc;
        }
        public int compare(Map.Entry<Class<? extends Steppable>, Double> e1,
                           Map.Entry<Class<? extends Steppable>, Double> e2) {
            if (e1.getValue() == null) {
                if (e2.getValue() == null) {
                    // null is equal to null
                    return 0;
                }
                else {
                    // null is less than everything else
                    return asc ? -1 : 1;
                }
            }
            else if (e2.getValue() == null) {
                // already know that e1 is not null
                // return 1 since null is less than everything else
                return asc ? 1 : -1;
            }
            else {
                // neither argument is null - compare them to each other
                if (asc) {
                    return e1.getValue().compareTo(e2.getValue());
                }
                else {
                    return e2.getValue().compareTo(e1.getValue());
                }
            }
        }
    }
    
    /**
     * A simulator factory class that allows command line arguments to be passed into
     * the simulator.
     */
    static class MakePGGSim implements sim.engine.MakesSimState {
        public SimState newInstance(long seed, String[] args) {
            // set up the argument parser
            net.sourceforge.argparse4j.inf.ArgumentParser parser =
                net.sourceforge.argparse4j.ArgumentParsers.newArgumentParser("PGGSim")
                    .description("A simulator for public goods games.")
                    .defaultHelp(true)
                    .version("${prog} 0.1");
            // version
            parser.addArgument("--version")
                .action(net.sourceforge.argparse4j.impl.Arguments.version());
            // version
            parser.addArgument("--verbose", "-v")
                .dest("verbose")
                .type(Boolean.class)
                .setDefault(false)
                .action(net.sourceforge.argparse4j.impl.Arguments.storeTrue())
                .help("run in verbose mode");
            // number of agents in the population
            parser.addArgument("-M")
                .metavar("numAgents")
                .type(Integer.class)
                .setDefault(100)
                .help("size of the agent population");
            // number of agents in a public goods game
            parser.addArgument("-N")
                .metavar("pggSize")
                .type(Integer.class)
                .setDefault(5)
                .help("size of public goods game");
            // gams per round
            parser.addArgument("-g")
                .metavar("games")
                .type(Integer.class)
                .setDefault(100)
                .help("number of games per round");
            // number of rounds
            parser.addArgument("-r")
                .metavar("rounds")
                .type(Integer.class)
                .setDefault(10000)
                .help("number of round");
            // output file
            parser.addArgument("-d")
                .metavar("directory")
                .type(net.sourceforge.argparse4j.impl.Arguments.fileType()
                      .verifyIsDirectory()
                      .verifyCanCreate())
                .setDefault(new java.io.File("."))
                .help("directory path for output files");
            // by default, run the simulator with the cooperator and defector strategies
            List<String> defaultSNames = new java.util.ArrayList<String>();
            defaultSNames.add(pgg.Cooperator.class.getName());
            defaultSNames.add(pgg.Defector.class.getName());
            // strategies
            parser.addArgument("-s")
                .metavar("strategy")
                .nargs("+")
                .type(String.class)
                .setDefault(defaultSNames)
                .help("list of strategies to include in simulation");
            // strategy proportions
            parser.addArgument("-p")
                .metavar("percent")
                .nargs("+")
                .type(Double.class)
                .help("initial percentage of population following each strategy");
            parser.addArgument("-ss")
                .metavar("strength")
                .type(Double.class)
                .setDefault(0.249)
                .help("selection strength measuring relative importance of base fitness " +
                      "and average payout in fitness calculation");
            // don't reset payoffs after each round
            parser.addArgument("-all")
                .type(Boolean.class)
                .setDefault(false)
                .action(net.sourceforge.argparse4j.impl.Arguments.storeTrue())
                .help("don't reset payoffs after each round");
            
            // parse the arguments
            net.sourceforge.argparse4j.inf.Namespace res = null;
            try {
                res = parser.parseArgs(args);
                System.out.println(res);
            }
            catch (net.sourceforge.argparse4j.internal.HelpScreenException ex) {
                System.exit(0);
            }
            catch (net.sourceforge.argparse4j.inf.ArgumentParserException ex) {
                System.out.println(ex);
                System.exit(1);
            }
            
            // read out the args to use when creating the simulation
            int numAgents = res.get("M");
            int pggSize = res.get("N");
            
            if (numAgents < pggSize) {
                numAgents = pggSize;
            }
            
            int gamesPerRound = res.get("g");
            int numRounds = res.get("r");
            double selectionStrength = res.get("ss");
            File oDir = res.get("d");
            boolean all = res.get("all");
            boolean verbose = res.get("verbose");
            List<String> sNames = res.get("s");
            List<Double> sProps = res.get("p");
            
            // fill in default proportions if not provided
            if (sProps == null) {
                int numStrategies = sNames.size();
                double equalProp = 1.0D/((double)numStrategies);
                sProps = new java.util.ArrayList<Double>();
                for (int i = 0; i < numStrategies; i++) {
                    sProps.add(equalProp);
                }
            }
            
            // make sure that the sNames and sProps lists have equal size
            if (sNames.size() != sProps.size()) {
                System.err.println("Number of proportions does not equal number of strategies");
                System.exit(1);
            }
            
            // make sure that the  specified proportios add up to 1
            double totalProp = 0;
            for (double prop : sProps) {
                totalProp += prop;
            }
            if (totalProp < 1.0D) {
                System.err.println("WARNING: Sum of strategy proportions (" + totalProp + ") is not equal to 1");
                //System.exit(1);
            }
            
            // collect list of strategies and their initial population proportions
            Map<Class<? extends Steppable>, Double> stratProps
                = new java.util.LinkedHashMap<Class<? extends Steppable>, Double>();;
            try {
                java.util.Iterator<Double> sPropsIter = sProps.iterator();
                for (String sName : sNames) {
                    Class<? extends Steppable> strategy = Class.forName(sName).asSubclass(Steppable.class);
                    stratProps.put(strategy, sPropsIter.next());
                }
                System.out.println(stratProps);
            }
            catch (ClassNotFoundException ex) {
                System.out.println(ex);
                System.exit(1);
            }
            catch (ClassCastException ex) {
                System.out.println(ex);
                System.exit(1);
            }
            
            // create a PGGSim
            return new PGGSim(seed, stratProps, numAgents, pggSize, gamesPerRound, numRounds,
                              selectionStrength, oDir, !all, verbose);
        }
        public Class simulationClass() {
            return PGGSim.class;
        }
    }
    
    /**
     * Run the simulation.
     */
    public static void main(String[] args) {
        doLoop(new MakePGGSim(), args);
        System.exit(0);
    }
}