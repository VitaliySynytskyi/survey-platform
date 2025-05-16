import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

// Налаштування тесту
export const options = {
  stages: [
    { duration: '30s', target: 10 }, // Поступове збільшення до 10 віртуальних користувачів
    { duration: '1m', target: 10 },  // Утримання 10 віртуальних користувачів протягом 1 хвилини
    { duration: '30s', target: 0 },  // Поступове зменшення до 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% запитів повинні завершуватися менше ніж за 500 мс
    http_req_failed: ['rate<0.05'],   // Менше 5% запитів можуть завершитися невдало
  },
};

// Базова URL API Gateway
const API_URL = 'http://localhost:8080/api/v1';

// Створення унікальних користувачів для тестування
function generateRandomUser() {
  const randomPart = randomString(8);
  return {
    email: `loadtest-${randomPart}@example.com`,
    password: 'Test@123456',
    name: `Load Test User ${randomPart}`,
  };
}

// Тестовий сценарій для автентифікації
export default function () {
  // 1. Реєстрація нового користувача
  const user = generateRandomUser();
  const registerRes = http.post(`${API_URL}/auth/register`, JSON.stringify(user), {
    headers: { 'Content-Type': 'application/json' },
  });

  check(registerRes, {
    'register status is 201': (r) => r.status === 201,
    'register response contains user_id': (r) => JSON.parse(r.body).user_id !== undefined,
  });

  sleep(1);

  // 2. Логін користувача
  const loginPayload = {
    email: user.email,
    password: user.password,
  };

  const loginRes = http.post(`${API_URL}/auth/login`, JSON.stringify(loginPayload), {
    headers: { 'Content-Type': 'application/json' },
  });

  check(loginRes, {
    'login status is 200': (r) => r.status === 200,
    'login response contains tokens': (r) => {
      const body = JSON.parse(r.body);
      return body.access_token !== undefined && body.refresh_token !== undefined;
    },
  });

  // Отримання токену автентифікації
  let accessToken;
  try {
    accessToken = JSON.parse(loginRes.body).access_token;
  } catch (e) {
    console.error('Failed to parse login response', e);
    return;
  }

  sleep(1);

  // 3. Запит інформації про користувача (me)
  const meRes = http.get(`${API_URL}/auth/me`, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${accessToken}`,
    },
  });

  check(meRes, {
    'me status is 200': (r) => r.status === 200,
    'me response contains user info': (r) => {
      const body = JSON.parse(r.body);
      return body.email === user.email && body.name === user.name;
    },
  });

  sleep(1);
} 