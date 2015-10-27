#!/bin/sh

if [ $RC -ne 0 ]; then
    echo "[err] Invalid RC : $RC"
fi
if [ $SD = "" ]; then
    echo "[err] Invalid SD : $SD"
fi
if [ $ED = "" ]; then
    echo "[err] Invalid ED : $ED"
fi
if [ $OUT != "" ]; then
    echo "[err] Invalid OUT : $OUT"
fi

exit 0
