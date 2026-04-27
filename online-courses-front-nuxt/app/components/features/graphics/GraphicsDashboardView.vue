<script setup lang="ts">
import type {
  EnrollmentSourceMetric,
  GraphicsDashboardResponse,
  ProgramAgeMetric,
  ProgramListenersMetric,
  ProgramRevenueMetric
} from '../../../types/graphics'
import AppCard from '../../ui/AppCard.vue'
import AppStatCard from '../../ui/AppStatCard.vue'
import AnalyticsAgeMatrix from './AnalyticsAgeMatrix.vue'
import AnalyticsHorizontalBars from './AnalyticsHorizontalBars.vue'
import AnalyticsSourceTimeline from './AnalyticsSourceTimeline.vue'

const props = defineProps<{
  data: GraphicsDashboardResponse
}>()

const compactNumber = new Intl.NumberFormat('ru-RU', {
  notation: 'compact',
  maximumFractionDigits: 1
})

const moneyFormatter = new Intl.NumberFormat('ru-RU', {
  style: 'currency',
  currency: 'RUB',
  maximumFractionDigits: 0
})

const formatCompact = (value: number) => compactNumber.format(value)
const formatMoney = (value: number) => moneyFormatter.format(value)

const sumBy = <T>(items: T[], getValue: (item: T) => number) =>
  items.reduce((total, item) => total + getValue(item), 0)

const buildDeltaMap = (
  baseline: ProgramListenersMetric[] | ProgramRevenueMetric[],
  accurate: ProgramListenersMetric[] | ProgramRevenueMetric[],
  getLabel: (item: ProgramListenersMetric | ProgramRevenueMetric) => string,
  getValue: (item: ProgramListenersMetric | ProgramRevenueMetric) => number
) => {
  const baselineMap = new Map(baseline.map(item => [getLabel(item), getValue(item)]))

  return new Map(
    accurate.map(item => {
      const label = getLabel(item)
      const delta = getValue(item) - (baselineMap.get(label) || 0)
      return [label, delta]
    })
  )
}

const countDeltaMap = computed(() =>
  buildDeltaMap(
    props.data.count,
    props.data.countAccurate,
    item => item.name_prof_education,
    item => item.listeners
  )
)

const worthDeltaMap = computed(() =>
  buildDeltaMap(
    props.data.worth,
    props.data.worthAccurate,
    item => item.name_prof_education,
    item => item.total_revenue
  )
)

const totalListeners = computed(() => sumBy(props.data.count, item => item.listeners))
const totalAccurateListeners = computed(() => sumBy(props.data.countAccurate, item => item.listeners))
const totalRevenue = computed(() => sumBy(props.data.worth, item => item.total_revenue))
const totalAccurateRevenue = computed(() => sumBy(props.data.worthAccurate, item => item.total_revenue))
const totalActiveGroups = computed(() => sumBy(props.data.groupLoad, item => item.active_enrolled))
const totalDivisionsLoad = computed(() => sumBy(props.data.divisionLoad, item => item.listeners))

const bestProgramByListeners = computed(() =>
  [...props.data.countAccurate].sort((a, b) => b.listeners - a.listeners)[0] || null
)

const bestProgramByRevenue = computed(() =>
  [...props.data.worthAccurate].sort((a, b) => b.total_revenue - a.total_revenue)[0] || null
)

const topPrograms = computed(() =>
  props.data.countAccurate.map(item => ({
    label: item.name_prof_education,
    value: item.listeners,
    hint: `+${countDeltaMap.value.get(item.name_prof_education) || 0} к базовому расчету`
  }))
)

const topRevenue = computed(() =>
  props.data.worthAccurate.map(item => ({
    label: item.name_prof_education,
    value: item.total_revenue,
    hint: `Разница ${formatMoney(worthDeltaMap.value.get(item.name_prof_education) || 0)}`
  }))
)

