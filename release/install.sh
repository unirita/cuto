#!/bin/sh

VERSION="V0.9.4L0"

echo -e "\n"
echo "***********************************************************"
echo "*                                                         *"
echo "*                                                         *"
echo "*           GoCUTO $VERSION Instaler                      *"
echo "*                         Last Update 2015/08/06          *"
echo "*                                           UNIRITA Inc.  *"
echo "*                                                         *"
echo "***********************************************************"

OSNAME=`uname`
OSNMAER=`uname -r`
HOSTNAME=`hostname`
CURRENT_DIR=`pwd`
CURRENT_USER=`whoami`

echo -e "\ncheack operating system..."
if [ $OSNAME = "Linux" ] ; then
    echo " OS = $OSNAME $OSNMAER"
elif [ $OSNAME = "Darwin" ] ; then
    echo " OS = $OSNAME $OSNMAER"
else
    echo "<error> don't support os"
    echo "...abort"
    exit
fi
echo "Ok..."

echo -e "\ncheack already installed file..."
if [ -s .installsed ] ; then
     echo "<error> Already instaled."
     echo "File [.installsed ] existed"
     echo "...abort"
     exit
fi
echo -e "Ok ...\n"

INSTALL_DIR=$CURRENT_DIR
BIND_ADDRESS=
LISTEN_PORT=
SILENT_MODE=
if [ "$1" == "-s" ] ; then
    SILENT_MODE="ON"
    BIND_ADDRESS="0.0.0.0"
    LISTEN_PORT="2015"
fi

if [ "$SILENT_MODE" != "ON" ] ; then
    YES_NO=
    while [ -z "$YES_NO" ] ; do
        echo "Do you want to install GoCuto $VERSION [ y/n ] ?"
        read YES_NO
        if [ "$YES_NO" = "y" ] ; then
            echo -e "\nStarting to install..."
        elif [ "$YES_NO" = "n" ] ; then
            echo -e "\n...canceled"
            exit
        else
            YES_NO=
        fi
    done
    
    echo -e "\nPlease enter bind-address name of the CUTO Servant"
    echo " [ Defalt  = 0.0.0.0 ]"
    echo "When you use the defalut value , please push an enter key as it is. "
    read BIND_ADDRESS
    if [ -z "$BIND_ADDRESS" ] ; then
    #    BIND_ADDRESS=$HOSTNAME
        BIND_ADDRESS="0.0.0.0"
        echo "Use defalt [ $BIND_ADDRESS ]"
    else
        echo "Node name of the CUTO Servant = $BIND_ADDRESS "
    fi
    
    echo -e "\nPlease enter port number of the CUTO Servant "
    echo " [ Defalt Port Number = 2015 ]"
    echo "When you use the defalut value , please push an enter key as it is. "
    read LISTEN_PORT
    if [ -z "$LISTEN_PORT" ] ; then
        LISTEN_PORT="2015"
        echo "Use defalt [ $LISTEN_PORT ]"
    else
        echo -e "Port number of the CUTO Servant [ $LISTEN_PORT ]\n"
    fi
fi

echo ""
echo "Install GoCUTO with the following setup information."
echo "****************************************************"
echo "  Install User                    = $CURRENT_USER"
echo "  Install Directory of the GoCuto = $INSTALL_DIR"
echo "  CUTO Servant bind-address       = $BIND_ADDRESS"
echo "  CUTO Servant Port Number        = $LISTEN_PORT"
echo "****************************************************"
echo ""

if [ "$SILENT_MODE" != "ON" ] ; then
    YES_NO=
    while [ -z "$YES_NO" ] ; do
        echo "OK? [ y/n ] "
        read YES_NO
        if [ "$YES_NO" = "y" ] ; then
            echo "Installing..."
        elif [ "$YES_NO" = "n" ] ; then
            echo "...abort"
            exit
        else
            YES_NO=
        fi
    done
fi

USEC=s/@CUTOUSER/`echo $CURRENT_USER | sed 's/\//\\\\\//g'`/g
echo $USEC >> .installsed

DIRC=s/@ROOT/`echo $INSTALL_DIR | sed 's/\//\\\\\//g'`/g
echo $DIRC >> .installsed

ADDC=s/@BIND_ADDRESS/$BIND_ADDRESS/g
echo $ADDC >> .installsed

POTC=s/@LISTEN_PORT/$LISTEN_PORT/g
echo $POTC >> .installsed


cd $CURRENT_DIR/shell
CUTO_SHELL="cutoenv.sh servant.sh servant_service.sh"
for z in $CUTO_SHELL ; do
    echo "changing $z ..."
    if [ -s $z ] ; then
        sed -f ../.installsed $z > $z.temp
        mv $z.temp $z
    else
        echo "<error> $z didn't exists"
    fi
done
chmod 744 ./*
chmod u-x cutoenv.sh


cd $CURRENT_DIR/bpmn
CUTO_BPMN="inst_test.csv"
for z in $CUTO_BPMN ; do
    echo "changing $z ..."
    if [ -s $z ] ; then
        sed -f ../.installsed $z > $z.temp
        mv $z.temp $z
    else
        echo "<error> $z didn't exists"
    fi
done
chmod 644 $CUTO_BPMN


cd $CURRENT_DIR/bin
CUTO_PARMS="master.ini servant.ini"
for z in $CUTO_PARMS ; do
    echo "changing $z ..."
    if [ -s $z ] ; then
        sed -f ../.installsed $z > $z.temp
        mv $z.temp $z
    else
        echo "<error> $z didn't exists"
    fi
done
chmod 644 $CUTO_PARMS

CUTO_BINARY="master servant show"
chmod 755 $CUTO_BINARY

cd $CURRENT_DIR
mkdir joblog log temp

echo "...completed !"
echo "Thank you for installing GoCUTO !"
