#!/bin/bash
function debug() {
    if [ "$1" == "enable" ]; then
        debug_flags_bak=$debug_flags
        debug_flags=1
    elif [ "$1" == "1" ]; then
        debug_flags_bak=$debug_flags
        debug_flags=1
    elif [ "$1" == "resume" ]; then
        debug_flags=$debug_flags_bak
    else
        debug_flags_bak=$debug_flags
        debug_flags=0
    fi
}

# 调试日志输出
function trace() {
    if [ -z $debug_flags ]; then
        return 0
    fi
    
    if [ $debug_flags -eq 1 ]; then
        echo -e $(FB 6)"$@"$(FB)
    fi
}

function error() {
    echo -e $(FB 1)"$@"$(FB)
}

function info() {
    echo -e $(FB 3)"$@"$(FB)
}