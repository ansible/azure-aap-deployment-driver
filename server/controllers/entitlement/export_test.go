package entitlement

import "server/util"

// Just exports setters (within package) for ease of test

func (c *EntitlementAPIController) TestOnlySetProductCode(productCode string) {
	c.productCode = productCode
}

func (c *EntitlementAPIController) TestOnlySetRequester(requester *util.HttpRequester) {
	c.httpRequester = requester
}
