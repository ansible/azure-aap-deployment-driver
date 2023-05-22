package engine

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
)

const (
	PARAMETERS_REGEXP = `\[parameters\('(\w+)'\)\]`
	KEYVAULTID_REGEXP = `\[extensionResourceId\(format\('/subscriptions/\{0\}/resourceGroups/\{1\}', subscription\(\)\.subscriptionId, resourceGroup\(\)\.name\), 'Microsoft\.KeyVault/vaults', parameters\('(\w+)'\)\)\]`
	// Another possible value for keyvault id reference regexp: `\[resourceId\('Microsoft\.KeyVault/vaults', parameters\('(\w+)'\)\)\]`
	KEYVAULTID_REPLACEMENT = "/subscriptions/%s/resourcegroups/%s/providers/Microsoft.KeyVault/vaults/%s"
	OUTPUTS_REGEXP         = `\[reference\(resourceId\('Microsoft\.Resources/deployments', '(\w+)'\).*\)\.outputs\.(\w+)\.value\]`
)

type Resolver struct {
	resourceGroup                         string
	subscriptionId                        string
	parametersReferenceRegexp             *regexp.Regexp
	parametersReferenceInKeyVaultIdRegexp *regexp.Regexp
	outputsReferenceRegex                 *regexp.Regexp
}

func NewResolver(subscriptionId, resourceGroup string) *Resolver {
	return &Resolver{
		subscriptionId:                        subscriptionId,
		resourceGroup:                         resourceGroup,
		parametersReferenceRegexp:             regexp.MustCompile(PARAMETERS_REGEXP),
		parametersReferenceInKeyVaultIdRegexp: regexp.MustCompile(KEYVAULTID_REGEXP),
		outputsReferenceRegex:                 regexp.MustCompile(OUTPUTS_REGEXP),
	}
}

func (resolver Resolver) ResolveReferencesToParameters(parameters map[string]interface{}, valueSource map[string]interface{}) {
	for k, v := range parameters {
		// k is the name of the parameter, v is the object containing either "value" or "reference" fields
		parameterName := getReferencedParameterName(v, resolver.parametersReferenceRegexp)
		if parameterName != "" {
			setParameterValue(parameters, k, parameterName, valueSource)
			continue
		}

		// Following block handles nested field "reference.keyVault.id" that can contain reference to a parameter
		keyvaultId := getReferencedKeyVaultId(v, resolver.parametersReferenceInKeyVaultIdRegexp)
		if keyvaultId != "" {
			setKeyVaultId(parameters, k, keyvaultId, valueSource, resolver.subscriptionId, resolver.resourceGroup)
		}
	}
}

func (resolver Resolver) ResolveDryRunParamsMap(params map[string]interface{}, outputs map[string]interface{}) map[string]interface{} {
	// Dry run wants just a map of key/value pairs {"access": "public"} for instance
	outMap := make(map[string]interface{})
	for k, v := range params {
		val, ok := outputs[k]
		if ok {
			outMap[k] = val.(map[string]interface{})["value"]
		} else {
			// Take default (empty) value since it will have correct type
			outMap[k] = v.(map[string]interface{})["value"]
		}
	}
	return outMap
}

func (resolver Resolver) ResolveReferencesToOutputs(parameters map[string]interface{}, outputs map[string]map[string]interface{}) error {
	for k, v := range parameters {

		// k is the name of the parameter, v is the object containing "value" field
		deploymentName, outputName := getReferencedDeploymentOutput(v, resolver.outputsReferenceRegex)
		if deploymentName != "" && outputName != "" {
			if err := setOutputValue(parameters, k, deploymentName, outputName, outputs); err != nil {
				return err
			}
		}
	}
	return nil

	// 	value, exists := v.(map[string]interface{})["value"]
	// 	valueString, isString := value.(string)
	// 	if exists && isString {
	// 		outputReference := resolver.outputsReferenceRegex.FindStringSubmatch(valueString)
	// 		if len(outputReference) >= 3 { // there must be at least 3, full match and two groups
	// 			deploymentName := outputReference[1]
	// 			outputName := outputReference[2]
	// 			deploymentOutputs, deploymentOutputsExist := outputs[deploymentName]
	// 			if deploymentOutputsExist {
	// 				output, outputExists := deploymentOutputs[outputName]
	// 				if outputExists {
	// 					parameters[k].(map[string]interface{})["value"] = output.(map[string]interface{})["value"]
	// 					log.Debugf("Configured value for parameter '%s' referencing output: %s.%s", k, deploymentName, outputName)
	// 				}
	// 			} else {
	// 				return fmt.Errorf("deployment named '%s' does not exist", deploymentName)
	// 			}
	// 		}
	// 	}
	// }
	// return nil
}

