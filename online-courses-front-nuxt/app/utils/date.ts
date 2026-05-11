const padDatePart = (value: string | number) => String(value).padStart(2, '0')

export const normalizeApiDate = (value?: string | null) => {
  if (!value) {
    return ''
  }

  const trimmed = String(value).trim()

  if (!trimmed) {
    return ''
  }

  const normalizedChunk = (((trimmed.split(' ')[0] || '').split('T')[0]) || '').replace(/\//g, '.')

  if (/^\d{4}-\d{2}-\d{2}$/.test(normalizedChunk)) {
    return normalizedChunk
  }

  if (/^\d{2}\.\d{2}\.\d{4}$/.test(normalizedChunk)) {
    const [day, month, year] = normalizedChunk.split('.')
    return day && month && year ? `${year}-${month}-${day}` : ''
  }

  if (/^\d{2}-\d{2}-\d{4}$/.test(normalizedChunk)) {
    const [day, month, year] = normalizedChunk.split('-')
    return day && month && year ? `${year}-${month}-${day}` : ''
  }

  const parsedDate = new Date(trimmed)

  if (Number.isNaN(parsedDate.getTime())) {
    return ''
  }

  return `${parsedDate.getFullYear()}-${padDatePart(parsedDate.getMonth() + 1)}-${padDatePart(parsedDate.getDate())}`
}

export const formatApiDate = (value?: string | null) => {
  const normalized = normalizeApiDate(value)

  if (!normalized) {
    return '—'
  }

  const [year, month, day] = normalized.split('-')
  return `${day}.${month}.${year}`
}

export const getDateSortValue = (value?: string | null) => {
  const normalized = normalizeApiDate(value)

  if (!normalized) {
    return Number.POSITIVE_INFINITY
  }

  return Number(normalized.replaceAll('-', ''))
}
