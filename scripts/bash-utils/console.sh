#!/bin/bash
. $(dirname "${BASH_SOURCE[0]}")/debug.sh

Color256() {
    # This program is free software. It comes without any warranty, to
    # the extent permitted by applicable law. You can redistribute it
    # and/or modify it under the terms of the Do What The Fuck You Want
    # To Public License, Version 2, as published by Sam Hocevar. See
    # http://sam.zoy.org/wtfpl/COPYING for more details.
    
    for fgbg in 38 48 ; do # Foreground / Background
        for color in {0..255} ; do # Colors
            # Display the color
            printf "\e[${fgbg};5;%sm  %3s  \e[0m" $color $color
            # Display 6 colors per lines
            if [ $((($color + 1) % 16)) == 4 ] ; then
                echo # New line
            fi
        done
        echo # New line
    done
}

#关闭所有属性
function RST() {
    echo -ne "\e[0m"
}

#设置粗体
function SET_BOLD() {
    echo -ne "\e[1m"
}

function CLR_BOLD() {
    echo -ne "\e[21m"
}

#设置暗淡
function SET_DIM() {
    echo -ne "\e[2m"
}

function CLR_DIM() {
    echo -ne "\e[22m"
}

#设置高亮度
function SET_LIGHT() {
    echo -e "\e[1m"
}

function CLR_LIGHT() {
    echo -ne "\e[21m"
}

#下划线
function SET_UNDERLINE() {
    echo -ne "\e[4m"
}

function CLR_UNDERLINE() {
    echo -ne "\e[24m"
}

#闪烁
function SET_BLINK() {
    echo -ne "\e[5m"
}

function CLR_BLINK() {
    echo -ne "\e[25m"
}

#反显，撞色显示，显示为白字黑底，或者显示为黑底白字
function SET_INV() {
    echo -ne "\e[7m"
}

function CLR_INV() {
    echo -e "\e[27m"
}

#消影，字符颜色将会与背景颜色相同
function SET_OPT() {
    echo -ne "\e[8m"
}

function CLR_OPT() {
    echo -ne "\e[28m"
}

#光标上移 n 行
function CUR_U() {
    echo -ne "\e[$1A"
}

#光标下移 n 行
function CUR_D() {
    echo -ne "\e[$1B"
}

#光标右移 n 行
function CUR_L() {
    echo -ne "\e[$1C "
}

#光标左移 n 行
function CUR_R() {
    echo -ne "\e[$1D "
}

#光标回到行首
function CUR_H() {
    echo -ne "\e[H"
}

#设置光标位置
function CUR_P() {
    echo -ne "\e[$2;$1H"
}

#删除光标后的内容
function DEL_A() {
    echo -ne "\e[0K"
}

#删除光标前的内容
function DEL_B() {
    echo -ne "\e[1K"
}

#删除光标所在行的内容
function DEL_L() {
    echo -ne "\e[2K"
}

#清屏
function SCR_CLR() {
    echo -ne "\e[2J"
}
#保存光标位置
function CUR_SAVE() {
    echo -ne "\e[s"
}

#恢复光标位置
function CUR_LOAD() {
    echo -ne "\e[u"
}

#隐藏光标
function CUR_HIDE() {
    echo -ne "\e[?25"
}

#显示光标
function CUR_SHOW() {
    echo -ne "\e[?25h"
}

#清除从光标到行尾的内容
function CUR_D2E() {
    echo -ne "\e[K"
}

