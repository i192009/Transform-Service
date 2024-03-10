#!/bin/bash
cdir=${0%/*}
echo $cdir
. $cdir/bash-utils/common.sh
. $cdir/bash-utils/install.sh

install python3
install_python_package grpcio
install_python_package grpcio-tools
install_python_package temporalio