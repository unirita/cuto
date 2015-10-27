#!/bin/sh

cd `dirname $0`
/bin/sh sleep.sh

echo $1
echo $2
echo $3
if [ $1 != "abc" ]; then
    echo "Invalid Argument\$1 - $1"
    exit 101
fi
if [ $2 != "\"e f g\"" ]; then
    echo "Invalid Argument\$2 - $2"
    exit 102
fi
if [ $3 != "h" ]; then
    echo "Invalid Argument\$3 - $3"
    exit 103
fi

exit 0
