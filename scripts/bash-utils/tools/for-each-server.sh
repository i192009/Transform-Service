#!/bin/bash

for host in $(cat nodes); do
    ssh $host $@
done
