# -*- coding: utf-8 -*-
"""
Utilities to execute the simulation using multiprocessing.

Created on Tue Mar 15 22:28:38 2016

@author: John Maloney
"""

import itertools
import multiprocessing

'''
Worker function used to play rounds using multiprocessing.
Inputs:
- args[0]: donor agent
- args[1]: recipient donor
- args[2]: cost
- args[2]: benfit
Outputs:
- return[0]: donor new overall payout
- return[1]: recipient new overall payout
- return[3]: joint payout from round
'''
def playrounds_worker(args):
    # get args
    donor = args[0]
    recipient = args[1]
    cost = args[2]
    benefit = args[3]
    # play the rounds
    total_payout = donor.playrounds(recipient, cost, benefit)
    # return payouts
    return (donor.payout, recipient.payout, total_payout)
    
# create a process pool
POOL = multiprocessing.Pool(multiprocessing.cpu_count());

'''
Play the required rounds of the IR game to complete the current generation.

Use multiprocessing to play the rounds.
'''
def playrounds_mp(t, cost, benefit):
    # build tuples for worker
    C = itertools.repeat(cost)
    B = itertools.repeat(benefit)
    args = itertools.izip(t.agents, C, B)
    # play rounds using multiprocessing
    rvals = POOL.map(playrounds_worker, args)
    return t.total_payouts;