const groupOccupancy = computed(() =>
  props.data.groupLoad.map(item => ({
    label: item.name_group,
    value: item.active_enrolled,
    hint: item.active_enrolled ? 'Есть активные слушатели' : 'Пока без активных слушателей'
  }))
)

const divisionOccupancy = computed(() =>
  props.data.divisionLoad.map(item => ({
    label: item.divisioneducation,
    value: item.listeners,
    hint: 'Суммарная активность по подразделению'
  }))
)

const ageRanges = computed(() => {
  const unique = Array.from(new Set(props.data.ageDiff.map(item => item.age_range)))

  return unique.sort((left, right) => {
    const leftStart = Number(left.split('-')[0]) || 0
    const rightStart = Number(right.split('-')[0]) || 0
    return leftStart - rightStart
  })
})

const ageDiffMatrix = computed(() => {
  const byProgram = new Map<string, ProgramAgeMetric[]>()

  props.data.ageDiff.forEach(item => {
    const bucket = byProgram.get(item.name_prof_education) || []
    bucket.push(item)
    byProgram.set(item.name_prof_education, bucket)
  })

  return Array.from(byProgram.entries()).map(([label, items]) => ({
    label,
    values: ageRanges.value.map(ageRange => ({
      ageRange,
      listeners: items.find(item => item.age_range === ageRange)?.listeners || 0
    }))
  }))
})

const sourcePalette = [
  'linear-gradient(135deg, #0f172a 0%, #1d4ed8 100%)',
  'linear-gradient(135deg, #0f766e 0%, #14b8a6 100%)',
  'linear-gradient(135deg, #9a3412 0%, #f97316 100%)',
  'linear-gradient(135deg, #7c2d12 0%, #fb7185 100%)',
  'linear-gradient(135deg, #4338ca 0%, #818cf8 100%)'
]

const sourceTimeline = computed(() => {
  const sources = Array.from(new Set(props.data.whoEnrolled.map(item => item.source)))
  const sourceColors = new Map(sources.map((source, index) => [source, sourcePalette[index % sourcePalette.length]]))
  const byMonth = new Map<string, EnrollmentSourceMetric[]>()

  props.data.whoEnrolled.forEach(item => {
    const bucket = byMonth.get(item.month) || []
    bucket.push(item)
    byMonth.set(item.month, bucket)
  })

  return Array.from(byMonth.entries())
    .sort(([left], [right]) => new Date(left).getTime() - new Date(right).getTime())
    .map(([month, items]) => {
      const total = sumBy(items, item => item.cnt)

      return {
        monthLabel: new Date(month).toLocaleDateString('ru-RU', {
          month: 'long',
          year: 'numeric'
        }),
        total,
        segments: items.map(item => ({
          label: item.source,
          value: item.cnt,
          color: sourceColors.get(item.source) || sourcePalette[0]
        }))
      }
    })
})
</script>

