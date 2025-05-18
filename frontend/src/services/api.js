import axios from 'axios'

// Create a configured axios instance
const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
})

// Auth endpoints
export const authApi = {
  register: (userData) => api.post('/v1/auth/register', userData),
  login: (credentials) => api.post('/v1/auth/login', credentials),
  refreshToken: (refreshToken) => api.post('/v1/auth/refresh', { refresh_token: refreshToken }),
  getCurrentUser: () => api.get('/v1/users/me'),
  updateCurrentUser: (userData) => api.put('/v1/users/me', userData)
}

// Survey endpoints
export const surveyApi = {
  getSurveys: () => api.get('/v1/surveys'),
  getSurvey: (id) => api.get(`/v1/surveys/${id}`),
  createSurvey: (surveyData) => api.post('/v1/surveys', surveyData),
  updateSurvey: (id, surveyData) => api.put(`/v1/surveys/${id}`, surveyData),
  deleteSurvey: (id) => api.delete(`/v1/surveys/${id}`),
  
  // Question related endpoints
  addQuestion: (surveyId, questionData) => api.post(`/v1/surveys/${surveyId}/questions`, questionData),
  updateQuestion: (questionId, questionData) => api.put(`/v1/questions/${questionId}`, questionData),
  deleteQuestion: (questionId) => api.delete(`/v1/questions/${questionId}`)
}

// Response endpoints
export const responseApi = {
  submitResponse: (responseData) => api.post('/v1/responses', responseData),
  getResponses: (surveyId) => api.get('/v1/responses', { params: { survey_id: surveyId } }),
  getResponseSummary: (surveyId) => api.get('/v1/responses/summary', { params: { survey_id: surveyId } })
}

export default api 