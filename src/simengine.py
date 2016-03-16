# -*- coding: utf-8 -*-
"""
Created on Sat Mar 12 12:17:25 2016

@author: John Maloney
"""

import tribe;
import random;
import simengine_mp;

RNGEN = random.SystemRandom();

class SimEngine:
    
    def __init__(self, num_tribes, num_agents_per_tribe):
        self.tribes = [ tribe.Tribe(i, num_agents_per_tribe) for i in xrange(num_tribes)];
        self.tribes_dict = {i : tribe.Tribe(i, num_agents_per_tribe) for i in xrange(num_tribes)};
        self.total_payouts = 0;

    '''
    Reset the simulations to prepare for participation in the next generation.
    '''
    def reset(self):
        self.total_payouts = 0;
        for t in self.tribes:
            t.reset();
    
    '''
    Play the required rounds of the IR game to complete the current generation.
    '''
    def playrounds(self, cost, benefit, USE_MP=False):
        if (USE_MP):
            simengine_mp.playrounds_mp(self, cost, benefit);
        else:
            for t in self.tribes:
                self.total_payouts += t.playrounds(cost, benefit);

        return self.total_payouts;
    
    '''
    Create the next generation by propagating action modules to the next
    generation based on the fitness those modules achieved.
    '''
    def create_next_gen(self):
        for t in self.tribes:
            t.create_next_gen();