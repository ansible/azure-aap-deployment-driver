package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

var parameters map[string]interface{}         // this will contain the parameters to test
var values map[string]interface{}             // this will contain values parameter will resolve to
var outputs map[string]map[string]interface{} // this will contain output values parameter resolve to as well

func TestResolveReferencesToParameters(t *testing.T) {
	resolver := NewResolver(SUB_ID, RG_NAME)
	resolver.ResolveReferencesToParameters(parameters, values)
	resolver.ResolveReferencesToOutputs(parameters, outputs)

	verifyParameterValue(t, parameters, "name", getExpectedParameterValue("aapName"))
	verifyParameterValue(t, parameters, "location", getExpectedParameterValue("location"))
	verifyParameterValue(t, parameters, "crossTenantRoleAssignment", getExpectedParameterValue("crossTenantRoleAssignment"))
	verifyParameterValue(t, parameters, "logAnalyticsWorkspaceResourceID",
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/log-%s-%s",
			SUB_ID, RG_NAME, getExpectedParameterValue("aapName"), getExpectedParameterValue("location")))
	verifyParameterValue(t, parameters, "userAssignedIdentityId",
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id-%s-%s",
			SUB_ID, RG_NAME, getExpectedParameterValue("aapName"), getExpectedParameterValue("location")))
	verifyParameterValue(t, parameters, "userAssignedIdentityPrincipalId", "some-principal-id-would-be-here")
	// all three keyvault references point to same keyvault, hence we check for the same value
	// the name of the keyvault seems shortened so the last part is only "-east"
	expectedKeyVualtId := fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/Microsoft.KeyVault/vaults/kv-%s-east",
		SUB_ID, RG_NAME, getExpectedParameterValue("aapName"))
	verifyKeyVaultReference(t, parameters, "servicePrincipalId", expectedKeyVualtId)
	verifyKeyVaultReference(t, parameters, "servicePrincipalSecret", expectedKeyVualtId)
	verifyKeyVaultReference(t, parameters, "tenantId", expectedKeyVualtId)
	// following parameters are getting their value from deployments
	verifyParameterValue(t, parameters, "privateDnsZoneId", getExpectedOutputValue("dnsAAPDeploy", "aksDnsZoneId"))
	verifyParameterValue(t, parameters, "vnetSubnetID", getExpectedOutputValue("networkingAAPDeploy", "aksSubnetId"))
}

func TestMain(m *testing.M) {
	if err := json.Unmarshal([]byte(PARAMETERS), &parameters); err != nil {
		log.Fatalf("Couldn't parse test parameters. %v", err)
	}
	if err := json.Unmarshal([]byte(PARAMETER_VALUES), &values); err != nil {
		log.Fatalf("Couldn't parse test parameter values. %v", err)
	}
	if err := json.Unmarshal([]byte(OUTPUT_VALUES), &outputs); err != nil {
		log.Fatalf("Couldn't parse test parameter values. %v", err)
	}
	os.Exit(m.Run())
}

func verifyParameterValue(t *testing.T, parameters map[string]interface{}, name string, expectedValue interface{}) {
	// since the parameters are maps to any type, need to cast them to correct type
	value, exists := parameters[name].(map[string]interface{})["value"]
	if !exists {
		t.Errorf("Parameter '%s' or its 'value' does not exist.", name)
	}
	if value != expectedValue {
		t.Errorf("Parameter '%s' value does not match expected: got: %v <> expected: %v", name, value, expectedValue)
	}
}

func verifyKeyVaultReference(t *testing.T, parameters map[string]interface{}, name string, expectedValue interface{}) {
	value, exists := parameters[name].(map[string]interface{})["reference"].(map[string]interface{})["keyVault"].(map[string]interface{})["id"]
	if !exists {
		t.Errorf("Parameter '%s' or its nested fields do not exist.", name)
	}
	if value != expectedValue {
		t.Errorf("Parameter '%s' value does not match expected: got: %v <> expected: %v", name, value, expectedValue)
	}
}

func getExpectedParameterValue(name string) interface{} {
	return values[name].(map[string]interface{})["value"]
}

func getExpectedOutputValue(deployment, output string) interface{} {
	return outputs[deployment][output].(map[string]interface{})["value"]
}
