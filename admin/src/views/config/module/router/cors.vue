<template>
  <a-form-item label="CORS Origins">
    <a-textarea v-model="data.cors_origins" :auto-size="{ minRows: 2, maxRows: 5 }" placeholder="https://example.com&#10;https://app.example.com" />
    <a-tooltip>
      <icon-question-circle class="ml-2 opacity-60" />
      <template #content>
        Allowed origins for CORS (one per line).<br/>
        Leave empty to only allow localhost for development.<br/>
        Example: https://www.example.com
      </template>
    </a-tooltip>
  </a-form-item>

  <a-divider />

  <a-alert type="info">
    <template #title>CORS Configuration Guide</template>
    <div>
      <p><strong>What is CORS?</strong></p>
      <p>CORS (Cross-Origin Resource Sharing) allows your frontend to communicate with this API from a different domain.</p>
      <br/>
      <p><strong>How to configure:</strong></p>
      <ul>
        <li>Enter each allowed origin on a new line</li>
        <li>Include the full URL (e.g., https://www.example.com)</li>
        <li>Do not include trailing slashes</li>
        <li>Localhost origins (http://localhost:*) are always allowed for development</li>
      </ul>
      <br/>
      <p><strong>Example:</strong></p>
      <pre>https://www.l9.lc
https://app.l9.lc</pre>
    </div>
  </a-alert>

  <a-divider />

  <a-alert type="warning">{{$t('needRestartApp')}}</a-alert>
</template>

<script setup>
import { inject } from 'vue'

const data = inject('data')
</script>
