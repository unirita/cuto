#/bin/sh

CUTOROOT=@ROOT

cd $CUTOROOT/bin
./servant -c servant.ini > $CUTOROOT/log/servant_console.log 2>&1 &
