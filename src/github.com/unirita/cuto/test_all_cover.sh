#!/bin/sh

OSNAME=`uname`
TESTROOT=$GOPATH/src/github.com/unirita/cuto
LOGFILE=$TESTROOT/cover_all_$OSNAME.txt
RETCODE=0


cd $TESTROOT/console
echo "github.com/unirita/cuto/console package tested..."
go test -coverprofile cover.out> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi


cd $TESTROOT/db
echo "github.com/unirita/cuto/db package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi
cd $TESTROOT/db/query
echo "github.com/unirita/cuto/db/query package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi

cd $TESTROOT/db/tx
echo "github.com/unirita/cuto/db/tx package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi



cd $TESTROOT/log
echo "github.com/unirita/cuto/log package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi


cd $TESTROOT/master
echo "github.com/unirita/cuto/master package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi
cd $TESTROOT/master/config
echo "github.com/unirita/cuto/master/config package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi

cd $TESTROOT/master/jobnet
echo "github.com/unirita/cuto/master/jobnet package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi
cd $TESTROOT/master/jobnet/parser
echo "github.com/unirita/cuto/master/jobnet/parser package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi


cd $TESTROOT/master/remote
echo "github.com/unirita/cuto/master/remote package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi



cd $TESTROOT/message
echo "github.com/unirita/cuto/message package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi


cd $TESTROOT/servant
echo "github.com/unirita/cuto/servant package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi
cd $TESTROOT/servant/config
echo "github.com/unirita/cuto/servant/config package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi

cd $TESTROOT/servant/job
echo "github.com/unirita/cuto/servant/job package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi

cd $TESTROOT/servant/remote
echo "github.com/unirita/cuto/servant/remote package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi



cd $TESTROOT/show
echo "github.com/unirita/cuto/show package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi
cd $TESTROOT/show/gen
echo "github.com/unirita/cuto/show/gen package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi



cd $TESTROOT/util
echo "github.com/unirita/cuto/util package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi



cd $TESTROOT/flowgen/converter
echo "github.com/unirita/cuto/flowgen/converter package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi



cd $TESTROOT/realtime/network
echo "github.com/unirita/cuto/realtime/network package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi



cd $TESTROOT/utctime
echo "github.com/unirita/cuto/utctime package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi


exit $RETCODE
