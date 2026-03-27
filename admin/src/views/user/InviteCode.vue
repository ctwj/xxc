<template>
  <div class="p-4">
    <a-space class="mb-4" size="medium">
      <a-button type="primary" @click="showCreateModal = true">
        <template #icon><icon-plus /></template>
        {{ t('create') }}
      </a-button>
      <a-button @click="loadData">
        <template #icon><icon-refresh /></template>
        {{ t('refresh') }}
      </a-button>
    </a-space>

    <a-table
      :columns="columns"
      :data="data"
      :loading="loading"
      :pagination="pagination"
      @pageChange="onPageChange"
      @pageSizeChange="onPageSizeChange"
    >
      <template #code="{ record }">
        <a-typography-text copyable>{{ record.code }}</a-typography-text>
      </template>
      <template #usage="{ record }">
        <span>{{ record.used_count }} / {{ record.max_uses || t('unlimited') }}</span>
      </template>
      <template #expire="{ record }">
        <template v-if="record.expire_at > 0">
          <a-tag v-if="record.expire_at <= Date.now() / 1000" color="red">{{ t('expired') }}</a-tag>
          <span v-else>{{ formatTime(record.expire_at) }}</span>
        </template>
        <span v-else>{{ t('neverExpire') }}</span>
      </template>
      <template #create_time="{ record }">
        <span>{{ formatTime(record.create_time) }}</span>
      </template>
      <template #action="{ record }">
        <a-popconfirm
          :content="t('confirmDelete')"
          @ok="onDelete(record)"
        >
          <a-button size="mini" type="outline" status="danger">
            {{ t('delete') }}
          </a-button>
        </a-popconfirm>
      </template>
    </a-table>

    <!-- 创建邀请码弹窗 -->
    <a-modal
      v-model:visible="showCreateModal"
      :title="t('createInviteCode')"
      @ok="onCreate"
      @cancel="resetForm"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item :label="t('maxUses')">
          <a-input-number
            v-model="form.max_uses"
            :min="0"
            :placeholder="t('maxUsesPlaceholder')"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item :label="t('expireDays')">
          <a-input-number
            v-model="form.expire_days"
            :min="0"
            :placeholder="t('expireDaysPlaceholder')"
            style="width: 100%"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { inviteCodeList, inviteCodeCount, inviteCodeCreate, inviteCodeDelete } from '@/api/user'
import { t } from '@/locale'

const loading = ref(false)
const data = ref([])
const showCreateModal = ref(false)
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showTotal: true,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50, 100]
})

const form = reactive({
  max_uses: 0,
  expire_days: 0
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: t('code'), slotName: 'code', width: 180 },
  { title: t('usage'), slotName: 'usage', width: 150 },
  { title: t('expireTime'), slotName: 'expire', width: 180 },
  { title: t('createTime'), slotName: 'create_time', width: 180 },
  { title: t('action'), slotName: 'action', width: 100 }
]

function formatTime(timestamp) {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      limit: pagination.pageSize,
      order: 'id desc'
    }
    const [listRes, countRes] = await Promise.all([
      inviteCodeList(params),
      inviteCodeCount()
    ])
    data.value = listRes || []
    pagination.total = countRes || 0
  } catch (e) {
    Message.error(e.message || 'Load failed')
  } finally {
    loading.value = false
  }
}

function onPageChange(page) {
  pagination.current = page
  loadData()
}

function onPageSizeChange(size) {
  pagination.pageSize = size
  pagination.current = 1
  loadData()
}

async function onCreate() {
  try {
    await inviteCodeCreate(form)
    Message.success(t('createSuccess'))
    showCreateModal.value = false
    resetForm()
    loadData()
  } catch (e) {
    Message.error(e.message || 'Create failed')
  }
}

async function onDelete(record) {
  try {
    await inviteCodeDelete(record.id)
    Message.success(t('deleteSuccess'))
    loadData()
  } catch (e) {
    Message.error(e.message || 'Delete failed')
  }
}

function resetForm() {
  form.max_uses = 0
  form.expire_days = 0
}

onMounted(() => {
  loadData()
})
</script>