# -*- coding: utf-8 -*-
"""
Generate the results for a single GPGG experiment.

Arguments:
  dir  - path to the director holding the results

Created on Wed May 11 16:20:15 2016

@author: John Maloney
"""

import argparse
import pandas as pd
import numpy as np
import sys
import os.path as path
import os
import matplotlib
import matplotlib.pyplot as plt
from matplotlib.font_manager import FontProperties
import json
import itertools

def main():
    desc = 'Produce results for a series of experiments'
    parser = argparse.ArgumentParser(description=desc)
    parser.add_argument('serdirs', type=str, nargs='+', help='list of directories that hold the series for the chart')
    parser.add_argument('-p', type=int, help='number of periods to include in the calculation', default=1000)
    parser.add_argument('-o', type=str, help='output file', default=None)
    parser.add_argument('-x', type=str, help='sim param that is x value', default="r")
    parser.add_argument('-s', type=str, help='sim param that is series value', default="z")
    parser.add_argument('-v', action='store_true')
    
    args = parser.parse_args()
    print args.senddirs
    print args
    
    run_script(args.serdirs, args.s, args.x, args.p, args.v)
# end main

# get the pstat.csv files and a JSON log file for a single experiment
# - expdir - the directory holding the experiment raw data
def get_exp_files(expdir, verbose):
    # structure for directory that holds the raw experimental results
    # expdir (name of root directory - passed into script as dname parameter)
    # -- graph<N> (holds results for Nth graph)
    # -- -- sim<M> (holds results for Mth simulation for Nth graph )
    # -- -- -- pstat.csv (holds population stats - strategy frequency)

    # for each pstat.csv file do the following
    # -- check if either strategy was eliminated in the last generation
    # -- -- if yes, then use the frequence from the last generation
    # -- otherwise, calc avg freq from last <periods> generations
    
    # Collect the pstat.csv files to process
    logfile=None
    files = []
    if (path.exists(expdir)):
        if (path.isdir(expdir)):
            for parent, ds, fs in os.walk(expdir):
                if ('pstat.csv' in fs):
                    files.append(path.join(parent, 'pstat.csv'))
                if ('graph01.log' in fs):
                    logfile=path.join(parent,'graph01.log')
        else:
            raise Exception('not a directory: %s\n' % expdir)            
    else:
        raise Exception('directory does not exist: %s\n' % expdir)

    # return list pstat.csv files
    if (verbose):
         sys.stderr.write('    # of pstat.csv files found: %d\n' % len(files))
         sys.stderr.write('    log file: %s\n' % logfile)
    return [files, logfile]
# end get_psfiles

def get_sim_result(csvfile, periods, verbose):    
    # load CSV data
    if (verbose):
        sys.stderr.write('    loading data from %s...\n' % csvfile)
    data = pd.read_csv(csvfile, skipinitialspace=True)
    
    # get cooperator data
    coop_data = data['Pc']

    # if one of the srategies was eliminated then return the last value
    last_value = coop_data[coop_data.shape[0]-1]
    if ((last_value<=0) or (last_value>=100)):
        if (verbose):
            if (last_value<=0):
                sys.stderr.write('      cooperators eliminated - return %5.2f...\n' % last_value)
            else:
                sys.stderr.write('      defectors eliminated - return %5.2f...\n' % last_value)
        return last_value

    # else - calculate the average for the last <periods> generations
    
    # negative periods argument means use all the data
    if ((periods < 0) or (data.shape[0] < periods)):
        periods = data.shape[0]
        if (verbose):
            sys.stderr.write('      calculate average using data for all %d periods...\n' % periods)
    else:
        if (verbose):
            sys.stderr.write('      calculate average using data for last %d periods...\n' % periods)
        start_idx = data.shape[0] - periods
        coop_data = coop_data[start_idx:]

    # calculate and return the mean of the column
    avg = coop_data.mean(axis=1)
    if (verbose):
        sys.stderr.write('      [%s] [%5.2f]\n' % (path.basename(csvfile), avg))
    return avg
# end get_sim_result

