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
    
    // Type into input and press enter
    cy.get('input[type="text"]').type(`${todoText}{enter}`)

    // Verify it appears in the list (scroll if needed)
    cy.contains(todoText).scrollIntoView().should('be.visible')
  })

  it('should toggle a todo status', () => {
    const todoText = 'Toggle Test ' + Date.now()
    
    // Add todo
    cy.get('input[type="text"]').type(`${todoText}{enter}`)
    
    // Find the todo item and click the check circle (button)
    cy.contains(todoText)
      .scrollIntoView()
      .parent()
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
    cy.get('input[type="text"]').type(`${activeText}{enter}`)
    cy.get('input[type="text"]').type(`${completedText}{enter}`)

    // Complete the second one
    cy.contains(completedText)
      .scrollIntoView()
      .parent()
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
    cy.get('input[type="text"]').type(`${todoText}{enter}`)
    
    // Hover over the todo item to make delete button visible (if using group-hover)
    // Cypress can force click hidden elements
    cy.contains(todoText)
      .scrollIntoView()
      .parent()
      .find('button')
      .last() // Assuming delete is the last button
      .click({ force: true })

    // Verify it's gone
    cy.contains(todoText).should('not.exist')
  })
})
