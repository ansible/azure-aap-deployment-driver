export function verifyRequiredEnvVariables() {
	let deploymentDriverUrl = Cypress.env('DEPLOYMENT_DRIVER_URL')
  let username = Cypress.env('USERNAME')
  let password = Cypress.env('PASSWORD')
  expect(deploymentDriverUrl).not.to.be.undefined
  // uncomment following when we implement SSO login
  //expect(username).not.to.be.undefined
  //expect(password).not.to.be.undefined
	return { deploymentDriverUrl, username, password }
}
