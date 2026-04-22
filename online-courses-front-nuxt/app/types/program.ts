export type ProgramListItem = {
  id_program_education: string
  name_prof_education: string
  time_education: number
  individual_price: number
  group_price: number
  campus_price: number
  price: number
  id_education_type: string
  id_divisions_education: string
}

export type ProgramPayload = {
  name_prof_education: string
  time_education: number
  individual_price: number
  group_price: number
  campus_price: number
  id_educationtype: string
  id_divisionseducation: string
}

export type ProgramUpdatePayload = {
  name_prof_education: string
  time_education: number
  individual_price: number
  group_price: number
  campus_price: number
  ID_EducationType: string
  ID_DivisionsEducation: string
}

export type ProgramFormState = {
  name_prof_education: string
  time_education: string
  individual_price: string
  group_price: string
  campus_price: string
  id_education_type: string
  id_divisions_education: string
}

export const createEmptyProgramFormState = (): ProgramFormState => ({
  name_prof_education: '',
  time_education: '',
  individual_price: '',
  group_price: '',
  campus_price: '',
  id_education_type: '',
  id_divisions_education: ''
})
