#!/bin/bash

file_name='cblog-service'

if [ -f $file_name  ];then
    echo "$file_name exists, start to stop original progress after 2 seconds!!!"
    sleep 1
else
    echo "$file_name doesn't exist, stop!!!"
    exit
fi

echo "stop process"
ps -ef | grep $file_name | grep -v grep | awk '{print $2}' | xargs kill
sleep 2

cmd="(./$file_name > /dev/null 2>&1) &"
eval $cmd  # 执行命令
pid=$!

sleep 2
echo "check $pid if exists"
ps -p $pid | awk '{ print $1 }' | grep $pid
sys_pid=$?
if [ $sys_pid -eq 0 ]; then
    printf "$cmd\nstarted process success (new pid: $pid)"
else
    echo "start failed!!!"
fi

echo $pid > pid.txt
