#!/bin/csh

if ( "$1" == "" ) then
        echo Nothing Arg#1
        exit 1
endif
echo $1
echo $GOPATH
exit 0

