<template>
  <a-card class="w-full" :title="$t('users')" :bordered="false">

    <a-grid :cols="{ xs: 2, sm: 2, md: 2, lg:2 }" :colGap="12" :rowGap="12">

      <a-grid-item>
        <a-statistic :title="$t('today')" :value="today" :value-from="0" show-group-separator animation>
          <template #suffix><span class="ml-1 text-gray-500"><icon-user /></span></template>
          <template #extra><a-spin v-if="loadingToday" /></template>
        </a-statistic>
      </a-grid-item>

      <a-grid-item>
        <a-statistic :title="$t('total')" :value="total" :value-from="0" show-group-separator animation>
          <template #suffix><span class="ml-1 text-gray-500"><icon-user-group /></span></template>
          <template #extra><a-spin v-if="loadingTotal" /></template>
        </a-statistic>
      </a-grid-item>

    </a-grid>

  </a-card>
</template>

<script setup>
import {useRequest} from "vue-request";
import {dashboardData} from "@/api/index.js";

const { data:total, loading:loadingTotal } = useRequest(dashboardData, {defaultParams:['userTotal']})
const { data:today, loading:loadingToday } = useRequest(dashboardData, {defaultParams:['userToday']})
</script>
