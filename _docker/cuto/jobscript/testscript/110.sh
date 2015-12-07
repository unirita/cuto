#!/bin/sh

if [ $1 != "123" ]; then
    echo "Invalid Argument - $1"
    exit 12
fi

exit 0
