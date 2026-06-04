#!/usr/bin/env bash
set -euo pipefail

go fmt ./...
go test ./...
go vet ./...

if find . \
  -path ./.git -prune -o \
  -type f \( -name "*.pem" -o -name "*.key" -o -name "*.p12" -o -name "*.pfx" -o -name ".env" \) \
  -print | grep -q .; then
  echo "发现禁止提交的敏感文件，请移除后再提交。"
  exit 1
fi

echo "检查完成。"

