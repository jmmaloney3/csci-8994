#!/bin/sh

#  Script.sh
#  
#
#  Created by John Maloney on 5/14/16.
#


# make the charts for the reports

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

# gpgg results dir
GPGG_RESULTS=$BIN/../gpgg-results

# pswitch results dir
PSWI_RESULTS=$BIN/../pswitch-results

# graph report figures dir
GRAPH_REP_DIR=$BIN/../graph-structure/fig/dynamic

# generate one chart for z=4 with w=9
# x-axis - r
# line label - g
SERIES=$PSWI_RESULTS/z4/w09
SERIES="$SERIES $GPGG_RESULTS/gtype1/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype3/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype4/z04"

OFILE=$GRAPH_REP_DIR/z4dynfix.png
echo "generating chart: z=4 (dymanic and fixed)"
$BIN/genchart.sh $SERIES -t "z = 4 and W = 9" -s "gtype" -o $OFILE

# generate one chart for z=6 with w=9
# x-axis - r
# line label - g
SERIES=$PSWI_RESULTS/z6/w09
SERIES="$SERIES $GPGG_RESULTS/gtype1/z06"
SERIES="$SERIES $GPGG_RESULTS/gtype3/z06"
SERIES="$SERIES $GPGG_RESULTS/gtype4/z06"

OFILE=$GRAPH_REP_DIR/z6dynfix.png
echo "generating chart: z=6 (dymanic and fixed)"
$BIN/genchart.sh $SERIES -t "z = 6 and W = 9" -s "gtype" -o $OFILE

# generate one chart for z=8 with w=9
# x-axis - r
# line label - g
SERIES=$PSWI_RESULTS/z8/w09
SERIES="$SERIES $GPGG_RESULTS/gtype1/z08"
SERIES="$SERIES $GPGG_RESULTS/gtype3/z08"
SERIES="$SERIES $GPGG_RESULTS/gtype4/z08"

OFILE=$GRAPH_REP_DIR/z8dynfix.png
echo "generating chart: z=8 (dymanic and fixed)"
$BIN/genchart.sh $SERIES -t "z = 8 and W = 9" -s "gtype" -o $OFILE

# generate one chart for z=10 with w=9
# x-axis - r
# line label - g
SERIES=$PSWI_RESULTS/z10/w09
SERIES="$SERIES $GPGG_RESULTS/gtype1/z10"
SERIES="$SERIES $GPGG_RESULTS/gtype3/z10"
SERIES="$SERIES $GPGG_RESULTS/gtype4/z10"

OFILE=$GRAPH_REP_DIR/z10dynfix.png
echo "generating chart: z=10 (dymanic and fixed)"
$BIN/genchart.sh $SERIES -t "z = 10 and W = 9" -s "gtype" -o $OFILE

