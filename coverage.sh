#!/bin/bash

echo "运行测试并检查覆盖率..."
if ! go test ./... -coverprofile=coverage.out -covermode=atomic -coverpkg=./...; then
  echo "::error::单元测试失败"
  exit 1
fi

# 获取覆盖率
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')

# 检查覆盖率
if [ -z "$COVERAGE" ]; then
  echo "::error::无法获取覆盖率数据"
  exit 1
fi

echo "覆盖率: $COVERAGE%"

if awk -v cov="$COVERAGE" 'BEGIN {exit !(cov >= 80)}'; then
  echo "✅ 覆盖率达标"
else
  echo "::error::覆盖率 $COVERAGE% 低于80%要求"
fi

gocov convert coverage.out | gocov-html -t kit > coverage.html
