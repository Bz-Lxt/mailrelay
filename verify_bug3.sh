#!/usr/bin/env bash
set -uo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-$ROOT/docker-compose.yml}"
# 并发脚本互不抢同一 compose 项目；不要写死 --platform，跟宿主 arch 走。
PROJECT="gogo-verify-3-$$"

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

# 段2 走容器对外入口复现：等健康，连续投递两封邮件到收件箱。
BASE="http://127.0.0.1:18110"

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

# 投递体不带 id 字段，交由服务自己生成。
first=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$BASE/v1/enqueue" \
  -H 'Content-Type: application/json' \
  -d '{"from":"a@local","to":"b@local","subject":"first","body":"x"}' || true)
second=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$BASE/v1/enqueue" \
  -H 'Content-Type: application/json' \
  -d '{"from":"a@local","to":"b@local","subject":"second","body":"y"}' || true)

# 段3 对外部可观测量断言：收件箱里的信封数必须等于投递数。
inbox=$(curl -s "$BASE/v1/mailboxes?box=inbound" 2>/dev/null || true)
count=$(printf '%s' "$inbox" | grep -o '"id"' | wc -l | tr -d ' ')

echo "[EXPECT] inbound envelope count == 2"
echo "[ACTUAL] inbound envelope count = $count"

# 复现到 bug -> exit 1
# 未复现（已修复或 main 干净） -> exit 0
# 环境失败 -> exit 2
if [[ "$count" -eq 2 && "$first" -eq 201 && "$second" -eq 201 ]]; then
  exit 0
fi
exit 1
