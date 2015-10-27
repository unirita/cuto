#!/bin/sh

RC=1
. ./cutoenv.sh

which sqlite3
if [ $? -ne 0 ] ; then
  echo Not found sqlite3
  exit $RC
fi

ISALIVE=`ps -u $LOGNAME | grep 'master' | grep -v grep | wc -l`
if [ $ISALIVE != 0 ] ; then
  echo Cuto master running now.
  exit $RC
fi

cd $CUTOROOT/data
if [ -s bk_jobnet.csv ] ; then
  rm bk_jobnet.csv
fi
if [ -s bk_job.csv ] ; then
  rm bk_job.csv
fi

sqlite3 cuto.sqlite < dbinit.query
if [ $? -eq 0 ] ; then
  echo Database initialize OK.
  RC=0
else
  echo Database initialize NG.
fi

exit $RC
