Cypress.Commands.add('logoutDeploymentEngineUI', () => {
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

});