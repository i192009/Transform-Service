#!/bin/bash

cdir=$(cd `dirname $0`; pwd)

if [ ! -d build ]; then
    mkdir -p $cdir/build
fi

cd $cdir/build
cmake ..
make

cd $cdir