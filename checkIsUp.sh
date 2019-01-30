#!/bin/sh

#Check if server is up

if [[ -z "${KAFKA_HOST}" ]]; then
  echo ERRO: KAFKA_HOST não informado
  exit 1
fi

if [[ -z "${KAFKA_PORT}" ]]; then
  echo ERRO: KAFKA_PORT não informado
  exit 1
fi

while [ "$(nc -z $KAFKA_HOST $KAFKA_PORT </dev/null; echo $?)" !=  "0" ];
do sleep 5;
echo "Waiting for KAFKA SERVER is UP and RESPONDING";
done

sleep 10;

./main
