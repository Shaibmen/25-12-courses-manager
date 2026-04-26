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
    throw new Error(getErrorMessage(error, 'РќРµ СѓРґР°Р»РѕСЃСЊ Р·Р°РіСЂСѓР·РёС‚СЊ РіСЂСѓРїРїС‹'))
  }
}

export const getGroupDetails = async (id: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<GroupApiItem>>(`group/details/${id}`)
    return normalizeGroupItem(response.data)
  } catch (error) {
    throw new Error(getErrorMessage(error, 'РќРµ СѓРґР°Р»РѕСЃСЊ Р·Р°РіСЂСѓР·РёС‚СЊ РіСЂСѓРїРїСѓ'))
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
    throw new Error(getErrorMessage(error, 'РќРµ СѓРґР°Р»РѕСЃСЊ СЃРѕР·РґР°С‚СЊ РіСЂСѓРїРїСѓ'))
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
    throw new Error(getErrorMessage(error, 'РќРµ СѓРґР°Р»РѕСЃСЊ РѕР±РЅРѕРІРёС‚СЊ РіСЂСѓРїРїСѓ'))
  }
}

export const deleteGroupRequest = async (id: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`group/${id}`, {
      method: 'DELETE'
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'РќРµ СѓРґР°Р»РѕСЃСЊ СѓРґР°Р»РёС‚СЊ РіСЂСѓРїРїСѓ'))
  }
}
