import { normalizeApiDate } from './date'

const RUSSIAN_PUBLIC_HOLIDAY_DATES = [
  '2026-01-01',
  '2026-01-02',
  '2026-01-03',
  '2026-01-04',
  '2026-01-05',
  '2026-01-06',
  '2026-01-07',
  '2026-01-08',
  '2026-01-09',
  '2026-02-23',
  '2026-03-08',
  '2026-03-09',
  '2026-05-01',
  '2026-05-09',
  '2026-05-11',
  '2026-06-12',
  '2026-11-04',
  '2026-12-31'
]

const holidayDateSet = new Set(RUSSIAN_PUBLIC_HOLIDAY_DATES)

export const isRussianPublicHoliday = (value?: string | null) => {
  const normalized = normalizeApiDate(value)
  return normalized ? holidayDateSet.has(normalized) : false
}

export const getRussianPublicHolidayDates = () => [...RUSSIAN_PUBLIC_HOLIDAY_DATES]
