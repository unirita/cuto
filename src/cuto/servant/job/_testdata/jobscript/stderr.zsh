#!/bin/zsh

if [ "$1" = "" ] ; then
        echo Invalid Arg#1.
        exit 12
fi
echo Argument1=$1 >&2
exit 0

