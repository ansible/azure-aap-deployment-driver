Cypress.Commands.add('loginWithRHAccount', () => {
    //cy.get('#pf-login-username-id').type('admin') - hardcoded by Deployment Engine UI
    // Enter the admin password for Deployment Engine UI login
    cy.get('#pf-login-password-id').type(Cypress.env('DEPLOYMENT_ENGINE_UI_PASSWORD'))
    cy.get('.pf-c-form__actions > .pf-c-button').click()

    // Check the UI elements on the Red Hat login dialog screen
    cy.get('.pf-c-modal-box__header > .pf-c-title')
      .contains('A valid subscription for Ansible Automation Platform in your Red Hat account is required')
      
    // Check the messages on the Red Hat login dialog screen
    cy.get('#pf-modal-part-2 > :nth-child(1)')
      .contains('Your Ansible Automation Platform deployment is underway') 
    cy.get('#pf-modal-part-2 > :nth-child(3)')
      .contains('To use Ansible Automation Platform on Azure, you MUST have a valid subscription for Ansible Automation Platform in your Red Hat account.')
    cy.get('#pf-modal-part-2 > :nth-child(5)')
      .contains('You can set up your Ansible Automation Platform subscription and your Red Hat account by clicking the button below. You will be redirected back to this page upon successful log in or account creation.')
    
    // Click the 'Log in with Red Hat account' button
    cy.get('.pf-c-modal-box__footer > .pf-c-button')
      .contains('Log in with Red Hat account')
      .should('have.focus')
      .click()

    // 'Log in with Red Hat account using user name and password
    cy.origin(Cypress.env('RH_SSO_URL'), () => {
      cy.on('uncaught:exception', (err, runnable) => {
        // returning false here prevents Cypress
        // inside the cy.origin() method from failing the test
        return false
      })
      cy.get('#username-verification', {timeout: 6000}).type(Cypress.env('RH_ACCOUNT_USERNAME'), {force: true})
      cy.get('#login-show-step2', {timeout: 6000})
        .should('have.text', 'Next')
        .click({force: true})
    cy.get('#password', {timeout: 6000}).type(Cypress.env('RH_ACCOUNT_PASSWORD'), {force: true})
    cy.get('button#rh-password-verification-submit-button.pf-c-button')
      .click({force: true})
    })

    // Check that the user is redirected to the Deployment Engine UI
    cy.get('h1')
      .contains('Ansible Automation Platform Deployment Engine')
      .should('be.visible')
 
  });
