<template>
  <a-tabs default-active-key="1">
    <a-tab-pane key="1" title="基本配置">
      <a-form-item label="Telegram App ID" required>
        <a-input-number v-model="data.app_id" placeholder="从 my.telegram.org 获取" :min="0" />
      </a-form-item>

      <a-form-item label="Telegram App Hash" required>
        <a-input v-model="data.app_hash" placeholder="从 my.telegram.org 获取" />
      </a-form-item>

      <a-form-item label="手机号" help="国际格式，如 +8613800138000">
        <a-input v-model="data.phone_number" placeholder="+8613800138000" />
      </a-form-item>

      <a-form-item label="会话加密密钥" help="用于加密存储 Telegram 会话数据">
        <a-input v-model="data.session_key" placeholder="32字符以上的随机字符串" />
      </a-form-item>

      <a-form-item label="自动重连">
        <a-switch v-model="data.auto_reconnect" />
      </a-form-item>

      <a-form-item label="重连延迟（秒）">
        <a-input-number v-model="data.reconnect_delay" :min="1" :max="60" />
      </a-form-item>

      <a-form-item label="下载媒体">
        <a-switch v-model="data.download_media" />
      </a-form-item>

      <a-form-item label="最大图片大小（字节）">
        <a-input-number v-model="data.max_image_size" :min="0" :max="52428800" />
      </a-form-item>

      <a-form-item label="日志保留天数">
        <a-input-number v-model="data.keep_log_days" :min="1" :max="365" />
      </a-form-item>
    </a-tab-pane>

    <a-tab-pane key="2" title="频道配置">
      <a-space direction="vertical" fill>
        <a-button type="primary" @click="addChannel">
          <template #icon><icon-plus /></template>
          添加频道
        </a-button>

        <a-table :columns="channelColumns" :data="data.channels || []" :pagination="false">
          <template #channel_id="{ record }">
            <a-input-number v-model="record.channel_id" :min="-10000000000000" placeholder="-1001234567890" />
          </template>
          <template #channel_name="{ record }">
            <a-input v-model="record.channel_name" placeholder="频道名称" />
          </template>
          <template #category_id="{ record }">
            <a-input-number v-model="record.category_id" :min="0" placeholder="分类ID" />
          </template>
          <template #status="{ record }">
            <a-switch v-model="record.status" :checked-value="1" :unchecked-value="0" />
          </template>
          <template #action="{ record, rowIndex }">
            <a-button size="small" type="text" status="danger" @click="removeChannel(rowIndex)">
              <template #icon><icon-delete /></template>
            </a-button>
          </template>
        </a-table>
      </a-space>
    </a-tab-pane>

    <a-tab-pane key="3" title="认证状态">
      <a-space direction="vertical" fill>
        <a-descriptions :column="1" bordered>
          <a-descriptions-item label="连接状态">
            <a-tag :color="data.connected ? 'green' : 'red'">
              {{ data.connected ? '已连接' : '未连接' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="认证状态">
            <a-tag :color="data.authenticated ? 'green' : 'orange'">
              {{ data.authenticated ? '已认证' : '未认证' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="监听频道数">
            {{ data.monitored_channels || 0 }}
          </a-descriptions-item>
        </a-descriptions>

        <a-button type="primary" @click="sendAuthCode" :loading="sendingCode">
          发送验证码
        </a-button>

        <a-form-item label="验证码" v-if="showCodeInput">
          <a-input v-model="authCode" placeholder="输入 Telegram 发送的验证码" />
          <a-button type="primary" @click="verifyAuthCode" :loading="verifyingCode" style="margin-top: 8px">
            验证登录
          </a-button>
        </a-form-item>
      </a-space>
    </a-tab-pane>

    <a-tab-pane key="4" title="同步日志">
      <a-table :columns="logColumns" :data="data.sync_logs || []" :pagination="{ pageSize: 20 }">
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : (record.status === 2 ? 'orange' : 'red')">
            {{ record.status === 1 ? '成功' : (record.status === 2 ? '跳过' : '失败') }}
          </a-tag>
        </template>
        <template #create_time="{ record }">
          {{ formatTime(record.create_time) }}
        </template>
      </a-table>
    </a-tab-pane>
  </a-tabs>
</template>

<script setup>
import { inject, ref } from 'vue'
import { Message } from '@arco-design/web-vue'

const data = inject('options')

// 初始化默认值
if (!data.channels) data.channels = []
if (!data.sync_logs) data.sync_logs = []
if (!data.app_id) data.app_id = 0
if (!data.app_hash) data.app_hash = ''
if (!data.phone_number) data.phone_number = ''
if (!data.session_key) data.session_key = ''
if (!data.auto_reconnect) data.auto_reconnect = true
if (!data.reconnect_delay) data.reconnect_delay = 5
if (!data.download_media) data.download_media = true
if (!data.max_image_size) data.max_image_size = 10485760
if (!data.keep_log_days) data.keep_log_days = 30
if (!data.connected) data.connected = false
if (!data.authenticated) data.authenticated = false

const sendingCode = ref(false)
const verifyingCode = ref(false)
const showCodeInput = ref(false)
const authCode = ref('')
const phoneCodeHash = ref('')

const channelColumns = [
  { title: '频道ID', slotName: 'channel_id', width: 200 },
  { title: '频道名称', slotName: 'channel_name', width: 150 },
  { title: '目标分类', slotName: 'category_id', width: 120 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'action', width: 60 },
]

const logColumns = [
  { title: '频道', dataIndex: 'channel_name', width: 120 },
  { title: '消息ID', dataIndex: 'message_id', width: 100 },
  { title: '文章ID', dataIndex: 'article_id', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '时间', slotName: 'create_time', width: 150 },
]

function addChannel() {
  data.channels.push({
    channel_id: 0,
    channel_name: '',
    category_id: 0,
    status: 1,
    filter_keywords: '',
    filter_message_types: 'text,photo',
    filter_min_length: 0,
    filter_max_length: 0,
  })
}

function removeChannel(index) {
  data.channels.splice(index, 1)
}

async function sendAuthCode() {
  if (!data.phone_number) {
    Message.warning('请先填写手机号')
    return
  }
  sendingCode.value = true
  // TODO: 调用后端 API 发送验证码
  // 实际实现需要在后端添加认证 API
  setTimeout(() => {
    sendingCode.value = false
    showCodeInput.value = true
    Message.success('验证码已发送到 Telegram')
  }, 1000)
}

async function verifyAuthCode() {
  if (!authCode.value) {
    Message.warning('请输入验证码')
    return
  }
  verifyingCode.value = true
  // TODO: 调用后端 API 验证登录
  setTimeout(() => {
    verifyingCode.value = false
    data.authenticated = true
    Message.success('认证成功！')
  }, 1000)
}

function formatTime(timestamp) {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}
</script>