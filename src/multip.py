# -*- coding: utf-8 -*-
"""
Created on Sun Mar 13 17:39:33 2016

@author: john
"""

import multiprocessing;
import simengine
import itertools

pool = multiprocessing.Pool(2);

class B:
    def __init__(self):
        self.x = 0;
    
class A:
    def __init__(self):
        self.b = B();
        self.y = 1;

def doit(args, conn=None):
    print "doit: %d %d" % (args[1], args[2])
    a = args[0]
    a.b.x = args[1]
    a.y = args[2]
    if (conn != None):
        conn.send(a)
        conn.close()

def printit(a1, a2):
    print "a1.b.x = %d" % a1.b.x;
    print "a1.y   = %d" % a1.y;
    
    print "a2.b.x = %d" % a2.b.x;
    print "a2.y   = %d" % a2.y;
   
def test():
    a1 = A();
    a2 = A();
    a3 = A();
    a4 = A();

    args1 = [[a1, 5, 6], [a2, 7, 8]];
    args2 = [[a3, 5, 6], [a4, 7, 8]];
    
    # execute without pool
    map(doit, args1);
    printit(a1, a2)
    
    # execute with a pool
    p_conn, c_conn = multiprocessing.Pipe()
    p = multiprocessing.Process(target=doit, args=(args2[0], c_conn))
    p.start()
    a3 = p_conn.recv()
    p.join()
    
    p_conn, c_conn = multiprocessing.Pipe()
    p = multiprocessing.Process(target=doit, args=(args2[1], c_conn))
    p.start()
    a4 = p_conn.recv()
    p.join()
    
    #pool.map(doit, args2)
    printit(a3, a4)

def tribe_doit(t, cost, benefit, conn=None):
    t.playrounds(cost, benefit)
    if (conn != None):
        conn.send(t)
        conn.close()

def tribe_queueit(t, cost, benefit, Q=None):
    t.playrounds(cost, benefit)
    if (Q != None):
        Q.put(t)

def tribe_poolit(args):
    t = args[0]
    cost = args[1]
    benefit = args[2]
    if (len(args) > 3):
        Q = args[3]
    else:
        Q = None
    if (len(args) > 4):
        V = args[4]
    else:
        V = None
    if (len(args) > 5):
        L = args[5]
    else:
        L = None

    tribe_queueit(t, cost, benefit, Q)
    
    with L:
        V.value += t.total_payouts

def testtribe():
    sim = simengine.SimEngine(10,5)
    
    for t in sim.tribes:
        print "tribe payout: %d (%d)" % (t.total_payouts, len(t.agents))
        for a in t.agents:
            print "  agent payout: %d" % a.payout

    Q = multiprocessing.Queue()
    
    processes = []
    for t in sim.tribes:
        processes.append(multiprocessing.Process(target=tribe_queueit, args=(t, 1, 3, Q)))

    for p in processes:
        p.start()
        
    for p in processes:
        p.join()

    while (not Q.empty()):
        t = Q.get()
        print "tribe payout: %d (%d)" % (t.total_payouts, len(t.agents))
        for a in t.agents:
            print "  agent payout: %d" % a.payout

def testpool():
    sim = simengine.SimEngine(2,2)
    
    for t in sim.tribes:
        print "tribe payout: %d (%d)" % (t.total_payouts, len(t.agents))
        for a in t.agents:
            print "  agent payout: %d" % a.payout

    # create a manager to manage the queue
    mgr = multiprocessing.Manager()
    Q = mgr.Queue()
    V = mgr.Value('i', 0)
    L = mgr.Lock()
    
    pool = multiprocessing.Pool()
    
    cost = itertools.repeat(1)
    benefit = itertools.repeat(3)
    queue = itertools.repeat(Q)
    value = itertools.repeat(V)
    lock = itertools.repeat(L)
    
    args = itertools.izip(sim.tribes, cost, benefit, queue, value, lock)
    
    pool.map(tribe_poolit, args)

    new_tribes = []
    while (not Q.empty()):
        t = Q.get()
        new_tribes.append(t)
        print "tribe payout: %d (%d)" % (t.total_payouts, len(t.agents))
        for a in t.agents:
            print "  agent payout: %d" % a.payout
    
    sim.tribes = new_tribes

    for t in sim.tribes:
        print "tribe payout: %d (%d)" % (t.total_payouts, len(t.agents))
        for a in t.agents:
            print "  agent payout: %d" % a.payout
    
    print 'total payout: %d' % V.value
 