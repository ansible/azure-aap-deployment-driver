package templates

import (
	"os"
	"path/filepath"
)

func DiscoverTemplateOrder(templateBasePath string) ([][]string, error) {
	dependenciesGraph := NewDependencyGraph()

	templateDirEntries, err := os.ReadDir(templateBasePath)

	if len(templateDirEntries) == 0 {
		return dependenciesGraph.GetAllDependenciesSorted(), nil
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
