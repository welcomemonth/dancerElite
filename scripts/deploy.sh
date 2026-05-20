#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP_NAME="${APP_NAME:-go-mountain}"
DEPLOY_DIR="${DEPLOY_DIR:-$ROOT_DIR/deploy}"
PORT="${PORT:-8080}"
DATABASE_DRIVER="${DATABASE_DRIVER:-sqlite3}"
DATABASE_DSN="${DATABASE_DSN:-$DEPLOY_DIR/data.db}"
JWT_SECRET="${JWT_SECRET:-}"

if [[ -z "$JWT_SECRET" ]]; then
  if command -v openssl >/dev/null 2>&1; then
    JWT_SECRET="$(openssl rand -hex 32)"
  else
    JWT_SECRET="$(date +%s%N)-go-mountain-default-secret-change-me"
  fi
fi

export SERVER_PORT="$PORT"
export DATABASE_DRIVER
export DATABASE_DSN
export JWT_SECRET

mkdir -p "$DEPLOY_DIR/bin" "$DEPLOY_DIR/configs" "$DEPLOY_DIR/uploads" "$DEPLOY_DIR/logs"

echo "构建前端并打包 Go 二进制..."
make -C "$ROOT_DIR" build
go build -o "$DEPLOY_DIR/bin/${APP_NAME}-migrate" "$ROOT_DIR/cmd/migrate/main.go"

PID_FILE="$DEPLOY_DIR/$APP_NAME.pid"
if [[ -f "$PID_FILE" ]]; then
  old_pid="$(cat "$PID_FILE")"
  if [[ -n "$old_pid" ]] && kill -0 "$old_pid" >/dev/null 2>&1; then
    echo "停止旧进程: $old_pid"
    kill "$old_pid"
    for _ in {1..20}; do
      if ! kill -0 "$old_pid" >/dev/null 2>&1; then
        break
      fi
      sleep 0.2
    done
    if kill -0 "$old_pid" >/dev/null 2>&1; then
      echo "旧进程未正常退出，强制停止: $old_pid"
      kill -9 "$old_pid"
    fi
  fi
fi

echo "同步部署文件到 $DEPLOY_DIR ..."
tmp_bin="$DEPLOY_DIR/bin/$APP_NAME.tmp"
cp "$ROOT_DIR/bin/go-mountain" "$tmp_bin"
mv "$tmp_bin" "$DEPLOY_DIR/bin/$APP_NAME"

if [[ ! -f "$DEPLOY_DIR/configs/config.yaml" ]]; then
  cat > "$DEPLOY_DIR/configs/config.yaml" <<EOF
server:
  port: $PORT

database:
  driver: $DATABASE_DRIVER
  dsn: $DATABASE_DSN

jwt:
  secret: $JWT_SECRET

wechat:
  app_id: ""
  secret: ""
  mch_id: ""
  mch_api_key: ""
  mch_serial_no: ""
  mch_private_key_path: ""
  notify_url: ""
EOF
  echo "已生成默认配置: $DEPLOY_DIR/configs/config.yaml"
else
  echo "保留已有配置: $DEPLOY_DIR/configs/config.yaml"
fi

echo "执行数据库迁移..."
(
  cd "$DEPLOY_DIR"
  "bin/${APP_NAME}-migrate"
)

echo "启动服务..."
(
  cd "$DEPLOY_DIR"
  nohup "bin/$APP_NAME" >"logs/$APP_NAME.log" 2>&1 &
  echo $! > "$PID_FILE"
)

new_pid="$(cat "$PID_FILE")"
for _ in {1..40}; do
  if curl -fsS "http://127.0.0.1:$PORT/api/ping" >/dev/null 2>&1; then
    echo "部署完成"
    echo "访问地址: http://127.0.0.1:$PORT/web/"
    echo "进程 PID: $new_pid"
    echo "日志文件: $DEPLOY_DIR/logs/$APP_NAME.log"
    exit 0
  fi
  sleep 0.25
done

echo "服务启动后健康检查失败，最近日志如下:"
tail -80 "$DEPLOY_DIR/logs/$APP_NAME.log" || true
exit 1
