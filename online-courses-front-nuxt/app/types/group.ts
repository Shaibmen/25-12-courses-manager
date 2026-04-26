export type GroupScheduleItem = {
  date: string
  theme: string
}

export type GroupItem = {
  group: string
  name_group: string
  rapspisanie: GroupScheduleItem[]
}

export type GroupPayload = {
  name_group: string
  raspisanie: GroupScheduleItem[]
}

export type GroupFormScheduleItem = GroupScheduleItem & {
  id: string
}

export type GroupFormState = {
  name_group: string
  schedule: GroupFormScheduleItem[]
}

const createScheduleId = () => `${Date.now()}-${Math.random().toString(36).slice(2, 10)}`

export const createEmptyGroupScheduleItem = (
  item?: Partial<GroupScheduleItem>
): GroupFormScheduleItem => ({
  id: createScheduleId(),
  date: item?.date || '',
  theme: item?.theme || ''
})

export const createEmptyGroupFormState = (): GroupFormState => ({
  name_group: '',
  schedule: [createEmptyGroupScheduleItem()]
})
