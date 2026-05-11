import type { ExecuterItem, ExecuterPayload } from '../../types/executer'

type ApiDataResponse<T> = {
  data: T
  message?: string
}

type ApiMessageResponse = {
  message?: string
}

const getErrorMessage = (error: unknown, fallback: string) => toUserErrorMessage(error, fallback)

export const getExecuters = async (page: number, filter: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<ExecuterItem[]>>('executer/', {
      query: { page, filter }
    })

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить список исполнителей'))
  }
}

export const createExecuterRequest = async (payload: ExecuterPayload) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>('executer/', {
      method: 'POST',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось создать исполнителя'))
  }
}

export const deleteExecuterRequest = async (id: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`executer/${id}`, {
      method: 'DELETE'
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось удалить исполнителя'))
  }
}
