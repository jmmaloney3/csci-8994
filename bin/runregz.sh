#!/bin/bash

# run the experiments for regular graphs with varying z
# example use: runregz.sh <dir> <numgraphs_per_experiment>

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
>&2 echo "  writing experiments for regular graphs with varying z results to $RESULTSDIR"

# get the number of experiments to execute
NUMGRAPHS=$1
>&2 echo "  using $NUMGRAPHS graph instances for each experiment"
shift # remove the number of graphs from list of arguments

for Z in 4 8 16 32 64;
do
  >&2 echo "  executing experiments for z=$Z..."
  # create directory to hold results for current value of z
  mkdir z`printf %02d $Z`
  cd z`printf %02d $Z`
  # generate series data for current value of z (all values of r)
  for R in 1 2 3 4 5;
  do
    >&2 echo "  executing experiment for r=$R..."
    $BIN/runexpgpgg.sh r$R $NUMGRAPHS -a 10000 -z $Z -g 10000 -r $R -c 1 -w 0
  done
  cd ..
done

cd ..
