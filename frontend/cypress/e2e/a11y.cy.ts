describe('Accessibility Testing', () => {
  beforeEach(() => {
    cy.visit('/')
    cy.injectAxe()
  })

  it('should have no detectable accessibility violations on load', () => {
    // Log violations to console
    const terminalLog = (violations) => {
      cy.task(
        'log',
        `${violations.length} accessibility violation${
          violations.length === 1 ? '' : 's'
        } ${violations.length === 1 ? 'was' : 'were'} detected`
      )
      // pluck specific keys to keep the table readable
      const violationData = violations.map(
        ({ id, impact, description, nodes }) => ({
          id,
          impact,
          description,
          nodes: nodes.length
        })
      )

      cy.task('table', violationData)
    }

    cy.checkA11y(undefined, {
      includedImpacts: ['critical', 'serious'],
      runOnly: {
        type: 'tag',
        values: ['wcag2a', 'wcag2aa'],
      },
    }, terminalLog)
  })
})
