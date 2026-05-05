<script setup lang="ts">
import type {
  DivisionListenersMetric,
  EnrollmentSourceMetric,
  GraphicsDashboardResponse,
  GroupEnrollmentMetric,
  ProgramAgeMetric,
  ProgramListenersMetric,
  ProgramPopularMetric,
  ProgramRevenueMetric
} from '../../../types/graphics'
import AppCard from '../../ui/AppCard.vue'
import AppStatCard from '../../ui/AppStatCard.vue'
import AnalyticsChart from './AnalyticsChart.vue'

type AnalyticsView = 'overview' | 'programs' | 'audience' | 'operations'

const props = defineProps<{
  data: GraphicsDashboardResponse
}>()

const activeView = ref<AnalyticsView>('overview')

const viewButtons: Array<{
  id: AnalyticsView
  label: string
  description: string
}> = [
  { id: 'overview', label: 'Сводка', description: 'Главные выводы и быстрый обзор' },
  { id: 'programs', label: 'Программы', description: 'Спрос, популярность и выручка' },
  { id: 'audience', label: 'Аудитория', description: 'Возрасты и каналы зачисления' },
  { id: 'operations', label: 'Нагрузка', description: 'Группы и подразделения' }
]

const numberFormatter = new Intl.NumberFormat('ru-RU')
const compactNumber = new Intl.NumberFormat('ru-RU', {
  notation: 'compact',
  maximumFractionDigits: 1
})
const moneyFormatter = new Intl.NumberFormat('ru-RU', {
  style: 'currency',
  currency: 'RUB',
  maximumFractionDigits: 0
})

const formatNumber = (value: number) => numberFormatter.format(value)
const formatMoney = (value: number) => moneyFormatter.format(value)
const formatCompact = (value: number) => compactNumber.format(value)

const wrapLabel = (value: string, maxLineLength = 18, maxLines = 3) => {
  const source = String(value || '').trim()

  if (!source) {
    return ''
  }

  const words = source.split(/\s+/)
  const lines: string[] = []
  let currentLine = ''

  for (const word of words) {
    const nextLine = currentLine ? `${currentLine} ${word}` : word

    if (nextLine.length <= maxLineLength) {
      currentLine = nextLine
      continue
    }

    if (currentLine) {
      lines.push(currentLine)
      currentLine = ''
    }

    if (word.length <= maxLineLength) {
      currentLine = word
      continue
    }

    const chunks = word.match(new RegExp(`.{1,${maxLineLength}}`, 'g')) || [word]

    for (const chunk of chunks) {
      if (lines.length < maxLines - 1) {
        lines.push(chunk)
      } else {
        lines.push(`${chunk.slice(0, Math.max(maxLineLength - 1, 1))}…`)
        return lines.join('\n')
      }
    }
  }

  if (currentLine) {
    lines.push(currentLine)
  }

  if (lines.length <= maxLines) {
    return lines.join('\n')
  }

  return `${lines.slice(0, maxLines).join('\n')}…`
}

const sumBy = <T>(items: T[], getter: (item: T) => number) =>
  items.reduce((total, item) => total + getter(item), 0)

const sortDesc = <T>(items: T[], getter: (item: T) => number) =>
  [...items].sort((left, right) => getter(right) - getter(left))

const totalListeners = computed(() => sumBy(props.data.count, item => item.listeners))
const totalAccurateListeners = computed(() => sumBy(props.data.countAccurate, item => item.listeners))
const totalExpectedRevenue = computed(() =>
  sumBy(props.data.worthAccurate, item => item.total_expected_revenue)
)
const totalGroupLoad = computed(() => sumBy(props.data.groupLoad, item => item.active_enrolled))
const totalDivisionLoad = computed(() => sumBy(props.data.divisionLoad, item => item.listeners))

const listenersSorted = computed(() =>
  sortDesc(props.data.countAccurate, item => item.listeners)
)

const popularSorted = computed(() =>
  sortDesc(props.data.popular, item => item.listeners)
)

const revenueSorted = computed(() =>
  sortDesc(props.data.worthAccurate, item => item.total_expected_revenue)
)

const groupsSorted = computed(() =>
  sortDesc(props.data.groupLoad, item => item.active_enrolled)
)

const divisionsSorted = computed(() =>
  sortDesc(props.data.divisionLoad, item => item.listeners)
)

