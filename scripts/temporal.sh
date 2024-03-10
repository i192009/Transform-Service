#!/bin/bash
echo ${0%/*}
. ${0%/*}/scripts/common.sh
. ${0%/*}/scripts/debug.sh
. ${0%/*}/scripts/console.sh

temporal=$HOME/.temporalio/bin/temporal
cdir=$(dirname $0)

function start() {
    nohup temporal server start-dev --ui-port 8233 --db-filename temporal.db 1>$cdir/temporal.log 2>&1 &
    # Try to find the Temporal server's process ID (PID)
    pid=$(pgrep -f 'temporal server start-dev')
    if [ ! -z "$pid" ]; then
        echo "temporal server started."
        echo "you can connect to temporal server by running: temporal cli or localhost:7233 in sdk."
        echo "you can visit http://localhost:8233 to see the temporal web."
    fi
}

function stop() {
    # Path to the log file where the Temporal server process ID might be stored
    local logfile="$1"

    if [ -z "$logfile" ]; then
        logfile=$cdir/temporal.log
    fi

    # Check if the log file exists
    if [ ! -f "$logfile" ]; then
        echo "Log file not found. Cannot determine the Temporal server process."
        exit 1
    fi

    # Try to find the Temporal server's process ID (PID)
    pid=$(pgrep -f 'temporal server start-dev')

    if [ -z "$pid" ]; then
        echo "Temporal server process not found."
        exit 1
    else
        # Terminate the Temporal server process
        kill -9 "$pid"
        if [ $? -eq 0 ]; then
            echo "Temporal server stopped successfully."
        else
            echo "Failed to stop Temporal server."
            exit 1
        fi
    fi
}

if [ ! -f "$temporal" ]; then
    echo temporal not install at $file
    curl -sSf https://temporal.download/cli.sh | sh
fi

append_to_path "$HOME/.temporalio/bin"

if [ "$1" == "start" ]; then
    start
    exit 0
elif [ "$1" == "stop" ]; then
    stop
    exit 0
fi

Menu \
    $(MenuItem "start temporal server") \
    $(MenuItem "stop temporal server")

case $? in
    0)
        start
        ;;
    1)
        stop
        ;;
    *)
        exit 1
        ;;
esac