# color text output
# TEXT content [fg] [bg]
function TEXT() {
    local ret=$1$(RST)
    if [[ $# > 1 ]]; then
        ret="\e[38;5;$2m"$ret
    fi

    if [[ $# > 2 ]]; then
        ret="\e[48;5;$3m"$ret
    fi

    echo -ne $ret
}

# set fg and bg
# FB [fg] [bg]
# if no any parameters, this function will reset the color to default
function FB() {
    case "$#" in
        "0")
            echo -ne "\e[0m"
            ;;
        "1")
            echo -ne "\e[38;5;$1m"
            ;;
        "2")
            echo -ne "\e[38;5;$1m\e[48;5;$2m"
            ;;
    esac
}

# set fg
# FG fg
function FG() {
    echo -ne "\e[38;5;$1m"
}

# set bg
# BG bg
function BG() {
    echo -ne "\e[48;5;$1m"
}

function Menu() {
    local choice=0

    local key=""
    local key_status=0
    local page=10
    local func=""
    
    OPTIND=1
    while getopts "f:p:" opt; do
        case $opt in
            f)
                func=$OPTARG
                ;;
            p)
                page=$OPTARG
                ;;
        esac
    done

    trace "optind = $OPTIND"
    shift $((OPTIND-1))
    local options=($@)

    trace "MenuItem callback function $func"
    trace "options=${#options[@]},${options[@]}"
    while true; do
        trace "key status = $key_status"
        if [[ $key_status == 0 ]]; then
            local start=$[choice / page * page]
            local thispage=(${options[@]:$start:$page})
            trace $start, $page, ${thispage[@]}
            # show menu items
            for idx in "${!thispage[@]}"; do
                if [ -n "$func" ]; then
                    trace "call $func ${thispage[$idx]}"
                    local option=`$func ${thispage[$idx]}`
                else
                    trace "show ${thispage[$idx]}"
                    local option=${thispage[$idx]//_/ }
                fi

                if [[ $idx == $[choice % page] ]]; then
                    echo -e "$(DEL_L)[$(FG 1)X$(FB)]" $(FB 15 7)"$idx - $option"$(FB)
                else
                    echo -e "$(DEL_L)[ ] $idx -" $option
                fi
            done

            # clear the rest of the page
            for (( idx= idx+1; idx < page; idx++ )); do
                echo `DEL_L`
            done

            echo "-------------------------------------"
            echo "$(DEL_L) Current choice is ${options[$choice]}"
            read -sn1 key
            if [ "$key" == 'q' ]; then
                trace "quit"
                return -1
            else
                CUR_U `expr $page + 2`
            fi
        fi

        # process the pressed key
        if [[ $key_status == 9 ]]; then
            if [ "$key" == "" ]; then
                trace 'pressed enter key'
                local ret=`expr index "0123456789" $[choice % page]`
                trace "ret = $ret"
                if (( $ret != 0 )); then
                    choice=$((start + ret - 1))
                    trace "choice = $choice"
                    key_status=0
                    CUR_D `expr $page + 2`
                    return $choice
                fi
            else
                trace 'pressed other key' $key
                local ret=`expr index "0123456789" $key`
                trace ret=$ret, $key

                if (( $ret >= 0 )); then
                    choice=$((start + ret - 1))
                    if (( $choice >= ${#options[@]} )); then
                        choice=`expr ${#options[@]} - 1`
                    fi

                    key_status=0
                    continue
                fi
            fi
        fi

        # echo 1 - $key_status
        if [[ $key_status == 0 && "$key" == $'\e' ]]; then
            key_status=1
            read -sn1 key
        else
            key_status=9
        fi

        # echo 2 - $key_status
        if [[ $key_status == 1 && "$key" == "[" ]]; then
            key_status=2
            read -sn1 key
        else
            key_status=9
        fi

        # echo 3 - $key_status
        if [[ $key_status == 2 ]]; then
            case $key in
                "A") # up key pressed
                    choice=$[choice - 1]
                    ;;
                "B") # down key pressed
                    choice=$[choice + 1]
                    ;;
                "C") # right key pressed
                    choice=$[choice + page]
                    ;;
                "D") # left key pressed
                    choice=$[choice - page]
                    ;;
                *)
                    ;;
            esac

            if (( $choice < 0 )); then
                choice=0
            elif (( $choice >= ${#options[@]} )); then
                choice=`expr ${#options[@]} - 1`
            fi

            trace "read key = $key; current choice is $choice"
            key_status=0
        fi
    done
}

function MenuItem() {
    echo -ne \"${@// /_}\"
}