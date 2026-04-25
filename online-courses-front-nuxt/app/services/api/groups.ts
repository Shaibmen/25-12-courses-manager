import type { GroupItem } from '../../types/group'

type ApiDataResponse<T> = {
  data: T
  message?: string
}

const getErrorMessage = (error: unknown, fallback: string) => toUserErrorMessage(error, fallback)

export const getGroups = async (page: number, filter: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<GroupItem[]>>('group/', {
      query: {
        page,
        filter
      }
    })

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить группы'))
  }
}

export const getGroupDetails = async (id: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<GroupItem>>(`group/details/${id}`)
    return response.data
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить группу'))
  }
}
