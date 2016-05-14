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
GRAPH_REP_DIR=$BIN/../graph-structure/fig/fixed

# generate one chart for each graph type
# x-axis - r
# line label - z

# generate chart for regular graphs
GTYPE=0
OFILE=$GRAPH_REP_DIR/regular.png
echo "generating chart: Regular ($GTYPE)"
$BIN/genchart.sh $GPGG_RESULTS/gtype$GTYPE/* -t "Regular" -s "z" -o $OFILE

# generate chart for homogeneous random
GTYPE=1
OFILE=$GRAPH_REP_DIR/homorand.png
echo "generating chart: Homogeneous Random ($GTYPE)"
$BIN/genchart.sh $GPGG_RESULTS/gtype$GTYPE/* -t "Homogeneous Random" -s "z" -o $OFILE

# generate chart for heterogeneous random
TITLE="Heterogeneous Random"
GTYPE=2
OFILE=$GRAPH_REP_DIR/heterand.png
echo "generating chart: Heterogeneous Random ($GTYPE)"
$BIN/genchart.sh $GPGG_RESULTS/gtype$GTYPE/* -t "Small World (p=1.0)" -s "z" -o $OFILE

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

# generate chart for small world (p=0.1)
GTYPE=5
OFILE=$GRAPH_REP_DIR/sw01.png
echo "generating chart: Small World (p=0.1) ($GTYPE)"
$BIN/genchart.sh $GPGG_RESULTS/gtype$GTYPE/* -t "Small World (p=0.1)" -s "z" -o $OFILE

# generate chart for small world (p=0.4)
GTYPE=6
OFILE=$GRAPH_REP_DIR/sw04.png
echo "generating chart: Small World (p=0.4) ($GTYPE)"
$BIN/genchart.sh $GPGG_RESULTS/gtype$GTYPE/* -t "Small World (p=0.4)" -s "z" -o $OFILE

# generate one chart for z=4 with all small world graph types
SERIES=$GPGG_RESULTS/gtype5/z04
SERIES="$SERIES $GPGG_RESULTS/gtype6/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype2/z04"

OFILE=$GRAPH_REP_DIR/smallworld.png
echo "generating chart: small world"
$BIN/genchart.sh $SERIES -t "Small World Networks (z=4)" -s "gtype" -o $OFILE

# generate one chart for z=4 with all scale free graph types
SERIES=$GPGG_RESULTS/gtype3/z04
SERIES="$SERIES $GPGG_RESULTS/gtype4/z04"

OFILE=$GRAPH_REP_DIR/scalefreeall.png
echo "generating chart: all scale free"
$BIN/genchart.sh $SERIES -t "Scale Free Networks (z=4)" -s "gtype" -o $OFILE

# generate one chart for z=4 with all graph types
# x-axis - r
# line label - g
SERIES=$GPGG_RESULTS/gtype0/z04
SERIES="$SERIES $GPGG_RESULTS/gtype1/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype5/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype6/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype2/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype3/z04"
SERIES="$SERIES $GPGG_RESULTS/gtype4/z04"

OFILE=$GRAPH_REP_DIR/z4all.png
echo "generating chart: z=4"
$BIN/genchart.sh $SERIES -t "z = 4" -s "gtype" -o $OFILE

# generate one chart for z=6 with all graph types
# x-axis - r
# line label - g
SERIES=$GPGG_RESULTS/gtype0/z06
SERIES="$SERIES $GPGG_RESULTS/gtype1/z06"
SERIES="$SERIES $GPGG_RESULTS/gtype5/z06"
SERIES="$SERIES $GPGG_RESULTS/gtype6/z06"
SERIES="$SERIES $GPGG_RESULTS/gtype2/z06"
SERIES="$SERIES $GPGG_RESULTS/gtype3/z06"
SERIES="$SERIES $GPGG_RESULTS/gtype4/z06"

OFILE=$GRAPH_REP_DIR/z6all.png
echo "generating chart: z=6"
$BIN/genchart.sh $SERIES -t "z = 6" -s "gtype" -o $OFILE

# generate one chart for z=8 with all graph types
# x-axis - r
# line label - g
SERIES=$GPGG_RESULTS/gtype0/z08
SERIES="$SERIES $GPGG_RESULTS/gtype1/z08"
SERIES="$SERIES $GPGG_RESULTS/gtype5/z08"
SERIES="$SERIES $GPGG_RESULTS/gtype6/z08"
SERIES="$SERIES $GPGG_RESULTS/gtype2/z08"
SERIES="$SERIES $GPGG_RESULTS/gtype3/z08"
SERIES="$SERIES $GPGG_RESULTS/gtype4/z08"

OFILE=$GRAPH_REP_DIR/z8all.png
echo "generating chart: z=8"
$BIN/genchart.sh $SERIES -t "z = 8" -s "gtype" -o $OFILE

# generate one chart for z=10 with all graph types
# x-axis - r
# line label - g
SERIES=$GPGG_RESULTS/gtype0/z10
SERIES="$SERIES $GPGG_RESULTS/gtype1/z10"
SERIES="$SERIES $GPGG_RESULTS/gtype5/z10"
SERIES="$SERIES $GPGG_RESULTS/gtype6/z10"
SERIES="$SERIES $GPGG_RESULTS/gtype2/z10"
SERIES="$SERIES $GPGG_RESULTS/gtype3/z10"
SERIES="$SERIES $GPGG_RESULTS/gtype4/z10"

OFILE=$GRAPH_REP_DIR/z10all.png
echo "generating chart: z=10"
$BIN/genchart.sh $SERIES -t "z = 10" -s "gtype" -o $OFILE