# get the result for an experiment
# the result is the fraction of the population that are cooperators
# - expdir - the directory holding the experiment raw data
# - periods - # of generations to use to calculate average
def get_exp_result(expdir, sparam, xparam, periods, verbose):
    if (verbose):
        sys.stderr.write('  calculate result for experiment %s...\n' % expdir)

    # get the files for the experiment
    expfiles = get_exp_files(expdir, verbose)
    psfiles = expfiles[0]
    logfile = expfiles[1]
    
    # get sim parameter from log file
    if (verbose):
        sys.stderr.write('    log file for experiment: %s\n' % logfile)
    fp = open(logfile, 'r')
    objs = json.load(fp)
    fp.close()
    sval = objs[0]['params'][sparam]
    if (verbose):
        sys.stderr.write('    series value for experiment %s: %f\n' % (sparam, sval))
    xval = objs[0]['params'][xparam]
    if (verbose):
        sys.stderr.write('    x value for experiment %s: %f\n' % (xparam, xval))
    
    # iterate through the psfiles and calculate the result
    results = []
    for f in psfiles:
        results.append(get_sim_result(f, periods, verbose))

    # return the average of the averages
    yval = np.mean(results)
    if (verbose):
        sys.stderr.write('    result for experiment %s: %5.2f\n' % (path.basename(expdir), yval))
    return [sval, xval, yval]
#end get_exp_result

# get the results for a series of experiments as a dict
# - serdir - the directory holding the series experiments
# - periods - # of generations to use to calculate average
def get_series_results(serdir, sparam, xparam, periods, verbose):
    if (verbose):
        sys.stderr.write('calculate results for series %s...\n' % serdir)

    # get results for the experiments in the series
    svals = []
    xvals = []
    yvals = []
    if (path.exists(serdir)):
        if (path.isdir(serdir)):
            fileinfo = os.walk(serdir).next()
            parent = fileinfo[0]
            for expdir in fileinfo[1]:
                expdir = path.join(parent, expdir)
                if (verbose):
                    sys.stderr.write('  experiment %s...\n' % expdir)
                result = get_exp_result(expdir,sparam,xparam,periods,verbose)
                svals.append(result[0])
                xvals.append(result[1])
                yvals.append(result[2])

    # return the series
    if (verbose):
        sys.stderr.write('s values for series %s: %s\n' % (path.basename(serdir), svals))
        sys.stderr.write('x values for series %s: %s\n' % (path.basename(serdir), xvals))
        sys.stderr.write('y values for series %s: %s\n' % (path.basename(serdir), yvals))
    return [svals[0], xvals, yvals]
# end get_series_results

# get the series data for a chart
# - serdirs - list of directories each holding data for one series (line)
# - serparam - simulation param used to label each series
# - xparam - simulation param used for x values
# - periods - number of generations to use for calculating results
def get_chart_series(serdirs, sparam, xparam, periods, verbose):
    if (verbose):
        sys.stderr.write('generate series data for chart...\n')
    
    series = []
    for serdir in serdirs:
        series.append(get_series_results(serdir, sparam, xparam, periods, verbose))
    
    return series

def make_plot(series, sparam, xlabel, ylabel, verbose):

    # set the font for the legend
    matplotlib.rcParams.update({'font.size': 8})
    fontP = FontProperties()
    fontP.set_size(8)

    # clear the plot
    #plt.figure(figsize=(8,3))
    plt.figure(figsize=(5,8))
    plt.clf()
    
    # plot the series
    slabels = []
    for s, m in zip(series, itertools.cycle('8s^x,*')):
        slabel = '%s = %d' % (sparam, s[0])
        slabels.append(slabel)
        xvals = s[1]
        yvals = s[2]
        if (verbose):
            sys.stderr.write('plot series %s:\n' % slabel)
            sys.stderr.write('  %s: %s\n' % (xlabel, xvals))
            sys.stderr.write('  %s: %s\n' % (ylabel, yvals))
        # plot the series
        plt.plot(xvals, yvals, linewidth=0.5, aa=True, marker=m)
    
    # set the axis labels
    plt.ylabel(ylabel)
    plt.xlabel(xlabel)
    
    # create the legend
    ax = plt.subplot(111)
    box = ax.get_position()
    ax.set_position([box.x0, box.y0 + box.height * 0.5,
                     box.width, box.height*0.25])
    plt.legend(slabels, loc='upper center', bbox_to_anchor=(0.5, -0.3),
               ncol=len(slabels), prop=fontP, frameon=False)

    # get file name for plot
    #idx = csvfile.find('.csv')
    #if (idx >= 0):
    #    pngfile = csvfile[0:idx] + '.png'
    #else:
    #    pngfile = csvfile + '.png'

    # write the plot to a file
    #plt.savefig(pngfile, bbox_inches='tight');

    plt.show()

def run_script(serdirs, sparam, xparam, periods, verbose):
    # get the results for the series
    series = get_chart_series(serdirs, sparam, xparam, periods, verbose)
    
    # generate the plot
    ylabel="frequency of cooperators"
    make_plot(series, sparam, xparam, ylabel, verbose)
# end run_script

# run main method when this file is run from command line
if __name__ == "__main__":
    main()