# -*- coding: utf-8 -*-
"""
Generate statistics for bit fixation.

Arguments:
  csvfile  - path to the CSV file (input)

Created on Tue Mar 22 20:40:17 2016

@author: John Maloney
"""

import argparse
import pandas as pd
import sys

def main():
    desc = 'Calculate the bit fixation statistics'
    parser = argparse.ArgumentParser(description=desc)
    parser.add_argument('csvfile', help='file that holds the data to be plotted')
    parser.add_argument('-p', type=int, help='number of periods to include in the calculation')
    parser.add_argument('-t', type=int, help='number of tribes in simulation')
    
    args = parser.parse_args()
    #print args.csvfile
    #print args
    
    if (not args.p):
        periods = -1
    else:
        periods = args.p

    if (not args.t):
        tribes = 64
    else:
        tribes = args.t

    run_script(args.csvfile, periods, tribes)
# end main

def get_result(percent):
    if (percent > 0.95):
        return '1'
    elif (percent < 0.05):
        return '0'
    else:
        return 'X'
# end get_result

def run_script(csvfile, periods, tribes):
    # load CSV data
    print 'Loading data from %s...' % csvfile
    data = pd.read_csv(csvfile)
    
    # bit columns
    assess_columns = ['n0','n1','n2','n3','n4','n5','n6','n7']
    assess_data = data[assess_columns]
    action_columns = ['a0','a1','a2','a3']
    action_data = data[action_columns]
    
    # negative periods argument means plot all the data
    if (periods < 0):
        print '  calculate statistics using data for all periods...'
    else:
        print '  calculate statistics using data for last %d periods...' % periods
        start_idx = data.shape[0] - periods
        assess_data = assess_data[start_idx:]
        action_data = action_data[start_idx:]

    # get numbe of tribes and agents
    num_tribes = data['t'][0]
    num_agents = data['a'][0]

    # calculate maximum count given specified number of periods
    max_assess = periods*num_tribes
    max_action = periods*num_agents
    
    # calculate assess results
    assess_percent = assess_data.sum()/max_assess
    assess_result = [ get_result(p) for p in assess_percent ]
    action_percent = action_data.sum()/max_action
    action_result = [ get_result(p) for p in action_percent ]
    
    # output result
    sys.stdout.write('assess: [')
    for ch in assess_result:
        sys.stdout.write(ch)
    sys.stdout.write(']')
    sys.stdout.write(' action:[')
    for ch in action_result:
        sys.stdout.write(ch)
    sys.stdout.write(']\n')
# end run_script

# run main method when this file is run from command line
if __name__ == "__main__":
    main()