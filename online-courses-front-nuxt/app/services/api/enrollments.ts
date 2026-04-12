import type {
  ContractorPayload,
  DocumentPayload,
  EnrollmentPayload,
  ListenerEnrollmentContext
} from '../../types/enrollment'

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

export const getListenerEnrollmentContext = async (listenerId: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<ListenerEnrollmentContext>>(`listener/details/${listenerId}`)
    return response.data
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить данные слушателя для записи'))
  }
}

export const upsertContractorRequest = async (listenerId: string, payload: ContractorPayload) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`contractor/${listenerId}`, {
      method: 'POST',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось сохранить заказчика'))
  }
}

export const deleteContractorRequest = async (contractorId: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`contractor/${contractorId}`, {
      method: 'DELETE'
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось удалить заказчика'))
  }
}

export const createEnrollmentRequest = async (payload: EnrollmentPayload) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>('enrollment/', {
      method: 'POST',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось создать запись на курс'))
  }
}

export const createEnrollmentDocumentRequest = async (payload: DocumentPayload) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>('document/', {
      method: 'POST',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось создать документы по записи'))
  }
}
