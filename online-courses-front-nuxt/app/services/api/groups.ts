import type { GroupItem, GroupPayload, GroupScheduleItem } from '../../types/group'

type ApiDataResponse<T> = {
  data: T
  message?: string
}

type ApiMessageResponse = {
  message?: string
}

type GroupApiItem = {
  group: string
  name_group: string
  rapspisanie?: GroupScheduleItem[]
  raspisanie?: GroupScheduleItem[]
}

const getErrorMessage = (error: unknown, fallback: string) => toUserErrorMessage(error, fallback)

const getAuthorizedHeaders = (): Record<string, string> => {
  const auth = useAuthState()
  const headers: Record<string, string> = {}

  if (auth.value.accessToken) {
    headers.Authorization = `Bearer ${auth.value.accessToken}`
  }

  return headers
}

const normalizeGroupItem = (item: GroupApiItem): GroupItem => ({
  group: item.group,
  name_group: item.name_group,
  rapspisanie: Array.isArray(item.rapspisanie)
    ? item.rapspisanie
    : Array.isArray(item.raspisanie)
      ? item.raspisanie
      : []
})

export const getGroups = async (page: number, filter: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<GroupApiItem[]>>('group/', {
      query: {
        page,
        filter
      }
    })

    return (response.data || []).map(normalizeGroupItem)
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить группы'))
  }
}

export const getGroupDetails = async (id: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<GroupApiItem>>(`group/details/${id}`)
    return normalizeGroupItem(response.data)
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить группу'))
  }
}

export const createGroupRequest = async (payload: GroupPayload) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>('group/', {
      method: 'POST',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось создать группу'))
  }
}

export const updateGroupRequest = async (id: string, payload: GroupPayload) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`group/${id}`, {
      method: 'PUT',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось обновить группу'))
  }
}

export const deleteGroupRequest = async (id: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`group/${id}`, {
      method: 'DELETE'
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось удалить группу'))
  }
}

export const exportGroupSchedule = async (id: string) => {
  const config = useRuntimeConfig()

  const response = await fetch(
    new URL(`group/export/${id}`, `${config.public.apiUrlCore}/`).toString(),
    { headers: getAuthorizedHeaders() }
  )

  if (!response.ok) {
    const message = await response.text().catch(() => '')
    throw new Error(message || `Не удалось скачать расписание (${response.status})`)
  }

  return await response.blob()
}
