#!/bin/bash

echo "运行测试并检查覆盖率..."
go test ./... -coverprofile=coverage.out -covermode=atomic -coverpkg=./...

COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')

echo "当前覆盖率: $COVERAGE%"

if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "错误: 覆盖率 $COVERAGE% 低于80%要求"
    #go tool cover -html=coverage.out
    gocov convert coverage.out | gocov-html -t kit > coverage.html
    exit 1
fi

echo "覆盖率达标"
#go tool cover -html=coverage.out
#gocov convert coverage.out | gocov-html > coverage.html
gocov convert coverage.out | gocov-html -t kit > coverage.html
