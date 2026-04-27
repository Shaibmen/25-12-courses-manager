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

const getErrorMessage = (error: unknown, fallback: string) => toUserErrorMessage(error, fallback)

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
    const response = await $fetch<unknown>(
      new URL(`exists?card-name=${encodeURIComponent(cardName)}`, `${config.public.apiUrlDoc}/`).toString(),
      {
        headers: getAuthorizedHeaders()
      }
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

const normalizeScanDiplomResponse = (response: unknown): string[] => {
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
    'files' in response &&
    Array.isArray(response.files)
  ) {
    return response.files as string[]
  }

  return []
}

export const buildListenerScanDiplomBaseName = (
  secondName?: string | null,
  firstName?: string | null,
  middleName?: string | null
) => {
  return [secondName, firstName, middleName]
    .map(value => String(value || '').trim())
    .filter(Boolean)
    .join('_')
}

export const getListenerScanDiplomFiles = async (baseName: string) => {
  const config = useRuntimeConfig()

  if (!baseName) {
    return []
  }

  try {
    const response = await $fetch<unknown>(
      new URL(`scan-diplom-exists?dogovor-name=${encodeURIComponent(baseName)}`, `${config.public.apiUrlDoc}/`).toString(),
      {
        headers: getAuthorizedHeaders()
      }
    )

    return normalizeScanDiplomResponse(response)
  } catch {
    return []
  }
}

export const downloadListenerScanDiplomFile = async (fileName: string) => {
  const config = useRuntimeConfig()

  const response = await fetch(
    new URL(`scan-diplom-download?dogovor-name=${encodeURIComponent(fileName)}`, `${config.public.apiUrlDoc}/`).toString(),
    {
      headers: getAuthorizedHeaders()
    }
  )

  if (!response.ok) {
    const message = await response.text().catch(() => '')

    throw new Error(message || `Не удалось скачать скан диплома (${response.status})`)
  }

  return await response.blob()
}

export const deleteListenerScanDiplomFile = async (fileName: string) => {
  const config = useRuntimeConfig()
  const normalizedName = fileName.toLowerCase().endsWith('.pdf')
    ? fileName.slice(0, -4)
    : fileName

  const response = await fetch(
    new URL(`scan-diplom-delete?dogovor-name=${encodeURIComponent(normalizedName)}`, `${config.public.apiUrlDoc}/`).toString(),
    {
      method: 'DELETE',
      headers: getAuthorizedHeaders()
    }
  )

  if (!response.ok) {
    const message = await response.text().catch(() => '')

    throw new Error(message || `Не удалось удалить скан диплома (${response.status})`)
  }

  return true
}

export const uploadListenerScanDiplomFile = async (scanName: string, file: File) => {
  const config = useRuntimeConfig()
  const authHeaders = getAuthorizedHeaders()
  const formData = new FormData()

  formData.append('photo', file, file.name)

  try {
    const response = await fetch(
      new URL(`scan/${encodeURIComponent(scanName)}?name-scan=${encodeURIComponent(scanName)}`, `${config.public.apiUrlCore}/`).toString(),
      {
        method: 'POST',
        headers: authHeaders,
        body: formData
      }
    )

    if (!response.ok) {
      const message = await response.text().catch(() => '')

      throw new Error(message || `Не удалось загрузить скан диплома (${response.status})`)
    }

    return await response.json().catch(() => ({}))
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить скан диплома'))
  }
}
