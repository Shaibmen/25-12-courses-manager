export type ApiListResponse<T> = {
  data: T[]
}

export type ProgramListenersMetric = {
  name_prof_education: string
  listeners: number
}

export type ProgramRevenueMetric = {
  name_prof_education: string
  total_revenue: number
}

export type ProgramAgeMetric = {
  name_prof_education: string
  age_range: string
  listeners: number
}

export type EnrollmentSourceMetric = {
  month: string
  source: string
  cnt: number
}

export type GroupEnrollmentMetric = {
  name_group: string
  active_enrolled: number
}

export type DivisionListenersMetric = {
  divisioneducation: string
  listeners: number
}

export type GraphicsDashboardResponse = {
  count: ProgramListenersMetric[]
  countAccurate: ProgramListenersMetric[]
  worth: ProgramRevenueMetric[]
  worthAccurate: ProgramRevenueMetric[]
  ageDiff: ProgramAgeMetric[]
  whoEnrolled: EnrollmentSourceMetric[]
  groupLoad: GroupEnrollmentMetric[]
  divisionLoad: DivisionListenersMetric[]
}
