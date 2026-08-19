#!/usr/bin/env bash
set -uo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-$ROOT/docker-compose.yml}"
# 并发脚本互不抢同一 compose 项目；不要写死 --platform，跟宿主 arch 走。
PROJECT="gogo-verify-9-$$"

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

# 段2 走容器对外入口复现：HTTP / CLI / 写文件 / 重启 / 再读。
# 禁止 go test、禁止宿主 go build / go run、禁止读隐藏判据。

# compose --wait 只看容器进程，go build 完成并绑定端口还要等一会，轮询到健康入口可用。
ready=""
for _ in $(seq 1 60); do
  if curl -sf -m 3 http://127.0.0.1:18110/health >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 1
done

# 先确认服务真的起来了，健康入口必须 200。
if [ -z "$ready" ]; then
  echo "[EXPECT] /health returns 200 once the service is up"
  echo "[ACTUAL] /health unreachable, service did not start"
  exit 2
fi

# 投递之前 /v1/stats 必须能正常返回；否则环境本身有问题。
pre_stats=$(curl -s -m 5 -o /tmp/gogo9_pre.out -w '%{http_code}' http://127.0.0.1:18110/v1/stats 2>/dev/null) || pre_stats="000"
if [ "$pre_stats" != "200" ]; then
  echo "[EXPECT] GET /v1/stats returns 200 before any mail is queued"
  echo "[ACTUAL] GET /v1/stats before enqueue http=$pre_stats"
  exit 2
fi

# 入队一封邮件，这是触发后续统计读取崩溃的前提。
enqueue_code=$(curl -s -m 5 -o /tmp/gogo9_eq.out -w '%{http_code}' \
  -X POST http://127.0.0.1:18110/v1/enqueue \
  -H 'Content-Type: application/json' \
  -d '{"from":"a@local","to":"b@local","subject":"hi","body":"x"}' 2>/dev/null) || enqueue_code="000"
if [ "$enqueue_code" != "201" ]; then
  echo "[EXPECT] POST /v1/enqueue returns 201"
  echo "[ACTUAL] POST /v1/enqueue http=$enqueue_code"
  exit 2
fi

# 段3 对外部可观测量断言：投递之后再读 /v1/stats。
# 修复前：统计入口跑进 runtime panic，连接被中断、拿到空响应。
# 修复后/干净 main：仍返回 200。
stats_code=$(curl -s -m 5 -o /tmp/gogo9_stats.out -w '%{http_code}' http://127.0.0.1:18110/v1/stats 2>/dev/null) || stats_code="000"
if [ "$stats_code" = "200" ]; then
  echo "[EXPECT] GET /v1/stats returns 200 after a mail is queued"
  echo "[ACTUAL] GET /v1/stats after enqueue http=200"
  exit 0
fi

echo "[EXPECT] GET /v1/stats returns 200 after a mail is queued"
echo "[ACTUAL] GET /v1/stats after enqueue http=$stats_code (empty/reset reply — runtime panic)"
exit 1
