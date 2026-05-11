import type {
  ProgramListItem,
  ProgramPayload,
  ProgramUpdatePayload
} from '../../types/program'

type ApiDataResponse<T> = {
  data: T
  message?: string
}

type ApiMessageResponse = {
  message?: string
}

const getErrorMessage = (error: unknown, fallback: string) => toUserErrorMessage(error, fallback)

export const getPrograms = async (page: number, filter: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<ProgramListItem[]>>('programeducation/', {
      query: {
        page,
        filter
      }
    })

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить программы обучения'))
  }
}

export const getProgramDetails = async (id: string) => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<ProgramListItem>>(`programeducation/${id}`)

    return response.data
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить программу обучения'))
  }
}

export const createProgramRequest = async (payload: ProgramPayload) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>('programeducation/', {
      method: 'POST',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось создать программу обучения'))
  }
}

export const updateProgramRequest = async (id: string, payload: ProgramUpdatePayload) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`programeducation/${id}`, {
      method: 'PUT',
      body: payload
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось обновить программу обучения'))
  }
}

export const deleteProgramRequest = async (id: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`programeducation/${id}`, {
      method: 'DELETE'
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось удалить программу обучения'))
  }
}
