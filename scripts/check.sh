#!/usr/bin/env bash
set -euo pipefail

go fmt ./...
go test ./...
go vet ./...

if {
  git ls-files
  git ls-files --others --exclude-standard
} | grep -E '(^|/)([^/]+\.(pem|key|p12|pfx)|\.env)$' >/dev/null; then
  echo "发现禁止提交的敏感文件，请移除后再提交。"
  exit 1
fi

echo "检查完成。"
