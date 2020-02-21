package editor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Filter is an interface which reads HCL and writes HCL
type Filter interface {
	// Filter reads HCL and writes HCL
	Filter(*hclwrite.File) (*hclwrite.File, error)
}
