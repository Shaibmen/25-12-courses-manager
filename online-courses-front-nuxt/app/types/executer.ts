export type ExecuterItem = {
  id_executor: string
  first_name: string
  second_name: string
  middle_name?: string
  status: string
  doverenost?: string
}

export type ExecuterPayload = {
  first_name: string
  second_name: string
  middle_name?: string
  status: string
  doverenost?: string
}

export type ExecuterFormState = {
  first_name: string
  second_name: string
  middle_name: string
  status: string
  doverenost: string
}

export const createEmptyExecuterFormState = (): ExecuterFormState => ({
  first_name: '',
  second_name: '',
  middle_name: '',
  status: '',
  doverenost: ''
})
