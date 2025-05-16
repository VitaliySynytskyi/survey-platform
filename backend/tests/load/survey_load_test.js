import http from 'k6/http';
import { check, sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

// Налаштування тесту
export const options = {
  stages: [
    { duration: '30s', target: 5 },  // Поступове збільшення до 5 віртуальних користувачів
    { duration: '1m', target: 5 },   // Утримання 5 віртуальних користувачів протягом 1 хвилини
    { duration: '30s', target: 0 },  // Поступове зменшення до 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<800'], // 95% запитів повинні завершуватися менше ніж за 800 мс
    http_req_failed: ['rate<0.05'],   // Менше 5% запитів можуть завершитися невдало
  },
};

// Базова URL API Gateway
const API_URL = 'http://localhost:8080/api/v1';

// Логін вже існуючого користувача
function login() {
  const credentials = {
    email: 'admin@example.com',
    password: 'Admin@123456',
  };

  const loginRes = http.post(`${API_URL}/auth/login`, JSON.stringify(credentials), {
    headers: { 'Content-Type': 'application/json' },
  });

  if (loginRes.status === 200) {
    try {
      return JSON.parse(loginRes.body).access_token;
    } catch (e) {
      console.error('Failed to parse login response', e);
    }
  }
  return null;
}

// Створення нового опитування
function createSurvey(token) {
  const uniqueID = uuidv4().substring(0, 8);
  const survey = {
    title: `Load Test Survey ${uniqueID}`,
    description: 'This is a survey created for load testing',
    questions: [
      {
        title: 'How would you rate our service?',
        type: 'scale',
        required: true,
        options: ['1', '2', '3', '4', '5'],
      },
      {
        title: 'What could we improve?',
        type: 'open-text',
        required: false,
      },
    ],
  };

  const createRes = http.post(`${API_URL}/surveys`, JSON.stringify(survey), {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
  });

  check(createRes, {
    'create survey status is 201': (r) => r.status === 201,
    'create survey returns survey details': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.id !== undefined && body.title === survey.title;
      } catch (e) {
        return false;
      }
    },
  });

  try {
    return JSON.parse(createRes.body).id;
  } catch (e) {
    console.error('Failed to parse create survey response', e);
    return null;
  }
}

// Тестовий сценарій для опитувань
export default function () {
  // 1. Логін
  const token = login();
  if (!token) {
    console.error('Login failed, skipping test');
    return;
  }

  sleep(1);

  // 2. Створення опитування
  const surveyId = createSurvey(token);
  if (!surveyId) {
    console.error('Creating survey failed, skipping test');
    return;
  }

  sleep(1);

  // 3. Отримання деталей опитування
  const surveyRes = http.get(`${API_URL}/surveys/${surveyId}`, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
  });

  check(surveyRes, {
    'get survey status is 200': (r) => r.status === 200,
    'get survey returns correct survey': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.id === surveyId;
      } catch (e) {
        return false;
      }
    },
  });

  sleep(1);

  // 4. Отримання всіх опитувань користувача
  const userSurveysRes = http.get(`${API_URL}/users/me/surveys`, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
  });

  check(userSurveysRes, {
    'get user surveys status is 200': (r) => r.status === 200,
  });

  sleep(1);
} 