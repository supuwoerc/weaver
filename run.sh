#!/bin/bash

APP_NAME="learn_gin_web"
DEPLOY_DIR="/var/www/learn-gin-web"
PID_FILE="$DEPLOY_DIR/$APP_NAME.pid"
APP_BINARY="$DEPLOY_DIR/$APP_NAME"
LOG_FILE="$DEPLOY_DIR/$APP_NAME.log"

start_app() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') - 启动新服务..."
  nohup env GIN_MODE=release $APP_BINARY > $LOG_FILE 2>&1 &
  NEW_PID=$!
  echo $NEW_PID > $PID_FILE
  echo "$(date '+%Y-%m-%d %H:%M:%S') - 服务启动成功: $NEW_PID"
}

graceful_restart() {
  if [ -f $PID_FILE ]; then
    OLD_PID=$(cat $PID_FILE)
    if [ -n "$OLD_PID" ] && kill -0 $OLD_PID 2>/dev/null; then
      echo "$(date '+%Y-%m-%d %H:%M:%S') - 优雅重启旧进程 $OLD_PID"
      kill -USR2 $OLD_PID
      echo "$(date '+%Y-%m-%d %H:%M:%S') - 已触发进程 $OLD_PID 优雅重启"

      # 等待新进程启动
      echo "$(date '+%Y-%m-%d %H:%M:%S') - 等待新进程启动..."
      TIMEOUT=30  # 最大等待时间（秒）
      ELAPSED=0
      while [ $ELAPSED -lt $TIMEOUT ]; do
        # 检查进程是否已启动
        NEW_PID=$(ps aux | grep "$APP_NAME" | grep -v grep | awk '{print $2}')
        if [ -n "$NEW_PID" ]; then
          echo "$(date '+%Y-%m-%d %H:%M:%S') - 新进程已启动: $NEW_PID"
          echo $NEW_PID > $PID_FILE
          echo "$(date '+%Y-%m-%d %H:%M:%S') - 新进程PID已更新到文件: $NEW_PID"
          return 0  # 成功启动，退出函数
        fi
        sleep 1  # 每秒检查一次
        ELAPSED=$((ELAPSED + 1))
      done
      echo "$(date '+%Y-%m-%d %H:%M:%S') - 超过最大等待时间，新进程未启动"
    else
      echo "$(date '+%Y-%m-%d %H:%M:%S') - 未找到旧进程或进程已停止，启动新服务"
      start_app
    fi
  else
    echo "$(date '+%Y-%m-%d %H:%M:%S') - 未发现PID文件，启动新服务"
    start_app
  fi
}

# 调用优雅重启
graceful_restart
