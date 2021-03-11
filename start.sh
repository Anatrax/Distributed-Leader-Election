#!/bin/bash

# Remove any existing output files or nodes from previous runs
rm output.txt
pkill -f leader_election

# Rebuild golang code
go build

# Get number of nodes to start up (Defaults to 7)
x=$1
if [ $# -eq 0 ]; then
    x=7
fi

# Let the user know how to exit this script properly
echo "TEST: Leaders automatically crash after sending out a few heartbeats"
echo "Press ENTER when you're ready to terminate the program"
sleep 3

# Disable Ctrl+C
trap '' 2

for ((i=0;i<$x;i++))
do
    echo "Initializing node ($((i+1))/$x)..."
	./leader_election $x $i >> output.txt &
    sleep 1
done

# Display output logs to the user
tail -f output.txt &

# Wait for the user to press ENTER before running the stop script to kill the nodes
read quit
./stop.sh

# Re-enable Ctrl+C
trap 2
