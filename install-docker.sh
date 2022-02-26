#!/usr/bin/env bash

if [ -z $server_name ]; then
  read -p "please enter server_name(default:server_center):" server_name
fi
if [ -z $server_name ]; then
  server_name="server_center"
fi

while :; do
  if [ ! -z $mysql_dsn ]; then
    break
  fi
  read -p "please enter mysql_dsn(required):" mysql_dsn
done

if [ -z $listen_port ]; then
  read -p "please enter listen port(default:7557):" listen_port
fi
if [ -z $listen_port ]; then
  listen_port="7557"
fi

echo
echo 'server_name: '$server_name
echo 'mysql_dsn: '$mysql_dsn
echo 'listen_port: '$listen_port
echo 'input any key go on, or control+c over'
read

echo 'create volume'
docker volume create log
echo 'create volume'
docker volume create server_center_data
echo 'stop container'
docker stop $server_name
echo 'remove container'
docker rm $server_name
echo 'remove image'
docker rmi $server_name
echo 'docker build'
docker build -t $server_name .
echo 'docker run'
docker run -d \
  --restart=always \
  --name $server_name \
  -v log:/log \
  -v server_center_data:/resource \
  -p $listen_port:7557 \
  -e server_name=$server_name \
  -e mysql_dsn=$mysql_dsn \
  $server_name

echo 'all finish'
