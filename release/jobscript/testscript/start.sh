#!/bin/sh

if [ $# -lt 2 ]; then
    if [ $# -lt 1 ]; then
        echo "Nothing ID."
        exit 101
    fi
    echo "Nothing start-date."
    exit 102
fi

CUR=`pwd`
cd `dirname $0`
cd ../../
VERCUR=`pwd`

# master environment $HOME
if [ $TESTENV1 != $HOME ]; then
    echo "Invalid TESTENV1[ $TESTENV1 ]"
    exit 103
fi

# servant environment $LANG
if [ $TESTENV2 != $LANG ]; then
    echo "Invalid TESTENV2[ $TESTENV2 ]"
    exit 104
fi

if [ $VERCUR != $CUR ]; then
    echo "Invalid Current-dir[ $CUR ]"
    exit 105
fi

exit 0
