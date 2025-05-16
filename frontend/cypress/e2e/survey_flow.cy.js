describe('Survey Creation and Taking Flow', () => {
  // Назва опитування для тестування
  const surveyTitle = `Test Survey ${Math.floor(Math.random() * 10000)}`
  let createdSurveyId;
  
  beforeEach(() => {
    // Використовуємо кастомну команду login для авторизації
    // Для створення опитування потрібен авторизований користувач
    cy.login(Cypress.env('regularUser').email, Cypress.env('regularUser').password)
  })
  
  it('should allow user to create a survey', () => {
    // Переходимо на сторінку створення опитування
    cy.visit('/surveys/create')
    
    // Заповнюємо форму створення опитування
    cy.get('[data-cy=input-survey-title]').type(surveyTitle)
    cy.get('[data-cy=input-survey-description]').type('This is a test survey description')
    
    // Додаємо перше питання (одиночний вибір)
    cy.get('[data-cy=btn-add-question]').click()
    cy.get('[data-cy=input-question-text]').last().type('What is your favorite color?')
    cy.get('[data-cy=select-question-type]').last().select('single-choice')
    cy.get('[data-cy=checkbox-required]').last().check()
    
    // Додаємо варіанти відповідей
    cy.get('[data-cy=btn-add-option]').click()
    cy.get('[data-cy=input-option-text]').eq(0).type('Red')
    cy.get('[data-cy=btn-add-option]').click()
    cy.get('[data-cy=input-option-text]').eq(1).type('Blue')
    cy.get('[data-cy=btn-add-option]').click()
    cy.get('[data-cy=input-option-text]').eq(2).type('Green')
    
    // Додаємо друге питання (текстове)
    cy.get('[data-cy=btn-add-question]').click()
    cy.get('[data-cy=input-question-text]').last().type('Please explain your choice')
    cy.get('[data-cy=select-question-type]').last().select('open-text')
    
    // Зберігаємо опитування
    cy.get('[data-cy=btn-save-survey]').click()
    
    // Перевіряємо перенаправлення на сторінку деталей опитування та успішне повідомлення
    cy.url().should('include', '/surveys/details/')
    cy.get('[data-cy=alert-success]').should('be.visible')
      .and('contain', 'Survey created successfully')
    
    // Зберігаємо ID створеного опитування для подальших тестів
    cy.url().then(url => {
      createdSurveyId = url.split('/').pop()
      cy.wrap(createdSurveyId).as('surveyId')
    })
  })
  
  it('should display created survey in the list', () => {
    // Переходимо на сторінку списку опитувань
    cy.visit('/surveys')
    
    // Перевіряємо, що наше опитування відображається у списку
    cy.get('[data-cy=survey-list]').should('be.visible')
    cy.get(`[data-cy=survey-title]:contains("${surveyTitle}")`).should('be.visible')
  })
  
  it('should allow user to edit survey', function() {
    // Використовуємо збережений ID опитування
    // Примітка: 'this.surveyId' доступний завдяки збереженню через cy.wrap().as()
    cy.visit(`/surveys/edit/${this.surveyId}`)
    
    // Змінюємо заголовок опитування
    const updatedTitle = `${surveyTitle} (Updated)`
    cy.get('[data-cy=input-survey-title]').clear().type(updatedTitle)
    
    // Додаємо ще одне питання
    cy.get('[data-cy=btn-add-question]').click()
    cy.get('[data-cy=input-question-text]').last().type('How satisfied are you with this product?')
    cy.get('[data-cy=select-question-type]').last().select('scale')
    
    // Налаштовуємо шкалу від 1 до 5
    cy.get('[data-cy=input-scale-min]').type('1')
    cy.get('[data-cy=input-scale-max]').type('5')
    cy.get('[data-cy=input-scale-min-label]').type('Not satisfied at all')
    cy.get('[data-cy=input-scale-max-label]').type('Very satisfied')
    
    // Зберігаємо зміни
    cy.get('[data-cy=btn-save-survey]').click()
    
    // Перевіряємо успішне оновлення
    cy.get('[data-cy=alert-success]').should('be.visible')
      .and('contain', 'Survey updated successfully')
  })
  
  it('should allow user to take the survey', function() {
    // Переходимо на сторінку проходження опитування
    cy.visit(`/surveys/take/${this.surveyId}`)
    
    // Перевіряємо, що заголовок опитування відображається
    cy.get('[data-cy=survey-title]').should('contain', surveyTitle)
    
    // Заповнюємо відповіді
    // Питання 1: одиночний вибір
    cy.get('[data-cy=question]').eq(0).within(() => {
      cy.get('[data-cy=option-radio]').eq(1).check() // Вибираємо "Blue"
    })
    
    // Питання 2: відкритий текст
    cy.get('[data-cy=question]').eq(1).within(() => {
      cy.get('[data-cy=text-answer]').type('I chose blue because it is my favorite color.')
    })
    
    // Питання 3: шкала (додане у попередньому тесті)
    cy.get('[data-cy=question]').eq(2).within(() => {
      cy.get('[data-cy=scale-option]').eq(3).check() // Вибираємо "4" із 5
    })
    
    // Відправляємо відповіді
    cy.get('[data-cy=btn-submit-response]').click()
    
    // Перевіряємо успішне завершення опитування
    cy.get('[data-cy=success-page]').should('be.visible')
    cy.get('[data-cy=success-message]').should('contain', 'Thank you for completing the survey')
  })
  
  it('should show survey responses in analytics', function() {
    // Переходимо на сторінку аналітики опитування 
    // (доступна для автора опитування)
    cy.visit(`/surveys/analytics/${this.surveyId}`)
    
    // Перевіряємо, що основна інформація відображається
    cy.get('[data-cy=survey-title]').should('contain', surveyTitle)
    cy.get('[data-cy=response-count]').should('contain', '1') // Одна відповідь від попереднього тесту
    
    // Перевіряємо, що графіки/аналітика для питань відображаються
    cy.get('[data-cy=question-analytics]').should('have.length.at.least', 3)
    
    // Перевіряємо дані для питання з одиночним вибором
    cy.get('[data-cy=question-analytics]').eq(0).within(() => {
      cy.get('[data-cy=chart]').should('be.visible')
      // Перевірка, що "Blue" має 100% відповідей (1 з 1)
      cy.get('[data-cy=option-percentage]').contains('Blue').should('contain', '100%')
    })
  })
}) 