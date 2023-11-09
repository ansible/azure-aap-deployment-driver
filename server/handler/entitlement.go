package handler

import (
	"net/http"
	"server/model"

	"gorm.io/gorm"
)

type EntitlementsCount struct {
	Count int64  `json:"count"`
	Error string `json:"error"`
}

func GetNumOfAzureMarketplaceEntitlements(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var entitlements []model.AzureMarketplaceEntitlement
	if tx := db.Model(&model.AzureMarketplaceEntitlement{}).Find(&entitlements); tx.Error != nil {
		respondError(w, http.StatusInternalServerError, tx.Error.Error())
		return
	}
	// iterate over the entitlements
	resp := EntitlementsCount{}
	for _, e := range entitlements {
		// error message here is from RH API call, we should not make it into our own error, hence 200 response
		if e.ErrorMessage != "" {
			resp.Error = e.ErrorMessage
			resp.Count = 0
			break
		}
		resp.Count = resp.Count + 1
	}
	respondJSON(w, http.StatusOK, &resp)
}
