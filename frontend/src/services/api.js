import api from '../utils/axiosConfig';

// Helper function to get auth headers dynamically
const getAuthHeaders = () => {
  const token = localStorage.getItem('token'); // Directly get token from localStorage
  if (token) {
    return { Authorization: `Bearer ${token}` };
  }
  return {};
};

// Auth endpoints
export const authApi = {
  register: (userData) => api.post('/api/v1/auth/register', userData),
  login: (credentials) => api.post('/api/v1/auth/login', credentials),
  refreshToken: (refreshToken) => api.post('/api/v1/auth/refresh', { refresh_token: refreshToken }),
  getCurrentUser: () => api.get('/api/v1/users/me', { headers: getAuthHeaders() }),
  updateCurrentUser: (userData) => api.put('/api/v1/users/me', userData, { headers: getAuthHeaders() })
}

// Survey endpoints
export const surveyApi = {
  getSurveys: () => api.get('/api/v1/surveys', { headers: getAuthHeaders() }),
  getSurvey: (id) => api.get(`/api/v1/surveys/${id}`, { headers: getAuthHeaders() }),
  createSurvey: (surveyData) => api.post('/api/v1/surveys', surveyData, { headers: getAuthHeaders() }),
  updateSurvey: (id, surveyData) => api.put(`/api/v1/surveys/${id}`, surveyData, { headers: getAuthHeaders() }),
  deleteSurvey: (id) => api.delete(`/api/v1/surveys/${id}`, { headers: getAuthHeaders() }),
  updateSurveyStatus: (id, isActive) => api.patch(`/api/v1/surveys/${id}/status`, { is_active: isActive }, { headers: getAuthHeaders() }),
  
  // Question related endpoints
  addQuestion: (surveyId, questionData) => api.post(`/api/v1/surveys/${surveyId}/questions`, questionData, { headers: getAuthHeaders() }),
  updateQuestion: (questionId, questionData) => api.put(`/api/v1/questions/${questionId}`, questionData, { headers: getAuthHeaders() }),
  deleteQuestion: (questionId) => api.delete(`/api/v1/questions/${questionId}`, { headers: getAuthHeaders() })
}

// Response endpoints
export const responseApi = {
  submitResponse: (responseData) => api.post('/api/v1/responses', responseData, { headers: getAuthHeaders() }),
  getResponses: (surveyId) => api.get('/api/v1/responses', { params: { survey_id: surveyId }, headers: getAuthHeaders() }),
  getResponseSummary: (surveyId) => api.get('/api/v1/responses/summary', { params: { survey_id: surveyId }, headers: getAuthHeaders() })
}

export default api 