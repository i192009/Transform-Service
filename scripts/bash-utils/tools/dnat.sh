#!/bin/bash

host=$1
port=$2
forward=$3

iptables -t nat -A POSTROUTING -p tcp -m tcp -d $host/32 --dport $forward -j SNAT --to-source 172.24.1.20
iptables -t nat -A PREROUTING -p tcp -m tcp --dport $port -j DNAT --to-destination $host:$forward
