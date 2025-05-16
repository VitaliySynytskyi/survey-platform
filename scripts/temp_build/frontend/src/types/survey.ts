export enum QuestionType {
  SingleChoice = 'single-choice',
  MultipleChoice = 'multiple-choice',
  OpenText = 'open-text',
  Scale = 'scale',
  MatrixSingle = 'matrix-single',
  MatrixMultiple = 'matrix-multiple'
}

export interface Option {
  value: string
  text: string
}

export interface ScaleSettings {
  min: number
  max: number
  minLabel: string
  maxLabel: string
}

export interface DisplayLogic {
  dependsOnQuestionId: string
  dependsOnAnswerValue: string
}

export interface Question {
  question_id: string
  text: string
  type: QuestionType
  is_required: boolean
  options?: Option[]
  scale_settings?: ScaleSettings
  matrix_rows?: string[]
  matrix_columns?: string[]
  display_logic?: DisplayLogic
}

export interface Survey {
  id: string
  title: string
  description?: string
  owner_id: string
  created_at: string
  updated_at: string
  questions: Question[]
}

export interface CreateSurveyRequest {
  title: string
  description?: string
  questions: Question[]
}

export interface UpdateSurveyRequest {
  title?: string
  description?: string
  questions?: Question[]
} 