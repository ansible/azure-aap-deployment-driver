package entitlement

import "server/util"

// Just exports setters for ease of test

func (c *EntitlementAPIController) TestSetProductCode(productCode string) {
	c.productCode = productCode
}

func (c *EntitlementAPIController) TestSetRequester(requester *util.HttpRequester) {
	c.httpRequester = requester
}