const bestProgramByListeners = computed(() => listenersSorted.value[0] || null)
const bestPopularProgram = computed(() => popularSorted.value[0] || null)
const bestRevenueProgram = computed(() => revenueSorted.value[0] || null)
const weakestRevenueProgram = computed(() => revenueSorted.value.at(-1) || null)
const busiestGroup = computed(() => groupsSorted.value[0] || null)

const countDeltaMap = computed(() => {
  const baseline = new Map(props.data.count.map(item => [item.name_prof_education, item.listeners]))

  return new Map(
    props.data.countAccurate.map(item => [
      item.name_prof_education,
      item.listeners - (baseline.get(item.name_prof_education) || 0)
    ])
  )
})

const educationTypeBreakdown = computed(() => {
  const totals = new Map<string, number>()

  props.data.popular.forEach(item => {
    totals.set(item.educationtype, (totals.get(item.educationtype) || 0) + item.listeners)
  })

  return Array.from(totals.entries())
    .map(([name, value]) => ({ name, value }))
    .sort((left, right) => right.value - left.value)
})

const ageRanges = computed(() => {
  const unique = Array.from(new Set(props.data.ageDiff.map(item => item.age_range)))
  return unique.sort((left, right) => {
    const leftStart = Number(left.split('-')[0]) || 0
    const rightStart = Number(right.split('-')[0]) || 0
    return leftStart - rightStart
  })
})

const programNamesForAge = computed(() =>
  Array.from(new Set(props.data.ageDiff.map(item => item.name_prof_education)))
)

const ageHeatmapOption = computed(() => {
  const maxValue = Math.max(...props.data.ageDiff.map(item => item.listeners), 0)

  return {
    tooltip: {
      trigger: 'item',
      formatter: (params: { data: [number, number, number] }) => {
        const [x, y, value] = params.data
        return `${programNamesForAge.value[y]}<br/>${ageRanges.value[x]}: ${formatNumber(value)}`
      }
    },
    grid: {
      top: 30,
      left: 190,
      right: 30,
      bottom: 60
    },
    xAxis: {
      type: 'category',
      data: ageRanges.value,
      axisLabel: {
        color: '#475569'
      }
    },
    yAxis: {
      type: 'category',
      data: programNamesForAge.value,
      axisLabel: {
        color: '#475569',
        width: 170,
        overflow: 'truncate'
      }
    },
    visualMap: {
      min: 0,
      max: maxValue || 1,
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      bottom: 0,
      inRange: {
        color: ['#e0f2fe', '#7dd3fc', '#0f766e']
      }
    },
    series: [
      {
        type: 'heatmap',
        data: props.data.ageDiff.map(item => [
          ageRanges.value.indexOf(item.age_range),
          programNamesForAge.value.indexOf(item.name_prof_education),
          item.listeners
        ]),
        label: {
          show: true,
          color: '#0f172a'
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowColor: 'rgba(15, 23, 42, 0.22)'
          }
        }
      }
    ]
  }
})

const listenersComparisonOption = computed(() => {
  const ordered = sortDesc(props.data.countAccurate, item => item.listeners)
  const labels = ordered.map(item => item.name_prof_education)
  const baseMap = new Map(props.data.count.map(item => [item.name_prof_education, item.listeners]))

  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      bottom: 0
    },
    grid: {
      top: 20,
      left: 250,
      right: 60,
      bottom: 50,
      containLabel: false
    },
    xAxis: {
      type: 'value',
      axisLabel: {
        color: '#64748b'
      }
    },
    yAxis: {
      type: 'category',
      data: labels,
      axisLabel: {
        color: '#334155',
        width: 220,
        lineHeight: 18,
        formatter: (value: string) => wrapLabel(value, 24, 3)
      }
    },
    series: [
      {
        name: 'Базовый расчет',
        type: 'bar',
        data: labels.map(label => baseMap.get(label) || 0),
        itemStyle: {
          color: '#94a3b8',
          borderRadius: [0, 10, 10, 0]
        }
      },
      {
        name: 'Точный расчет',
        type: 'bar',
        data: labels.map(label => ordered.find(item => item.name_prof_education === label)?.listeners || 0),
        itemStyle: {
          color: '#0f766e',
          borderRadius: [0, 10, 10, 0]
        }
      }
    ]
  }
})

const popularProgramsOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'shadow'
    },
    formatter: (params: Array<{ axisValue: string, seriesName: string, value: number }>) => {
      const lines = params.map(item => `${item.seriesName}: ${formatNumber(item.value)}`)
      return `${params[0]?.axisValue || ''}<br/>${lines.join('<br/>')}`
    }
  },
  grid: {
    top: 30,
    left: 250,
    right: 150,
    bottom: 24,
    containLabel: false
  },
  xAxis: {
    type: 'value',
    axisLabel: {
      color: '#64748b'
    }
  },
  yAxis: {
    type: 'category',
    data: popularSorted.value.map(item => item.name_prof_education),
    axisLabel: {
      color: '#334155',
      width: 220,
      lineHeight: 18,
      formatter: (value: string) => wrapLabel(value, 24, 3)
    }
  },
  series: [
    {
      name: 'Слушатели',
      type: 'bar',
      data: popularSorted.value.map(item => ({
        value: item.listeners,
        itemStyle: {
          color: '#1d4ed8',
          borderRadius: [0, 10, 10, 0]
        }
      })),
      label: {
        show: true,
        position: 'right',
        distance: 12,
        width: 120,
        overflow: 'break',
        lineHeight: 18,
        formatter: (params: { dataIndex: number, value: number }) =>
          `${popularSorted.value[params.dataIndex]?.educationtype || ''} · ${formatNumber(params.value)}`
      }
    }
  ]
}))

const revenueProgramsOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'shadow'
    },
    formatter: (params: Array<{ dataIndex: number, value: number }>) => {
      const item = revenueSorted.value[params[0]?.dataIndex || 0]
      if (!item) {
        return ''
      }

      return `${item.name_prof_education}<br/>${item.educationtype}<br/>${formatMoney(item.total_expected_revenue)}`
    }
  },
  grid: {
    top: 20,
    left: 250,
    right: 120,
    bottom: 24,
    containLabel: false
  },
  xAxis: {
    type: 'value',
    axisLabel: {
      color: '#64748b',
      formatter: (value: number) => formatCompact(value)
    }
  },
  yAxis: {
    type: 'category',
    data: revenueSorted.value.map(item => item.name_prof_education),
    axisLabel: {
      color: '#334155',
      width: 220,
      lineHeight: 18,
      formatter: (value: string) => wrapLabel(value, 24, 3)
    }
  },
  series: [
    {
      type: 'bar',
      data: revenueSorted.value.map(item => ({
        value: item.total_expected_revenue,
        itemStyle: {
          color: '#0f766e',
          borderRadius: [0, 10, 10, 0]
        }
      })),
      label: {
        show: true,
        position: 'right',
        distance: 12,
        width: 92,
        overflow: 'break',
        formatter: (params: { value: number }) => formatCompact(params.value)
      }
    }
  ]
}))

const educationTypeOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    formatter: (params: { name: string, value: number, percent: number }) =>
      `${params.name}<br/>${formatNumber(params.value)} слушателей<br/>${params.percent}%`
  },
  legend: {
    type: 'scroll',
    orient: 'vertical',
    right: 0,
    top: 20,
    bottom: 20,
    width: 180,
    itemGap: 14,
    formatter: (value: string) => wrapLabel(value, 18, 3)
  },
  series: [
    {
      type: 'pie',
      radius: ['45%', '72%'],
      center: ['30%', '50%'],
      avoidLabelOverlap: true,
      itemStyle: {
        borderColor: '#ffffff',
        borderWidth: 2
      },
      label: {
        show: false
      },
      labelLine: {
        show: false
      },
      emphasis: {
        label: {
          show: true,
          formatter: (params: { name: string, value: number, percent: number }) =>
            `${wrapLabel(params.name, 16, 3)}\n${formatNumber(params.value)} · ${params.percent}%`,
          lineHeight: 18,
          fontWeight: 600
        }
      },
      data: educationTypeBreakdown.value.map((item, index) => ({
        ...item,
        itemStyle: {
          color: ['#1d4ed8', '#0f766e', '#f97316', '#7c3aed', '#e11d48'][index % 5]
        }
      }))
    }
  ]
}))

const enrollmentSourcesOption = computed(() => {
  const months = Array.from(new Set(props.data.whoEnrolled.map(item => item.month)))
    .sort((left, right) => new Date(left).getTime() - new Date(right).getTime())
  const monthLabels = months.map(month =>
    new Date(month).toLocaleDateString('ru-RU', { month: 'short', year: 'numeric' })
  )
  const sources = Array.from(new Set(props.data.whoEnrolled.map(item => item.source)))
  const colors = ['#0f172a', '#0f766e', '#f97316', '#7c3aed', '#e11d48']

  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      bottom: 0
    },
    grid: {
      top: 30,
      left: 50,
      right: 30,
      bottom: 60
    },
    xAxis: {
      type: 'category',
      data: monthLabels,
      axisLabel: {
        color: '#475569'
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: '#64748b'
      }
    },
    series: sources.map((source, index) => ({
      name: source,
      type: 'bar',
      stack: 'total',
      emphasis: {
        focus: 'series'
      },
      itemStyle: {
        color: colors[index % colors.length],
        borderRadius: [6, 6, 0, 0]
      },
      data: months.map(month =>
        props.data.whoEnrolled.find(item => item.month === month && item.source === source)?.cnt || 0
      )
    }))
  }
})

