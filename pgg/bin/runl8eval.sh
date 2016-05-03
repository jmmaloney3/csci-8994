#!/bin/bash

# find the directory that holds the script
# - see http://stackoverflow.com/questions/59895/can-a-bash-script-tell-what-directory-its-stored-in/246128#246128
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
SOURCE="$(readlink "$SOURCE")"
[[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
BIN="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

# echo $BIN

# sub-directories
ROBUSTDIR=robust
STABLEDIR=stable
VIABLEDIR=viable
DEFECTDIR=defect
COOPERDIR=cooper

L8CLASS=$1
CLASSES=$L8CLASS
EXTRACLASS=$3
ROUNDS=$2
RESULTDIR="$(echo $L8CLASS | tr '[:upper:]' '[:lower:]')"
RESULTDIR=${RESULTDIR#pgg.}
RESULTDIR=${RESULTDIR%agent}
if [ ! -z "$3" ]
  then
    SUFFIX="$(echo $EXTRACLASS | tr '[:upper:]' '[:lower:]')"
    SUFFIX=${SUFFIX#pgg.}
    SUFFIX=${SUFFIX%agent}
    RESULTDIR=$RESULTDIR-$SUFFIX
    CLASSES="$CLASSES $EXTRACLASS"
    echo $RESULTDIR
fi
# echo $RESULTDIR

echo "Evaluating $CLASSES..."

# evaluate the specified L8 agent

# create results directory
mkdir $RESULTDIR
cd ./$RESULTDIR

# evaluate robustness
echo "  evaluating robustness..."
mkdir $ROBUSTDIR
cd ./$ROBUSTDIR
$BIN/pggsim.sh -s pgg.Defector pgg.Cooperator $CLASSES -ss 0.6 -r $ROUNDS &> ./sim.log
$BIN/lineplot.sh ./scounts.csv
cd ..

# evaluate stability
echo "  evaluating stability..."
mkdir $STABLEDIR
cd ./$STABLEDIR
$BIN/pggsim.sh -s pgg.Defector pgg.Cooperator $CLASSES -p 0 0 1 -ss 0.6 -r $ROUNDS &> ./sim.log
$BIN/lineplot.sh ./scounts.csv
cd ..

# evaluate initial viability
echo "  evaluating viability..."
mkdir $VIABLEDIR
cd ./$VIABLEDIR

# -- evaluate against defector
mkdir $DEFECTDIR
cd ./$DEFECTDIR
$BIN/pggsim.sh -s pgg.Defector pgg.Cooperator $CLASSES -p 1 0 0 -ss 0.6 -r $ROUNDS &> ./sim.log
$BIN/lineplot.sh ./scounts.csv
cd ..

# -- evaluate against cooperator
mkdir $COOPERDIR
cd ./$COOPERDIR
$BIN/pggsim.sh -s pgg.Defector pgg.Cooperator $CLASSES -p 0 1 0 -ss 0.6 -r $ROUNDS &> ./sim.log
$BIN/lineplot.sh ./scounts.csv
cd ..

# done
cd ..