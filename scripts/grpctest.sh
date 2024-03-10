#!/bin/bash
cdir=${0%/*}

. ${0%/*}/scripts/common.sh
. ${0%/*}/scripts/debug.sh
. ${0%/*}/scripts/console.sh
. ${0%/*}/scripts/install.sh

if [ ! -f "$HOME/go/bin/grpcurl" ]; then
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
    if [ $? -ne 0 ]; then
        echo "grpcurl install failed."
        exit 1
    fi
fi

install jq

command=$(which grpcurl)
if [ -f "$HOME/go/bin/grpcurl" ]; then
    command=$HOME/go/bin/grpcurl
elif [ -z "$command" ]; then
    echo "grpcurl not found."
    exit 1
fi

# send grpc request
# usage: grpc_request server method path/file.json casename
function grpc_request() {
    local server=$1
    local method=$2
    local casefile=$3
    local casename=$4

    if [ -z "$server" ]; then
        return 1
    fi

    if [ -z "$method" ]; then
        return 1
    fi

    if [ -z "$casefile" || ! -f "$casefile" ]; then
        return 1
    fi

    if [ -z "$casename" ]; then
        return 1
    fi

    jq ".$casename" $casefile | $command -plaintext -d @ $server $method $cdir/cases.json
}

function main() {
    Menu \
        $(MenuItem "CreateJob") \
        $(MenuItem "CancelJob") \
        $(MenuItem "GetJobDetails") \
        $(MenuItem "GetJobList")


    local choice=$?
    case $choice in
        0)
            echo "Choice CreateJob"
            CreateJob localhost:9742 TransformService2.TransformV2/CreateJob CreateJob[0]
            ;;
        1)
            echo "Choice CancelJob"
            CancelJob
            ;;
        2)
            echo "Choice GetJobDetails"
            GetJobDetails
            ;;
        3)
            echo "Choice GetJobList"
            GetJobList
            ;;
    esac
}

function CreateJob() {
    echo "grpcurl -plaintext -d "${body//"/\\"}" localhost:9742 TransformService2.TransformV2/CreateJob"
}

# debug enable
main