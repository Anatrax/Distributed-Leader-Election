#!/bin/bash
rm output.txt
pkill -f better_grpc_base

go build

x=$1
if [ $# -eq 0 ]; then
    x=10
fi

echo "TEST: Leaders automatically crash after sending out a few heartbeats"
echo "Press ENTER when you're ready to terminate the program"
sleep 3

for ((i=0;i<$x;i++))
do
    echo "Initializing node ($((i+1))/$x)..."
	./better_grpc_base $x $i >> output.txt &
    sleep 1
done

tail -f output.txt &

read quit
./stop.sh
