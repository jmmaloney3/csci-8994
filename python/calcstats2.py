# -*- coding: utf-8 -*-
"""
Generate final statistics for bit fixation.

Created on Sat Mar 26 10:47:41 2016

@author: John Maloney
"""

import argparse
import pandas as pd
import os.path as path
import sys

def main():
    desc = 'Calculate the FINAL bit fixation statistics'
    parser = argparse.ArgumentParser(description=desc)
    parser.add_argument('csvfile', help='file that holds the data produced by calcstats')
    parser.add_argument('-o', type=str, help='output file')
    parser.add_argument('-v', action='store_true')
    
    args = parser.parse_args()
    #print args.csvfile
    #print args
    
    run_script(args.csvfile, args.o, args.v)
# end main

def run_script(csvfile, ofile_name, verbose):
    # make sure input file exists and is a regular file
    if (not path.exists(csvfile)):
        sys.stderr.write('input file %s does not exist\n' % csvfile)
        return
        if (path.isdir(csvfile)):
            sys.stderr.write('input file %s is a directory\n' % csvfile)
            return

    # get csv writer
    if (ofile_name is None):
        calc_stats(csvfile, sys.stdout, verbose)
    else:    
        if (path.exists(ofile_name)):
            sys.stderr.write('output file %s exists\n' % ofile_name)
            return
        else:
            with open(ofile_name, 'wb') as ofile:
                calc_stats(csvfile, ofile, verbose)
# end run_script

def get_result(data, column):
    counts = data[column].value_counts()
    # bit value fixated at 1
    if ('1' in counts):
        b1 = counts['1']
    elif (1 in counts):
        b1 = counts[1]
    else:
        b1 = 0
    # bit value fixated at 0
    if ('0' in counts):
        b0 = counts['0']
    elif (0 in counts):
        b0 = counts[0]
    else:
        b0 = 0
    # bit value did not fixate
    if ('X' in counts):
        bX = counts['X']
    elif (-1 in counts):
        bX = counts[-1]
    else:
        bX = 0

    if (b1 > (b0 + bX)):
        return '1'
    elif (b0 > (b1 + bX)):
        return '0'
    else:
        return 'X'
# end get_result

def calc_stats(csvfile, ofile, verbose):    
    # define column names
    assess_columns = ['n0','n1','n2','n3','n4','n5','n6','n7']
    action_columns = ['a0','a1','a2','a3']

    # load CSV data
    if (verbose):
        sys.stderr.write('Loading data from %s...\n' % csvfile)
    data = pd.read_csv(csvfile, skipinitialspace=True)
    
    # calculate results
    assess_result = [ get_result(data, column) for column in assess_columns ]
    action_result = [ get_result(data, column) for column in action_columns ]
    
    # output results
    ofile.write('assess: [%s]' % ','.join(bit for bit in assess_result))
    ofile.write(' action: [%s]\n' % ','.join(bit for bit in action_result))
# end process_file

# run main method when this file is run from command line
if __name__ == "__main__":
    main()