/// <reference types="cypress" />

describe('Deployment driver web UI', () => {
  beforeEach(() => {
    cy.viewport(1920, 1080)
    cy.loginWithRHAccount()
  })

  it('main view contains label, steps and status area with cancel button', () => {
    cy.get('main#primary-app-container').as('main')
    cy.get('@main').contains('Ansible Automation Platform Deployment Engine')
    cy.get('@main').contains('Deployment Steps')
    cy.get('@main').get('div.deployProgress').as('statusArea')
    cy.get('@statusArea').get('button.cancelButton').contains('Cancel Deployment')

    cy.get('.pf-c-brand').should('be.visible')
    cy.get('.pf-c-alert__title')
      .contains('You currently have a subscription to Ansible Automation Platform')
    cy.get('p').contains('To manage or setup new subscription, visit the')
    cy.get('p > .pf-c-button')
      .contains('Red Hat Hybrid Cloud Console')
      .should('have.attr', 'href', 'https://console.redhat.com/')

    // Check the H1 title
    cy.get('h1')
      .contains('Ansible Automation Platform Deployment Engine')
      .should('be.visible') 
  })

  it('left navigation has expected navigation links', () => {
    cy.get('#nav-toggle').click()
    cy.get('div#page-sidebar').as('navigation')

    cy.get('@navigation').contains('Deployment')
      .should('be.visible')
      .should('have.attr','href').and('include','/')

    cy.get('@navigation').contains('Documentation')
      .should('be.visible')
      .should('have.attr','href').and('include', '/documentation')

    cy.get('@navigation').contains('Logout')
      .should('be.visible')    
  })

  it('list of deployment steps contains expected steps', () => {
    // Check the title of Deployment Steps
    cy.get('.pf-c-title')
      .contains('Deployment Steps')
      .should('be.visible')
    
    // Check the detailed deployment step's names in the Deployment Steps panel
    cy.get('div.deploy-step>ul').as('steps')
    cy.get('@steps').contains('VNET and_Subnets')
    cy.get('@steps').contains('Private DNS')
    cy.get('@steps').contains('AAP Repository')
    cy.get('@steps').contains('Database Server and Databases')
    cy.get('@steps').contains('AKS Cluster')
    cy.get('@steps').contains('AAP Operators')
    cy.get('@steps').contains('AAP_Applications')
    cy.get('@steps').contains('Application Ingress')
    cy.get('@steps').contains('Seeded Content')
    cy.get('@steps').contains('Billing')
    cy.get('@steps').contains('Deployment Cleanup')

  })


  it('clicking cancel button brings up a dialog that can be closed', () => {
    // Check the Cancel Deployment button on the main screen
    cy.get('div.deployProgress').as('statusArea')
    cy.get('@statusArea').get('button.cancelButton').click()

    // Check the UI elements on the pupped up dialog after clicking the Cancel Deployment button from the main screen
    cy.get('div[role="dialog"]').as('dialog')
    cy.get('@dialog', {timeout: 4000}).contains('Cancel Deployment')
    cy.get('@dialog').get('button.pf-m-primary').contains('Confirm')
    cy.get('@dialog').get('button.pf-m-link').as('closeButton')
    cy.get('@closeButton').contains('No').should('not.be.selected')
    cy.get('@closeButton').click()
    cy.get('@dialog').should('not.exist')   
  })

  it('Checking the Deployment Engine UI logout', () => {
    // Select the Logout from the navigation menu  
  cy.get('#nav-toggle').click()
  cy.get('#Logout-2').click()

  //Check the messages on the popped up Logout screen
  cy.get('.pf-c-modal-box__title-text')
    .contains('Logout')
    .should('be.visible')
  cy.get('#pf-modal-part-3')
    .contains('Are you sure you want Logout?')
    .should('be.visible')
 
  // Check the 'Cancel' button is unselected
  cy.get('.pf-c-modal-box__footer > .pf-m-link')
    .contains('Cancel')
    .should('not.be.selected')

  //Check the 'Confirm' button is selected 
  cy.get('button#primary-loading-button.pf-c-button.pf-m-primary.pf-m-progress', {timeout: 4000})
    .contains('Confirm')
    .should('be.visible')
    .click()
  })
})
