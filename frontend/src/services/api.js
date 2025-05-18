import api from '../utils/axiosConfig';

// Auth endpoints
export const authApi = {
  register: (userData) => api.post('/api/v1/auth/register', userData),
  login: (credentials) => api.post('/api/v1/auth/login', credentials),
  refreshToken: (refreshToken) => api.post('/api/v1/auth/refresh', { refresh_token: refreshToken }),
  getCurrentUser: () => api.get('/api/v1/users/me'),
  updateCurrentUser: (userData) => api.put('/api/v1/users/me', userData)
}

// Survey endpoints
export const surveyApi = {
  getSurveys: () => api.get('/api/v1/surveys'),
  getSurvey: (id) => api.get(`/api/v1/surveys/${id}`),
  createSurvey: (surveyData) => api.post('/api/v1/surveys', surveyData),
  updateSurvey: (id, surveyData) => api.put(`/api/v1/surveys/${id}`, surveyData),
  deleteSurvey: (id) => api.delete(`/api/v1/surveys/${id}`),
  updateSurveyStatus: (id, isActive) => api.patch(`/api/v1/surveys/${id}/status`, { is_active: isActive }),
  
  // Question related endpoints
  addQuestion: (surveyId, questionData) => api.post(`/api/v1/surveys/${surveyId}/questions`, questionData),
  updateQuestion: (questionId, questionData) => api.put(`/api/v1/questions/${questionId}`, questionData),
  deleteQuestion: (questionId) => api.delete(`/api/v1/questions/${questionId}`)
}

// Response endpoints
export const responseApi = {
  submitResponse: (responseData) => api.post('/api/v1/responses', responseData),
  getResponses: (surveyId) => api.get('/api/v1/responses', { params: { survey_id: surveyId } }),
  getResponseSummary: (surveyId) => api.get('/api/v1/responses/summary', { params: { survey_id: surveyId } })
}

export default api 