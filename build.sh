#/bin/sh

rm $GOPATH/bin/master
rm $GOPATH/bin/servant
rm $GOPATH/bin/show

# *****************
# Unit test
# *****************
cd $GOPATH/src/cuto
. ./test_all_cover.sh
if [ "$?" -ne "0" ] ; then
    echo "unit test NG."
    exit 1
fi

# *****************
# All build
# *****************
echo "master building..."
go install cuto/master
if [ "$?" -ne "0" ] ; then
    echo "master build NG."
    exit 1
fi

echo "servant building..."
go install cuto/servant
if [ "$?" -ne "0" ] ; then
    echo "servant build NG."
    exit 1
fi

echo "show utility building..."
go install cuto/show
if [ "$?" -ne "0" ] ; then
    echo "show build NG."
    exit 1
fi

exit 0

