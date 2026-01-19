describe('Todo App E2E', () => {
  beforeEach(() => {
    // Visit the app
    // Assuming the app runs on localhost:34115 (Vite default or configured)
    // Actually we should check where it runs. 
    // Usually Wails dev runs on random port but we can test the web part separately if running `npm run dev`
    // Let's assume we run `npm run dev` for E2E on localhost:5173 (Vite default)
    cy.visit('http://localhost:5173')
  })

  it('can add a new todo', () => {
    cy.contains('button', 'Add').click()
    cy.get('input[placeholder="What needs to be done?"]').type('New E2E Task')
    cy.contains('button', 'Add').click()
    cy.contains('New E2E Task').should('be.visible')
  })

  it('can complete a todo', () => {
    // Add a task first
    cy.contains('button', 'Add').click()
    cy.get('input[placeholder="What needs to be done?"]').type('Task to Complete')
    cy.contains('button', 'Add').click()

    // Find the task and click check
    cy.contains('Task to Complete').parent().parent().find('button').first().click()
    
    // Verify strikethrough class
    cy.contains('Task to Complete').should('have.class', 'line-through')
  })

  it('can create a project', () => {
    cy.get('aside').contains('button', '+').click()
    cy.get('input[placeholder="Name"]').type('E2E Project')
    cy.contains('button', 'Create').click()
    
    cy.get('aside').contains('E2E Project').should('be.visible')
  })

  it('can filter tasks by project', () => {
    // Create project
    cy.get('aside').contains('button', '+').click()
    cy.get('input[placeholder="Name"]').type('Filter Project')
    cy.contains('button', 'Create').click()

    // Select project
    cy.get('aside').contains('Filter Project').click()

    // Add task to project
    cy.contains('button', 'Add').click()
    cy.get('input[placeholder="What needs to be done?"]').type('Project Task')
    cy.contains('button', 'Add').click()

    // Verify task is visible
    cy.contains('Project Task').should('be.visible')

    // Switch to Inbox (should not see it)
    cy.contains('button', 'Inbox').click()
    cy.contains('Project Task').should('not.exist')
  })
})
