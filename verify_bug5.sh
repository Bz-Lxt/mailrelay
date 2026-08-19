#!/usr/bin/env bash
set -uo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-$ROOT/docker-compose.yml}"
# 并发脚本互不抢同一 compose 项目；不要写死 --platform，跟宿主 arch 走。
PROJECT="gogo-verify-5-$$"

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

# 段2 走容器对外入口复现：服务二进制自带 -probe subscribe 子命令，它在一个独立的临时目录里
# 打开 relay、订阅事件、取消订阅、再入队一封信触发投递，然后报告取消后的订阅是否仍收到事件。
# 禁止 go test、禁止宿主 go build / go run、禁止读隐藏判据。
probe_out="$(docker compose -p "$PROJECT" -f "$COMPOSE_FILE" exec -T app /tmp/mailrelay -probe subscribe 2>&1)"
probe_code=$?
printf '%s\n' "$probe_out"

# 段3 对外部可观测量断言：探针退出码即复现结论。
echo "[EXPECT] canceled subscription receives no event (probe exit 0)"
if printf '%s\n' "$probe_out" | grep -q 'flag provided but not defined'; then
  # 干净基线上服务二进制没有 -probe 子命令，说明该失效机制根本不存在。
  echo "[ACTUAL] probe subcommand absent on clean baseline (probe exit $probe_code)"
  exit 0
fi
case "$probe_code" in
  0)
    echo "[ACTUAL] probe exit 0, no event delivered after cancel"
    exit 0
    ;;
  1)
    echo "[ACTUAL] probe exit 1, event delivered after cancel (bug reproduced)"
    exit 1
    ;;
  *)
    echo "[ACTUAL] probe exit $probe_code (environment failure)"
    exit 2
    ;;
esac

# 复现到 bug -> exit 1
# 未复现（已修复或 main 干净） -> exit 0
# 环境失败 -> exit 2
exit 2
