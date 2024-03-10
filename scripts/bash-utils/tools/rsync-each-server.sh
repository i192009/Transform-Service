#!/bin/bash

if [ $# -lt 2 ]; then
	echo "using: rsync-each-server src dst [host ...]"
	exit 1
fi

if [ "" != "${*:3}" ]; then
	hosts="${*:3}"
fi

if [ "" == "$hosts" ]; then
	echo "hosts is empty!"
	exit 1
fi

echo "sync files to these hosts"
for host in $hosts; do
	echo "rsync -rav $1 $host:$2"
	rsync -rav $1 $host:$2
done


