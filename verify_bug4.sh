#!/usr/bin/env bash
set -uo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-$ROOT/docker-compose.yml}"
# 并发脚本互不抢同一 compose 项目；不要写死 --platform，跟宿主 arch 走。
PROJECT="gogo-verify-4-$$"
BODY_FILE="/tmp/gogo-verify-4-resp-$$"

cleanup() {
  docker compose -p "$PROJECT" -f "$COMPOSE_FILE" down -v >/dev/null 2>&1 || true
  rm -f "$BODY_FILE" 2>/dev/null || true
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

BASE="http://127.0.0.1:18110"

# 段2 走容器对外入口复现：等健康，再投递。
ready=0
for _ in $(seq 1 60); do
  if curl -sf "$BASE/health" >/dev/null 2>&1; then
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

# 对照：收件人正常的入队必须成功（201），证明服务与端点正常。
ok=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$BASE/v1/enqueue" \
  -H 'Content-Type: application/json' \
  -d '{"from":"a@local","to":"b@local","subject":"ctrl","body":"y"}' || true)
if [[ "$ok" != "201" ]]; then
  echo "[EXPECT] enqueue with a well-formed recipient returns 201"
  echo "[ACTUAL] HTTP $ok"
  exit 2
fi

# 复现：收件人不像一个地址（没有 @），本应被当成 4xx 拒掉，不该冒服务端错误。
bad=$(curl -s -o "$BODY_FILE" -w '%{http_code}' -X POST "$BASE/v1/enqueue" \
  -H 'Content-Type: application/json' \
  -d '{"to":"garbage","subject":"hi","body":"x"}' || true)

# 段3 对外部可观测量断言。
echo "[EXPECT] POST /v1/enqueue with a malformed recipient (no @) does not 5xx"
echo "[ACTUAL] HTTP $bad body=$(cat "$BODY_FILE" 2>/dev/null | tr -d '\n')"

# 复现到 bug -> exit 1
# 未复现（已修复或 main 干净） -> exit 0
# 环境失败 -> exit 2
case "$bad" in
  5*)
    exit 1
    ;;
  4*)
    exit 0
    ;;
  2*)
    exit 0
    ;;
  *)
    exit 2
    ;;
esac
