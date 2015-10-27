#!/bin/csh

if ( "$TESTENV1" != "ENVENVENV" ) then
    echo Undefined TESTENV1.
    exit 12
endif
echo $TESTENV1
if ( "$1" != "" )  then
        echo Invalid Arg#1.
        exit 12
endif

exit 0

