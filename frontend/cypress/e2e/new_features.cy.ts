describe('New Features E2E', () => {
  beforeEach(() => {
    cy.visit('/')
    // Ensure English mode for consistent text matching
    cy.get('body').then(($body) => {
      if ($body.text().includes('添加')) {
         cy.get('button[title="语言"]').click()
         cy.contains('Todo App').should('be.visible')
      }
    })
  })

  it('can set remind time for a task', () => {
    const taskTitle = 'Task with Reminder ' + Date.now()
    
    // Open Add Modal
    cy.contains('button', 'Add').click()
    
    // Fill Title
    cy.get('input[placeholder="What needs to be done?"]').type(taskTitle)
    
    // Set Remind At
    // Find the label "Remind At" and get the input following it
    // Or since we know the structure:
    // Due Date is the first datetime-local, Remind At is the second
    cy.get('input[type="datetime-local"]').eq(1).then($input => {
        // Set a future date. Format: YYYY-MM-DDTHH:mm
        const now = new Date()
        now.setMinutes(now.getMinutes() + 30) // 30 mins later
        now.setMinutes(now.getMinutes() - now.getTimezoneOffset()) // Adjust for local
        const dateStr = now.toISOString().slice(0, 16)
        $input.val(dateStr)
        // Trigger input event to update v-model
        cy.wrap($input).trigger('input')
    })
    
    // Click Add
    cy.get('div[class*="fixed"]').contains('button', 'Add').click()
    
    // Verify
    cy.contains(taskTitle).should('be.visible')
    // Verify Bell icon is present in the task card
    // The bell icon is rendered when remind_at is present
    cy.contains(taskTitle).parents('.group').find('.text-indigo-600').should('exist')
  })

  it('can set repeat rule for a task', () => {
    const taskTitle = 'Repeating Task ' + Date.now()
    
    // Open Add Modal
    cy.contains('button', 'Add').click()
    
    // Fill Title
    cy.get('input[placeholder="What needs to be done?"]').type(taskTitle)
    
    // Set Repeat to "Daily"
    // Find label "Repeat"
    cy.contains('label', 'Repeat').parent().find('button').first().click()
    // Select "Daily" from dropdown
    cy.contains('button', 'Daily').click()
    
    // Click Add
    cy.get('div[class*="fixed"]').contains('button', 'Add').click()
    
    // Verify
    cy.contains(taskTitle).should('be.visible')
    // Verify Repeat icon/text is present
    // We look for the purple badge which indicates repeat
    cy.contains(taskTitle).parents('.group').find('.text-purple-600').should('contain', 'Daily')
  })
})
