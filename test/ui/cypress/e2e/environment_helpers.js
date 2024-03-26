export function verifyRequiredEnvVariables() {
  expect(Cypress.env('DEPLOYMENT_DRIVER_URL')).not.to.be.undefined
  expect(Cypress.env('DEPLOYMENT_ENGINE_UI_PASSWORD')).not.to.be.undefined
  expect(Cypress.env('RH_SSO_URL')).not.to.be.undefined
  expect(Cypress.env('RH_ACCOUNT_USERNAME')).not.to.be.undefined
  expect(Cypress.env('RH_ACCOUNT_PASSWORD')).not.to.be.undefined
  return 
}
