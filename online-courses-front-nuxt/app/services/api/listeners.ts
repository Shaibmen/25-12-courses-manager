import type {
  EnrollmentProgramDetails,
  LevelEducationItem,
  ListenerDetailsResponse,
  ListenerFormPayload,
  ListenerListItem
} from '../../types/listener'

type ApiDataResponse<T> = {
  data: T
  message?: string
}

type ApiMessageResponse = {
  message?: string
}

const getErrorMessage = (error: unknown, fallback: string) => {
  if (
    typeof error === 'object' &&
    error !== null &&
    'data' in error &&
    typeof error.data === 'object' &&
    error.data !== null &&
    'message' in error.data &&
    typeof error.data.message === 'string'
  ) {
    return error.data.message
  }

  if (error instanceof Error && error.message) {
    return error.message
  }

  return fallback
}

const getAuthorizedHeaders = (): Record<string, string> => {
  const auth = useAuthState()
  const headers: Record<string, string> = {}

  if (auth.value.accessToken) {
    headers.Authorization = `Bearer ${auth.value.accessToken}`
  }

  return headers
}

export const getListeners = async (page: number, filter: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<ListenerListItem[]>>('listener/', {
      query: {
        page,
        filter
      }
    })

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить список слушателей'))
  }
}

export const getListenerDetails = async (id: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<ListenerDetailsResponse>>(`listener/details/${id}`)

    return response.data
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить данные слушателя'))
  }
}

export const createListenerRequest = async (payload: ListenerFormPayload) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>('listener/', {
      method: 'POST',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось создать слушателя'))
  }
}

export const updateListenerRequest = async (id: string, payload: ListenerFormPayload) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`listener/${id}`, {
      method: 'PUT',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось обновить данные слушателя'))
  }
}

export const deleteListenerRequest = async (id: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`listener/${id}`, {
      method: 'DELETE'
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось удалить слушателя'))
  }
}

export const getEducationLevels = async () => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<LevelEducationItem[]>>('leveleducation/', {
      query: {
        filter: ''
      }
    })

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить уровни образования'))
  }
}

export const getListenerEnrollments = async (id: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<EnrollmentProgramDetails[]>>(`enrollment/details/${id}`)

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить курсы слушателя'))
  }
}

export const getListenerFiles = async (cardName: string) => {
  const config = useRuntimeConfig()

  if (!cardName) {
    return []
  }

  try {
    const response = await $fetch<unknown>(new URL(`exists?card-name=${encodeURIComponent(cardName)}`, `${config.public.apiUrlDoc}/`).toString(), {
      headers: getAuthorizedHeaders()
    })

    if (Array.isArray(response)) {
      return response as string[]
    }

    if (
      typeof response === 'object' &&
      response !== null &&
      'data' in response &&
      Array.isArray(response.data)
    ) {
      return response.data as string[]
    }

    if (
      typeof response === 'object' &&
      response !== null &&
      'data' in response &&
      typeof response.data === 'object' &&
      response.data !== null &&
      'files' in response.data &&
      Array.isArray(response.data.files)
    ) {
      return response.data.files as string[]
    }

    return []
  } catch {
    return []
  }
}

export const downloadListenerFile = async (fileName: string) => {
  const config = useRuntimeConfig()

  const response = await fetch(
    new URL(`download?card-name=${encodeURIComponent(fileName)}`, `${config.public.apiUrlDoc}/`).toString(),
    {
      headers: getAuthorizedHeaders()
    }
  )

  if (!response.ok) {
    const message = await response.text().catch(() => '')

    throw new Error(message || `Не удалось скачать файл (${response.status})`)
  }

  return await response.blob()
}
