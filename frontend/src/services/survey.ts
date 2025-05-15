import apiClient from './axios'
import { Survey, CreateSurveyRequest, UpdateSurveyRequest } from '@/types/survey'

interface SurveyListResponse {
  surveys: Survey[]
  totalCount: number
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
  }
}

export default surveyService 