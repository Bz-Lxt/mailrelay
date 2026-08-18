#!/usr/bin/env bash
set -uo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-$ROOT/docker-compose.yml}"
# 并发脚本互不抢同一 compose 项目；不要写死 --platform，跟宿主 arch 走。
PROJECT="gogo-verify-2-$$"

cleanup() {
  docker compose -p "$PROJECT" -f "$COMPOSE_FILE" down -v >/dev/null 2>&1 || true
}
trap cleanup EXIT

# 段1 用仓库自带 compose 起服务。构建/启动失败必须 exit 2，不能当成 bug 复现。
# 镜像必须是官方多架构（golang / alpine 等），禁止 --platform linux/amd64 或 linux/arm64。
if ! docker compose -p "$PROJECT" -f "$COMPOSE_FILE" up -d --wait --build; then
  echo "[EXPECT] docker compose up succeeded"
  echo "[ACTUAL] compose up failed"
  exit 2
fi

command -v curl >/dev/null 2>&1 || {
  echo "[EXPECT] curl available on host"
  echo "[ACTUAL] curl missing"
  exit 2
}

# 段2 走容器对外入口复现：等健康后，查 outbound 信箱摘要。
ready=0
for _ in $(seq 1 60); do
  if curl -sf http://127.0.0.1:18110/health >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 0.5
done
if [[ "$ready" -ne 1 ]]; then
  echo "[EXPECT] /health reachable"
  echo "[ACTUAL] /health not reachable"
  exit 2
fi

# 段3 对外部可观测量断言：GET /v1/mailbox?box=outbound 的状态码。
code=$(curl -s -o /dev/null -w '%{http_code}' "http://127.0.0.1:18110/v1/mailbox?box=outbound")

echo "[EXPECT] GET /v1/mailbox?box=outbound returns 200, not a crashed connection"
echo "[ACTUAL] http_code=$code"

# 复现到 bug -> exit 1
# 未复现（已修复或 main 干净） -> exit 0
# 环境失败 -> exit 2
if [[ "$code" == "000" ]]; then
  exit 1
fi
exit 0
