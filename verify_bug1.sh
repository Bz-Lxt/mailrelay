#!/usr/bin/env bash
set -uo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-$ROOT/docker-compose.yml}"
# 并发脚本互不抢同一 compose 项目；不要写死 --platform，跟宿主 arch 走。
PROJECT="gogo-verify-1-$$"

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

# 段2 走容器对外入口复现：等健康，并发投递，再读收件箱。
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

# 投递体不带 id 字段，交由服务自己生成，这样并发时才会撞同一处生成逻辑。
N=50
pids=()
for _ in $(seq 1 "$N"); do
  curl -s -o /dev/null -X POST http://127.0.0.1:18110/v1/enqueue \
    -H 'Content-Type: application/json' \
    -d '{"from":"a@local","to":"b@local","subject":"c","body":"d"}' &
  pids+=("$!")
done
for p in "${pids[@]}"; do
  wait "$p"
done

# 段3 对外部可观测量断言：并发投递 N 次后，收件箱里的信封数。
body=$(curl -s http://127.0.0.1:18110/v1/mailboxes?box=inbound) || body=""
count=$(printf '%s' "$body" | grep -o '"id"' | wc -l | tr -d ' ')

echo "[EXPECT] inbound envelope count == $N"
echo "[ACTUAL] inbound envelope count = $count"

# 复现到 bug -> exit 1
# 未复现（已修复或 main 干净） -> exit 0
# 环境失败 -> exit 2
if [[ "$count" -eq "$N" ]]; then
  exit 0
fi
exit 1
