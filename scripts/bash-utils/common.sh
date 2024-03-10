#!/bin/bash

function pause() {
    if (( $# == 0 )); then
        read -n1 -p 'press any to continue ...' pressed
    elif (( $# == 1 )); then
        read -n1 -t $1 -p 'press any to continue ...' pressed
    elif (( $# == 2 )); then
        read -n1 -p "$1" -t $2 pressed
    fi
    echo # 空一行
}

function get_ip() {
    echo $(ifconfig | grep '\<inet\>'| grep -v '127.0.0.1' | awk '{ print $2}' | awk 'NR==1')
}

# 获取网络地址掩码
# get_ip_mask ip [8|16|24|32]
function get_ip_mask() {
    if (( $# < 2 )); then
        return 1
    fi

    case "$2" in
        "8")
            echo $(echo $1 | awk -v FS='.' '{print $1".0.0.0/8"}')
            ;;
        "16")
            echo $(echo $1 | awk -v FS='.' -v OFS='.' '{print $1,$2".0.0/16"}')
            ;;
        "24")
            echo $(echo $1 | awk -v FS='.' -v OFS='.' '{print $1,$2,$3".0/24"}')
            ;;
        "32")
            echo $1/32
            ;;
        *)
            echo ""
            return 2
            ;;
    esac

    return 0
}

function update_hosts() {
    local IP="$1"
    local HOSTNAME="$2"
    local HOSTS_FILE="/etc/hosts"

    if [ -z "${IP}" ] || [ -z "${HOSTNAME}" ]; then
        echo "Usage: update_hosts ip hostname"
        exit 1
    fi

    if grep -q ${HOSTNAME} "${HOSTS_FILE}"; then
        sed -i'.bak' "/${HOSTNAME}/d" "${HOSTS_FILE}"
    fi

    echo "${IP} ${HOSTNAME}" | tee -a "${HOSTS_FILE}" > /dev/null
}

# Function to check if a path is in the PATH environment variable
function in_search_path() {
    local desired_path="$1"

    if [[ ":$PATH:" == *":$desired_path:"* ]]; then
        echo "The path '$desired_path' is in the PATH environment variable."
        return 0  # Return 0 for true/success
    else
        echo "The path '$desired_path' is not in the PATH environment variable."
        return 1  # Return 1 for false/failure
    fi
}

function append_to_path() {
    local desired_path="$1"

    if in_search_path "$desired_path"; then
        return 0  # Return 0 for true/success
    fi

    export PATH="$PATH:$desired_path"
    return 0  # Return 1 for false/failure
}