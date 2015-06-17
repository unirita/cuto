#!/bin/bash
#
# GoCuto        Starts syslogd/klogd.
#
# chkconfig: 345 99 01
# description: GoCuto servant for Linux is the facility run servant.
### BEGIN INIT INFO
# Provides: $CUTO
### END INIT INFO

# Source function library.
. /etc/init.d/functions


CUTOUSER=@CUTOUSER
CUTOROOT=@ROOT

start() {

#servant start
        ISALIVE=`ps -u $CUTOUSER | grep 'servant' | grep -v grep | wc -l`
        if [ $ISALIVE != 0 ] ; then
                echo "#### cuto servant already Started  ####"
                exit 1
        else
                echo "#### cuto servant Start .. ####"

                su - $CUTOUSER -c "$CUTOROOT/shell/servant.sh > $CUTOROOT/log/servant_service.log 2>&1 &"

                sleep 10

                nowtime=`date +"%Y/%m/%d %H:%M:%S"`
                ISALIVE=`ps -u $CUTOUSER | grep 'servant' | grep -v grep | wc -l`
                if [ $ISALIVE != 0 ] ; then
                    echo "$nowtime Servant process start successful." >> $CUTOROOT/log/servant_chklog.log
                else
                    echo "$nowtime Not found Servant process." >> $CUTOROOT/log/servant_chklog.log
                    exit 1
                fi
        fi
        exit 0
}

stop() {

#servant stop
        echo "servant Process Stopped ..."

        ps -ef | grep $CUTOUSER | grep servant | grep -v grep
        if [[ $? = 0 ]] ; then
                nowtime=`date +"%Y/%m/%d %H:%M:%S"`

                echo `ps -ef | grep $CUTOUSER | grep 'servant' | grep -v grep |awk '{print $2}'`
                export killid=`ps -ef | grep $CUTOUSER | grep 'servant' | grep -v grep |awk '{print $2}'`
                kill -15 $killid

                echo "$nowtime Servant process terminated." >> $CUTOROOT/log/servant_chklog.log
        else
                echo "### servant already terminated ###"
        fi

        exit 0

        return 0
}
rhstatus() {
        return 0
}
restart() {
        return 0
}
reload()  {
    return 1
}


case "$1" in
  start)
        start
        ;;
  stop)
        stop
        ;;
  *)
        echo $"Usage: $0 {start|stop}"
        exit 2
esac

exit $?
