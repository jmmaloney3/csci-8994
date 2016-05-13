#!/bin/bash

# run the experiments for specified graph type with varying z
# example use: rungtype.sh <dir> <Z> <numgraphs_per_experiment> <num_sims_per_graph>

# find the directory that holds the script
# - see http://stackoverflow.com/questions/59895/can-a-bash-script-tell-what-directory-its-stored-in/246128#246128
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
SOURCE="$(readlink "$SOURCE")"
[[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
# assumption: this bash script is in the $BIN directory
BIN="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

# echo $BIN

# get the name to use for the directory that holds results
RESULTSDIR=$1
shift # remove the directory name from the list of arguments

# make sure the result directory doesn't exist
if [ -e $RESULTSDIR ];
then
  # echo to stderr
  >&2 echo "file or directory $RESULTSDIR already exists"
  exit 1
fi

# run the runsim go program to execute simulation
# assumption: the go program is in the $BIN directory
mkdir $RESULTSDIR
cd $RESULTSDIR

# get the value for Z
Z=$1
shift # remove the Z from list of arguments
>&2 echo "  writing experiments for z=$Z to $RESULTSDIR"

# get the number of experiments to execute
NUMGRAPHS=$1
>&2 echo "  using $NUMGRAPHS graph instances for each experiment"
shift # remove the number of graphs from list of arguments

# get the number of simulations to execute per graph
NUMSIMS=$1
>&2 echo "  using $NUMSIMS for each graph"
shift # remove the number of sims from list of arguments

# Constants
NGENS=101000
N=1000

for W in 1 3 5 7 9;
do
  >&2 echo "  executing experiments for w=$W..."
  # create directory to hold results for current value of w
  mkdir w`printf %02d $W`
  cd w`printf %02d $W`
  # generate series data for current value of z (all values of r)
  for R in 2 3 4 5 6 7 8 9;
  do
    >&2 echo "  executing experiment for r=$R..."
    time $BIN/runexpgpgg.sh r$R $NUMGRAPHS -s $NUMSIMS -a $N -z $Z -g $NGENS -r $R -c 1 -w $W -gtype 1 &
  done
  wait
  cd ..
done

cd ..
