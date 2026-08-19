#!/usr/bin/env bash
set -uo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-$ROOT/docker-compose.yml}"
# 并发脚本互不抢同一 compose 项目；不要写死 --platform，跟宿主 arch 走。
PROJECT="gogo-verify-10-$$"

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

BASE="http://127.0.0.1:18110"

# 段2 走容器对外入口复现：等健康，再投递一封 subject 超长（带开头标记）的信封。
ready=0
for _ in $(seq 1 90); do
  if curl -sf "$BASE/health" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 0.5
done
if [[ "$ready" -ne 1 ]]; then
  echo "[EXPECT] /health reachable after warm-up"
  echo "[ACTUAL] /health not reachable"
  exit 2
fi

# 对照：短 subject 入队必须成功（201），证明服务与端点正常。
ctrl=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$BASE/v1/enqueue" \
  -H 'Content-Type: application/json' \
  -d '{"from":"a@local","to":"b@local","subject":"ctrl","body":"y"}' || true)
if [[ "$ctrl" != "201" ]]; then
  echo "[EXPECT] enqueue with a short subject returns 201"
  echo "[ACTUAL] HTTP $ctrl"
  exit 2
fi

# 复现：subject 共 250 字节，开头是 HEAD 标记，后面接 246 个 a。
TAIL_A="$(head -c 246 /dev/zero | tr '\0' 'a')"
SUBJECT="HEAD${TAIL_A}"
RESP="$(curl -s -X POST "$BASE/v1/enqueue" \
  -H 'Content-Type: application/json' \
  -d "{\"from\":\"a@local\",\"to\":\"b@local\",\"subject\":\"${SUBJECT}\",\"body\":\"x\"}" || true)"

# 段3 对外部可观测量断言：
#  - subject 字段应保留开头 HEAD 标记（字段被截到 200 字节时保留开头）。
#  - id 字段应是完整长度的摘要（64 个十六进制字符），而不是只剩尾段。
subj_head=0
if printf '%s' "$RESP" | grep -q '"subject":"HEAD'; then
  subj_head=1
fi
id_full=0
if printf '%s' "$RESP" | grep -q '"id":"[0-9a-f]\{64\}"'; then
  id_full=1
fi

echo "[EXPECT] subject keeps HEAD and id is 64-char digest"
echo "[ACTUAL] subject_head=${subj_head} id_full=${id_full} resp=$(printf '%s' "$RESP" | tr -d '\n')"

# 复现到 bug -> exit 1
# 未复现（已修复或 main 干净） -> exit 0
# 环境失败 -> exit 2
if [[ "$subj_head" -eq 1 && "$id_full" -eq 1 ]]; then
  exit 0
fi
exit 1
