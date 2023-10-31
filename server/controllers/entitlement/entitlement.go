package entitlement

import (
	"context"
	"encoding/json"
	"server/config"
	"server/model"
	"server/persistence"
	"server/util"
	"sync"

	log "github.com/sirupsen/logrus"
)

type EntitlementAPIController struct {
	ctx            context.Context
	httpRequester  *util.HttpRequester
	apiUrl         string
	productCode    string
	subscriptionId string
	database       *persistence.Database
}

type APIFilter struct {
	VendorProductCode   string `json:"vendorProductCode,omitempty"`
	AzureSubscriptionId string `json:"azureSubscriptionId,omitempty"`
	AzureTenantId       string `json:"azureTenantId,omitempty"`
}

type APIResponseContent struct {
	SourcePartner     string
	RhAccountId       string
	Status            string
	PartnerIdentities map[string]string
	RhEntitlements    []map[string]string
}

type APIResponse struct {
	Content []APIResponseContent
	Page    map[string]float64
}

var (
	once                    sync.Once
	entitlementCtrlInstance *EntitlementAPIController
)

func NewEntitlementController(context context.Context, db *persistence.Database) *EntitlementAPIController {
	once.Do(func() {

		cert := config.GetEnvironment().SW_SUB_API_CERTIFICATE
		key := config.GetEnvironment().SW_SUB_API_PRIVATEKEY
		url := config.GetEnvironment().SW_SUB_API_URL
		code := config.GetEnvironment().SW_SUB_VENDOR_PRODUCT_CODE
		subs := config.GetEnvironment().SUBSCRIPTION
		var requester *util.HttpRequester

		if cert == "" || key == "" {
			log.Warn("Entitlements controller will not be initialized because certificate or key are not provided.")
		} else {
			var err error
			requester, err = util.NewHttpRequesterWithCertificate(cert, key)
			if err != nil {
				log.Warnf("Could not initialize entitlements controller. %v\n", err)
			}
		}

		entitlementCtrlInstance = &EntitlementAPIController{
			ctx:            context,
			httpRequester:  requester,
			apiUrl:         url,
			productCode:    code,
			subscriptionId: subs,
			database:       db,
		}
	})
	return entitlementCtrlInstance
}

func (controller *EntitlementAPIController) FetchSubscriptions() {
	// no need to wait for this one, its not long running and http request uses context
	go func() {
		if controller.httpRequester != nil {
			// in the future we might need to handle pagination... maybe...
			resp, err := controller.httpRequester.MakeRequestWithJSONBody(
				controller.ctx,
				"POST",
				controller.apiUrl,
				nil,
				APIFilter{
					VendorProductCode:   controller.productCode,
					AzureSubscriptionId: controller.subscriptionId,
					//AzureTenantId:       "",
				},
			)
			if err != nil {
				log.Warnf("Failed to get response from subscription API: %v", err)
				return
			}

			response := APIResponse{}
			if err := json.Unmarshal(resp.Body, &response); err != nil {
				log.Errorf("Couldn't unmarshal JSON response. %v", err)
				return
			}

			if err := storeEntitlements(controller.database, &response); err != nil {
				log.Errorf("Couldn't persist entitlements in database. %v", err)
			}

			return
		}
		log.Warn("Entitlements can not be fetched, entitlement controller was not initialized.")
	}()
}

func storeEntitlements(db *persistence.Database, data *APIResponse) error {
	if len(data.Content) == 0 {
		return nil
	}

	for _, c := range data.Content {
		// supporting only Azure marketplace entitlements for now
		if c.SourcePartner == "azure_marketplace" {

			var azSubId, azCustId string
			var exists bool
			if azSubId, exists = c.PartnerIdentities["azureSubscriptionId"]; !exists {
				azSubId = ""
			}
			if azCustId, exists = c.PartnerIdentities["azureCustomerId"]; !exists {
				azCustId = ""
			}
			entitlement := model.AzureMarketplaceEntitlement{
				AzureSubscriptionId: azSubId,
				AzureCustomerId:     azCustId,
				RHEntitlements:      make([]model.RedHatEntitlements, 0),
				RedHatAccountId:     c.RhAccountId,
				Status:              c.Status,
			}
			for _, rhe := range c.RhEntitlements {
				var sku, subNum string
				var exists bool
				if sku, exists = rhe["sku"]; !exists {
					sku = ""
				}
				if subNum, exists = rhe["subscriptionNumber"]; !exists {
					subNum = ""
				}
				entitlement.RHEntitlements = append(entitlement.RHEntitlements, model.RedHatEntitlements{
					Sku:                sku,
					SubscriptionNumber: subNum,
				})
			}
			log.Printf("%v", entitlement)
			db.Instance.Save(&entitlement)
		}
	}
	return nil
}
