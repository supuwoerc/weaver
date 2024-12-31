#!/bin/bash

APP_NAME="learn_gin_web"
DEPLOY_DIR="/var/www/learn-gin-web"
PID_FILE="$DEPLOY_DIR/$APP_NAME.pid"
APP_BINARY="$DEPLOY_DIR/$APP_NAME"
LOG_FILE="$DEPLOY_DIR/$APP_NAME.log"

# 获取进程名称
get_process_name() {
  ps -p $1 -o comm= 2>/dev/null
}

start_app() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') - 启动新服务..."
  nohup env GIN_MODE=release $APP_BINARY > $LOG_FILE 2>&1 &
  NEW_PID=$!
  echo "$(date '+%Y-%m-%d %H:%M:%S') - 服务启动成功: $NEW_PID"
  echo $NEW_PID > $PID_FILE
}

graceful_restart() {
  if [ -f $PID_FILE ]; then
    OLD_PID=$(cat $PID_FILE)
    if [ -n "$OLD_PID" ]; then
      # 检查进程是否存在，并且检查进程名称
      PROCESS_NAME=$(get_process_name $OLD_PID)
      if [ "$PROCESS_NAME" == "$APP_NAME" ]; then
        echo "$(date '+%Y-%m-%d %H:%M:%S') - 优雅重启旧进程 $OLD_PID"
        kill -USR2 $OLD_PID
        echo "$(date '+%Y-%m-%d %H:%M:%S') - 已触发进程 $OLD_PID 优雅重启"
      else
        echo "$(date '+%Y-%m-%d %H:%M:%S') - PID $OLD_PID 不属于应用 $APP_NAME，启动新服务"
        start_app
      fi
    else
      echo "$(date '+%Y-%m-%d %H:%M:%S') - PID 文件无效，启动新服务"
      start_app
    fi
  else
    echo "$(date '+%Y-%m-%d %H:%M:%S') - 未发现PID文件，启动新服务"
    start_app
  fi
}

# 调用优雅重启
graceful_restart
