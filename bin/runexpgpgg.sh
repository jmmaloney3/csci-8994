#!/bin/bash

# run the simulations required for an experiment for graph PGG
# example use: runexpgpgg.sh <dir> <num_graphs> <simengine_params>

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
  >&2 echo "file or directory $RESULTSDIR already exist"
  exit 1
fi

# run the runsim go program to execute simulation
# assumption: the go program is in the $BIN directory
mkdir $RESULTSDIR
cd $RESULTSDIR
>&2 echo "  writing experiment results to $RESULTSDIR"

# get the number of experiments to execute
NUMEXPS=$1
>&2 echo "  running simulations for $NUMEXPS graph instances"
shift # remove the number of sims from list of arguments

# echo the command that will be used to execute each simulation
>&2 echo "  execute following for each graph instance:"
>&2 echo "    $BIN/runsimgpgg $@ -d graph<num> &> graph<num>.log"

for i in `seq -f "%02g" 1 $NUMEXPS`;
do
  >&2 echo "  executing simulations for graph instance $i..."
  $BIN/runsimgpgg "$@" -d "graph$i" &> "graph$i.log"
done

cd ..
