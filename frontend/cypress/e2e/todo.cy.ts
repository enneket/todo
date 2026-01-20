describe('Todo App E2E', () => {
  beforeEach(() => {
    // Visit the app
    cy.visit('/')
    
    // Ensure we are in English mode for testing
    // Check if the title is Chinese "待办事项" or English "Todo App"
    // The button title changes too: "语言" (zh) or "Language" (en)
    cy.get('body').then(($body) => {
      // If we see "添加" (Add in zh), we need to switch
      if ($body.text().includes('添加')) {
         // Find language button. 
         // In zh mode, title is "语言". In en mode, title is "Language".
         // We can find by the translate icon if we can target it, or just use the title "语言"
         cy.get('button[title="语言"]').click()
         cy.contains('Todo App').should('be.visible')
      }
    })
  })

  it('can add a new todo', () => {
    // Open Add Modal
    cy.contains('button', 'Add').click()
    
    // Type title
    cy.get('input[placeholder="What needs to be done?"]').type('New E2E Task')
    
    // Click Add
    // The button in modal says "Add"
    cy.get('div[class*="fixed"]').contains('button', 'Add').click()
    
    // Verify
    cy.contains('New E2E Task').should('be.visible')
  })

  it('can complete a todo', () => {
    // Add a task first
    cy.contains('button', 'Add').click()
    cy.get('input[placeholder="What needs to be done?"]').type('Task to Complete')
    cy.get('div[class*="fixed"]').contains('button', 'Add').click()

    // Find the task and click check
    // The check circle is inside a button
    cy.contains('Task to Complete').parent().parent().find('button').first().click()
    
    // Verify strikethrough class (line-through)
    cy.contains('Task to Complete').should('have.class', 'line-through')
  })

  it('can create a project', () => {
    // Click + button for projects
    cy.contains('h3', 'Projects').parent().find('button').click()
    
    // Modal opens. Input has no placeholder, but is inside the modal.
    // The modal title is "Add Project"
    cy.contains('h2', 'Add Project').parent().parent().find('input[type="text"]').type('E2E Project')
    
    cy.contains('button', 'Create').click()
    
    // Verify project appears in sidebar
    cy.get('aside').contains('E2E Project').should('be.visible')
  })

  it('can filter tasks by project', () => {
    // Create project
    cy.contains('h3', 'Projects').parent().find('button').click()
    
    const projName = 'Filter Project ' + Date.now()
    cy.contains('h2', 'Add Project').parent().parent().find('input[type="text"]').type(projName)
    
    cy.contains('button', 'Create').click()

    // Select project (Click on it in sidebar)
    cy.get('aside').contains(projName).click()

    // Add task to project
    cy.contains('button', 'Add').click()
    cy.get('input[placeholder="What needs to be done?"]').type('Project Task')
    cy.get('div[class*="fixed"]').contains('button', 'Add').click()

    // Verify task is visible
    cy.contains('Project Task').should('be.visible')

    // Switch to Inbox (should not see it)
    cy.contains('button', 'Inbox').click()
    cy.contains('Project Task').should('not.exist')
    
    // Switch back to All Tasks (should see it)
    cy.contains('button', 'All Tasks').click()
    cy.contains('Project Task').should('be.visible')
  })
})
