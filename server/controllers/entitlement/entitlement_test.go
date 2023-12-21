package entitlement_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"server/controllers/entitlement"
	"server/model"
	"server/persistence"
	"server/test"
	"server/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEntitlementController(t *testing.T) {
	// Create server
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Unable to read body from entitlement req: %v", err)
		}
		var filter entitlement.APIFilter
		err = json.Unmarshal(body, &filter)
		if err != nil {
			log.Printf("Testing error, could not unmarshal json: %v", err)
			os.Exit(1)
		}
		//if filter.VendorProductCode == "normal" {
		w.WriteHeader(http.StatusOK)
		ent := make(map[string]interface{})
		ent["rhAccountId"] = "123456"
		ent["sourcePartner"] = "azure_marketplace"
		ent["partnerIdentities"] = map[string]string{
			"azureSubscriptionId": "7890ab",
			"azureTenantId":       "cdef01",
		}
		ent["rhEntitlements"] = []map[string]string{
			{
				"sku":                "AnsibleSku",
				"subscriptionNumber": "234567",
			},
		}
		ent["purchase"] = map[string]interface{}{
			"vendorProductCode": "rhmaap",
			"contracts":         []string{},
		}
		ent["status"] = "SUBSCRIBED"
		ent["entitlementDates"] = map[string]string{
			"startDate": "2023-01-01T00:00:00.000000Z",
		}
		if filter.VendorProductCode == "rhmaap" {
			ent["purchase"] = map[string]interface{}{
				"vendorProductCode": "rhmaap",
				"contracts":         []string{},
			}
		} else {
			// Induce json unmarshal error
			ent["sourcePartner"] = map[string]string{"THISWILLFAIL": "NOTCONSTRUCTEDPROPERLY"}
		}
		resp := make(map[string]interface{})
		resp["Content"] = []map[string]interface{}{ent}
		resp["Page"] = map[string]float64{"one": 1}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Cant marshal server response to json: %v", err)
		}
		_, err = w.Write(jsonResp)
		if err != nil {
			t.Errorf("Unable to write server response: %v", err)
		}
	}))

	os.Setenv("SW_SUB_VENDOR_PRODUCT_CODE", "rhmaap")
	os.Setenv("SUBSCRIPTION", "234567")
	os.Setenv("SW_SUB_API_URL", fmt.Sprintf("%s/%s", server.URL, "entitlements"))
	os.Setenv("SW_SUB_API_CERTIFICATE", "dummy")
	os.Setenv("SW_SUB_API_PRIVATEKEY", "dummy")
	test.SetEnvironment()

	db := persistence.NewInMemoryDB()
	ec := entitlement.NewEntitlementController(context.Background(), db)
	req := util.NewHttpRequester()
	req.TestOnlySetClient(server.Client())
	ec.TestSetRequester(req)
	ec.FetchSubscriptions()
	time.Sleep(1 * time.Second) // Time for goroutine to finish
	ent := []model.AzureMarketplaceEntitlement{}
	db.Instance.Find(&ent)

	assert.Len(t, ent, 1)
	ec.TestSetProductCode("BOGUS")
	ec.FetchSubscriptions()
	time.Sleep(1 * time.Second)
	db.Instance.Find(&ent)
	assert.Len(t, ent, 2)
	assert.NotEqual(t, ent[1].ErrorMessage, "")
}
