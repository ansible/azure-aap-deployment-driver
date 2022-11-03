// ==== This file is for test data only ====
package engine

const (
	SUB_ID     = "some-sub-id-would-be-here" // make sure this matches the value used in PARAMETER_VALUES
	RG_NAME    = "mrg-aapSomething"          // make sure this matches the value used in PARAMETER_VALUES
	PARAMETERS = `
{
  "name": {
    "value": "[parameters('aapName')]"
  },
  "location": {
    "value": "[parameters('location')]"
  },
  "logAnalyticsWorkspaceResourceID": {
    "value": "[parameters('logWorkspaceId')]"
  },
  "privateDnsZoneId": {
    "value": "[reference(resourceId('Microsoft.Resources/deployments', 'dnsAAPDeploy')).outputs.aksDnsZoneId.value]"
  },
  "userAssignedIdentityId": {
    "value": "[parameters('userAssignedIdentityId')]"
  },
  "userAssignedIdentityPrincipalId": {
    "value": "[parameters('userAssignedIdentityPrincipalId')]"
  },
  "vnetSubnetID": {
    "value": "[reference(resourceId('Microsoft.Resources/deployments', 'networkingAAPDeploy')).outputs.aksSubnetId.value]"
  },
  "servicePrincipalId": {
    "reference": {
      "keyVault": {
        "id": "[extensionResourceId(format('/subscriptions/{0}/resourceGroups/{1}', subscription().subscriptionId, resourceGroup().name), 'Microsoft.KeyVault/vaults', parameters('keyVaultName'))]"
      },
      "secretName": "service-principal-client-id"
    }
  },
  "servicePrincipalSecret": {
    "reference": {
      "keyVault": {
        "id": "[extensionResourceId(format('/subscriptions/{0}/resourceGroups/{1}', subscription().subscriptionId, resourceGroup().name), 'Microsoft.KeyVault/vaults', parameters('keyVaultName'))]"
      },
      "secretName": "service-principal-secret"
    }
  },
  "tenantId": {
    "reference": {
      "keyVault": {
        "id": "[extensionResourceId(format('/subscriptions/{0}/resourceGroups/{1}', subscription().subscriptionId, resourceGroup().name), 'Microsoft.KeyVault/vaults', parameters('keyVaultName'))]"
      },
      "secretName": "tenant-id"
    }
  },
  "crossTenantRoleAssignment": {
    "value": "[parameters('crossTenantRoleAssignment')]"
  }
}
`
	PARAMETER_VALUES = `
{
  "aapName": {
    "type": "String",
    "value": "aapuniquename123"
  },
  "access": {
    "type": "String",
    "value": "public"
  },
  "crossTenantRoleAssignment": {
    "type": "Bool",
    "value": false
  },
  "location": {
    "type": "String",
    "value": "eastus"
  },
  "keyVaultName": {
    "type": "String",
    "value": "kv-aapuniquename123-east"
  },
  "logWorkspaceId": {
    "type": "String",
    "value": "/subscriptions/some-sub-id-would-be-here/resourceGroups/mrg-aapSomething/providers/Microsoft.OperationalInsights/workspaces/log-aapuniquename123-eastus"
  },
  "userAssignedIdentityId": {
    "type": "String",
    "value": "/subscriptions/some-sub-id-would-be-here/resourceGroups/mrg-aapSomething/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id-aapuniquename123-eastus"
  },
  "userAssignedIdentityPrincipalId": {
    "type": "String",
    "value": "some-principal-id-would-be-here"
  },
  "publicDNSZoneName": {
    "type": "String",
    "value": "aapuniquename123.some.domain.name"
  },
  "vnetConfig": {
    "type": "Object",
    "value": {
      "name": "vnet01_doNotEditName",
      "resourceGroup": "resourceGroup",
      "addressPrefixes": [
        "10.33.33.0/24"
      ],
      "addressPrefix": "10.33.33.0/24",
      "newOrExisting": "new",
      "subnets": {
        "aks": {
          "name": "cluster_doNotEditName",
          "addressPrefix": "10.33.33.0/26",
          "startAddress": "10.33.33.4"
        },
        "appgw": {
          "name": "appgw_doNotEditName",
          "addressPrefix": "10.33.33.64/28",
          "startAddress": "10.33.33.68"
        },
        "postgres": {
          "name": "database_doNotEditName",
          "addressPrefix": "10.33.33.80/28",
          "startAddress": "10.33.33.84"
        },
        "plink": {
          "name": "private_link_doNotEditName",
          "addressPrefix": "10.33.33.96/28",
          "startAddress": "10.33.33.100"
        }
      }
    }
  }
}
`
	OUTPUT_VALUES = `
{
  "dnsAAPDeploy":{
    "aksDnsZoneId":{
      "value": "/subscriptions/some-sub-id-would-be-here/resourceGroups/mrg-aapSomething/providers/Microsoft.Network/privateDnsZones/privatelink.eastus.azmk8s.io"
    },
    "blobDnsZoneId":{
      "value":""
    },
    "blobDnsZoneName":{
      "value":""
    },
    "psqlDnsZoneId":{
      "value":""
    }
  },
  "networkingAAPDeploy":{
    "aksSubnetId": {
      "value":"/subscriptions/some-sub-id-would-be-here/resourceGroups/mrg-aapSomething/providers/Microsoft.Network/virtualNetworks/vnet-aapuniquename123-eastus/subnets/snet-aapuniquename123-eastus-aks"
    }
  }
}
`
)
