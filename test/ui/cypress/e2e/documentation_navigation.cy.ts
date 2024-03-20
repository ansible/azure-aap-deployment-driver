/// <reference types="cypress" />

describe('Documentation Navigation Link spec', () => {
    it('Check the UI elements for the Documentation link', () => {
      cy.loginWithRHAccount()

      cy.get('#nav-toggle').click()
      cy.get('#Documentation-1').click()

      // Check the title of the Documentation screen
      cy.get('.pf-c-title')
        .should('have.text','Documentation')
        .should('be.visible')

        cy.get('.pf-c-empty-state__body > .pf-c-button')
        .contains('Red Hat Ansible Automation Platform on Microsoft Azure Guide')
        .should('have.attr','href',
        'https://access.redhat.com/documentation/en-us/ansible_on_clouds/2.x/html/red_hat_ansible_automation_platform_on_microsoft_azure_guide/index')          
    })
  })