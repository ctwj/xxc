<template>
  <a-form-item :label="$t('enable')">
    <a-switch type="round" v-model="data.enable" />
  </a-form-item>

  <a-tabs type="line" default-active-key="api" :style="{ height: '600px' }">
    <a-tab-pane key="api" title="API 配置">
      <a-space direction="vertical" :size="16" fill>
        <a-form-item label="API URL">
          <a-input v-model="data.api_url" placeholder="https://api.openai.com/v1" />
          <template #extra>
            <a-text type="secondary">AI API 地址，兼容 OpenAI 格式</a-text>
          </template>
        </a-form-item>

        <a-form-item label="API Key">
          <a-input-password v-model="data.api_key" placeholder="sk-..." />
          <template #extra>
            <a-text type="secondary">API 密钥，用于调用 AI 服务</a-text>
          </template>
        </a-form-item>

        <a-form-item label="Model">
          <a-input v-model="data.model" placeholder="gpt-3.5-turbo" />
          <template #extra>
            <a-text type="secondary">使用的 AI 模型名称</a-text>
          </template>
        </a-form-item>
      </a-space>
    </a-tab-pane>

    <a-tab-pane key="keywords" title="关键词配置">
      <a-space direction="vertical" :size="16" fill>
        <a-form-item label="Min Keywords">
          <a-input-number v-model="data.min_keywords" :min="1" :max="20" :style="{ width: '200px' }" />
          <template #extra>
            <a-text type="secondary">最少提取的关键词数量</a-text>
          </template>
        </a-form-item>

        <a-form-item label="Max Keywords">
          <a-input-number v-model="data.max_keywords" :min="1" :max="20" :style="{ width: '200px' }" />
          <template #extra>
            <a-text type="secondary">最多提取的关键词数量</a-text>
          </template>
        </a-form-item>

        <a-form-item label="Min Long Tail">
          <a-input-number v-model="data.min_long_tail" :min="0" :max="10" :style="{ width: '200px' }" />
          <template #extra>
            <a-text type="secondary">最少长尾关键词数量（3-5 个词的组合）</a-text>
          </template>
        </a-form-item>
      </a-space>
    </a-tab-pane>

    <a-tab-pane key="tags" title="标签配置">
      <a-space direction="vertical" :size="16" fill>
        <a-form-item label="Min Tags">
          <a-input-number v-model="data.min_tags" :min="0" :max="10" :style="{ width: '200px' }" />
          <template #extra>
            <a-text type="secondary">最少生成的标签数量</a-text>
          </template>
        </a-form-item>

        <a-form-item label="Max Tags">
          <a-input-number v-model="data.max_tags" :min="0" :max="10" :style="{ width: '200px' }" />
          <template #extra>
            <a-text type="secondary">最多生成的标签数量</a-text>
          </template>
        </a-form-item>
      </a-space>
    </a-tab-pane>

    <a-tab-pane key="content" title="内容配置">
      <a-space direction="vertical" :size="16" fill>
        <a-form-item label="Enable Rewrite">
          <a-switch type="round" v-model="data.enable_rewrite" />
          <template #extra>
            <a-text type="secondary">启用文章内容改写功能，保留格式和图片，优化关键词表述</a-text>
          </template>
        </a-form-item>

        <a-form-item label="Auto Publish">
          <a-switch type="round" v-model="data.auto_publish" />
          <template #extra>
            <a-text type="secondary">SEO 处理完成后自动将文章状态修改为已发布</a-text>
          </template>
        </a-form-item>

        <a-form-item label="Other Category Name">
          <a-input v-model="data.other_category_name" placeholder="其他软件" />
          <template #extra>
            <a-text type="secondary">需要重新分类的分类名称</a-text>
          </template>
        </a-form-item>
      </a-space>
    </a-tab-pane>

    <a-tab-pane key="advanced" title="高级配置">
      <a-space direction="vertical" :size="16" fill>
        <a-form-item label="Batch Size">
          <a-input-number v-model="data.batch_size" :min="1" :max="100" :style="{ width: '200px' }" />
          <template #extra>
            <a-text type="secondary">每次处理的文章数量</a-text>
          </template>
        </a-form-item>

        <a-form-item label="Force Regenerate">
          <a-switch type="round" v-model="data.force_regenerate" />
          <template #extra>
            <a-text type="warning">强制重新生成已处理的文章（会覆盖已有的 SEO 数据）</a-text>
          </template>
        </a-form-item>

        <a-form-item label="Skip Published">
          <a-switch type="round" v-model="data.skip_published" />
          <template #extra>
            <a-text type="secondary">跳过已发布的文章，只处理未发布的文章</a-text>
          </template>
        </a-form-item>
      </a-space>
    </a-tab-pane>
  </a-tabs>
</template>

<script setup>
import { inject } from "vue";
const data = inject("options");
</script>