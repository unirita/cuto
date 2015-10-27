#!/bin/bash

if [ "$TESTENV1" != "ENVENVENV" ] ; then
    echo Undefined TESTENV1.
    exit 12
fi
if [ "$1" != "X" ] ; then
        echo Invalid Arg#1.
        exit 12
fi
echo $TESTENV1 $1
if [ "$2" == "" ] ; then
        echo Invalid Arg#2
        exit 12
fi
exit $2

