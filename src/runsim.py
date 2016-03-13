# -*- coding: utf-8 -*-
"""
Run the simulation with the specified arguments.

Arguments:
  tribes  - number of tribes
  agents  - number of agents per tribe
  cost    - cost c to donate
  benefit - benefit b received from donation

Created on Sun Mar 13 16:00:53 2016

@author: John Maloney
"""

import argparse
import simengine

def main():
    desc = 'Run the indirect reciprocity simulation.'
    parser = argparse.ArgumentParser(description=desc)
    parser.add_argument('-t', type=int, help='number of tribes')
    parser.add_argument('-a', type=int, help='number of agents per tribe')
    parser.add_argument('-c', type=int, help='cost c to donate')
    parser.add_argument('-b', type=int, help='benefit b received from donation')
    parser.add_argument('-g', type=int, help='number of generations to simulate')
    
    args = parser.parse_args()
    #print args.csvfile
    #print args
    
    if (not args.t):
        tribes = 64
    else:
        tribes = args.t

    if (not args.a):
        agents = 64
    else:
        agents = args.a

    if (not args.c):
        cost = 1
    else:
        cost = args.c

    if (not args.b):
        benefit = 3
    else:
        benefit = args.b

    if (not args.g):
        gens = 10
    else:
        gens = args.g

    run_script(tribes, agents, cost, benefit, gens)
# end main

def run_script(tribes, agents, cost, benefit, gens):
    # create the simulation
    sim = simengine.SimEngine(tribes, agents)
    
    # simuate the request number of generations
    for i in xrange(gens):
        sim.playrounds(cost, benefit)
        print sim.total_payouts
        sim.create_next_gen()
        sim.reset()
# end run_script

# run main method when this file is run from command line
if __name__ == "__main__":
    main()