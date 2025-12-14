#!/bin/bash

# 完整测试流程：注册 -> 登录 -> 修改密码 -> 用新密码登录

BASE_URL="http://localhost:1004/usercenter/v1"
MOBILE="13900139999"
OLD_PASSWORD="123456"
NEW_PASSWORD="654321"

echo "========================================="
echo "步骤1: 注册新用户"
echo "========================================="
REGISTER_RESP=$(curl -s -X POST ${BASE_URL}/user/register \
  -H "Content-Type: application/json" \
  -d "{
    \"mobile\": \"${MOBILE}\",
    \"password\": \"${OLD_PASSWORD}\"
  }")

echo "注册响应: ${REGISTER_RESP}"
echo ""

# 等待1秒
sleep 1

echo "========================================="
echo "步骤2: 登录获取 Token"
echo "========================================="
LOGIN_RESP=$(curl -s -X POST ${BASE_URL}/user/login \
  -H "Content-Type: application/json" \
  -d "{
    \"mobile\": \"${MOBILE}\",
    \"password\": \"${OLD_PASSWORD}\"
  }")

echo "登录响应: ${LOGIN_RESP}"
echo ""

# 提取 token (使用 jq 或 grep)
if command -v jq &> /dev/null; then
    TOKEN=$(echo ${LOGIN_RESP} | jq -r '.data.accessToken')
else
    # 简单的文本提取（如果没有 jq）
    TOKEN=$(echo ${LOGIN_RESP} | grep -o '"accessToken":"[^"]*' | grep -o '[^"]*$')
fi

echo "提取的 Token: ${TOKEN}"
echo ""

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    echo "❌ 获取 Token 失败，请检查登录是否成功"
    exit 1
fi

# 等待1秒
sleep 1

echo "========================================="
echo "步骤3: 修改密码"
echo "========================================="
CHANGE_RESP=$(curl -s -X POST ${BASE_URL}/user/changePassword \
  -H "Content-Type: application/json" \
  -H "Authorization: ${TOKEN}" \
  -d "{
    \"oldPassword\": \"${OLD_PASSWORD}\",
    \"newPassword\": \"${NEW_PASSWORD}\",
    \"confirmPassword\": \"${NEW_PASSWORD}\"
  }")

echo "修改密码响应: ${CHANGE_RESP}"
echo ""

# 等待1秒
sleep 1

echo "========================================="
echo "步骤4: 用旧密码登录（应该失败）"
echo "========================================="
OLD_LOGIN_RESP=$(curl -s -X POST ${BASE_URL}/user/login \
  -H "Content-Type: application/json" \
  -d "{
    \"mobile\": \"${MOBILE}\",
    \"password\": \"${OLD_PASSWORD}\"
  }")

echo "旧密码登录响应: ${OLD_LOGIN_RESP}"
echo ""

# 等待1秒
sleep 1

echo "========================================="
echo "步骤5: 用新密码登录（应该成功）"
echo "========================================="
NEW_LOGIN_RESP=$(curl -s -X POST ${BASE_URL}/user/login \
  -H "Content-Type: application/json" \
  -d "{
    \"mobile\": \"${MOBILE}\",
    \"password\": \"${NEW_PASSWORD}\"
  }")

echo "新密码登录响应: ${NEW_LOGIN_RESP}"
echo ""

echo "========================================="
echo "测试完成！"
echo "========================================="
