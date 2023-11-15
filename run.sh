#!/bin/sh

case $1 in 
	start)
		nohup ./cms 2>&1 >> /mooc/mooc/cms/logs/log.log 2> /dev/null &
		echo "服务已启动..."
		sleep 1
	;;
	stop)
		killall cms
		echo "服务已停止..."
		sleep 1
	;;
	restart)
		killall cms
		sleep 1
		nohup ./cms 2>&1 >> /mooc/mooc/cms/logs/log.log 2> /dev/null &
		echo "服务已重启..."
		sleep 1
	;;
	*) 
		echo "$0 {start|stop|restart}"
		exit 4
	;;
esac

