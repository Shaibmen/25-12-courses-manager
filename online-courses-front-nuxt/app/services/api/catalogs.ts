import type { DivisionItem, EducationTypeItem } from '../../types/catalogs'
import type { LevelEducationItem } from '../../types/listener'

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

export const getDivisions = async (filter = '') => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<DivisionItem[]>>('divisions/', {
      query: { filter }
    })

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить подразделения'))
  }
}

export const createDivisionRequest = async (divisions: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>('divisions/', {
      method: 'POST',
      body: { divisions }
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось создать подразделение'))
  }
}

export const updateDivisionRequest = async (id: string, divisions: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`divisions/${id}`, {
      method: 'PUT',
      body: { divisions }
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось обновить подразделение'))
  }
}

export const deleteDivisionRequest = async (id: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`divisions/${id}`, {
      method: 'DELETE'
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось удалить подразделение'))
  }
}

export const getEducationTypes = async (filter = '') => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<EducationTypeItem[]>>('educationtype/', {
      query: { filter }
    })

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить типы обучения'))
  }
}

export const createEducationTypeRequest = async (typeName: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>('educationtype/', {
      method: 'POST',
      body: { type_name: typeName }
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось создать тип обучения'))
  }
}

export const updateEducationTypeRequest = async (id: string, typeName: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`educationtype/${id}`, {
      method: 'PUT',
      body: { type_name: typeName }
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось обновить тип обучения'))
  }
}

export const deleteEducationTypeRequest = async (id: string) => {
  const api = useApiClient()

  try {
    return await api.core<ApiMessageResponse>(`educationtype/${id}`, {
      method: 'DELETE'
    })
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось удалить тип обучения'))
  }
}

export const getEducationLevelsCatalog = async (filter = '') => {
  const api = useApiClient()

  try {
    const response = await api.core<ApiDataResponse<LevelEducationItem[]>>('leveleducation/', {
      query: { filter }
    })

    return response.data || []
  } catch (error) {
    throw new Error(getErrorMessage(error, 'Не удалось загрузить уровни обучения'))
  }
}
