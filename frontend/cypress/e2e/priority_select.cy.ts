describe('Priority Selector Test', () => {
  beforeEach(() => {
    // Visit the app root
    cy.visit('/')
  })

  it('should display the priority selector with correct default value', () => {
    // Check if the BaseSelect button exists and contains "Medium" (default)
    // Note: The text might depend on the current locale, assuming default is 'en'
    // But since we fixed locale to 'en' in core memories, we check for English text or the component structure
    
    // Check if the button exists inside the BaseSelect component
    cy.get('button').contains('Medium').should('exist')
    
    // Click the selector to open dropdown
    cy.get('button').contains('Medium').click()
    
    // Check if options are visible
    cy.contains('High').should('be.visible')
    cy.contains('Low').should('be.visible')
  })

  it('should allow changing priority', () => {
    // Open selector
    cy.get('button').contains('Medium').click()
    
    // Select 'High'
    cy.contains('High').click()
    
    // Verify the button text changed to 'High'
    cy.get('button').contains('High').should('exist')
  })
})
