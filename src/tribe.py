# -*- coding: utf-8 -*-
"""
Created on Sat Mar 12 13:30:44 2016

@author: John Maloney
"""

import agent;
import random;

RNGEN = random.SystemRandom();
def randbool():
    return bool(RNGEN.randint(0,1));

class Tribe:
    rngen = random.SystemRandom();
    agents = [];
    total_payouts = 0;
    
    def __init__(self, num_agents):
        self.agents = [agent.Agent() for j in xrange(num_agents)];

    '''
    Reset the tribe's agents to prepare for participation in the next generation.
    '''
    def reset(self):
        self.total_payouts = 0;
        for a in self.agents:
            a.reset();

    '''
    Play the required rounds of the IR game to complete the current generation.
    '''
    def playrounds(self, cost, benefit):
        for i in xrange(len(self.agents)):
            for j in xrange(i+1, len(self.agents)):
                # randomly assign the agents to roles
                if (randbool()):
                    # agent i is donor
                    self.total_payouts += self.agents[i].playround(self.agents[j], cost, benefit);
                else:
                    # agent j is donor
                    self.total_payouts += self.agents[j].playround(self.agents[i], cost, benefit);
        # return the total payouts for use by the sim engine
        return self.total_payouts;
    
    '''
    Randomly select an agent from the local population.  The chance that an
    agent is selected is proportional to its fitness.
    '''
    def select_local_parent(self):
        ri = RNGEN.randint(0, self.total_payouts);
        thresh = 0;
        for a in self.agents:
            thresh += a.payout
            if (ri <= thresh):
                return a;

    '''
    Create the next generation by propagating action modules to the next
    generation based on the fitness those modules achieved.
    '''
    def create_next_gen(self):
        for a in self.agents:
            parent = self.select_local_parent();
            a.actm = parent.actm;