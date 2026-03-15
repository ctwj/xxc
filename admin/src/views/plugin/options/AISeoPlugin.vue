<template>
  <a-form-item :label="$t('enable')">
    <a-switch type="round" v-model="data.enable" />
  </a-form-item>

  <a-form-item label="API URL">
    <a-input v-model="data.api_url" placeholder="https://api.openai.com/v1" />
  </a-form-item>

  <a-form-item label="API Key">
    <a-input-password v-model="data.api_key" placeholder="sk-..." />
  </a-form-item>

  <a-form-item label="Model">
    <a-input v-model="data.model" placeholder="gpt-3.5-turbo" />
  </a-form-item>

  <a-divider />

  <a-form-item label="Min Keywords">
    <a-input-number v-model="data.min_keywords" :min="1" :max="20" :style="{ width: '200px' }" />
  </a-form-item>

  <a-form-item label="Max Keywords">
    <a-input-number v-model="data.max_keywords" :min="1" :max="20" :style="{ width: '200px' }" />
  </a-form-item>

  <a-form-item label="Min Long Tail">
    <a-input-number v-model="data.min_long_tail" :min="0" :max="10" :style="{ width: '200px' }" />
  </a-form-item>

  <a-divider />

  <a-form-item label="Min Tags">
    <a-input-number v-model="data.min_tags" :min="0" :max="10" :style="{ width: '200px' }" />
  </a-form-item>

  <a-form-item label="Max Tags">
    <a-input-number v-model="data.max_tags" :min="0" :max="10" :style="{ width: '200px' }" />
  </a-form-item>

  <a-divider />

  <a-form-item label="Enable Rewrite">
    <a-switch type="round" v-model="data.enable_rewrite" />
  </a-form-item>

  <a-form-item label="Other Category Name">
    <a-input v-model="data.other_category_name" placeholder="其他软件" />
  </a-form-item>

  <a-divider />

  <a-form-item label="Batch Size">
    <a-input-number v-model="data.batch_size" :min="1" :max="100" :style="{ width: '200px' }" />
  </a-form-item>

  <a-form-item label="Force Regenerate">
    <a-switch type="round" v-model="data.force_regenerate" />
    <template #extra>
      <a-text type="warning">强制重新生成已处理的文章</a-text>
    </template>
  </a-form-item>
</template>

<script setup>
import { inject } from "vue";
const data = inject("options");
</script>