import type {
  ApiListResponse,
  DivisionListenersMetric,
  EnrollmentSourceMetric,
  GraphicsDashboardResponse,
  GroupEnrollmentMetric,
  ProgramAgeMetric,
  ProgramListenersMetric,
  ProgramPopularMetric,
  ProgramRevenueMetric
} from '../../types/graphics'

export const getGraphicCount = async () => {
  const api = useApiClient()

  return await api.core<ApiListResponse<ProgramListenersMetric>>('graphic/count')
}

export const getGraphicCountAccurate = async () => {
  const api = useApiClient()

  return await api.core<ApiListResponse<ProgramListenersMetric>>('graphic/count/accurate')
}

export const getGraphicPopular = async () => {
  const api = useApiClient()

  return await api.core<ApiListResponse<ProgramPopularMetric>>('graphic/popular')
}

export const getGraphicWorthAccurate = async () => {
  const api = useApiClient()

  return await api.core<ApiListResponse<ProgramRevenueMetric>>('graphic/worth/accurate')
}

export const getGraphicAgeDiff = async () => {
  const api = useApiClient()

  return await api.core<ApiListResponse<ProgramAgeMetric>>('graphic/agediff')
}

export const getGraphicWhoEnrolled = async () => {
  const api = useApiClient()

  return await api.core<ApiListResponse<EnrollmentSourceMetric>>('graphic/whoenrolled')
}

export const getGraphicGroupLoad = async () => {
  const api = useApiClient()

  return await api.core<ApiListResponse<GroupEnrollmentMetric>>('graphic/group')
}

export const getGraphicDivisionLoad = async () => {
  const api = useApiClient()

  return await api.core<ApiListResponse<DivisionListenersMetric>>('graphic/division')
}

export const getGraphicsDashboard = async (): Promise<GraphicsDashboardResponse> => {
  const [
    count,
    countAccurate,
    popular,
    worthAccurate,
    ageDiff,
    whoEnrolled,
    groupLoad,
    divisionLoad
  ] = await Promise.all([
    getGraphicCount(),
    getGraphicCountAccurate(),
    getGraphicPopular(),
    getGraphicWorthAccurate(),
    getGraphicAgeDiff(),
    getGraphicWhoEnrolled(),
    getGraphicGroupLoad(),
    getGraphicDivisionLoad()
  ])

  return {
    count: count.data,
    countAccurate: countAccurate.data,
    popular: popular.data,
    worthAccurate: worthAccurate.data,
    ageDiff: ageDiff.data,
    whoEnrolled: whoEnrolled.data,
    groupLoad: groupLoad.data,
    divisionLoad: divisionLoad.data
  }
}
