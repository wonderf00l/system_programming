#!/bin/bash

N=1000

for i in $(seq 1 $N); do
    echo "Запуск номер $i"
    # go run main.go
    ./exec
done