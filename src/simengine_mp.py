# -*- coding: utf-8 -*-
"""
Utilities to execute the simulation using multiprocessing.

Created on Mon Mar 14 22:07:27 2016

@author: John Maloney
"""

import itertools
import multiprocessing

'''
Worker function used to play rounds using multiprocessing.
- arg[0]: tribe
- arg[1]: cost
- arg[2]: benfit
- arg[3]: managed queue
- arg[4]: managed int for payouts
- arg[5]: lock for managing access to payouts
'''
def playrounds_worker(args):
    # get args
    t = args[0]
    cost = args[1]
    benefit = args[2]
#    queue = args[3]
#    payouts = args[4]
#    lock = args[5]
    # play the rounds
    p = t.playrounds(cost, benefit)
    return t
    # put tribe on the queue
#    queue.put(t)
    # add payout to total payouts from all tribes
#    with lock:
#        payouts.value += p;
    
# create a process pool, managed queue, shared int and a lock
POOL = multiprocessing.Pool(multiprocessing.cpu_count());
MGR = multiprocessing.Manager();
QUEUE = MGR.Queue();
PAYOUTS = MGR.Value('i',0);
PAYOUTS_LOCK = MGR.Lock();

'''
Play the required rounds of the IR game to complete the current generation.

Use multiprocessing to play the rounds.
'''
def playrounds_mp(sim, cost, benefit):
    # build tuples for worker
    C = itertools.repeat(cost)
    B = itertools.repeat(benefit)
#    Q = itertools.repeat(QUEUE)
#    P = itertools.repeat(PAYOUTS)
#    L = itertools.repeat(PAYOUTS_LOCK)
#    args = itertools.izip(sim.tribes, C, B, Q, P, L)
    args = itertools.izip(sim.tribes, C, B)
    # play rounds using multiprocessing
    new_tribes = POOL.map(playrounds_worker, args)
    # copy tribes back to simengine
#    new_tribes = []
#    while (not QUEUE.empty()):
#        new_tribes.append(QUEUE.get())
    sim.tribes = new_tribes
    sim.total_payouts = sum([t.total_payouts for t in sim.tribes])
    # update total_payout from all tribes and reset PAYOUTS to zero
#    with PAYOUTS_LOCK:
#        sim.total_payouts = PAYOUTS.value;
#        PAYOUTS.value = 0
    # return total_payouts
    return sim.total_payouts;