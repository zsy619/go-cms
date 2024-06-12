#!/bin/sh

current_path=$(pwd)
echo "The current path is: $current_path"

dir_path="$current_path/logs"

# 判断路径是否存在
if [ ! -d "$dir_path" ]; then
	echo "Directory does not exist. Creating now..."
	mkdir -p "$dir_path"
	echo "Directory created."
else
	echo "Directory already exists."
fi

case $1 in 
	start)
		nohup ./cms 2>&1 >> $dir_path/log.log 2> /dev/null &
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
		nohup ./cms 2>&1 >> $dir_path/log.log 2> /dev/null &
		echo "服务已重启..."
		sleep 1
	;;
	*) 
		echo "$0 {start|stop|restart}"
		exit 4
	;;
esac

