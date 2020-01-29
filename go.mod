module github.com/minamijoyo/hcledit

go 1.13

require (
	github.com/goreleaser/goreleaser v0.123.3
	github.com/hashicorp/hcl/v2 v2.3.1-0.20200103191330-7990d6e9a2c9
	github.com/hashicorp/logutils v1.0.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/lint v0.0.0-20191125180803-fdd1cda4f05f
)

replace github.com/hashicorp/hcl/v2 => github.com/minamijoyo/hcl/v2 v2.0.1-0.20200129060650-436ebf7cac4f
