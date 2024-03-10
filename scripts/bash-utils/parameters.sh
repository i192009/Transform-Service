#!/bin/bash
. $(dirname "${BASH_SOURCE[0]}")/console.sh

# 读取参数
# ReadParameter $name $default
function ReadParameter() {
    trace name=$1, $(eval "echo \$$1")

    local default=$(eval "echo \$$1")

    if [ -n "$2" ]; then
        default=$2
    fi

    read -p "$1($default): " input
    if [ -n "$input" ]; then 
        export $1=$input
    else
        export $1=$default
    fi

    trace $(eval "echo \$$1")
}

# 列出所有参数
function ListParameters() {
    trace "parameters=${parameters[@]}"

    for ele in ${parameters[@]}; do
        echo $ele=$(eval echo \$"$ele")
    done
}

# 根据参数名显示参数值
function ShowParamerter() {
    trace "parameter=$1"
    echo -e "$(FB 45)$1$(FB 15)=$(eval "echo \$$1")"
}

# 检查必要的参数，若参数为空则要求用户输入。
function ValidateParameters() {
    local parameters=($@)
    for parameter in ${parameters[@]}; do
        local value=$(eval "echo \$$parameter")
        if [ -z $value ]; then
            echo "$parameter is empty."
            return 1
        fi
    done
}

# 使用菜单的方式设置变量值
function setting() {
    local parameters=($@)

    while true; do
        Menu -f ShowParamerter ${parameters[@]}
        local choice=$?
        if (( $choice == 255 )); then
            break
        fi

        ReadParameter ${parameters[$choice]}
    done
}

# 检查必要的参数，若参数为空则要求用户输入。
function setting_validate() {
    local parameters=($@)
    for parameter in ${parameters[@]}; do
        local value=$(eval "echo \$$parameter")
        if [ -z $value ]; then
            read -p "$parameter=" $parameter
        fi
    done

    ListParameters ${parameters[@]}
}