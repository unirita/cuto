#/bin/sh

OSNAME=`uname`
TESTROOT=$GOPATH/src/cuto
LOGFILE=$TESTROOT/cover_all_$OSNAME.txt
RETCODE=0


cd $TESTROOT/console
echo "console package tested..."
go test -coverprofile cover.out> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi


cd $TESTROOT/db
echo "db package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi
cd $TESTROOT/db/query
echo "db/query package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi

cd $TESTROOT/db/tx
echo "db/tx package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi



cd $TESTROOT/log
echo "log package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi


cd $TESTROOT/master
echo "master package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi
cd $TESTROOT/master/config
echo "master/config package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi

cd $TESTROOT/master/jobnet
echo "master/jobnet package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi
cd $TESTROOT/master/jobnet/parser
echo "master/jobnet/parser package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi


cd $TESTROOT/master/remote
echo "master/remote package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi



cd $TESTROOT/message
echo "message package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi


cd $TESTROOT/servant
echo "servant package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi
cd $TESTROOT/servant/config
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi

cd $TESTROOT/servant/job
echo "servant/job package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi

cd $TESTROOT/servant/remote
echo "servant/remote package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi



cd $TESTROOT/show
echo "show package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi
cd $TESTROOT/show/gen
echo "show/gen package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi



cd $TESTROOT/util
echo "util package tested..."
go test -coverprofile cover.out>> $LOGFILE
if [ "$?" -ne "0" ] ; then
  echo "NG."
  RETCODE=1
fi


exit $RETCODE
