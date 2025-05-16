const { defineConfig } = require('cypress')

module.exports = defineConfig({
  e2e: {
    baseUrl: 'http://localhost:8080',
    specPattern: 'cypress/e2e/**/*.cy.{js,jsx,ts,tsx}',
    viewportWidth: 1280,
    viewportHeight: 720,
    defaultCommandTimeout: 10000,
    video: false,
    screenshotOnRunFailure: true,
    setupNodeEvents(on, config) {
      // Тут можна додати налаштування плагінів або власних подій
      return config
    },
  },
  
  env: {
    // Змінні середовища для тестів
    apiUrl: 'http://localhost:3000/api/v1',
    adminUser: {
      email: 'admin@example.com',
      password: 'admin_password',
    },
    regularUser: {
      email: 'user@example.com',
      password: 'user_password',
    },
  },
}) 