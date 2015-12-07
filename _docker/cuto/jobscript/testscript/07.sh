#!/bin/sh

#if [ $PATH != "/home/cuto" ]; then
#    echo "[err] Invalid path. - $PATH"
#    exit 1
#fi

if [ $# -ne 1 ]; then
    echo "[err] Invalid Argument Count. - $#"
    exit 2
fi

if [ $1 != "a b" ]; then
    echo "[err] Invalid Argument\$1 - $1"
    exit 3
fi

exit 0
