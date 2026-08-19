#!/usr/bin/env bash
set -uo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-$ROOT/docker-compose.yml}"
# 并发脚本互不抢同一 compose 项目；不要写死 --platform，跟宿主 arch 走。
PROJECT="gogo-verify-7-$$"
BODY_FILE="/tmp/gogo-verify-7-stats-$$"

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

# 段2 走容器对外入口复现：等健康，再连投 N 封信。
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

# 对照：单封入队必须 201，证明端点与服务正常。
probe=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$BASE/v1/enqueue" \
  -H 'Content-Type: application/json' \
  -d '{"from":"a@local","to":"b@local","subject":"probe","body":"x"}' || true)
if [[ "$probe" != "201" ]]; then
  echo "[EXPECT] POST /v1/enqueue returns 201"
  echo "[ACTUAL] HTTP $probe"
  exit 2
fi

# 复现：再连投 4 封（含对照那封共 5 封）。投递计数应在 stats 里累加成 5。
N=5
for _ in $(seq 1 4); do
  curl -s -o /dev/null -X POST "$BASE/v1/enqueue" \
    -H 'Content-Type: application/json' \
    -d '{"from":"a@local","to":"b@local","subject":"batch","body":"y"}' || true
done

# 段3 对外部可观测量断言：/v1/stats 返回的 enqueued 计数。
body=$(curl -s "$BASE/v1/stats" || true)
printf '%s' "$body" >"$BODY_FILE"
# JSON 形如 {"dispatched":0,"enqueued":5}，取 enqueued 字段后的整数。
count=$(printf '%s' "$body" | sed -n 's/.*"enqueued":[[:space:]]*\([0-9][0-9]*\).*/\1/p' | head -n1)
[[ -z "$count" ]] && count=-1

echo "[EXPECT] /v1/stats enqueued == $N"
echo "[ACTUAL] enqueued = $count body=$(cat "$BODY_FILE" 2>/dev/null | tr -d '\n')"

# 复现到 bug -> exit 1
# 未复现（已修复或 main 干净） -> exit 0
# 环境失败 -> exit 2
if [[ "$count" -eq "$N" ]]; then
  exit 0
fi
exit 1
