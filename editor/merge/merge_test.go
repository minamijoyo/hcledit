package merge

import (
	"bytes"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestMergeBlocks(t *testing.T) {
	tests := []struct {
		name        string
		sourceHCL   string
		targetHCL   string
		ignoreAttrs []string
		setAttrs    map[string]string
		expectedHCL string
	}{
		{
			name: "update existing block and preserve target env",
			sourceHCL: `
module "network" {
  version = "v2.0.0"
  env     = "dev"
}`,
			targetHCL: `
module "network" {
  version = "v1.0.0"
  env     = "prod"
}`,
			ignoreAttrs: []string{"env"},
			setAttrs:    nil,
			expectedHCL: `
module "network" {
  version = "v2.0.0"
  env     = "prod"
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceFile, diags := hclwrite.ParseConfig([]byte(tt.sourceHCL), "source.tf", hcl.InitialPos)
			if diags.HasErrors() {
				t.Fatalf("failed to parse source: %s", diags.Error())
			}

			targetFile, diags := hclwrite.ParseConfig([]byte(tt.targetHCL), "target.tf", hcl.InitialPos)
			if diags.HasErrors() {
				t.Fatalf("failed to parse target: %s", diags.Error())
			}

			Blocks(sourceFile, targetFile, "module", tt.ignoreAttrs, tt.setAttrs)

			resultBytes := hclwrite.Format(targetFile.Bytes())

			expectedFile, diags := hclwrite.ParseConfig([]byte(tt.expectedHCL), "expected.tf", hcl.InitialPos)
			if diags.HasErrors() {
				t.Fatalf("failed to parse expected: %s", diags.Error())
			}
			expectedBytes := hclwrite.Format(expectedFile.Bytes())

			if !bytes.Equal(resultBytes, expectedBytes) {
				t.Errorf("Merge output mismatch.\nGot:\n%s\nExpected:\n%s", string(resultBytes), string(expectedBytes))
			}
		})
	}
}