func getReferencedParameterName(parameterValue interface{}, referenceRegexp *regexp.Regexp) string {
	parameterName := ""
	// Following block handles simple "value" field that can contain reference to a parameter
	valueField, valueFieldExists := parameterValue.(map[string]interface{})["value"]
	valueString, valueFieldIsString := valueField.(string)
	if valueFieldExists && valueFieldIsString {
		// check if the value is a parameter reference (as opposed to other values or references)
		paramReference := referenceRegexp.FindStringSubmatch(valueString)
		if len(paramReference) >= 2 { // there must be at least 2, one for full match and one for the group
			parameterName = paramReference[1] // using index 1 because 0 is the full match, we need
		}
	}
	return parameterName
}

func setParameterValue(parameters map[string]interface{}, name string, sourceName string, valueSource map[string]interface{}) {
	for sk, sv := range valueSource {
		// sk is the parameter name, sk is the object containing "value" field
		if sk == sourceName {
			parameters[name].(map[string]interface{})["value"] = sv.(map[string]interface{})["value"]
			log.Debugf("Configured value for parameter '%s' referencing another parameter named: %s", name, sourceName)
			break // bail once value is set
		}
	}
}

func getReferencedKeyVaultId(parameterValue interface{}, referenceRegexp *regexp.Regexp) string {
	keyVaultId := ""
	referenceField, referenceFieldExists := parameterValue.(map[string]interface{})["reference"]
	if referenceFieldExists {
		keyVaultField, keyVaultFieldExists := referenceField.(map[string]interface{})["keyVault"]
		if keyVaultFieldExists {
			idField, idFieldExists := keyVaultField.(map[string]interface{})["id"]
			idString, idFieldIsString := idField.(string)
			if idFieldExists && idFieldIsString {
				// check if the value is a parameter reference (as opposed to other values or references)
				paramReference := referenceRegexp.FindStringSubmatch(idString)
				if len(paramReference) >= 2 { // there must be at least 2, one for full match and one for the group
					keyVaultId = paramReference[1] // using index 1 because 0 is the full match, we need only the group
				}
			}
		}
	}
	return keyVaultId
}

func setKeyVaultId(parameters map[string]interface{}, name string, sourceName string, valueSource map[string]interface{}, subscriptionId string, resourceGroup string) {
	for sk, sv := range valueSource {
		// sk is the parameter name, sk is the object containing "value" field
		if sk == sourceName {
			resolvedValue := fmt.Sprintf(KEYVAULTID_REPLACEMENT, subscriptionId, resourceGroup, sv.(map[string]interface{})["value"])
			parameters[name].(map[string]interface{})["reference"].(map[string]interface{})["keyVault"].(map[string]interface{})["id"] = resolvedValue
			log.Debugf("Configured key vault id reference in parameter '%s' referencing another parameter named '%s' with value: %s", name, sourceName, resolvedValue)
			break
		}
	}
}

func getReferencedDeploymentOutput(parameterValue interface{}, referenceRegexp *regexp.Regexp) (string, string) {
	deploymentName := ""
	outputName := ""
	value, exists := parameterValue.(map[string]interface{})["value"]
	valueString, isString := value.(string)
	if exists && isString {
		outputReference := referenceRegexp.FindStringSubmatch(valueString)
		if len(outputReference) >= 3 { // there must be at least 3, full match and two groups
			deploymentName = outputReference[1]
			outputName = outputReference[2]
		}
	}
	return deploymentName, outputName
}

func setOutputValue(parameters map[string]interface{}, name string, deploymentName string, outputName string, outputs map[string]map[string]interface{}) error {
	deploymentOutputs, deploymentOutputsExist := outputs[deploymentName]
	if deploymentOutputsExist {
		output, outputExists := deploymentOutputs[outputName]
		if outputExists {
			parameters[name].(map[string]interface{})["value"] = output.(map[string]interface{})["value"]
			log.Debugf("Configured value for parameter '%s' referencing output: %s.%s", name, deploymentName, outputName)
		} else {
			return fmt.Errorf("output name '%s' from deployment '%s' does not exist", outputName, deploymentName)
		}
	} else {
		return fmt.Errorf("deployment named '%s' does not exist", deploymentName)
	}
	return nil
}
