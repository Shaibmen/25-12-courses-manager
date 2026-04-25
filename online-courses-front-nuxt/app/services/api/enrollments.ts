import type {
  AccurateEnrollmentItem,
  ContractorPayload,
  DocumentPayload,
  EnrollmentItem,
  EnrollmentPayload,
  EnrollmentUpdatePayload,
  ListenerEnrollmentContext
} from '../../types/enrollment'

type ApiDataResponse<T> = {
  data: T
  message?: string
}

type ApiMessageResponse = {
  message?: string
}

const getErrorMessage = (error: unknown, fallback: string) => toUserErrorMessage(error, fallback)

export const getEnrollments = async (page: number, filter: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<EnrollmentItem[]>>('enrollment/', {
      query: {
        page,
        filter
      }
    })

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить записи на курсы'))
  }
}

export const getEnrollmentByProgram = async (programId: string, page: number) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<EnrollmentItem[]>>(`enrollment/${programId}`, {
      query: {
        page
      }
    })

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить записи по программе'))
  }
}

export const getEnrollmentDetails = async (listenerId: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<EnrollmentItem[]>>(`enrollment/details/${listenerId}`)

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить записи слушателя'))
  }
}

export const getAccurateEnrollments = async () => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<AccurateEnrollmentItem[]>>('enrollment/accurate')

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить точные записи'))
  }
}

export const deleteEnrollment = async (listenerId: string, programId: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`enrollment/${listenerId}/${programId}`, {
      method: 'DELETE'
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось удалить запись'))
  }
}

export const updateEnrollmentRequest = async (
  listenerId: string,
  programId: string,
  payload: EnrollmentUpdatePayload
) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`enrollment/${listenerId}/${programId}`, {
      method: 'PUT',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось обновить запись'))
  }
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
