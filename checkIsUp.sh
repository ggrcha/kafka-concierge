#!/bin/sh

#Check if server is up

while [ "$(nc -z $KAFKA_HOST $KAFKA_PORT </dev/null; echo $?)" !=  "0" ];
do sleep 5;
echo "Waiting for KAFKA SERVER is UP and RESPONDING";
done

sleep 10;

./main
