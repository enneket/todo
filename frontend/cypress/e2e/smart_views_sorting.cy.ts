describe('Smart Views and Sorting', () => {
  // Fixed "Now": June 15, 2024 12:00 PM
  const NOW = new Date('2024-06-15T12:00:00').getTime()

  beforeEach(() => {
    // Set fixed system time
    cy.clock(NOW)
    cy.visit('/')

    // Ensure we are in English mode
    cy.get('body').then(($body) => {
      if ($body.text().includes('添加')) {
         cy.get('button[aria-label="Language"]').click()
         cy.contains('Todo App').should('be.visible')
      }
    })
  })

  // Helper to create task
  const createTask = (title: string, date: string, priority: string) => {
    // Open Add Modal
    // If floating button exists, use it. Or text "Add".
    // The hero button has text "Add" inside.
    cy.contains('button', 'Add').click()
    
    // Type Title
    cy.get('input[placeholder="What needs to be done?"]').type(title)
    
    // Set Priority
    // Find the label "Priority" and get the dropdown button next to it (in the same container div)
    cy.contains('label', 'Priority').parent().find('button').first().click()
    // Select the option
    const priorityLabel = priority.charAt(0).toUpperCase() + priority.slice(1)
    // The dropdown is in a portal/absolute div, but since we just clicked, it should be visible.
    // We look for the button with the label inside the dropdown
    cy.get('div[class*="absolute"]').contains('button', priorityLabel).click()

    // Set Due Date
    if (date) {
        // datetime-local input
        cy.contains('label', 'Due Date').parent().find('input').type(date)
    }

    // Click Save/Add
    // In modal footer
    cy.get('div[class*="fixed"]').contains('button', 'Add').click()
    
    // Verify it appeared (wait for list update)
    cy.contains(title).should('be.visible')
  }

  it('filters tasks by Smart Views', () => {
    // 1. Create Data
    // Today: 2024-06-15
    createTask('Task Today', '2024-06-15T14:00', 'medium')
    
    // Upcoming (Tomorrow): 2024-06-16
    createTask('Task Upcoming', '2024-06-16T14:00', 'low')
    
    // Overdue (Yesterday): 2024-06-14
    createTask('Task Overdue', '2024-06-14T14:00', 'high')
    
    // Future (Next Month): 2024-07-01
    createTask('Task Future', '2024-07-01T14:00', 'medium')

    // 2. Verify "Today" View
    cy.contains('button', 'Today').click()
    // Should see:
    cy.contains('Task Today').should('be.visible')
    // Should NOT see:
    cy.contains('Task Upcoming').should('not.exist')
    cy.contains('Task Overdue').should('not.exist')
    cy.contains('Task Future').should('not.exist')

    // 3. Verify "Upcoming" View
    cy.contains('button', 'Upcoming').click()
    // Should see Today and Upcoming (Next 7 days)
    cy.contains('Task Today').should('be.visible')
    cy.contains('Task Upcoming').should('be.visible')
    // Should NOT see Overdue or Far Future
    cy.contains('Task Overdue').should('not.exist')
    cy.contains('Task Future').should('not.exist')

    // 4. Verify "Overdue" View
    cy.contains('button', 'Overdue').click()
    // Should see Overdue
    cy.contains('Task Overdue').should('be.visible')
    // Should NOT see others
    cy.contains('Task Today').should('not.exist')
    cy.contains('Task Upcoming').should('not.exist')
  })

  it('sorts tasks by Priority and Due Date', () => {
    const prefix = 'SortTest '
    // Create tasks with specific attributes for sorting
    // High Priority, Mid Date
    createTask(prefix + 'High', '2024-06-15T12:00', 'high')
    // Medium Priority, Late Date
    createTask(prefix + 'Medium', '2024-06-15T18:00', 'medium')
    // Low Priority, Early Date
    createTask(prefix + 'Low', '2024-06-15T08:00', 'low')

    // Switch to All Tasks to see them all
    cy.contains('button', 'All Tasks').click()

    // 1. Sort by Priority (High -> Medium -> Low)
    cy.get('button[aria-label="Sort by"]').click()
    cy.contains('button', 'Priority').click()
    
    // Verify Order
    cy.get('main h3').should($items => {
        const texts = [...$items].map(i => i.innerText).filter(t => t.startsWith(prefix))
        expect(texts).to.deep.equal([prefix + 'High', prefix + 'Medium', prefix + 'Low'])
    })

    // 2. Sort by Due Date Ascending (Early -> Mid -> Late)
    cy.get('button[aria-label="Sort by"]').click()
    cy.contains('button', 'Due Date (Ascending)').click()
    
    cy.get('main h3').should($items => {
        const texts = [...$items].map(i => i.innerText).filter(t => t.startsWith(prefix))
        expect(texts).to.deep.equal([prefix + 'Low', prefix + 'High', prefix + 'Medium'])
    })

    // 3. Sort by Due Date Descending (Late -> Mid -> Early)
    cy.get('button[aria-label="Sort by"]').click()
    cy.contains('button', 'Due Date (Descending)').click()
    
    cy.get('main h3').should($items => {
        const texts = [...$items].map(i => i.innerText).filter(t => t.startsWith(prefix))
        expect(texts).to.deep.equal([prefix + 'Medium', prefix + 'High', prefix + 'Low'])
    })
  })
})
