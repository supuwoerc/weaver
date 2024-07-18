#!/bin/bash

# Define the application name and directories
APP_NAME="learn_gin_web"
DEPLOY_DIR="/var/www/learn-gin-web"
PID_FILE="$DEPLOY_DIR/$APP_NAME.pid"
APP_BINARY="$DEPLOY_DIR/$APP_NAME"

# Function to start the application
start_app() {
  echo "Starting the new application..."
  nohup GIN_MODE=release $APP_BINARY &> /dev/null &
  NEW_PID=$!
  echo $NEW_PID > $PID_FILE
  echo "New application started with PID $NEW_PID."
}

# Check if the application is already running
if [ -f $PID_FILE ]; then
  OLD_PID=$(cat $PID_FILE)
  if [ -n "$OLD_PID" ] && kill -0 $OLD_PID 2>/dev/null; then
    echo "Stopping the old process with PID $OLD_PID..."
    kill -9 $OLD_PID
    echo "Old process stopped."
  else
    echo "No old process found or process already stopped."
  fi
else
  echo "No PID file found."
fi

# Start the new application
start_app
