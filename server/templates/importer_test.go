package templates

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test DiscoverTemplateOrder with a dir path with no children dirs
func TestDiscoverTemplateOrder(t *testing.T) {
	templatePathWithNoTemplates, _ := os.Getwd()
	result, _ := DiscoverTemplateOrder(templatePathWithNoTemplates)
	assert.Empty(t, result, "returned ordered templates should be empty if no templates are within the supplied path.")
}
