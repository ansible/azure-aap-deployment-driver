/// <reference types="cypress" />
// ***********************************************
// This example commands.ts shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })
//
import './auth'
import './logout'

declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * Use AAP admin credential to login Deployment Engine UI, then complete the 
       * Red Hat account login from the Deployment Engine modal diagol. The user will be
       * redirected to the Deployment Engine UI after completing the Red Hat account login
       * successfully.
       * @example
       * cy.loginWithRHAccount()
       */
      loginWithRHAccount(): Chainable

      /**
       * Logout Deployment Engine UI via the Logout navigation link
       * @example
       * cy.logoutDeploymentEngineUI()
       */
      logoutDeploymentEngineUI(): Chainable

      /**
       * Check if the required variables are set       * 
       * @param requiredVariables placeholders for the required variables
       * @example
       * cy.requiredVariablesAreSet(['ENV_VAR1','ENV_VAR2','ENV_VAR3'])
       */
      requiredVariablesAreSet(requiredVariables: string[]): Chainable
    }
  }
}