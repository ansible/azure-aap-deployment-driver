package templates

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const MainTemplateName = "mainTemplate"

// gets the main template and parameters
//
//	returns: mainTemplate, mainParameters, error
func GetMainTemplateAndParameters(templateBasePath string) (map[string]any, map[string]any, error) {
	mainTemplate, err := readJSON(filepath.Join(templateBasePath, MainTemplateName+".json"))
	if err != nil {
		return nil, nil, err
	}
	mainParameters, err := readJSON(filepath.Join(templateBasePath, MainTemplateName+".parameters.json"))
	if err != nil {
		return nil, nil, err
	}

	return mainTemplate, mainParameters, nil
}

func DiscoverTemplateOrder(templateBasePath string) ([][]string, error) {
	log.Infof("Starting deployment template discovery in location: %s", templateBasePath)

	dependenciesGraph := NewDependencyGraph()

	templateDirEntries, err := os.ReadDir(templateBasePath)
	entryCount := len(templateDirEntries)

	if entryCount == 0 {
		log.Infof("%d deployment templates found in location [%s]", entryCount, templateBasePath)
	}

	if err != nil {
		return nil, err
	}

	for _, entry := range templateDirEntries {
		// expecting only directories
		if entry.IsDir() {
			name := entry.Name()
			dependencyFileName := filepath.Join(templateBasePath, name, name+".dependencies.json")
			// read dependencies file
			fileContent, err := readDependencyJSON(dependencyFileName)
			if err != nil {
				return nil, err
			}
			// only entries with dependencies are added, those without dependencies don't need to be added
			for _, entryValue := range fileContent {
				dependenciesGraph.AddDependency(name, entryValue)
			}
		}
	}
	return dependenciesGraph.GetAllDependenciesSorted(), nil
}

func ReadJSONTemplate(templateBasePath string, templateName string) (map[string]interface{}, error) {
	return readJSON(filepath.Join(templateBasePath, templateName, templateName+".json"))
}

func ReadJSONTemplateParameters(templateBasePath string, templateName string) (map[string]interface{}, error) {
	return readJSON(filepath.Join(templateBasePath, templateName, templateName+".parameters.json"))
}
