#!/bin/bash
. $(dirname "${BASH_SOURCE[0]}")/debug.sh
. $(dirname "${BASH_SOURCE[0]}")/common.sh
. $(dirname "${BASH_SOURCE[0]}")/install.sh

function MysqlExec() {
    trace "mysql execute paths=($@)"
    for path in $@; do
        trace "current path = $path"
        if [ -d $path ]; then
            trace "$path is a directory"
            if [ -e $path/.mysqlfiles ]; then
                MysqlExec $(sed "s/^/${path//\//\\\/}\/&/g" $path/.mysqlfiles)
            else
                MysqlExec $(find $path | grep -E ".*\.sql$")
            fi
        elif [ -f $path ]; then
            trace "$path is a file"
            more $path
            read -n1 -p "execute this file?[$path](Y/n)" answer
            echo #空一行
            if [[ "$answer" == "" || "$answer" == "y" || "$answer" == "Y" ]]; then
                trace "run command: cat $path | mycli -h $MYSQL_LOCAL_HOST -P $MYSQL_LOCAL_PORT -u root -p$MYSQL_ROOT_PASSWORD -e \"source $path\" --"
                mycli -h $MYSQL_LOCAL_HOST -P $MYSQL_LOCAL_PORT -u root -p$MYSQL_ROOT_PASSWORD -e "source $path" --
                if [ 0 -ne $? ]; then
                    error "execute sql file failed, file = $path; command = cat $path | mycli -h $MYSQL_LOCAL_HOST -P $MYSQL_LOCAL_PORT -u root -p$MYSQL_ROOT_PASSWORD -e --"
                    return 1
                fi
            fi
        else
            echo "file not found, path = $path"
        fi
    done
}

function InitMysql() {
    if [ ! -f "/usr/bin/mycli" ]; then
        install mycli
        if [ ! -f "/usr/bin/mycli" ]; then
            echo "mycli install failed."
            return 1
        fi
    fi

    MysqlExec $@
}

function MongoExec() {
    trace "mongo execute paths=\($@\)"
    for path in $@; do
        if [ -d $path ]; then
            if [ -e $path/.mongofiles ]; then
                MongoExec $(sed "s/^/${path//\//\\\/}\/&/g" $path/.mongofiles)
            else
                MongoExec $(find $path | grep -E ".*\.js$")
            fi
        elif [ -f $path ]; then
            trace "$path is a file"
            more $path
            read -n1 -p "execute this file?[$path](Y/n)" answer
            echo #空一行
            if [[ "$answer" == "" || "$answer" == "y" || "$answer" == "Y" ]]; then
                info "command = mongosh -f $path mongodb://$MONGO_ROOT_USER:$MONGO_ROOT_PASSWORD@$MONGO_LOCAL_HOST:$MONGO_LOCAL_PORT?authSource=admin"
                mongosh -f "$path" "mongodb://$MONGO_ROOT_USER:$MONGO_ROOT_PASSWORD@$MONGO_LOCAL_HOST:$MONGO_LOCAL_PORT?authSource=admin"
                if [ 0 -ne $? ]; then
                    error "execute mongo script file failed, file = $path"
                    pause
                    return 1
                fi
            fi
        else
            echo "file not found, path = $path"
            pause
        fi
    done
}

function InitMongo() {
    install wget
    if [ ! -f "/usr/bin/mongosh" ]; then
        wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | apt-key add -
        echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/6.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-6.0.list
        apt update
        install mongodb-mongosh

        if [ ! -f "/usr/bin/mongosh" ]; then
            echo "mongosh install failed."
            return 1
        fi
    fi

    export MONGO_ROOT_USER=$MONGO_ROOT_USER
    export MONGO_ROOT_PASSWORD=$MONGO_ROOT_PASSWORD
    export MONGO_LOCAL_HOST=$MONGO_LOCAL_HOST
    export MONGO_LOCAL_PORT=$MONGO_LOCAL_PORT
    MongoExec $@
    export -n MONGO_ROOT_USER
    export -n MONGO_ROOT_PASSWORD
    export -n MONGO_LOCAL_HOST
    export -n MONGO_LOCAL_PORT
}