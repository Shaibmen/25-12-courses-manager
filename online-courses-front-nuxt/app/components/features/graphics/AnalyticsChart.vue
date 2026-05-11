<script setup lang="ts">
import { BarChart, HeatmapChart, LineChart, PieChart } from 'echarts/charts'
import {
  GridComponent,
  LegendComponent,
  TooltipComponent,
  VisualMapComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { use } from 'echarts/core'
import VChart from 'vue-echarts'

use([
  BarChart,
  HeatmapChart,
  LineChart,
  PieChart,
  GridComponent,
  LegendComponent,
  TooltipComponent,
  VisualMapComponent,
  CanvasRenderer
])

withDefaults(defineProps<{
  option: Record<string, unknown>
  height?: string
}>(), {
  height: '420px'
})
</script>

<template>
  <ClientOnly>
    <VChart
      class="analytics-chart"
      :style="{ height }"
      :option="option"
      autoresize
    />
  </ClientOnly>
</template>

<style scoped>
.analytics-chart {
  width: 100%;
}
</style>
