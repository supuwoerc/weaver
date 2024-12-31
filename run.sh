#!/bin/bash

# Define the application name and directories
APP_NAME="learn_gin_web"
DEPLOY_DIR="/var/www/learn-gin-web"
PID_FILE="$DEPLOY_DIR/$APP_NAME.pid"
APP_BINARY="$DEPLOY_DIR/$APP_NAME"
LOG_FILE="$DEPLOY_DIR/$APP_NAME.log"

# Function to start the application
start_app() {
  echo "启动新服务..."
  nohup env GIN_MODE=release $APP_BINARY > $LOG_FILE 2>&1 &
  NEW_PID=$!
  echo $NEW_PID > $PID_FILE
  echo "服务启动成功: $NEW_PID"
}

# Function to perform graceful restart
graceful_restart() {
  if [ -f $PID_FILE ]; then
    OLD_PID=$(cat $PID_FILE)
    if [ -n "$OLD_PID" ] && kill -0 $OLD_PID 2>/dev/null; then
      echo "优雅重启旧进程 $OLD_PID"
      kill -USR2 $OLD_PID
      echo "已触发进程 $OLD_PID 优雅重启"
      sleep 5
    else
      echo "未找到旧进程或进程已停止"
    fi
  else
    echo "未发现PID文件,启动新服务"
  fi
}

# Perform graceful restart if application is running
graceful_restart

# Start the new application
start_app
