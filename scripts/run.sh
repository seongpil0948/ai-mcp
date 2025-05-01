#!/bin/bash
# Path: scripts/run.sh

# 환경 변수 로드
if [ -f .env ]; then
    echo "환경 변수 로드 중..."
    export $(grep -v '^#' .env | xargs)
fi

# 인자 처리
MODE="http"
PORT="8080"
DEBUG=false
CONFIG_PATH="configs/config.yaml"

# 명령줄 인자 처리
while [[ $# -gt 0 ]]; do
    case $1 in
        -m|--mode)
            MODE="$2"
            shift 2
            ;;
        -p|--port)
            PORT="$2"
            shift 2
            ;;
        -d|--debug)
            DEBUG=true
            shift
            ;;
        -c|--config)
            CONFIG_PATH="$2"
            shift 2
            ;;
        *)
            echo "알 수 없는 옵션: $1"
            exit 1
            ;;
    esac
done

# 디버그 모드 활성화 옵션
DEBUG_OPT=""
if [ "$DEBUG" = true ]; then
    DEBUG_OPT="--debug"
fi

# 서버 로그 디렉토리 생성
mkdir -p logs

echo "theshop-ai-mcp-server 실행 중..."
echo "모드: $MODE"
echo "포트: $PORT"
echo "설정 파일: $CONFIG_PATH"

# 서버 실행
if [ "$MODE" = "stdio" ]; then
    ./bin/mcpserver --mode stdio --config "$CONFIG_PATH" $DEBUG_OPT
else
    ./bin/mcpserver --mode http --port "$PORT" --config "$CONFIG_PATH" $DEBUG_OPT 2>&1 | tee -a logs/server.log
fi