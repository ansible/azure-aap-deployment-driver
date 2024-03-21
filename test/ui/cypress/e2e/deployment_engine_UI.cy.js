/// <reference types="cypress" />
describe('Deployment Engine UI spec', () => {
  beforeEach(() => {
    cy.loginWithRHAccount()
  })
  // The following afterEach() caused the test case execution failure. 
  // Need to do some more researh.
  //afterEach(() => {
  //  cy.logoutDeploymentEngineUI()
  //})

  // Check the navigation links
  it('Check the navigation links', () => {
    cy.get('#nav-toggle').click()

    // Check the Deployment navigation link
    cy.get('#Deployment-0')
    .contains('Deployment')
    .should('be.visible')
    .should('have.attr','href').and('include','/')
    
    // Check the Documentation navigation link
    cy.get('#Documentation-1')
      .contains('Documentation')
      .should('be.visible')
      .should('have.attr','href').and('include', '/documentation')
    
    // Check the Logout navigation link
    cy.get('#Logout-2')
      .contains('Logout')
      .should('be.visible')
  })

  it('Check the titles and button on the main screen', () => {
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

  it('Check the Deployment Steps', () => {
    // Check the title of Deployment Steps
    cy.get('.pf-c-title')
      .contains('Deployment Steps')
      .should('be.visible')
    
    // Check the detailed deployment step's names in the Deployment Steps panel
    cy.fixture('deployment_steps.json').then((deploymentStepsData) => {
      let total_steps = deploymentStepsData.steps.length
      cy.log('Total Deployment Steps: '+total_steps)
      
      cy.get('div.deploy-step>ul').as('steps')
      for (let i=0; i<total_steps; i++) {
        cy.log('Deployment Step '+ deploymentStepsData.steps[i].stepNo +": " + deploymentStepsData.steps[i].Name)
        
        cy.get('@steps')
          .contains(`${deploymentStepsData.steps[`${i}`].Name}`)
      }
    })
  })

  it('Check the Cancel Deployment button and its related messages', () => {
    // Check the Cancel Deployment button on the main screen
    cy.get('.pf-l-bullseye > .pf-c-button')
      .should('have.text', 'Cancel Deployment')
      .click()

    // Check the messages on the pupped up screen after clicking the Cancel Deployment button
    cy.get('.pf-c-modal-box__title-text')
      .contains('Cancel Deployment')
      .should('be.visible')
    cy.get('#pf-modal-part-4')
      .contains("Are you sure you want to cancel your deployment? If so, click the 'Confirm' Button or press the 'No' Button to return to your Deployment.")
      .should('be.visible')
    
    // Check the "Confirm" button has focus by default
    cy.get('.pf-c-modal-box__footer > .pf-m-primary', {timeout: 4000})
      .contains('Confirm')
      .should('be.visible')
    
    // Check there are options to give up in canceling the deployment
    cy.get('#pf-modal-part-2 > .pf-m-plain')
      .should('be.visible')
    cy.get('.pf-c-modal-box__footer > .pf-m-link')
      .contains('No')
      .should('not.be.selected')
      .click()
  })

  it('Logout the Deployment Engine UI', () => {
    cy.logoutDeploymentEngineUI()
  })

})