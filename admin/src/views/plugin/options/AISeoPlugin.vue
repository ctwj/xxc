<template>
  <a-form-item :label="$t('enable')">
    <a-switch type="round" v-model="data.enable" />
  </a-form-item>

  <a-tabs type="line" default-active-key="functions">
    <a-tab-pane key="functions" title="功能开关">
      <!-- 核心优化 - 紧凑网格布局 -->
      <div class="mb-3">
        <div class="text-xs text-gray-500 mb-2">核心优化（推荐全开）</div>
        <div class="grid grid-cols-2 gap-x-6 gap-y-1">
          <div class="flex items-center justify-between py-1">
            <span class="text-sm">标题优化</span>
            <a-switch type="round" v-model="data.enable_title_optimize" size="small" />
          </div>
          <div class="flex items-center justify-between py-1">
            <span class="text-sm">关键词提取</span>
            <a-switch type="round" v-model="data.enable_keywords" size="small" />
          </div>
          <div class="flex items-center justify-between py-1">
            <span class="text-sm">描述优化</span>
            <a-switch type="round" v-model="data.enable_description" size="small" />
          </div>
          <div class="flex items-center justify-between py-1">
            <span class="text-sm">标签生成</span>
            <a-switch type="round" v-model="data.enable_tags" size="small" />
          </div>
        </div>
      </div>

      <!-- 高级优化 -->
      <div class="mb-3">
        <div class="text-xs text-gray-500 mb-2">高级优化（按需开启）</div>
        <div class="grid grid-cols-2 gap-x-6 gap-y-1">
          <div class="flex items-center justify-between py-1">
            <span class="text-sm">内容改写 <a-tooltip content="深度优化内容，消耗较多 Token"><icon-question-circle class="text-gray-400 text-xs" /></a-tooltip></span>
            <a-switch type="round" v-model="data.enable_rewrite" size="small" />
          </div>
          <div class="flex items-center justify-between py-1">
            <span class="text-sm">分类推荐 <a-tooltip content="仅对未分类文章生效"><icon-question-circle class="text-gray-400 text-xs" /></a-tooltip></span>
            <a-switch type="round" v-model="data.enable_category_recommend" size="small" />
          </div>
        </div>
      </div>

      <!-- 其他设置 -->
      <div class="flex items-center gap-6 pt-2 border-t border-gray-100">
        <div class="flex items-center gap-2">
          <span class="text-sm">自动发布</span>
          <a-switch type="round" v-model="data.auto_publish" size="small" />
        </div>
      </div>
    </a-tab-pane>

    <a-tab-pane key="api" title="API 配置">
      <a-space direction="vertical" :size="12" fill>
        <a-form-item label="AI 类型">
          <a-select v-model="data.ai_type" @change="onAITypeChange" :style="{ width: '280px' }">
            <a-option value="openai">OpenAI 兼容</a-option>
            <a-option value="nvidia">NVIDIA NIM</a-option>
            <a-option value="zhipu">智普 GLM</a-option>
          </a-select>
        </a-form-item>

        <a-form-item label="API URL">
          <a-input v-model="data.api_url" placeholder="https://api.openai.com/v1" />
        </a-form-item>

        <a-form-item label="API Key">
          <a-input-password v-model="data.api_key" placeholder="sk-..." />
        </a-form-item>

        <a-form-item label="Model">
          <a-input v-model="data.model" placeholder="llama-3.1-8b-instruct" />
        </a-form-item>
      </a-space>
    </a-tab-pane>

    <a-tab-pane key="params" title="参数配置">
      <div class="grid grid-cols-3 gap-x-6 gap-y-3">
        <a-form-item label="最少关键词">
          <a-input-number v-model="data.min_keywords" :min="1" :max="20" />
        </a-form-item>

        <a-form-item label="最多关键词">
          <a-input-number v-model="data.max_keywords" :min="1" :max="20" />
        </a-form-item>

        <a-form-item label="最少长尾词">
          <a-input-number v-model="data.min_long_tail" :min="0" :max="10" />
        </a-form-item>

        <a-form-item label="最少标签">
          <a-input-number v-model="data.min_tags" :min="0" :max="10" />
        </a-form-item>

        <a-form-item label="最多标签">
          <a-input-number v-model="data.max_tags" :min="0" :max="10" />
        </a-form-item>

        <a-form-item label="其他分类">
          <a-input v-model="data.other_category_name" placeholder="其他软件" />
        </a-form-item>
      </div>
    </a-tab-pane>

    <a-tab-pane key="advanced" title="高级配置">
      <div class="grid grid-cols-2 gap-x-6 gap-y-3">
        <a-form-item label="批量数量">
          <a-input-number v-model="data.batch_size" :min="1" :max="100" />
        </a-form-item>

        <a-form-item label="指定文章ID">
          <a-input-number v-model="data.article_id" :min="0" placeholder="0=批量" />
        </a-form-item>
      </div>

      <div class="flex items-center gap-6 pt-2">
        <div class="flex items-center gap-2">
          <span class="text-sm">强制重新生成</span>
          <a-switch type="round" v-model="data.force_regenerate" size="small" />
        </div>
        <div class="flex items-center gap-2">
          <span class="text-sm">跳过已发布</span>
          <a-switch type="round" v-model="data.skip_published" size="small" />
        </div>
      </div>
    </a-tab-pane>
  </a-tabs>
</template>

<script setup>
import { inject } from "vue";
const data = inject("options");

// AI 类型对应的默认配置
const aiTypeConfig = {
  openai: {
    url: "https://api.openai.com/v1",
    model: "gpt-3.5-turbo"
  },
  nvidia: {
    url: "https://integrate.api.nvidia.com/v1",
    model: "abacusai/dracarys-llama-3.1-70b-instruct"
  },
  zhipu: {
    url: "https://open.bigmodel.cn/api/paas/v4",
    model: "glm-4.7-flash"
  },
};

// AI 类型变更时自动填充 URL 和 Model
function onAITypeChange(value) {
  if (aiTypeConfig[value]) {
    data.value.api_url = aiTypeConfig[value].url;
    data.value.model = aiTypeConfig[value].model;
  }
}

// 初始化：如果 ai_type 为空，根据当前 api_url 反推类型
if (!data.value.ai_type && data.value.api_url) {
  for (const [type, config] of Object.entries(aiTypeConfig)) {
    if (data.value.api_url.includes(config.url.replace("https://", "").replace("http://", ""))) {
      data.value.ai_type = type;
      break;
    }
  }
}
// 默认值
if (!data.value.ai_type) {
  data.value.ai_type = "openai";
}
</script>