const groupLoadOption = computed(() => buildSingleSeriesOption(
  groupsSorted.value,
  item => item.name_group,
  item => item.active_enrolled,
  '#f97316'
))

const divisionLoadOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    formatter: (params: { name: string, value: number, percent: number }) =>
      `${params.name}<br/>${formatNumber(params.value)} слушателей<br/>${params.percent}%`
  },
  legend: {
    bottom: 0
  },
  series: [
    {
      type: 'pie',
      radius: ['40%', '74%'],
      data: divisionsSorted.value.map((item, index) => ({
        name: item.divisioneducation,
        value: item.listeners,
        itemStyle: {
          color: ['#1d4ed8', '#0f766e', '#f97316', '#7c3aed', '#e11d48'][index % 5]
        }
      })),
      label: {
        formatter: '{b}\n{d}%'
      }
    }
  ]
}))

function buildSingleSeriesOption<T>(
  items: T[],
  labelGetter: (item: T) => string,
  valueGetter: (item: T) => number,
  color: string
) {
  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      top: 20,
      left: 220,
      right: 88,
      bottom: 24,
      containLabel: false
    },
    xAxis: {
      type: 'value',
      axisLabel: {
        color: '#64748b'
      }
    },
    yAxis: {
      type: 'category',
      data: items.map(labelGetter),
      axisLabel: {
        color: '#334155',
        width: 190,
        lineHeight: 18,
        formatter: (value: string) => wrapLabel(value, 22, 3)
      }
    },
    series: [
      {
        type: 'bar',
        data: items.map(item => ({
          value: valueGetter(item),
          itemStyle: {
            color,
            borderRadius: [0, 10, 10, 0]
          }
        })),
        label: {
          show: true,
          position: 'right',
          distance: 10,
          width: 72,
          overflow: 'break',
          formatter: (params: { value: number }) => formatNumber(params.value)
        }
      }
    ]
  }
}

const insightCards = computed(() => [
  {
    title: 'Лидер по спросу',
    value: bestProgramByListeners.value?.name_prof_education || 'Нет данных',
    note: bestProgramByListeners.value
      ? `${formatNumber(bestProgramByListeners.value.listeners)} слушателей`
      : 'Пока нет данных'
  },
  {
    title: 'Лидер по популярности',
    value: bestPopularProgram.value?.educationtype || 'Нет данных',
    note: bestPopularProgram.value
      ? `${bestPopularProgram.value.name_prof_education} · ${formatNumber(bestPopularProgram.value.listeners)}`
      : 'Пока нет данных'
  },
  {
    title: 'Максимальная выручка',
    value: bestRevenueProgram.value?.name_prof_education || 'Нет данных',
    note: bestRevenueProgram.value
      ? formatMoney(bestRevenueProgram.value.total_expected_revenue)
      : 'Пока нет данных'
  },
  {
    title: 'Самая загруженная группа',
    value: busiestGroup.value?.name_group || 'Нет данных',
    note: busiestGroup.value
      ? `${formatNumber(busiestGroup.value.active_enrolled)} активных зачислений`
      : 'Пока нет данных'
  }
])
</script>

