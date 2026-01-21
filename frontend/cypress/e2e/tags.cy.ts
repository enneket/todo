describe('Tag System E2E', () => {
  beforeEach(() => {
    cy.visit('/')
    
    // Ensure English mode
    cy.get('body').then(($body) => {
      if ($body.text().includes('添加')) {
         cy.get('button[title="语言"]').click()
      }
    })
  })

  it('can add tags to a todo', () => {
    // Open Add Modal
    cy.contains('button', 'Add').click()
    
    // Type title and tags
    cy.get('input[placeholder="What needs to be done?"]').type('Tagged Task')
    cy.get('input[placeholder="Tags (comma separated)"]').type('work, urgent')
    
    // Click Add
    cy.get('div[class*="fixed"]').contains('button', 'Add').click()
    
    // Verify tags are displayed on the card
    cy.contains('Tagged Task').parents('.group').within(() => {
      cy.contains('#work').should('be.visible')
      cy.contains('#urgent').should('be.visible')
    })
  })

  it('can filter tasks by clicking tag in sidebar', () => {
    // 1. Create a task with a specific tag
    const uniqueTag = 'tag-' + Date.now()
    cy.contains('button', 'Add').click()
    cy.get('input[placeholder="What needs to be done?"]').type('Task with Unique Tag')
    cy.get('input[placeholder="Tags (comma separated)"]').type(uniqueTag)
    cy.get('div[class*="fixed"]').contains('button', 'Add').click()

    // 2. Verify tag appears in sidebar
    cy.get('aside').contains(uniqueTag).should('be.visible')

    // 3. Click "Inbox" or "All Tasks" first to ensure we are in a broader view
    cy.contains('button', 'All Tasks').click()

    // 4. Click the tag in sidebar
    cy.get('aside').contains(uniqueTag).click()

    // 5. Verify the view updates
    // The header should show "# <tagname>"
    cy.get('header').contains(`# ${uniqueTag}`).should('be.visible')
    
    // The task should be visible
    cy.contains('Task with Unique Tag').should('be.visible')
    
    // Other tasks should not be visible (create a dummy one to verify if needed, 
    // but verifying the header and presence is good enough for now)
  })

  it('can filter tasks by clicking tag on the card', () => {
    const cardTag = 'click-me'
    
    // Create task
    cy.contains('button', 'Add').click()
    cy.get('input[placeholder="What needs to be done?"]').type('Card Tag Task')
    cy.get('input[placeholder="Tags (comma separated)"]').type(cardTag)
    cy.get('div[class*="fixed"]').contains('button', 'Add').click()

    // Switch to All Tasks
    cy.contains('button', 'All Tasks').click()

    // Find the tag on the card and click it
    cy.contains('Card Tag Task').parents('.group').find('span').contains(`#${cardTag}`).click()

    // Verify view changed
    cy.get('header').contains(`# ${cardTag}`).should('be.visible')
    cy.contains('Card Tag Task').should('be.visible')
  })
})
