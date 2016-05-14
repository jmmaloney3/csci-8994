#!/bin/sh

#  make-z4all-chart.sh
#  
#
#  Created by John Maloney on 5/13/16.
#

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

# graph report figures dir
GRAPH_REP_DIR=$BIN/../graph-structure/fig

# generate one chart for z=4 with all graph types
# x-axis - r
# line label - g
SERIES=$GPGG_RESULTS/gtype0/z04
SERIES="$SERIES $GPGG_RESULTS/gtype1/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype2/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype3/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype4/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype5/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype6/z04"

OFILE=$GRAPH_REP_DIR/z4all.png
echo "generating chart: z=4"
$BIN/genchart.sh $SERIES -t "z = 4" -s "gtype" -o $OFILE