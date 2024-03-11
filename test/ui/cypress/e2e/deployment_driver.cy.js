import { verifyRequiredEnvVariables } from './environment_helpers'

describe('Deployment driver web UI', () => {
  let deploymentDriverUrl, username, password

  before(() => {
    ({deploymentDriverUrl, username, password} = verifyRequiredEnvVariables())
  })

  beforeEach(() => {
    cy.viewport(1920,1080)
    cy.visit(deploymentDriverUrl)
  })

  it('main view contains label, steps and status area with cancel button', () => {
    cy.get('main#primary-app-container').as('main')
    cy.get('@main').contains('Ansible Automation Platform Deployment Engine')
    cy.get('@main').contains('Deployment Steps')
    cy.get('@main').get('div.deployProgress').as('statusArea')
    cy.get('@statusArea').get('button.cancelButton').contains('Cancel Deployment')
  })

  it('left navigation has expected navigation links', () => {
    cy.get('div#page-sidebar').as('navigation')
    cy.get('@navigation').contains('Deployment')
    cy.get('@navigation').contains('Documentation')
    cy.get('@navigation').contains('Logout')
  })

  it('list of deployment steps contains expected steps', () => {
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
  })

  it('clicking cancel button brings up a dialog that can be closed', () => {
    cy.get('div.deployProgress').as('statusArea')
    cy.get('@statusArea').get('button.cancelButton').click()
    cy.get('div[role="dialog"]').as('dialog')
    cy.get('@dialog').contains('Cancel Deployment')
    cy.get('@dialog').get('button.pf-m-primary').contains('Confirm')
    cy.get('@dialog').get('button.pf-m-link').as('closeButton')
    cy.get('@closeButton').contains('No')
    cy.get('@closeButton').click()
    cy.get('@dialog').should('not.exist')
  })
})
