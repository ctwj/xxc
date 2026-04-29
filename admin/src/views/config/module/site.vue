<template>
  <a-form-item :label="$t('name')">
    <a-input v-model="data.name" class="w-64" />
  </a-form-item>

  <a-form-item :label="$t('url')">
    <a-input  class="w-64" v-model="data.url" placeholder="http://www.xxx.com" />
  </a-form-item>

  <a-form-item :label="$t('title')" style="max-width: 600px">
    <a-input v-model="data.title" />
  </a-form-item>

  <a-form-item :label="$t('keywords')" field="keywords" style="max-width: 600px" extra="keywords">
    <a-textarea v-model="data.keywords" class="w-full" :auto-size="{ minRows: 3, maxRows: 3 }"/>
  </a-form-item>

  <a-form-item :label="$t('description')" field="description" style="max-width: 600px" extra="description">
    <a-textarea v-model="data.description" class="w-full" :auto-size="{ minRows: 3, maxRows: 3 }"/>
  </a-form-item>

  <a-divider />

  <a-typography-title :heading="6">Webhook 配置 (ISR Revalidation)</a-typography-title>
  <a-form-item label="Webhook URL" style="max-width: 600px">
    <a-input v-model="data.webhook_url" placeholder="https://your-frontend.com/api/revalidate" />
    <template #extra>
      前端 ISR revalidation API 地址，用于文章发布后自动更新前端页面
    </template>
  </a-form-item>

  <a-form-item label="Webhook Secret" style="max-width: 600px">
    <a-input v-model="data.webhook_secret" placeholder="your-secret-key" />
    <template #extra>
      用于验证 webhook 请求的密钥，需与前端 REVALIDATE_SECRET 环境变量一致
    </template>
  </a-form-item>

  <a-divider />
</template>

<script setup>
  import {inject,watch} from 'vue'

  const data = inject('data')

  watch(()=>data.value.url, (val)=>{
    if(!val) return
    if(!val.startsWith("http") && !val.startsWith("//")) data.value.url = "http://" + data.value.url
  })

</script>