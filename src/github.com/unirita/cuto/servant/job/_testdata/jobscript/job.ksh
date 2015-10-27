#!/bin/ksh

if [ "$1" != "A B" ] ; then
        echo Invalid Arg#1.
        exit 12
fi
echo $TESTENV1 $1
exit 0

