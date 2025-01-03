#!/bin/bash

APP_NAME="learn_gin_web"
DEPLOY_DIR=$(cd "$(dirname "$0")" && pwd)
APP_BINARY="$DEPLOY_DIR/$APP_NAME"
LOG_FILE="$DEPLOY_DIR/$APP_NAME.log"

# 标题
function title() {
  echo -e "\033[33m$1\033[0m"
}

# 错误
function error() {
  echo -e "\033[31mError: $1\033[0m" >&2
}

# 获取程序的PID
function get_app_pid() {
    local app="$1"
    pgrep -f "${app}"
}

# 检查应用是否有执行权限
function check_app_permission() {
    local app="$1"
    if [ ! -x "${app}" ]; then
        error "${app} - 缺少执行权限..."
        exit 1
    fi
}

# 启动新的程序
function start_app() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - 启动新服务..."
    nohup env GIN_MODE=release $APP_BINARY > $LOG_FILE 2>&1 &
    NEW_PID=$!
    echo "$(date '+%Y-%m-%d %H:%M:%S') - 服务启动成功: $NEW_PID"
}

# 使用grace http重启程序
function graceful_restart_app() {
    local pid="$1"
    kill -USR2 "${pid}"
    echo "$(date '+%Y-%m-%d %H:%M:%S') - 已发送重启信号到 ${pid}"
}

# 部署程序
deploy() {
    # 检查能不能执行程序
    check_app_permission $APP_BINARY
    # 获取当前正在运行的程序对应的PID
    # shellcheck disable=SC2155
    local pid=$(get_app_pid $APP_BINARY)
    if [ -n "${pid}" ]; then
        graceful_restart_app "${pid}"
    else
        start_app
    fi
}

deploy
