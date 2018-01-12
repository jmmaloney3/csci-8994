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
import os.path as path
import os
import csv

def main():
    desc = 'Calculate the bit fixation statistics'
    parser = argparse.ArgumentParser(description=desc)
    parser.add_argument('csvfile', help='file that holds the raw generation data')
    parser.add_argument('-p', type=int, help='number of periods to include in the calculation')
    parser.add_argument('-t', type=percent_type, help='threshold for ALLD/ALLC strategy types')
    parser.add_argument('-o', type=str, help='output file')
    parser.add_argument('-v', action='store_true')
    
    args = parser.parse_args()
    #print args.csvfile
    #print args
    
    if (not args.p):
        periods = -1
    else:
        periods = args.p
    
    if (not args.t):
        allcd_t = 0.1
    else:
        allcd_t = args.t

    if (not args.o):
        ofile = None
    else:
        ofile = args.o

    run_script(args.csvfile, periods, allcd_t, ofile, args.v)
# end main

def percent_type(value):
    fvalue = float(value)
    if ((fvalue <0) or (fvalue > 1)):
        msg = "%s is an invalid percent value" % value
        raise argparse.ArgumentTypeError(msg)
    return fvalue

def get_result(percent):
    if (percent > 0.95):
        return '1'
    elif (percent < 0.05):
        return '0'
    else:
        return 'X'
# end get_result

def run_script(csvfile, periods, allcd_t, ofile_name, verbose):
    # Collect the files to process
    files = []
    if (path.exists(csvfile)):
        if (path.isdir(csvfile)):
            fnames = os.walk(csvfile).next()[2]
            for fn in fnames:
                files.append(path.join(csvfile, fn))
        else:
            files.append(csvfile)

    # get output file handle and calculate stats
    if (ofile_name is None):
        calc_stats(files, periods, allcd_t, sys.stdout, verbose)
    else:    
        if (path.exists(ofile_name)):
            sys.stderr.write('output file %s exists\n' % ofile_name)
            return
        else:
            with open(ofile_name, 'wb') as ofile:
                calc_stats(files, periods, allcd_t, ofile, verbose)
# end run_script
    
def calc_stats(files, periods, allcd_t, ofile, verbose):
    # create csv writer
    csv_writer = csv.writer(ofile)
    
    # define column names
    assess_columns = ['n0','n1','n2','n3','n4','n5','n6','n7']
    action_columns = ['a00','a01','a02','a03','a04','a05','a06','a07','a08','a09','a10','a11','a12','a13','a14','a15']
    po_columns     = ['po', 'minpo', 'maxpo']
    
    # write headers to output file
    csv_writer.writerow(assess_columns+action_columns)

    # process the files and calculate statistics
    for ifile in files:
        fname, fext = os.path.splitext(ifile)
        if (fext == '.csv'):
            process_file(ifile, periods, csv_writer, assess_columns, action_columns, po_columns, allcd_t, verbose)
# end calc_stats

def process_file(csvfile, periods, csv_writer, assess_columns, action_columns, po_columns, allcd_t, verbose):    
    # load CSV data
    if (verbose):
        sys.stderr.write('Loading data from %s...\n' % csvfile)
    data = pd.read_csv(csvfile, skipinitialspace=True)
    
    allc = 'a15'
    alld = 'a00'
    
    # get bit column data
    assess_data = data[assess_columns]
    action_data = data[action_columns]
    po_data     = data[po_columns]
    allcd_data  = data[[allc, alld]]

    # negative periods argument means use all the data
    if ((periods < 0) or (data.shape[0] < periods)):
        periods = data.shape[0]
        if (verbose):
            sys.stderr.write('  calculate statistics using data for all %d periods...\n' % periods)
    else:
        if (verbose):
            sys.stderr.write('  calculate statistics using data for last %d periods...\n' % periods)
        start_idx = data.shape[0] - periods
        assess_data = assess_data[start_idx:]
        action_data = action_data[start_idx:]
        po_data     = po_data[start_idx]
        allcd_data  = allcd_data[start_idx:]

    # get number of tribes and agents
    num_tribes = data['t'][0]
    num_agents = data['a'][0]

    # calculate maximum count given specified number of periods
    max_assess = periods*num_tribes
    max_action = periods*(num_agents*num_tribes)
    
    # check allc/alld threshold
    allcd_percent = allcd_data.sum()/max_action
    if (allcd_percent[allc] > allcd_t):
        sys.stderr.write('  [%s] ALLC prevelance (%6.4f) exceeds %4.2f threshold\n' % (path.basename(csvfile), allcd_percent[allc], allcd_t))
        return
    if (allcd_percent[alld] > allcd_t):
        sys.stderr.write('  [%s] ALLD prevelance (%6.4f) exceeds %4.2f threshold\n' % (path.basename(csvfile), allcd_percent[alld], allcd_t))
        return

    # calculate results
    assess_percent = assess_data.sum()/max_assess
    assess_result = [ get_result(p) for p in assess_percent ]
    action_percent = action_data.sum()/max_action
    # action_result = [ get_result(p) for p in action_percent ]
    action_result = [ str(p) for p in action_percent ]
    
    # calculate payout percent of maximum possible payout
    po_percent = data['po']/data['pomax']
    
    # output percentages
    assess_str = ','.join('%4.2f' % n for n in assess_percent)
    action_str = ','.join('%4.2f' % n for n in action_percent)
    sys.stderr.write('  [%s] [%s] [%s]\n' % (path.basename(csvfile), assess_str, action_str))
    # output result
    csv_writer.writerow(assess_result+action_result)
# end process_file

# run main method when this file is run from command line
if __name__ == "__main__":
    main()