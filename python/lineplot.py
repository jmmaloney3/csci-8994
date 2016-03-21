# -*- coding: utf-8 -*-
"""
Generate a line plot for the specified attributes.

Arguments:
  csvfile  - path to the CSV file (input)

Created on Thu Dec  4 07:13:00 2015

@author: John Maloney
"""

import csvutils as cu
import numpy as np
import matplotlib
import matplotlib.pyplot as plt
from matplotlib.font_manager import FontProperties
import argparse

def main():
    desc = 'Create a line plot using the data in the specified CSV file'
    parser = argparse.ArgumentParser(description=desc)
    parser.add_argument('csvfile', help='file that holds the data to be plotted')
    parser.add_argument('-p', type=int, help='number of periods to include in the plot')
    parser.add_argument('-g', type=int, help='size of gap beween plotted data points')
    
    args = parser.parse_args()
    #print args.csvfile
    #print args
    
    if (not args.p):
        periods = -1
    else:
        periods = args.p

    if (not args.g):
        gap = 1
    else:
        gap = args.g

    run_script(args.csvfile, periods, gap)
# end main

def run_script(csvfile, periods, gap):
    # load CSV data
    print 'Loading matrix from %s...' % csvfile
    matrix, headers = cu.load_matrix(csvfile, True)
    
    # negative periods argument means plot all the data
    if (periods < 0):
        periods = matrix.shape[0]

    print '  plot data for %d periods...' % periods
    print '  plot every %dth data point...' % gap

    # the matrix is structured as follows:
    # -- the first row contains the column headers
    # -- the first column contains the round numbers
    # -- each additional column contains the data for one strategy

    # set the font for the legend
    matplotlib.rcParams.update({'font.size': 8})
    fontP = FontProperties()
    fontP.set_size(8)

    # clear the plot
    plt.figure(figsize=(8,3))
    plt.clf()
    
    #colors = {'b0' : 'b', 'b1': 'r',
    #          'NonParticipant' : 'y', 'Punisher' : 'g',
    #          'L8JudgingAgent': 'c',
    #          'L8BBGAgent' : 'c',
    #          'L8BGBAgent' : 'c',
    #          'L8BGGAgent' : 'c',
    #          'L8GBBAgent' : 'c',
    #          'L8GBGAgent' : 'c',
    #          'L8StandingAgent' : 'c',
    #          'L8GGGAgent' : 'c'
    #          }
    for i in range(1,matrix.shape[1]-1):
        indices = np.arange(0,periods)
        indices = indices[::gap]
        #plt.plot(indices, matrix[indices,i], linewidth=0.5, color=colors[headers[i].strip()], aa=True)
        plt.plot(indices, matrix[indices,i], linewidth=0.5, aa=True)

    # label the axes
    #if (csvfile.find('scounts') >=0):
    #    ylab = "strategy frequencies (%)"
    #elif (csvfile.find('spayouts') >= 0):
    #    ylab = "average payouts"
    #elif (csvfile.find('sfitness') >= 0):
    #    ylab = "average fitness"
    #else:
    #    ylab = "unknown data"
    ylab = "num of tribes"
    
    plt.ylabel(ylab)
    plt.xlabel("generations")
    
    # create the legend
    ax = plt.subplot(111)
    box = ax.get_position()
    ax.set_position([box.x0, box.y0 + box.height * 0.5,
                     box.width, box.height*0.5])
    plt.legend(headers[1::], loc='upper center', bbox_to_anchor=(0.5, -0.3),
               ncol=matrix.shape[1]-1, prop=fontP, frameon=False)

    # get file name for plot
    idx = csvfile.find('.csv')
    if (idx >= 0):
        pngfile = csvfile[0:idx] + '.png'
    else:
        pngfile = csvfile + '.png'

    # write the plot to a file
    plt.savefig(pngfile, bbox_inches='tight');

    # plt.show()
# end run_script

# run main method when this file is run from command line
if __name__ == "__main__":
    main()