<template>
  <section class="graphics-stack">
    <section class="graphics-hero">
      <div class="graphics-hero__copy">
        <p class="graphics-hero__eyebrow">Аналитика и графики</p>
        <h2>Сводка по курсам, выручке, возрастам и каналам зачисления</h2>
        <p>
          Раздел собирает все ключевые графики из API и показывает, где растет набор,
          какие программы приносят больше выручки и как распределяется аудитория.
        </p>
      </div>

      <div class="graphics-hero__focus">
        <span>Лидер по слушателям</span>
        <strong>{{ bestProgramByListeners?.name_prof_education || 'Нет данных' }}</strong>
        <span v-if="bestProgramByListeners">
          {{ bestProgramByListeners.listeners }} слушателей в точном расчете
        </span>
      </div>
    </section>

    <div class="graphics-kpis">
      <AppStatCard
        title="Слушатели"
        :value="totalAccurateListeners"
        :description="`Базовый расчет: ${totalListeners}`"
      />
      <AppStatCard
        title="Выручка"
        :value="formatMoney(totalAccurateRevenue)"
        :description="`Без уточнения: ${formatMoney(totalRevenue)}`"
      />
      <AppStatCard
        title="Активные группы"
        :value="totalActiveGroups"
        :description="`Суммарная нагрузка подразделений: ${totalDivisionsLoad}`"
      />
      <AppStatCard
        title="Топ по доходу"
        :value="bestProgramByRevenue ? formatCompact(bestProgramByRevenue.total_revenue) : '0'"
        :description="bestProgramByRevenue?.name_prof_education || 'Нет данных'"
      />
    </div>

    <div class="graphics-grid">
      <AppCard title="Популярность программ">
        <p class="chart-caption">
          Точный расчет по числу слушателей. Подпись справа показывает разницу с базовым подсчетом.
        </p>
        <AnalyticsHorizontalBars
          :items="topPrograms"
          color-from="#0f172a"
          color-to="#38bdf8"
        />
      </AppCard>

      <AppCard title="Выручка по программам">
        <p class="chart-caption">
          Приоритетно показываем точную выручку, чтобы было проще сравнивать финансовую отдачу программ.
        </p>
        <AnalyticsHorizontalBars
          :items="topRevenue"
          color-from="#0f766e"
          color-to="#2dd4bf"
          value-suffix=" ₽"
        />
      </AppCard>
    </div>

    <div class="graphics-grid graphics-grid--secondary">
      <AppCard title="Возрастные группы по программам">
        <p class="chart-caption">
          Чем насыщеннее ячейка, тем больше слушателей в возрастной группе.
        </p>
        <AnalyticsAgeMatrix
          :rows="ageDiffMatrix"
          :age-ranges="ageRanges"
        />
      </AppCard>

      <AppCard title="Источники зачисления по месяцам">
        <p class="chart-caption">
          Каждый месяц разбит на сегменты по источникам. Так проще увидеть структуру набора.
        </p>
        <AnalyticsSourceTimeline :items="sourceTimeline" />
      </AppCard>
    </div>

    <div class="graphics-grid">
      <AppCard title="Нагрузка по группам">
        <p class="chart-caption">
          Быстрый срез по текущим учебным группам и числу активных зачислений.
        </p>
        <AnalyticsHorizontalBars
          :items="groupOccupancy"
          color-from="#7c2d12"
          color-to="#fb7185"
        />
      </AppCard>

      <AppCard title="Нагрузка по подразделениям">
        <p class="chart-caption">
          Показывает, какие подразделения сейчас ведут больше всего слушателей.
        </p>
        <AnalyticsHorizontalBars
          :items="divisionOccupancy"
          color-from="#312e81"
          color-to="#818cf8"
        />
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
    radial-gradient(circle at top right, rgba(56, 189, 248, 0.22), transparent 30%),
    radial-gradient(circle at bottom left, rgba(20, 184, 166, 0.2), transparent 26%),
    linear-gradient(135deg, #0f172a 0%, #111827 46%, #0f766e 100%);
  color: #e2e8f0;
  box-shadow: 0 28px 80px rgba(15, 23, 42, 0.22);
}

.graphics-hero__copy h2 {
  margin: 0;
  font-size: clamp(1.8rem, 3vw, 2.6rem);
}

.graphics-hero__copy p:last-child {
  margin-bottom: 0;
  color: rgba(226, 232, 240, 0.84);
}

.graphics-hero__eyebrow {
  margin: 0 0 0.8rem;
  font-size: 0.76rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #7dd3fc;
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
  color: rgba(226, 232, 240, 0.8);
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

.graphics-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.graphics-grid--secondary {
  align-items: start;
}

.chart-caption {
  margin: 0 0 1rem;
  color: #64748b;
}

@media (max-width: 1100px) {
  .graphics-kpis {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 900px) {
  .graphics-hero,
  .graphics-grid,
  .graphics-kpis {
    grid-template-columns: 1fr;
  }
}
</style>
