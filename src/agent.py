# -*- coding: utf-8 -*-
"""
Created on Sat Mar 12 12:11:12 2016

@author: John Maloney
"""

import random
import rep

RNGEN = random.SystemRandom();
def randbool():
    return bool(RNGEN.randint(0,1));

class Agent:
    rep = rep.GOOD;
    payout = 0;
    actm = None;
    num_games = 0;
    
    def __init__(self):
        self.actm = ActionModule();

    '''
    Reset the agent's internal state to prepare for participation in the
    next generation.
    '''
    def reset(self):
        self.payout = 0;
        self.num_games = 0;
        
    '''
    Returns True if the agent chooses to donate and False otherwise.
    '''
    def chooses_donate(self, recipient):
        return self.actm.choose_donate(self, recipient);
    
    '''
    Play a round of the IR game with this agent playing the role of the
    donor agent.  The total payout earned by both agents is returned.
    '''
    def playround(self, recipient, cost, benefit):
        # increase number of games played by agents
        self.num_games += 1;
        recipient.num_games += 1;
        
        # keep track of total payout earned by both agents
        total_payout = 0;
        
        if (self.chooses_donate(recipient)):
            # donor donates
            # -- recipient receives benefit
            recipient.payout += benefit;
            # -- donor pays cost
            self.payout -= cost;
            # update total payout
            total_payout += (benefit - cost);

        # to prevent negative payout, each donor receives cost
        self.payout += cost;
        recipient.payout += cost;
        total_payout += (2*cost);
        
        # return the total payut earned
        return total_payout;
    
class ActionModule:
    bits = [];
    
    '''
    Create an action module with the specified actions.  The argument
    defines the action to take in each of the four possible situations:
    
      Donor  Recipient  Action
      -----  ---------  ------
      GOOD   GOOD       bits[0] (False = no donation, else donate)
      GOOD   BAD        bits[1] (False = no donation, else donate)
      BAD    GOOD       bits[2] (False = no donation, else donate)
      BAD    BAD        bits[3] (False = no donation, else donate)
      
    If bits is not provided then a random action module is created.
    '''
    def __init__(self, bits=None):
        if (bits == None):
            bits = [randbool(),
                    randbool(),
                    randbool(),
                    randbool()];
        
        self.bits = bits;

    '''
    Returns True if the agent should donate and False otherwise.
    '''
    def choose_donate(self, donor, recipient):
        if (donor.rep == rep.GOOD):
            if (recipient.rep == rep.GOOD):
                # GOOD GOOD
                return self.bits[0];
            else:
                # GOOD BAD
                return self.bits[1];
        else:
            if (recipient.rep == rep.GOOD):
                # BAD GOOD
                return self.bits[2];
            else:
                # BAD BAD
                return self.bits[3];