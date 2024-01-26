package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestAccProviderFrameworkSatisfaction(t *testing.T) {
	t.Parallel()

	var _ provider.Provider = &GoPackagerProvider{}
}
