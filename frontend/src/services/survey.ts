import apiClient from './axios'
import { Survey, CreateSurveyRequest, UpdateSurveyRequest } from '@/types/survey'

interface SurveyListResponse {
  surveys: Survey[]
  totalCount: number
}

interface SurveyResponse {
  survey_id: string
  respondent_id?: string
  anonymous_id?: string
  answers: Array<{
    question_id: string
    value: string | string[] | number
  }>
}

// New interfaces for analytics
interface SurveyResults {
  survey_id: string
  title: string
  description: string
  total_responses: number
  question_results: QuestionResult[]
  completion_rate?: number
  average_time_spent?: string
}

interface QuestionResult {
  question_id: string
  title: string
  type: string
  response_data: any
  total_answers: number
  analytics: {
    option_counts?: Record<string, {
      option_id: string
      text: string
      count: number
      percent: number
    }>
    average?: number
    median?: number
    min?: number
    max?: number
    word_frequency?: Record<string, number>
  }
}

interface IndividualResponsesResult {
  survey_id: string
  total_count: number
  responses: SurveyResponse[]
  current_page: number
  total_pages: number
  items_per_page: number
}

const surveyService = {
  /**
   * Get all surveys for a user
   */
  async getUserSurveys(userId: string, page: number = 1, limit: number = 10): Promise<SurveyListResponse> {
    const response = await apiClient.get<SurveyListResponse>(
      `/users/${userId}/surveys?page=${page}&limit=${limit}`
    )
    return response.data
  },

  /**
   * Get a survey by ID
   */
  async getSurvey(surveyId: string): Promise<Survey> {
    const response = await apiClient.get<Survey>(`/surveys/${surveyId}`)
    return response.data
  },

  /**
   * Get a public survey by ID for taking
   */
  async getPublicSurvey(surveyId: string): Promise<Survey> {
    const response = await apiClient.get<Survey>(`/surveys/${surveyId}/public`)
    return response.data
  },

  /**
   * Submit a response to a survey
   */
  async submitSurveyResponse(surveyId: string, responseData: SurveyResponse): Promise<void> {
    await apiClient.post(`/surveys/${surveyId}/responses`, responseData)
  },

  /**
   * Create a new survey
   */
  async createSurvey(survey: CreateSurveyRequest): Promise<Survey> {
    const response = await apiClient.post<Survey>('/surveys', survey)
    return response.data
  },

  /**
   * Update an existing survey
   */
  async updateSurvey(surveyId: string, survey: UpdateSurveyRequest): Promise<Survey> {
    const response = await apiClient.put<Survey>(`/surveys/${surveyId}`, survey)
    return response.data
  },

  /**
   * Delete a survey
   */
  async deleteSurvey(surveyId: string): Promise<void> {
    await apiClient.delete(`/surveys/${surveyId}`)
  },

  /**
   * Get survey results
   */
  async getSurveyResults(surveyId: string): Promise<SurveyResults> {
    const response = await apiClient.get<SurveyResults>(`/surveys/${surveyId}/results`)
    return response.data
  },

  /**
   * Get individual responses for a survey
   */
  async getIndividualResponses(
    surveyId: string, 
    page: number = 1, 
    limit: number = 10, 
    startDate?: string, 
    endDate?: string
  ): Promise<IndividualResponsesResult> {
    let url = `/surveys/${surveyId}/responses?page=${page}&limit=${limit}`
    
    if (startDate) {
      url += `&startDate=${encodeURIComponent(startDate)}`
    }
    
    if (endDate) {
      url += `&endDate=${encodeURIComponent(endDate)}`
    }
    
    const response = await apiClient.get<IndividualResponsesResult>(url)
    return response.data
  }
}

export default surveyService 