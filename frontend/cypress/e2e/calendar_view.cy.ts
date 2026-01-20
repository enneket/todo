describe('Calendar View E2E', () => {
  beforeEach(() => {
    cy.visit('/')
    
    // Ensure English
    cy.get('body').then(($body) => {
      if ($body.text().includes('添加')) {
         cy.get('button[title="语言"]').click()
      }
    })
    cy.contains('Todo App').should('be.visible')
  })

  it('can navigate to calendar view', () => {
    // Click Calendar in sidebar
    cy.contains('button', 'Calendar').click()
    
    // Verify Calendar view is active (Month view by default)
    cy.contains('button', 'Month').should('have.class', 'bg-white')
    cy.contains('button', 'Week').should('not.have.class', 'bg-white')
    
    // Check for days of week header
    cy.contains('Sun').should('be.visible')
    cy.contains('Mon').should('be.visible')
  })

  it('displays tasks in calendar', () => {
    // Add a task for today
    const taskTitle = 'Calendar Task ' + Date.now()
    
    cy.contains('button', 'Add').click()
    cy.get('input[placeholder="What needs to be done?"]').type(taskTitle)
    
    // Set date to today (default is tomorrow)
    // We need to input current date-time into the datetime-local input
    // Format: YYYY-MM-DDThh:mm
    const now = new Date()
    // Adjust to local ISO string
    const offset = now.getTimezoneOffset() * 60000
    const localISOTime = (new Date(now.getTime() - offset)).toISOString().slice(0, 16)
    
    // Find the due date input. It is the first datetime-local input
    cy.get('input[type="datetime-local"]').first().type(localISOTime)
    
    cy.get('div[class*="fixed"]').contains('button', 'Add').click()

    // Navigate to Calendar
    cy.contains('button', 'Calendar').click()
    
    // Verify task is visible in calendar
    cy.contains(taskTitle).should('be.visible')
  })

  it('can switch between month and week views', () => {
    cy.contains('button', 'Calendar').click()
    
    // Switch to Week view
    cy.contains('button', 'Week').click()
    
    // Verify Week view is active
    cy.contains('button', 'Week').should('have.class', 'bg-white')
    
    // Grid should have 7 columns (days) and multiple rows, but the implementation is
    // v-for="(day, index) in calendarDays" inside a grid container.
    // In Week view, calendarDays has length 7.
    // We can check if we see 7 day cells.
    // The day cells have class 'border-b border-r' etc.
    // Let's count them.
    cy.get('.grid-rows-6 > div').should('have.length', 7)
    
    // Switch back to Month view
    cy.contains('button', 'Month').click()
    // Month view has 42 cells
    cy.get('.grid-rows-6 > div').should('have.length', 42)
  })
})
