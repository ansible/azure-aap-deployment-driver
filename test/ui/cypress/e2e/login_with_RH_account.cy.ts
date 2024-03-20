/// <reference types="cypress" />

describe('Deployment Engine Login spec', () => {
  it('Login with Red Hat SSO', () => {
    cy.loginWithRHAccount()
  })
})