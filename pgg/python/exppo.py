# -*- coding: utf-8 -*-
"""
Calculate the expected payouts for cooperators and defectors if the cooperators
defect when the percentage of defectors in the group exceeds a threshold.

Created on Sat Dec 12 16:06:33 2015

@author: John Maloney
"""
import scipy.stats as stats

import matplotlib.pyplot as plt

def main():

    run_script(100, 5, 5, 3)
# end main

def run_script(M, n, N, r):
    # M - total agents in the population
    # n - total number of defectors in the population
    # N - number of agents to be selected

    # construct the required hypergeometric distribution
    # -- assume that a cooperator has alredy been selected
    #    evaluate payouts based on the rest of the players
    #    that are selected
    rv = stats.hypergeom(M-1, n, N-1)
    
    cpayouts = []
    dpayouts = []
    
    print 'N = %d' % N
    
    for T in xrange(N):
        cpo = float(0)
        dpo = float(0)
        for k in xrange(T+1):
            # calculate base payout when there are k defectors
            base_payout = (float(N-k)/float(N))*float(r)
            print 'base payout (k=%d): %f' % (k, base_payout)

            # get probability that k defectors are selected
            p = rv.pmf(k)
            cpo = cpo + (p * (base_payout - 1))
            if (k > 0):
                dpo = dpo + (p * base_payout)

        # collect caculated payouts
        cpayouts.append(cpo)
        dpayouts.append(dpo)
        # end for k
    # end for T
    
    print cpayouts
    print dpayouts
    
    # plot the two payout curves
    plt.clf()
    plt.plot(cpayouts, linewidth=0.5, color='b', aa=True)
    plt.plot(dpayouts, linewidth=0.5, color='r', aa=True)
    
    plt.show()

# end run_script

# run main method when this file is run from command line
if __name__ == "__main__":
    main()