<template>
  <section class="graphics-stack">
    <section class="graphics-hero">
      <div class="graphics-hero__copy">
        <p class="graphics-hero__eyebrow">Аналитический модуль</p>
        <h2>Графики по программам, подразделениям, выручке и нагрузке</h2>
        <p>
          Все ключевые срезы собраны в одном месте. Удобная аналитика для быстрого понимания ситуации и принятия решений по развитию программ и оптимизации ресурсов. (базовый посчет это активнные данные, а общий суммируется по всем записям за все время)
        </p>
      </div>

      <div class="graphics-hero__focus">
        <span>Программа-лидер</span>
        <strong>{{ bestProgramByListeners?.name_prof_education || 'Нет данных' }}</strong>
        <span v-if="bestProgramByListeners">
          +{{ countDeltaMap.get(bestProgramByListeners.name_prof_education) || 0 }} к базовому расчету
        </span>
      </div>
    </section>

    <div class="graphics-kpis">
      <AppStatCard
        title="Слушатели"
        :value="totalAccurateListeners"
        :description="`Базовый подсчет: ${totalListeners}`"
      />
      <AppStatCard
        title="Ожидаемая выручка"
        :value="formatMoney(totalExpectedRevenue)"
        :description="bestRevenueProgram?.educationtype || 'Нет типа обучения'"
      />
      <AppStatCard
        title="Нагрузка групп"
        :value="totalGroupLoad"
        :description="`Суммарная нагрузка подразделений: ${totalDivisionLoad}`"
      />
      <AppStatCard
        title="Топ доход"
        :value="bestRevenueProgram ? formatCompact(bestRevenueProgram.total_expected_revenue) : '0'"
        :description="bestRevenueProgram?.name_prof_education || 'Нет данных'"
      />
    </div>

    <AppCard title="Разделы аналитики">
      <div class="view-switcher">
        <button
          v-for="view in viewButtons"
          :key="view.id"
          type="button"
          class="view-switcher__button"
          :class="{ 'view-switcher__button--active': activeView === view.id }"
          @click="activeView = view.id"
        >
          <strong>{{ view.label }}</strong>
          <span>{{ view.description }}</span>
        </button>
      </div>
    </AppCard>

    <div v-if="activeView === 'overview'" class="graphics-stack">
      <div class="insight-grid">
        <article
          v-for="card in insightCards"
          :key="card.title"
          class="insight-card"
        >
          <span class="insight-card__title">{{ card.title }}</span>
          <strong>{{ card.value }}</strong>
          <span>{{ card.note }}</span>
        </article>
      </div>

      <div class="chart-grid">
        <AppCard title="Сравнение записей текущих записей и записи за все время">
          <p class="chart-caption">
            Двойной горизонтальный график помогает быстро понять, где точный расчет заметно выше базового.
          </p>
          <AnalyticsChart :option="listenersComparisonOption" height="430px" />
        </AppCard>

        <AppCard title="Распределение по форматам обучения">
          <p class="chart-caption">
            Кольцевая диаграмма показывает, какой формат обучения сейчас дает больше всего слушателей.
          </p>
          <AnalyticsChart :option="educationTypeOption" height="430px" />
        </AppCard>
      </div>
    </div>

    <div v-else-if="activeView === 'programs'" class="chart-grid">
      <AppCard title="Популярные программы">
        <p class="chart-caption">
          Здесь видно, какие программы сильнее по числу слушателей и в каком формате они обучаются.
        </p>
        <AnalyticsChart :option="popularProgramsOption" height="520px" />
      </AppCard>

      <AppCard title="Ожидаемая выручка по программам">
        <p class="chart-caption">
          Горизонтальный график помогает сразу увидеть лидеров по ожидаемой выручке.
        </p>
        <AnalyticsChart :option="revenueProgramsOption" height="520px" />
      </AppCard>
    </div>

    <div v-else-if="activeView === 'audience'" class="chart-grid">
      <AppCard title="Возрастная матрица">
        <p class="chart-caption">
          Тепловая карта показывает, в каких возрастных сегментах каждая программа выражена сильнее.
        </p>
        <AnalyticsChart :option="ageHeatmapOption" height="520px" />
      </AppCard>

      <AppCard title="Источники зачисления по месяцам">
        <p class="chart-caption">
          Стековая диаграмма позволяет быстро оценить структуру набора по каналам поступления.
        </p>
        <AnalyticsChart :option="enrollmentSourcesOption" height="520px" />
      </AppCard>
    </div>

    <div v-else class="chart-grid">
      <AppCard title="Нагрузка по группам">
        <p class="chart-caption">
          Чем длиннее полоса, тем выше текущая нагрузка группы по активным зачислениям.
        </p>
        <AnalyticsChart :option="groupLoadOption" height="480px" />
      </AppCard>

      <AppCard title="Нагрузка по подразделениям">
        <p class="chart-caption">
          Кольцевая диаграмма показывает долю каждого подразделения в общей текущей нагрузке.
        </p>
        <AnalyticsChart :option="divisionLoadOption" height="480px" />
      </AppCard>
    </div>

    <div class="summary-grid">
      <AppCard title="Лучше по выручке">
        <div class="summary-box summary-box--good">
          <strong>{{ bestRevenueProgram?.name_prof_education || 'Нет данных' }}</strong>
          <span v-if="bestRevenueProgram">{{ formatMoney(bestRevenueProgram.total_expected_revenue) }}</span>
        </div>
      </AppCard>

      <AppCard title="Слабее по выручке">
        <div class="summary-box summary-box--warn">
          <strong>{{ weakestRevenueProgram?.name_prof_education || 'Нет данных' }}</strong>
          <span v-if="weakestRevenueProgram">{{ formatMoney(weakestRevenueProgram.total_expected_revenue) }}</span>
        </div>
      </AppCard>
    </div>
  </section>
