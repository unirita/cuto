#!/bin/sh

CUTOROOT=@ROOT;export CUTOROOT

cd $CUTOROOT/bin
./servant -c servant.ini
