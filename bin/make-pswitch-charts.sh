#!/bin/sh

#  make-pswitch-charts.sh
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

# pswitch results dir
PSWI_RESULTS=$BIN/../pswitch-results

# graph report figures dir
GRAPH_REP_DIR=$BIN/../graph-structure/fig

# generate charts for partner switching
# x-axis - r
# line label - w

# generate chart for partner switching with Z=2
#Z=2
#OFILE=$GRAPH_REP_DIR/pswitch-z2.png
#echo "generating chart: z=$Z"
#$BIN/genchart.sh $PSWI_RESULTS/z$Z/* -t "z = 2" -s "W" -o $OFILE

# generate chart for partner switching with Z=4
Z=4
OFILE=$GRAPH_REP_DIR/pswitch-z4.png
echo "generating chart: z=$Z"
$BIN/genchart.sh $PSWI_RESULTS/z$Z/* -t "z = 4" -s "W" -o $OFILE

# generate chart for partner switching with Z=6
Z=6
OFILE=$GRAPH_REP_DIR/pswitch-z6.png
echo "generating chart: z=$Z"
$BIN/genchart.sh $PSWI_RESULTS/z$Z/* -t "z = 6" -s "W" -o $OFILE

# generate chart for partner switching with Z=8
Z=8
OFILE=$GRAPH_REP_DIR/pswitch-z8.png
echo "generating chart: z=$Z"
$BIN/genchart.sh $PSWI_RESULTS/z$Z/* -t "z = 8" -s "W" -o $OFILE

# generate chart for partner switching with Z=10
Z=10
OFILE=$GRAPH_REP_DIR/pswitch-z10.png
echo "generating chart: z=$Z"
$BIN/genchart.sh $PSWI_RESULTS/z$Z/* -t "z = 10" -s "W" -o $OFILE

# generate one chart for w=1 with all values of z
# x-axis - r
# line label - g
SERIES=$PSWI_RESULTS/z4/w01
SERIES="$SERIES $PSWI_RESULTS/z6/w01"
SERIES="$SERIES $PSWI_RESULTS/z8/w01"
SERIES="$SERIES $PSWI_RESULTS/z10/w01"

OFILE=$GRAPH_REP_DIR/w1all.png
echo "generating chart: W=1"
$BIN/genchart.sh $SERIES -t "W = 1" -s "z" -o $OFILE

# generate one chart for w=3 with all values of z
# x-axis - r
# line label - g
SERIES=$PSWI_RESULTS/z4/w03
SERIES="$SERIES $PSWI_RESULTS/z6/w03"
SERIES="$SERIES $PSWI_RESULTS/z8/w03"
SERIES="$SERIES $PSWI_RESULTS/z10/w03"

OFILE=$GRAPH_REP_DIR/w3all.png
echo "generating chart: W=3"
$BIN/genchart.sh $SERIES -t "W = 3" -s "z" -o $OFILE

# generate one chart for w=5 with all values of z
# x-axis - r
# line label - g
SERIES=$PSWI_RESULTS/z4/w05
SERIES="$SERIES $PSWI_RESULTS/z6/w05"
SERIES="$SERIES $PSWI_RESULTS/z8/w05"
SERIES="$SERIES $PSWI_RESULTS/z10/w05"

OFILE=$GRAPH_REP_DIR/w5all.png
echo "generating chart: W=5"
$BIN/genchart.sh $SERIES -t "W = 5" -s "z" -o $OFILE

# generate one chart for w=7 with all values of z
# x-axis - r
# line label - g
SERIES=$PSWI_RESULTS/z4/w07
SERIES="$SERIES $PSWI_RESULTS/z6/w07"
SERIES="$SERIES $PSWI_RESULTS/z8/w07"
SERIES="$SERIES $PSWI_RESULTS/z10/w07"

OFILE=$GRAPH_REP_DIR/w7all.png
echo "generating chart: W=7"
$BIN/genchart.sh $SERIES -t "W = 7" -s "z" -o $OFILE

# generate one chart for w=9 with all values of z
# x-axis - r
# line label - g
SERIES=$PSWI_RESULTS/z4/w09
SERIES="$SERIES $PSWI_RESULTS/z6/w09"
SERIES="$SERIES $PSWI_RESULTS/z8/w09"
SERIES="$SERIES $PSWI_RESULTS/z10/w09"

OFILE=$GRAPH_REP_DIR/w9all.png
echo "generating chart: W=9"
$BIN/genchart.sh $SERIES -t "W = 9" -s "z" -o $OFILE

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


