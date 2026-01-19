describe('Todo App', () => {
  beforeEach(() => {
    // Clear todos before each test (optional, but good for isolation)
    // For now we assume a fresh start or just test the flow
    cy.visit('/')
    // Ensure we are in English mode for testing
    // Check if the title is NOT "Todo App" (which means it's likely "待办事项")
    cy.get('h1').then(($h1) => {
      if (!$h1.text().includes('Todo App')) {
        // Click language toggle button to switch to English
        cy.get('button[title="语言"]').click()
        // Wait for the title to change to ensure app is ready
        cy.contains('Todo App').should('be.visible')
      }
    })
  })

  it('should display the header', () => {
    cy.get('h1').should('contain', 'Todo App')
  })

  it('should add a new todo', () => {
    const todoText = 'New Test Todo ' + Date.now()
    
    // Open modal
    cy.contains('button', 'Add').click()
    // Type into title
    cy.get('input[placeholder="What needs to be done?"]').type(todoText)
    // Click Add in modal
    cy.get('.fixed button').contains('Add').click()

    // Verify it appears in the list (scroll if needed)
    cy.contains(todoText).scrollIntoView().should('be.visible')
  })

  it('should toggle a todo status', () => {
    const todoText = 'Toggle Test ' + Date.now()
    
    // Add todo
    cy.contains('button', 'Add').click()
    cy.get('input[placeholder="What needs to be done?"]').type(todoText)
    cy.get('.fixed button').contains('Add').click()
    
    // Find the todo item and click the check circle (button)
    cy.contains(todoText)
      .scrollIntoView()
      .parents('.group')
      .find('button')
      .first()
      .click()

    // Verify line-through style (indicating completion)
    cy.contains(todoText).scrollIntoView().should('have.class', 'line-through')
  })

  it('should filter todos', () => {
    const activeText = 'Active Task ' + Date.now()
    const completedText = 'Completed Task ' + Date.now()

    // Add two tasks
    cy.contains('button', 'Add').click()
    cy.get('input[placeholder="What needs to be done?"]').type(activeText)
    cy.get('.fixed button').contains('Add').click()

    cy.contains('button', 'Add').click()
    cy.get('input[placeholder="What needs to be done?"]').type(completedText)
    cy.get('.fixed button').contains('Add').click()

    // Complete the second one
    cy.contains(completedText)
      .scrollIntoView()
      .parents('.group')
      .find('button')
      .first()
      .click()

    // Wait for the completion to be reflected in UI
    cy.contains(completedText).should('have.class', 'line-through')

    // Click "Active" filter
    cy.contains('button', 'Active').click()
    cy.contains(activeText).scrollIntoView().should('be.visible')
    cy.contains(completedText).should('not.exist')

    // Click "Completed" filter
    cy.contains('button', 'Completed').click()
    cy.contains(completedText).scrollIntoView().should('be.visible')
    cy.contains(activeText).should('not.exist')

    // Click "All" filter
    cy.contains('button', 'All').click()
    cy.contains(activeText).scrollIntoView().should('be.visible')
    cy.contains(completedText).scrollIntoView().should('be.visible')
  })

  it('should delete a todo', () => {
    const todoText = 'Delete Test ' + Date.now()
    
    // Add todo
    cy.contains('button', 'Add').click()
    cy.get('input[placeholder="What needs to be done?"]').type(todoText)
    cy.get('.fixed button').contains('Add').click()
    
    // Hover over the todo item to make delete button visible (if using group-hover)
    // Cypress can force click hidden elements
    cy.contains(todoText)
      .scrollIntoView()
      .parents('.group')
      .find('button')
      .last() // Assuming delete is the last button
      .click({ force: true })

    // Verify it's gone
    cy.contains(todoText).should('not.exist')
  })
})
