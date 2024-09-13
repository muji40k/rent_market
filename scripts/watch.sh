#! /bin/bash

if [ "function" != "$(LC_ALL=C type -t task)" ]; then
    echo "Task not set (or not exported)" 1>&2
    exit 1
fi

if [ -z "${ROOT+x}" ]; then
    echo "Root not set, set default as './'" 1>&2
    ROOT=./
fi

function get_modification_time {
    stat --format="%Y" "$1"
}

function get_last_modified {
    dirs=($(find $ROOT -newer "$1" -type d))

    if [ 0 -eq ${#dirs[@]} ]; then
        echo $1
    else
        get_last_modified ${dirs[0]}
    fi
}

function print_log {
    date +"# [%R] $1" 1>&2
}

print_log "Start watcher in directory '$ROOT'"
task
modified=$(get_last_modified $ROOT)
time=$(get_modification_time "$modified")

while sleep 1
do
    cmodified=$(get_last_modified "$modified")

    if [[ "$cmodified" != "$modified" ]]; then
        print_log "Update"
        task
        modified=$cmodified
        time=$(get_modification_time $modified)
    elif [ ! -e $modified ]; then
        print_log "Update"
        task
        modified=$(get_last_modified $ROOT)
        time=$(get_modification_time $modified)
    elif [ $(get_modification_time $modified) -gt $time ]; then
        print_log "Update"
        task
        time=$(get_modification_time $modified)
    fi
done

