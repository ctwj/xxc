<template>
  <a-space class="mb-2" size="medium">
    <a-input-search
      v-model="keyword"
      :placeholder="t('search')"
      style="width: 300px"
      @search="onSearch"
      @press-enter="onSearch"
    />
    <a-button @click="loadData">
      <template #icon><icon-refresh /></template>
      {{ t('refresh') }}
    </a-button>
  </a-space>

  <a-table ref="table" row-key="id"
           :style="{height:tableHeight + 'px'}"
           :columns="columns"
           :data="data"
           :bordered="false"
           :loading="loading"
           filter-icon-align-left
           :scroll="{x:'100%', y:tableHeight-78}"
           v-model:pagination="pagination"
           @pageChange="pageChange"
           @pageSizeChange="pageSizeChange">

    <template #th>
      <th style="color:var(--color-text-3)"></th>
    </template>
    <template #td="{ record }">
      <td style="color:var(--color-text-1)" class="border-opacity-30" />
    </template>

    <template #pagination-left><div><a-spin v-if="loadingCount" /></div></template>

    <template #status="{ record }">
      <a-tag :color="record.status === 1 ? 'green' : 'red'" style="cursor: pointer" @click="toggleStatus(record)">
        {{ record.status === 1 ? t('active') : t('disabled') }}
      </a-tag>
    </template>

    <template #emailVerified="{ record }">
      <a-tag :color="record.email_verified ? 'blue' : 'orange'">
        {{ record.email_verified ? t('verified') : t('unverified') }}
      </a-tag>
    </template>

    <template #time="{ record, column }">
      <n-time v-if="record[column.dataIndex] > 0" :time="record[column.dataIndex]" :to="Date.now()/1000" type="relative" unix />
      <span v-else> - </span>
    </template>

    <template #action="{ record }">
      <a-space>
        <a-popconfirm :content="t('confirmDelete')" @ok="onDelete(record)">
          <a-button size="mini" type="outline" status="danger">{{ t('delete') }}</a-button>
        </a-popconfirm>
      </a-space>
    </template>
  </a-table>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { Message } from '@arco-design/web-vue'
import { userList, userCount, userEnable, userDisable, userDelete } from '@/api/user'
import { t } from '@/locale'
import { useWindowSize } from '@vueuse/core'

const { height } = useWindowSize()
const tableHeight = computed(() => height.value - 180)

const loading = ref(false)
const loadingCount = ref(false)
const data = ref([])
const keyword = ref('')

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showTotal: true,
  showPageSize: true,
  pageSizeOptions: [20, 50, 100, 200]
})

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: 80,
    sortable: { sortDirections: ['ascend', 'descend'] }
  },
  {
    title: t('username'),
    dataIndex: 'username',
    width: 150,
    ellipsis: true,
    tooltip: true
  },
  {
    title: t('email'),
    dataIndex: 'email',
    width: 200,
    ellipsis: true,
    tooltip: true
  },
  {
    title: t('nickname'),
    dataIndex: 'nickname',
    width: 120,
    ellipsis: true,
    tooltip: true
  },
  {
    title: t('emailVerified'),
    slotName: 'emailVerified',
    width: 100,
    align: 'center'
  },
  {
    title: t('status'),
    slotName: 'status',
    width: 80,
    align: 'center'
  },
  {
    title: t('createTime'),
    dataIndex: 'create_time',
    slotName: 'time',
    width: 120,
    align: 'right'
  },
  {
    title: t('lastLogin'),
    dataIndex: 'last_login',
    slotName: 'time',
    width: 120,
    align: 'right'
  },
  {
    title: t('action'),
    slotName: 'action',
    width: 100,
    align: 'center'
  }
]

async function loadData() {
  loading.value = true
  loadingCount.value = true
  try {
    const params = {
      page: pagination.current,
      limit: pagination.pageSize,
      order: 'id desc',
      keyword: keyword.value
    }
    const [listRes, countRes] = await Promise.all([
      userList(params),
      userCount({ keyword: keyword.value })
    ])
    data.value = listRes || []
    pagination.total = countRes || 0
  } catch (e) {
    Message.error(e.message || 'Load failed')
  } finally {
    loading.value = false
    loadingCount.value = false
  }
}

function onSearch() {
  pagination.current = 1
  loadData()
}

function pageChange(page) {
  pagination.current = page
  loadData()
}

function pageSizeChange(size) {
  pagination.pageSize = size
  pagination.current = 1
  loadData()
}

async function toggleStatus(record) {
  try {
    if (record.status === 1) {
      await userDisable(record.id)
      Message.success(t('operationSuccess'))
    } else {
      await userEnable(record.id)
      Message.success(t('operationSuccess'))
    }
    loadData()
  } catch (e) {
    Message.error(e.message || 'Operation failed')
  }
}

async function onDelete(record) {
  try {
    await userDelete(record.id)
    Message.success(t('deleteSuccess'))
    loadData()
  } catch (e) {
    Message.error(e.message || 'Delete failed')
  }
}

onMounted(() => {
  loadData()
})
</script>