</template>

<style scoped>
.graphics-stack {
  display: grid;
  gap: 1rem;
}

.graphics-hero {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(280px, 0.8fr);
  gap: 1rem;
  padding: 1.5rem;
  border-radius: 1.6rem;
  background:
    radial-gradient(circle at top right, rgba(56, 189, 248, 0.2), transparent 30%),
    radial-gradient(circle at bottom left, rgba(20, 184, 166, 0.18), transparent 26%),
    linear-gradient(135deg, #0f172a 0%, #111827 46%, #0f766e 100%);
  color: #e2e8f0;
  box-shadow: 0 28px 80px rgba(15, 23, 42, 0.22);
}

.graphics-hero__eyebrow {
  margin: 0 0 0.8rem;
  font-size: 0.76rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #7dd3fc;
}

.graphics-hero__copy h2 {
  margin: 0;
  font-size: clamp(1.8rem, 3vw, 2.6rem);
}

.graphics-hero__copy p:last-child {
  margin-bottom: 0;
  color: rgba(226, 232, 240, 0.84);
}

.graphics-hero__focus {
  display: grid;
  align-content: end;
  gap: 0.45rem;
  padding: 1.25rem;
  border-radius: 1.3rem;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(226, 232, 240, 0.12);
}

.graphics-hero__focus span {
  color: rgba(226, 232, 240, 0.82);
}

.graphics-hero__focus strong {
  font-size: 1.2rem;
  color: #f8fafc;
}

.graphics-kpis {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 1rem;
}

.view-switcher {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.9rem;
}

.view-switcher__button {
  display: grid;
  gap: 0.3rem;
  text-align: left;
  padding: 1rem;
  border-radius: 1.15rem;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.95), rgba(241, 245, 249, 0.88));
  color: #0f172a;
  cursor: pointer;
  transition:
    transform 180ms ease,
    box-shadow 220ms ease,
    border-color 220ms ease;
}

.view-switcher__button:hover {
  transform: translateY(-2px);
  box-shadow: 0 16px 32px rgba(15, 23, 42, 0.08);
}

.view-switcher__button span {
  color: #64748b;
  font-size: 0.85rem;
}

.view-switcher__button--active {
  border-color: rgba(45, 212, 191, 0.35);
  background: linear-gradient(135deg, #0f172a 0%, #0f766e 100%);
  color: #f8fafc;
  box-shadow: 0 20px 42px rgba(15, 118, 110, 0.22);
}

.view-switcher__button--active span {
  color: rgba(226, 232, 240, 0.84);
}

.insight-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 1rem;
}

.insight-card {
  display: grid;
  gap: 0.45rem;
  padding: 1.1rem;
  border-radius: 1.2rem;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(240, 249, 255, 0.9));
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.06);
}

.insight-card strong {
  color: #0f172a;
  font-size: 1.02rem;
}

.insight-card span:last-child,
.insight-card__title {
  color: #64748b;
}

.chart-grid,
.summary-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.chart-caption {
  margin: 0 0 1rem;
  color: #64748b;
}

.summary-box {
  display: grid;
  gap: 0.4rem;
  padding: 1rem;
  border-radius: 1rem;
}

.summary-box strong {
  color: #0f172a;
}

.summary-box span {
  color: #475569;
}

.summary-box--good {
  background: linear-gradient(180deg, rgba(236, 253, 245, 0.96), rgba(167, 243, 208, 0.78));
}

.summary-box--warn {
  background: linear-gradient(180deg, rgba(255, 247, 237, 0.96), rgba(253, 186, 116, 0.76));
}

@media (max-width: 1200px) {
  .graphics-kpis,
  .insight-grid,
  .view-switcher {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 900px) {
  .graphics-hero,
  .graphics-kpis,
  .view-switcher,
  .insight-grid,
  .chart-grid,
  .summary-grid {
    grid-template-columns: 1fr;
  }
}
</style>
