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
        <a-space>
          <a-button type="primary" @click="addChannel">
            <template #icon><icon-plus /></template>
            添加频道
          </a-button>
          <a-button type="outline" @click="fetchUserChannels" :loading="loadingChannels" :disabled="!data.authenticated">
            <template #icon><icon-download /></template>
            从 Telegram 获取频道列表
          </a-button>
        </a-space>

        <a-alert v-if="!data.authenticated" type="warning" style="margin-top: 8px">
          请先完成 Telegram 认证后才能获取频道列表
        </a-alert>

        <a-modal
          v-model:visible="channelSelectVisible"
          title="选择 Telegram 频道"
          :width="700"
          @ok="confirmChannelSelection"
          unmount-on-close
          :modal-style="{ maxHeight: '80vh' }"
        >
          <div class="channel-select-content">
            <div class="channel-search">
              <a-input v-model="channelSearchKeyword" placeholder="搜索频道名称或用户名..." allow-clear>
                <template #prefix><icon-search /></template>
              </a-input>
              <span class="channel-count">
                共 {{ availableChannels.length }} 个频道，{{ filteredAvailableChannels.length }} 个匹配
              </span>
            </div>
            <a-table
              :columns="selectColumns"
              :data="filteredAvailableChannels"
              :pagination="paginationProps"
              :row-selection="rowSelectionProps"
              v-model:selectedKeys="selectedChannelKeys"
              row-key="id"
              :scroll="{ y: 300 }"
            >
              <template #title="{ record }">
                <a-space>
                  <span>{{ record.title }}</span>
                  <a-tag v-if="record.username" size="small">@{{ record.username }}</a-tag>
                </a-space>
              </template>
              <template #type="{ record }">
                <a-tag :color="record.is_broadcast ? 'blue' : 'green'">
                  {{ record.is_broadcast ? '频道' : '群组' }}
                </a-tag>
              </template>
            </a-table>
          </div>
        </a-modal>

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
          <template #filter="{ record }">
            <a-button size="small" type="text" @click="editFilter(record)">
              <template #icon><icon-filter /></template>
              {{ record.filter_keywords ? '已配置' : '配置' }}
            </a-button>
          </template>
          <template #status="{ record }">
            <a-switch v-model="record.status" :checked-value="1" :unchecked-value="0" />
          </template>
          <template #action="{ record, rowIndex }">
            <a-space>
              <a-button size="small" type="text" status="danger" @click="removeChannel(rowIndex)">
                <template #icon><icon-delete /></template>
              </a-button>
            </a-space>
          </template>
        </a-table>
      </a-space>
    </a-tab-pane>

    <a-tab-pane key="3" title="认证状态">
      <a-space direction="vertical" fill>
        <!-- 会话状态警告 -->
        <a-alert v-if="sessionStatus.need_reauth" type="warning" style="margin-bottom: 16px">
          <template #title>会话需要重新认证</template>
          {{ sessionStatus.message }}
        </a-alert>

        <a-descriptions :column="1" bordered>
          <a-descriptions-item label="手机号">
            {{ data.phone_number || '未设置' }}
          </a-descriptions-item>
          <a-descriptions-item label="连接状态">
            <a-tag :color="data.connected ? 'green' : 'red'">
              {{ data.connected ? '已连接' : '未连接' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="认证状态">
            <a-tag :color="sessionStatus.need_reauth ? 'red' : (data.authenticated ? 'green' : 'orange')">
              {{ sessionStatus.need_reauth ? '会话已过期' : (data.authenticated ? '已认证' : '未认证') }}
            </a-tag>
            <span v-if="sessionStatus.message && sessionStatus.status !== 'ok'" style="margin-left: 8px; color: var(--color-text-3);">
              {{ sessionStatus.message }}
            </span>
          </a-descriptions-item>
          <a-descriptions-item label="监听频道数">
            {{ data.monitored_channels || 0 }}
          </a-descriptions-item>
        </a-descriptions>

        <a-space v-if="sessionStatus.need_reauth || !data.authenticated">
          <a-button type="primary" @click="sendAuthCode" :loading="sendingCode">
            发送验证码
          </a-button>
          <a-button type="outline" @click="checkSession" :loading="checkingSession">
            检查会话状态
          </a-button>
        </a-space>

        <a-space v-else>
          <a-button type="primary" status="danger" @click="showClearConfirm = true">
            清除认证
          </a-button>
          <a-button type="outline" @click="checkSession" :loading="checkingSession">
            检查会话状态
          </a-button>
        </a-space>

        <a-form-item label="验证码" v-if="showCodeInput">
          <a-input v-model="authCode" placeholder="输入 Telegram 发送的验证码" />
          <a-button type="primary" @click="verifyAuthCode" :loading="verifyingCode" style="margin-top: 8px">
            验证登录
          </a-button>
        </a-form-item>

        <a-modal v-model:visible="showClearConfirm" title="确认清除认证" @ok="clearAuth" @cancel="showClearConfirm = false">
          <p>确定要清除 Telegram 认证吗？清除后需要重新验证手机号。</p>
          <p style="color: var(--color-danger);">此操作不可恢复！</p>
        </a-modal>
      </a-space>
    </a-tab-pane>

    <a-tab-pane key="4" title="同步日志">
      <a-space direction="vertical" fill>
        <a-button type="primary" @click="refreshLogs" :loading="loadingLogs">
          <template #icon><icon-refresh /></template>
          刷新日志
        </a-button>

        <a-table :columns="logColumns" :data="syncLogs" :pagination="{ pageSize: 20 }">
          <template #channel_id="{ record }">
            <span>{{ getChannelName(record.channel_id) }}</span>
          </template>
          <template #status="{ record }">
            <a-tag :color="record.status === 1 ? 'green' : (record.status === 2 ? 'orange' : 'red')">
              {{ record.status === 1 ? '成功' : (record.status === 2 ? '跳过' : '失败') }}
            </a-tag>
          </template>
          <template #create_time="{ record }">
            {{ formatTime(record.create_time) }}
          </template>
        </a-table>
      </a-space>
    </a-tab-pane>
  </a-tabs>

  <!-- 过滤规则编辑模态框 -->
  <a-modal v-model:visible="filterModalVisible" title="编辑过滤规则" :width="500" @ok="saveFilter">
    <a-form :model="filterForm" auto-label-width>
      <a-form-item label="消息类型">
        <a-checkbox-group v-model="filterForm.messageTypes">
          <a-checkbox value="text">文本</a-checkbox>
          <a-checkbox value="photo">图片</a-checkbox>
          <a-checkbox value="video">视频</a-checkbox>
        </a-checkbox-group>
      </a-form-item>
      <a-form-item label="最小长度">
        <a-input-number v-model="filterForm.minLength" :min="0" placeholder="0表示不限制" />
      </a-form-item>
      <a-form-item label="最大长度">
        <a-input-number v-model="filterForm.maxLength" :min="0" placeholder="0表示不限制" />
      </a-form-item>
      <a-form-item label="关键词规则">
        <a-radio-group v-model="filterForm.keywordType">
          <a-radio value="none">不过滤</a-radio>
          <a-radio value="whitelist">白名单</a-radio>
          <a-radio value="blacklist">黑名单</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item v-if="filterForm.keywordType !== 'none'" label="关键词列表">
        <a-textarea v-model="filterForm.keywordsText" placeholder="每行一个关键词" :auto-size="{ minRows: 3, maxRows: 6 }" />
        <template #extra>
          <span v-if="filterForm.keywordType === 'whitelist'">只同步包含这些关键词的消息</span>
          <span v-else-if="filterForm.keywordType === 'blacklist'">不同步包含这些关键词的消息</span>
        </template>
      </a-form-item>
      <a-form-item v-if="filterForm.keywordType !== 'none'" label="匹配模式">
        <a-switch v-model="filterForm.matchAll" />
        <template #extra>
          <span v-if="filterForm.keywordType === 'whitelist'">{{ filterForm.matchAll ? '必须匹配所有关键词' : '匹配任一关键词即可' }}</span>
          <span v-else>黑名单模式下匹配任一关键词即过滤</span>
        </template>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { inject, ref, computed, watch } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from "@/api/axios"

const data = inject('options')

// 初始化默认值
function initDefaults() {
  if (data.value == null) data.value = {}

  // 处理 channels_json -> channels 转换
  if (data.value.channels == null) {
    if (data.value.channels_json && data.value.channels_json !== '[]') {
      try {
        data.value.channels = JSON.parse(data.value.channels_json)
      } catch (e) {
        console.error('解析 channels_json 失败:', e)
        data.value.channels = []
      }
    } else {
      data.value.channels = []
    }
  }

  if (data.value.sync_logs == null) data.value.sync_logs = []
  if (data.value.app_id == null) data.value.app_id = 0
  if (data.value.app_hash == null) data.value.app_hash = ''
  if (data.value.phone_number == null) data.value.phone_number = ''
  if (data.value.session_key == null) data.value.session_key = ''
  if (data.value.auto_reconnect == null) data.value.auto_reconnect = true
  if (data.value.reconnect_delay == null) data.value.reconnect_delay = 5
  if (data.value.download_media == null) data.value.download_media = true
  if (data.value.max_image_size == null) data.value.max_image_size = 10485760
  if (data.value.keep_log_days == null) data.value.keep_log_days = 30
  if (data.value.connected == null) data.value.connected = false
  if (data.value.authenticated == null) data.value.authenticated = false
}

// 监听 data 变化，当数据加载时重新初始化
watch(() => data.value, (newVal) => {
  if (newVal) {
    initDefaults()
  }
}, { immediate: true, deep: true })

// 监听 channels 变化，同步更新 channels_json
watch(() => data.value?.channels, (newChannels) => {
  if (newChannels) {
    data.value.channels_json = JSON.stringify(newChannels)
  }
}, { deep: true })

// 立即初始化
initDefaults()

// 自动检查会话状态
checkSession()

const sendingCode = ref(false)
const verifyingCode = ref(false)
const showCodeInput = ref(false)
const authCode = ref('')
const phoneCodeHash = ref('')
const loadingChannels = ref(false)
const channelSelectVisible = ref(false)
const availableChannels = ref([])
const selectedChannelKeys = ref([])
const showClearConfirm = ref(false)
const channelSearchKeyword = ref('')
const checkingSession = ref(false)
const sessionStatus = ref({
  status: 'unknown',
  message: '',
  need_reauth: false
})

// 日志和状态监控
const loadingLogs = ref(false)
const syncLogs = ref([])

// 过滤规则相关
const filterModalVisible = ref(false)
const currentFilterChannel = ref(null)
const filterForm = ref({
  messageTypes: ['text', 'photo'],
  minLength: 0,
  maxLength: 0,
  keywordType: 'none',
  keywordsText: '',
  matchAll: false,
  caseSensitive: false
})

// 搜索过滤后的频道列表
const filteredAvailableChannels = computed(() => {
  if (!channelSearchKeyword.value || channelSearchKeyword.value.trim() === '') {
    return availableChannels.value
  }
  const keyword = channelSearchKeyword.value.toLowerCase().trim()
  return availableChannels.value.filter(ch => {
    const titleMatch = ch.title?.toLowerCase().includes(keyword)
    const usernameMatch = ch.username?.toLowerCase().includes(keyword)
    return titleMatch || usernameMatch
  })
})

const channelColumns = [
  { title: '频道ID', slotName: 'channel_id', width: 180 },
  { title: '频道名称', slotName: 'channel_name', width: 120 },
  { title: '目标分类', slotName: 'category_id', width: 100 },
  { title: '过滤规则', slotName: 'filter', width: 100 },
  { title: '状态', slotName: 'status', width: 70 },
  { title: '操作', slotName: 'action', width: 80 },
]

const selectColumns = [
  { title: '频道名称', slotName: 'title', width: 200 },
  { title: '频道ID', dataIndex: 'id', width: 120 },
  { title: '类型', slotName: 'type', width: 80 },
]

const paginationProps = {
  pageSize: 10,
  showTotal: true,
}

const rowSelectionProps = {
  type: 'checkbox',
  showCheckedAll: true,
  onlyCurrent: false,
}

const logColumns = [
  { title: '频道', slotName: 'channel_id', width: 150 },
  { title: '消息ID', dataIndex: 'message_id', width: 100 },
  { title: '文章ID', dataIndex: 'article_id', width: 80 },
  { title: '标题', dataIndex: 'message_title', ellipsis: true, tooltip: true },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '时间', slotName: 'create_time', width: 150 },
]

function addChannel() {
  data.value.channels.push({
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
  data.value.channels.splice(index, 1)
}

async function sendAuthCode() {
  const phone = data.value.phone_number
  console.log('发送验证码 - 手机号:', phone, '类型:', typeof phone)

  if (!phone || phone.trim() === '') {
    Message.warning('请先填写手机号')
    return
  }

  // 检查 API 配置
  if (!data.value.app_id || !data.value.app_hash) {
    Message.warning('请先配置 Telegram App ID 和 App Hash')
    return
  }

  sendingCode.value = true

  try {
    const response = await axios.post('/plugin/telegram/sendCode', {
      phone_number: phone.trim()
    })

    if (response.data.success) {
      phoneCodeHash.value = response.data.data?.phone_code_hash || ''
      showCodeInput.value = true
      Message.success('验证码已发送到 Telegram，请查收')
    } else {
      Message.error(response.data.message || '发送验证码失败')
    }
  } catch (error) {
    console.error('发送验证码错误:', error)
    Message.error(error.response?.data?.message || error.message || '发送验证码失败')
  } finally {
    sendingCode.value = false
  }
}

async function verifyAuthCode() {
  if (!authCode.value) {
    Message.warning('请输入验证码')
    return
  }

  verifyingCode.value = true

  try {
    const response = await axios.post('/plugin/telegram/verifyCode', {
      code: authCode.value.trim()
    })

    if (response.data.success) {
      data.value.authenticated = true
      data.value.connected = true
      Message.success('认证成功！')
      showCodeInput.value = false
      authCode.value = ''
    } else {
      Message.error(response.data.message || '验证失败')
    }
  } catch (error) {
    console.error('验证码验证错误:', error)
    Message.error(error.response?.data?.message || error.message || '验证失败')
  } finally {
    verifyingCode.value = false
  }
}

async function clearAuth() {
  try {
    const response = await axios.post('/plugin/telegram/clearAuth')

    if (response.data.success) {
      data.value.authenticated = false
      data.value.connected = false
      sessionStatus.value = {
        status: 'no_session',
        message: '认证已清除，请重新认证',
        need_reauth: true
      }
      showClearConfirm.value = false
      Message.success('认证已清除')
    } else {
      Message.error(response.data.message || '清除失败')
    }
  } catch (error) {
    console.error('清除认证错误:', error)
    Message.error(error.response?.data?.message || error.message || '清除失败')
  }
}

async function checkSession() {
  checkingSession.value = true

  try {
    const response = await axios.get('/plugin/telegram/checkSession')

    if (response.data.success) {
      sessionStatus.value = response.data.data || {
        status: 'unknown',
        message: '',
        need_reauth: false
      }

      if (sessionStatus.value.need_reauth) {
        Message.warning(sessionStatus.value.message || '会话需要重新认证')
      } else if (sessionStatus.value.status === 'ok') {
        Message.success('会话状态正常')
      }
    } else {
      Message.error(response.data.message || '检查失败')
    }
  } catch (error) {
    console.error('检查会话状态错误:', error)
    Message.error(error.response?.data?.message || error.message || '检查失败')
  } finally {
    checkingSession.value = false
  }
}

function formatTime(timestamp) {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

async function fetchUserChannels() {
  loadingChannels.value = true

  try {
    const response = await axios.get('/plugin/telegram/channels')

    if (response.data.success) {
      availableChannels.value = response.data.data || []
      selectedChannelKeys.value = []
      channelSelectVisible.value = true
      Message.success(`获取到 ${availableChannels.value.length} 个频道`)
    } else {
      Message.error(response.data.message || '获取频道列表失败')
    }
  } catch (error) {
    console.error('获取频道列表错误:', error)
    Message.error(error.response?.data?.message || error.message || '获取频道列表失败')
  } finally {
    loadingChannels.value = false
  }
}

function confirmChannelSelection() {
  const selectedChannels = availableChannels.value.filter(ch =>
    selectedChannelKeys.value.includes(ch.id)
  )

  for (const ch of selectedChannels) {
    // 检查是否已存在
    const exists = data.value.channels.some(c => c.channel_id === ch.id)
    if (!exists) {
      data.value.channels.push({
        channel_id: ch.id,
        channel_name: ch.title,
        category_id: 0,
        status: 1,
        filter_keywords: '',
        filter_message_types: 'text,photo',
        filter_min_length: 0,
        filter_max_length: 0,
      })
    }
  }

  channelSelectVisible.value = false
  Message.success(`已添加 ${selectedChannels.length} 个频道`)
}

// 编辑过滤规则
function editFilter(record) {
  currentFilterChannel.value = record

  // 解析现有的过滤配置
  const types = (record.filter_message_types || 'text,photo').split(',').filter(t => t)
  filterForm.value.messageTypes = types

  filterForm.value.minLength = record.filter_min_length || 0
  filterForm.value.maxLength = record.filter_max_length || 0

  // 解析关键词配置
  if (record.filter_keywords) {
    try {
      const kwConfig = JSON.parse(record.filter_keywords)
      filterForm.value.keywordType = kwConfig.type || 'none'
      filterForm.value.keywordsText = (kwConfig.keywords || []).join('\n')
      filterForm.value.matchAll = kwConfig.match_all || false
      filterForm.value.caseSensitive = kwConfig.case_sensitive || false
    } catch (e) {
      filterForm.value.keywordType = 'none'
      filterForm.value.keywordsText = ''
    }
  } else {
    filterForm.value.keywordType = 'none'
    filterForm.value.keywordsText = ''
  }

  filterModalVisible.value = true
}

// 保存过滤规则
function saveFilter() {
  if (!currentFilterChannel.value) return

  const record = currentFilterChannel.value

  // 保存消息类型
  record.filter_message_types = filterForm.value.messageTypes.join(',')

  // 保存长度限制
  record.filter_min_length = filterForm.value.minLength
  record.filter_max_length = filterForm.value.maxLength

  // 保存关键词配置
  if (filterForm.value.keywordType !== 'none' && filterForm.value.keywordsText.trim()) {
    const keywords = filterForm.value.keywordsText.split('\n').map(k => k.trim()).filter(k => k)
    const kwConfig = {
      type: filterForm.value.keywordType,
      keywords: keywords,
      match_all: filterForm.value.matchAll,
      case_sensitive: filterForm.value.caseSensitive
    }
    record.filter_keywords = JSON.stringify(kwConfig)
  } else {
    record.filter_keywords = ''
  }

  filterModalVisible.value = false
  Message.success('过滤规则已保存')
}

// 刷新日志
async function refreshLogs() {
  loadingLogs.value = true

  try {
    const response = await axios.get('/plugin/telegram/logs')

    if (response.data.success) {
      syncLogs.value = response.data.data || []
      Message.success(`获取到 ${syncLogs.value.length} 条日志`)
    } else {
      Message.error(response.data.message || '获取日志失败')
    }
  } catch (error) {
    console.error('获取日志错误:', error)
    Message.error(error.response?.data?.message || error.message || '获取日志失败')
  } finally {
    loadingLogs.value = false
  }
}

// 刷新状态
async function refreshStatus() {
  try {
    const response = await axios.get('/plugin/telegram/debug')

    if (response.data.success) {
      const status = response.data.data || {}
      // 更新 data 中的状态
      data.value.connected = status.connected
      data.value.authenticated = status.authenticated
      data.value.monitored_channels = status.monitored_channels
      data.value.reconnect_attempts = status.reconnect_attempts
      data.value.max_reconnect = status.max_reconnect
      Message.success('状态已刷新')
    } else {
      Message.error(response.data.message || '获取状态失败')
    }
  } catch (error) {
    console.error('获取状态错误:', error)
    Message.error(error.response?.data?.message || error.message || '获取状态失败')
  }
}

// 根据频道 ID 获取频道名称
function getChannelName(channelId) {
  if (!channelId) return '-'
  const channel = data.value.channels?.find(ch => ch.channel_id === channelId)
  return channel?.channel_name || `ID: ${channelId}`
}

// 初始化时加载日志
refreshLogs()
</script>

<style scoped>
.channel-select-content {
  display: flex;
  flex-direction: column;
}

.channel-search {
  margin-bottom: 12px;
}

.channel-count {
  display: block;
  margin-top: 8px;
  color: var(--color-text-3);
  font-size: 12px;
}
</style>
