package merge

import (
	"bytes"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestBlocks(t *testing.T) {
	tests := []struct {
		name        string
		sourceHCL   string
		targetHCL   string
		blockType   string
		ignoreAttrs []string
		setAttrs    map[string]string
		expectedHCL string
	}{
		// ---------------------------------------------------------
		// HAPPY PATHS
		// ---------------------------------------------------------
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
			blockType:   "module",
			ignoreAttrs: []string{"env"},
			expectedHCL: `
module "network" {
  version = "v2.0.0"
  env     = "prod"
}`,
		},
		{
			name: "add missing block and inject target env",
			sourceHCL: `module "db" {
  version = "v1.5.0"
  env     = "dev"
}`,
			targetHCL: `module "network" {
  version = "v1.0.0"
}`,
			blockType:   "module",
			ignoreAttrs: []string{"env"},
			setAttrs:    map[string]string{"env": "prod"},
			expectedHCL: `module "network" {
  version = "v1.0.0"
}
module "db" {
  version = "v1.5.0"
  env     = "prod"
}`,
		},
		// ---------------------------------------------------------
		// EDGE CASES
		// ---------------------------------------------------------
		{
			name:      "edge case: empty source file leaves target unmodified",
			sourceHCL: ``,
			targetHCL: `
module "network" {
  version = "v1.0.0"
}`,
			blockType: "module",
			expectedHCL: `
module "network" {
  version = "v1.0.0"
}`,
		},
		{
			name: "edge case: ignores unrelated block types in source",
			sourceHCL: `
resource "aws_s3_bucket" "b" {
  bucket = "my-bucket"
}

module "network" {
  version = "v2.0.0"
}`,
			targetHCL: `
module "network" {
  version = "v1.0.0"
}`,
			blockType: "module", // Only merge modules, ignore the s3 bucket
			expectedHCL: `
module "network" {
  version = "v2.0.0"
}`,
		},
		{
			name: "edge case: handles blocks with multiple labels (resources)",
			sourceHCL: `
resource "aws_instance" "web" {
  ami  = "ami-new"
  type = "t3.large"
}`,
			targetHCL: `
resource "aws_instance" "web" {
  ami  = "ami-old"
  type = "t2.micro"
}

resource "aws_instance" "db" {
  ami = "ami-old"
}`,
			blockType:   "resource",
			ignoreAttrs: []string{"type"}, // Protect the target's instance type
			expectedHCL: `
resource "aws_instance" "web" {
  ami  = "ami-new"
  type = "t2.micro"
}

resource "aws_instance" "db" {
  ami = "ami-old"
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

			// Execute the engine using the dynamic blockType
			Blocks(sourceFile, targetFile, tt.blockType, tt.ignoreAttrs, tt.setAttrs)

			// Format bytes to normalize HCL whitespace before comparing
			resultBytes := hclwrite.Format(targetFile.Bytes())

			expectedFile, diags := hclwrite.ParseConfig([]byte(tt.expectedHCL), "expected.tf", hcl.InitialPos)
			if diags.HasErrors() {
				t.Fatalf("failed to parse expected: %s", diags.Error())
			}
			expectedBytes := hclwrite.Format(expectedFile.Bytes())

			if !bytes.Equal(bytes.TrimSpace(resultBytes), bytes.TrimSpace(expectedBytes)) {
				t.Errorf("\n=== Output Mismatch ===\nGot:\n%s\nExpected:\n%s", string(resultBytes), string(expectedBytes))
			}
		})
	}
}
