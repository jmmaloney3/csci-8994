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

# pswitch results dir
PSWI_RESULTS=$BIN/../pswitch-results

# graph report figures dir
GRAPH_REP_DIR=$BIN/../graph-structure/fig

# generate charts for partner switching
# x-axis - r
# line label - w

# generate chart for partner switching with Z=2
Z=2
OFILE=$GRAPH_REP_DIR/pswitch-z2.png
echo "generating chart: z=$Z"
$BIN/genchart.sh $PSWI_RESULTS/z$Z/* -t "z = 2" -s "W" -o $OFILE

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
