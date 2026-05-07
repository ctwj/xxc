# API Contracts: Telegram Channel Sync

**Date**: 2026-04-30
**Feature**: 002-telegram-channel-sync

## 概述

本文档定义 Telegram 频道同步插件的后端 API 接口契约。

**架构说明**: API 端点通过插件内部实现，注册到 Moss 的路由系统中。所有逻辑封装在 `main/plugins/` 目录，不创建独立的 controller 文件。

## 基础信息

- **Base URL**: `/admin/api/telegram-sync`
- **认证**: 需要 Admin JWT Token
- **响应格式**: JSON

## 通用响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 错误响应

```json
{
  "code": 1001,
  "message": "错误描述",
  "data": null
}
```

## API 端点

### 1. 频道管理

#### 1.1 获取频道列表

```
GET /admin/api/telegram-sync/channels
```

**Query Parameters**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| limit | int | 否 | 每页数量，默认 20 |
| status | int | 否 | 状态过滤：1=启用, 0=禁用 |

**Response**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 10,
    "page": 1,
    "limit": 20,
    "items": [
      {
        "id": 1,
        "channel_id": -1001234567890,
        "channel_name": "技术频道",
        "channel_link": "https://t.me/tech_channel",
        "status": 1,
        "category_id": 5,
        "last_sync_time": 1714473600,
        "total_sync_count": 150,
        "create_time": 1714387200
      }
    ]
  }
}
```

#### 1.2 添加频道

```
POST /admin/api/telegram-sync/channels
```

**Request Body**:

```json
{
  "channel_id": -1001234567890,
  "channel_name": "技术频道",
  "channel_link": "https://t.me/tech_channel",
  "status": 1,
  "category_id": 5,
  "article_status": true,
  "filter_keywords": "{\"type\":\"whitelist\",\"keywords\":[\"技术\"]}",
  "filter_message_types": "text,photo",
  "filter_min_length": 50,
  "remark": "技术文章同步"
}
```

**Response**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "channel_id": -1001234567890,
    "channel_name": "技术频道"
  }
}
```

#### 1.3 更新频道

```
PUT /admin/api/telegram-sync/channels/:id
```

**Request Body**: 同添加频道

#### 1.4 删除频道

```
DELETE /admin/api/telegram-sync/channels/:id
```

**Response**:

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

#### 1.5 切换频道状态

```
POST /admin/api/telegram-sync/channels/:id/toggle
```

**Response**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "status": 0
  }
}
```

### 2. 同步日志

#### 2.1 获取同步日志

```
GET /admin/api/telegram-sync/logs
```

**Query Parameters**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| limit | int | 否 | 每页数量，默认 50 |
| channel_id | int64 | 否 | 频道 ID 过滤 |
| status | int | 否 | 状态过滤：0=失败, 1=成功, 2=跳过 |

**Response**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 500,
    "page": 1,
    "limit": 50,
    "items": [
      {
        "id": 1,
        "channel_id": -1001234567890,
        "message_id": 12345,
        "article_id": 1001,
        "status": 1,
        "message_title": "Go 语言入门教程",
        "create_time": 1714473600
      }
    ]
  }
}
```

#### 2.2 获取日志详情

```
GET /admin/api/telegram-sync/logs/:id
```

### 3. 插件配置

#### 3.1 获取配置

```
GET /admin/api/telegram-sync/config
```

**Response**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "app_id": 12345,
    "app_hash": "********",
    "phone_number": "+8613800138000",
    "auto_reconnect": true,
    "reconnect_delay": 5,
    "sync_delay": 1,
    "download_media": true,
    "max_image_size": 10485760,
    "log_level": "info",
    "keep_log_days": 30
  }
}
```

#### 3.2 更新配置

```
PUT /admin/api/telegram-sync/config
```

**Request Body**:

```json
{
  "app_id": 12345,
  "app_hash": "your_app_hash",
  "phone_number": "+8613800138000",
  "auto_reconnect": true,
  "download_media": true
}
```

### 4. 认证管理

#### 4.1 发送验证码

```
POST /admin/api/telegram-sync/auth/send-code
```

**Request Body**:

```json
{
  "phone_number": "+8613800138000"
}
```

**Response**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "phone_code_hash": "abc123..."
  }
}
```

#### 4.2 验证登录

```
POST /admin/api/telegram-sync/auth/verify
```

**Request Body**:

```json
{
  "phone_code": "12345",
  "phone_code_hash": "abc123..."
}
```

**Response**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "authenticated": true,
    "user": {
      "id": 123456789,
      "first_name": "User",
      "phone": "+8613800138000"
    }
  }
}
```

#### 4.3 检查认证状态

```
GET /admin/api/telegram-sync/auth/status
```

**Response**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "authenticated": true,
    "user": {
      "id": 123456789,
      "first_name": "User",
      "phone": "+8613800138000"
    },
    "session_valid": true
  }
}
```

#### 4.4 登出

```
POST /admin/api/telegram-sync/auth/logout
```

### 5. 状态监控

#### 5.1 获取连接状态

```
GET /admin/api/telegram-sync/status
```

**Response**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "connected": true,
    "authenticated": true,
    "last_heartbeat": 1714473600,
    "monitored_channels": 5,
    "active_channels": 4,
    "uptime_seconds": 86400
  }
}
```

#### 5.2 获取统计数据

```
GET /admin/api/telegram-sync/stats
```

**Query Parameters**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| days | int | 否 | 统计天数，默认 7 |

**Response**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_synced": 1500,
    "total_success": 1425,
    "total_failed": 30,
    "total_skipped": 45,
    "success_rate": 95.0,
    "daily_stats": [
      {
        "date": "2026-04-30",
        "synced": 50,
        "success": 48,
        "failed": 1,
        "skipped": 1
      }
    ]
  }
}
```

## 错误码定义

| 错误码 | 说明 |
|--------|------|
| 1001 | 频道已存在 |
| 1002 | 频道不存在 |
| 1003 | 无效的频道 ID |
| 1004 | 未认证 Telegram |
| 1005 | 认证失败 |
| 1006 | 验证码错误 |
| 1007 | 验证码过期 |
| 1008 | 网络连接失败 |
| 1009 | Telegram API 错误 |
| 1010 | 配置无效 |
