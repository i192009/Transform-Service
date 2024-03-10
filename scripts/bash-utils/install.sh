#!/bin/bash
. $(dirname "${BASH_SOURCE[0]}")/debug.sh
. $(dirname "${BASH_SOURCE[0]}")/common.sh
. $(dirname "${BASH_SOURCE[0]}")/parameters.sh

function installed() {
    local pkg=$1
    dpkg -l $pkg > /dev/null 2>&1
    return $?
}

function install() {
    if [ "$debug_flags" == "1" ]; then
        echo apt install -y $@
        return 0
    else
        for pkg in $@; do
            installed $pkg
            if [ "$?" == "0" ]; then
                info "$pkg is already installed"
                continue
            else
                apt install -y $@
            fi
        done
        return $?
    fi
}

function install_python_package() {
    local retry=0
    local package=$1
    while [ $retry -lt 3 ]; do
        # Try to install the package
        echo "pip install \"$package\""
        pip install "$package"
        status=$?

        # Check if installation was successful
        if [ $status -eq 0 ]; then
            echo "Installation successful."
            break
        fi

        # Wait for a bit before retrying
        read -t 5 -n 1 -p "Installation failed, retrying in 5 seconds..." key

        # Check for Esc key press
        if [ $? -eq 0 ]; then
            echo "Installation aborted."
            break
        fi
        retry = $((retry + 1))
    done
}

function install_nfs_server() {
    # 创建nfs挂载目录
    if [ ! -e $NFS_SERVER_PATH ]; then
        mkdir -p $NFS_SERVER_PATH
        if [ 0 -ne $? ]; then 
            error "make dir error, path = $NFS_SERVER_PATH"
            return 1; fi

        chmod a+w $NFS_SERVER_PATH
        if [ 0 -ne $? ]; then
            error "change directory mode error, path = $NFS_SERVER_PATH"
            return 2;
        fi
    fi

    # install nfs server
    install nfs-kernel-server
    if [ 0 -ne $? ]; then
        error "install nfs-kernel-server error"
        return 3; 
    fi

    # search and delete whole line
    sed -i "/^${NFS_SERVER_PATH//\//\\\/}\\t\\t${NFS_SERVER_EXPORT//\//\\\/}.*/d" /etc/exports
    # append new line
    echo -e "${NFS_SERVER_PATH}\t\t${NFS_SERVER_EXPORT}(rw,no_root_squash,insecure,sync)" >> /etc/exports

    trace "---------------------------------------"
    trace "$(cat /etc/exports)"

    if [ "enabled" == "$(systemctl is-enabled nfs-server)" ]; then
        info "restart nfs-server"
        systemctl restart nfs-server
    else
        info "enable nfs-server"
        systemctl enable nfs-server
        systemctl start nfs-server
    fi
}

function remote_command() {
    local command="$@"
    for client in $(cat $K8S_NODE_LIST) ; do
        trace "ssh -i $K8S_NODE_PRIVATE_KEY $K8S_NODE_USER@$client" $command
        ssh -i $K8S_NODE_PRIVATE_KEY $K8S_NODE_USER@$client $command
    done
}

function remote_install() {
    remote_command 'apt-get install -y '$@
}