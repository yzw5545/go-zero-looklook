#!/bin/bash

# 启动 usercenter 服务（RPC + API）

echo "==> 检查端口占用..."
lsof -i :2004 && echo "端口 2004 已被占用，请先关闭" && exit 1
lsof -i :1004 && echo "端口 1004 已被占用，请先关闭" && exit 1

echo "==> 启动 usercenter-rpc (端口 2004)..."
go run app/usercenter/cmd/rpc/usercenter.go -f app/usercenter/cmd/rpc/etc/usercenter.yaml &
RPC_PID=$!

sleep 2

echo "==> 启动 usercenter-api (端口 1004)..."
go run app/usercenter/cmd/api/usercenter.go -f app/usercenter/cmd/api/etc/usercenter.yaml &
API_PID=$!

echo ""
echo "✅ 服务已启动："
echo "   - usercenter-rpc: 端口 2004 (PID: $RPC_PID)"
echo "   - usercenter-api: 端口 1004 (PID: $API_PID)"
echo ""
echo "测试命令："
echo "   curl -X POST http://localhost:1004/usercenter/v1/user/login -H 'Content-Type: application/json' -d '{\"mobile\":\"13900139002\",\"password\":\"123456\"}'"
echo ""
echo "停止服务："
echo "   kill $RPC_PID $API_PID"

wait
