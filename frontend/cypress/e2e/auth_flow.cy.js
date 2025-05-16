describe('Authentication Flow', () => {
  beforeEach(() => {
    // Перед кожним тестом очищаємо localStorage/cookies
    cy.clearLocalStorage()
    cy.clearCookies()
  })

  it('should allow a user to sign up', () => {
    // Випадковий email для запобігання конфліктам під час повторних запусків
    const randomEmail = `test${Math.floor(Math.random() * 10000)}@example.com`
    const password = 'Password123!'

    cy.visit('/signup')
    
    // Перевіряємо, що форма реєстрації існує
    cy.get('[data-cy=signup-form]').should('be.visible')
    
    // Заповнюємо форму
    cy.get('[data-cy=input-username]').type('testuser')
    cy.get('[data-cy=input-email]').type(randomEmail)
    cy.get('[data-cy=input-password]').type(password)
    cy.get('[data-cy=input-password-confirm]').type(password)
    
    // Підтверджуємо правила та умови
    cy.get('[data-cy=checkbox-terms]').check()
    
    // Відправляємо форму
    cy.get('[data-cy=btn-signup]').click()
    
    // Очікуємо, що користувача буде перенаправлено на сторінку входу
    // з повідомленням про успішну реєстрацію
    cy.url().should('include', '/login')
    cy.get('[data-cy=alert-success]').should('be.visible')
      .and('contain', 'Registration successful')
  })

  it('should show validation errors for invalid input', () => {
    cy.visit('/signup')
    
    // Введення недійсних даних
    cy.get('[data-cy=input-email]').type('invalid-email')
    cy.get('[data-cy=input-password]').type('short')
    
    // Спроба відправити форму
    cy.get('[data-cy=btn-signup]').click()
    
    // Очікуємо побачити повідомлення про помилки валідації
    cy.get('[data-cy=error-email]').should('be.visible')
      .and('contain', 'valid email')
    cy.get('[data-cy=error-password]').should('be.visible')
      .and('contain', 'at least 8 characters')
  })

  it('should allow a user to log in and log out', () => {
    // Використовуємо тестового користувача з налаштувань
    const { email, password } = Cypress.env('regularUser')
    
    cy.visit('/login')
    
    // Перевіряємо, що форма логіну існує
    cy.get('[data-cy=login-form]').should('be.visible')
    
    // Заповнюємо форму
    cy.get('[data-cy=input-email]').type(email)
    cy.get('[data-cy=input-password]').type(password)
    
    // Відправляємо форму
    cy.get('[data-cy=btn-login]').click()
    
    // Перевіряємо, що користувача перенаправлено на dashboard
    cy.url().should('include', '/dashboard')
    
    // Перевіряємо, що меню для авторизованого користувача доступне
    cy.get('[data-cy=user-menu]').should('be.visible')
    cy.get('[data-cy=user-menu]').click()
    
    // Натискаємо на "Вийти"
    cy.get('[data-cy=btn-logout]').click()
    
    // Перевіряємо, що користувач розлогінився і перенаправлений на головну сторінку
    cy.url().should('eq', Cypress.config().baseUrl + '/')
    cy.get('[data-cy=nav-login]').should('be.visible')
  })

  it('should show error for invalid credentials', () => {
    cy.visit('/login')
    
    // Вводимо неправильні дані
    cy.get('[data-cy=input-email]').type('wrong@example.com')
    cy.get('[data-cy=input-password]').type('wrongpassword')
    
    // Відправляємо форму
    cy.get('[data-cy=btn-login]').click()
    
    // Перевіряємо, що відображається повідомлення про помилку
    cy.get('[data-cy=alert-error]').should('be.visible')
      .and('contain', 'Invalid email or password')
    
    // Перевіряємо, що ми залишились на сторінці логіну
    cy.url().should('include', '/login')
  })
}) 