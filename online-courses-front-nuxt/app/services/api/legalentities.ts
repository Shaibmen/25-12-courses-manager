import type {
  LegalEntityDetailsResponse,
  LegalEntityListItem,
  LegalEntityPayload
} from '../../types/legalentity'

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

export const getLegalEntities = async (page: number, filter: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<LegalEntityListItem[]>>('legalentity/', {
      query: { page, filter }
    })

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить юридические лица'))
  }
}

export const getLegalEntityDetails = async (id: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<LegalEntityDetailsResponse>>(`legalentity/details/${id}`)
    return response.data
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить юридическое лицо'))
  }
}

export const createLegalEntityRequest = async (payload: LegalEntityPayload) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>('legalentity/', {
      method: 'POST',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось создать юридическое лицо'))
  }
}

export const updateLegalEntityRequest = async (id: string, payload: LegalEntityPayload) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`legalentity/${id}`, {
      method: 'PUT',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось обновить юридическое лицо'))
  }
}

export const deleteLegalEntityRequest = async (id: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`legalentity/${id}`, {
      method: 'DELETE'
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось удалить юридическое лицо'))
  }
}

export const getLegalEntityFiles = async (companyName: string) => {
  const config = useRuntimeConfig()

  if (!companyName) {
    return []
  }

  try {
    const response = await $fetch<unknown>(
      new URL(`exists?card-name=${encodeURIComponent(companyName)}`, `${config.public.apiUrlDoc}/`).toString(),
      { headers: getAuthorizedHeaders() }
    )

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

export const downloadLegalEntityFile = async (fileName: string) => {
  const config = useRuntimeConfig()

  const response = await fetch(
    new URL(`download?card-name=${encodeURIComponent(fileName)}`, `${config.public.apiUrlDoc}/`).toString(),
    { headers: getAuthorizedHeaders() }
  )

  if (!response.ok) {
    const message = await response.text().catch(() => '')
    throw new Error(message || `Не удалось скачать файл (${response.status})`)
  }

  return await response.blob()
}
