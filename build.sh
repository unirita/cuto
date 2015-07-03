#!/bin/sh

rm $GOPATH/bin/master
rm $GOPATH/bin/servant
rm $GOPATH/bin/show
rm $GOPATH/bin/flowgen

# *****************
# Unit test
# *****************
cd $GOPATH/src/cuto
/bin/sh test_all_cover.sh
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

echo "show utility building..."
go install cuto/flowgen
if [ "$?" -ne "0" ] ; then
    echo "flowgen build NG."
    exit 1
fi

chmod a+x $GOPATH/bin/*

rm $GOPATH/cutoroot/bin/master
rm $GOPATH/cutoroot/bin/servant
rm $GOPATH/cutoroot/bin/show
rm $GOPATH/cutoroot/bin/flowgen

cp $GOPATH/bin/master $GOPATH/cutoroot/bin
cp $GOPATH/bin/servant $GOPATH/cutoroot/bin
cp $GOPATH/bin/show $GOPATH/cutoroot/bin
cp $GOPATH/bin/flowgen $GOPATH/cutoroot/bin

exit 0

