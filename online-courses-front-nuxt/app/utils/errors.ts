const extractStatusCode = (error: unknown) => {
  if (
    typeof error === 'object' &&
    error !== null &&
    'statusCode' in error &&
    typeof error.statusCode === 'number'
  ) {
    return error.statusCode
  }

  if (
    typeof error === 'object' &&
    error !== null &&
    'status' in error &&
    typeof error.status === 'number'
  ) {
    return error.status
  }

  return null
}

export const toUserErrorMessage = (error: unknown, fallback: string) => {
  if (
    typeof error === 'object' &&
    error !== null &&
    'data' in error &&
    typeof error.data === 'object' &&
    error.data !== null &&
    'message' in error.data &&
    typeof error.data.message === 'string' &&
    error.data.message.trim()
  ) {
    const code = extractStatusCode(error)
    return code ? `Код ${code}: ${error.data.message}` : error.data.message
  }

  const code = extractStatusCode(error)

  if (code) {
    if (code === 404) {
      return `Код ${code}: не удалось найти данные`
    }

    if (code === 401 || code === 403) {
      return `Код ${code}: недостаточно прав для выполнения действия`
    }

    if (code >= 500) {
      return `Код ${code}: ошибка сервера, попробуйте позже`
    }

    return `Код ${code}: ${fallback}`
  }

  if (error instanceof Error && /fetch failed|Failed to fetch|NetworkError/i.test(error.message)) {
    return `Код сети: ${fallback}`
  }

  return fallback
}
