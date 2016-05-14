#!/bin/sh

#  gen-scalefree.sh
#  
#
#  Created by John Maloney on 5/13/16.
#

#!/bin/sh

#  make-charts.sh
#
#
#  Created by John Maloney on 5/13/16.
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

# graph report figures dir
GRAPH_REP_DIR=$BIN/../graph-structure/fig

# generate one chart for each graph type
# x-axis - r
# line label - z

# generate chart for scale free
GTYPE=3
OFILE=$GRAPH_REP_DIR/sfree.png
echo "generating chart: Scale Free ($GTYPE)"
$BIN/genchart.sh $GPGG_RESULTS/gtype$GTYPE/* -t "Scale Free" -s "z" -o $OFILE

# generate chart for scale free
GTYPE=4
OFILE=$GRAPH_REP_DIR/usfree.png
echo "generating chart: Uniform Scale Free ($GTYPE)"
$BIN/genchart.sh $GPGG_RESULTS/gtype$GTYPE/* -t "Uniform Scale Free" -s "z" -o $OFILE

# generate one chart for z=4 with all scale free graph types
SERIES=$GPGG_RESULTS/gtype3/z04
SERIES="$SERIES $GPGG_RESULTS/gtype4/z04"

OFILE=$GRAPH_REP_DIR/scalefreeall.png
echo "generating chart: all scale free"
$BIN/genchart.sh $SERIES -t "Scale Free Networks (z=4)" -s "gtype" -o $OFILE
