#!/bin/sh

cd `dirname $0`
echo $TEST05
if [ $TEST05 != "あ,い,う" ]; then
    echo "Invalid TEST05."
    exit 12
fi

exit 0
