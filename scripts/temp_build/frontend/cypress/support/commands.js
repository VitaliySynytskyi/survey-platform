// frontend/cypress/support/commands.js

// Кастомна команда для логіну
Cypress.Commands.add('login', (email, password) => {
  // Перевіряємо, чи вже авторизований (щоб уникнути зайвих запитів)
  cy.window().its('localStorage').then((localStorage) => {
    // Якщо у localStorage вже є auth_token, не потрібно виконувати логін
    if (localStorage.getItem('auth_token')) {
      return;
    }

    // Виконуємо логін через API для швидкості (замість UI)
    cy.request({
      method: 'POST',
      url: `${Cypress.env('apiUrl')}/auth/login`,
      body: {
        email: email,
        password: password,
      },
    }).then((response) => {
      // Перевіряємо успішність запиту
      expect(response.status).to.eq(200);
      expect(response.body).to.have.property('access_token');
      
      // Зберігаємо токени в localStorage
      localStorage.setItem('auth_token', response.body.access_token);
      localStorage.setItem('refresh_token', response.body.refresh_token);
      
      // Зберігаємо інформацію про користувача
      if (response.body.user) {
        localStorage.setItem('user', JSON.stringify(response.body.user));
      }
    });
  });
});

// Кастомна команда для створення тестового опитування через API
Cypress.Commands.add('createTestSurvey', (title, description) => {
  // Формуємо базовий шаблон опитування
  const surveyData = {
    title: title || `Test Survey ${Math.floor(Math.random() * 10000)}`,
    description: description || 'This is an automated test survey',
    questions: [
      {
        text: 'What is your favorite color?',
        type: 'single-choice',
        is_required: true,
        options: [
          { text: 'Red' },
          { text: 'Blue' },
          { text: 'Green' }
        ]
      },
      {
        text: 'Please provide any additional feedback',
        type: 'open-text',
        is_required: false
      }
    ]
  };

  // Надсилаємо запит на створення опитування
  cy.request({
    method: 'POST',
    url: `${Cypress.env('apiUrl')}/surveys`,
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
    },
    body: surveyData
  }).then((response) => {
    // Перевіряємо успішність запиту
    expect(response.status).to.eq(201);
    // Повертаємо ID створеного опитування
    return response.body.id;
  });
});

// Кастомна команда для видалення тестового опитування через API
Cypress.Commands.add('deleteTestSurvey', (surveyId) => {
  cy.request({
    method: 'DELETE',
    url: `${Cypress.env('apiUrl')}/surveys/${surveyId}`,
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
    }
  }).then((response) => {
    expect(response.status).to.eq(204);
  });
}); 