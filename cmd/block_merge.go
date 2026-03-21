package cmd

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor/merge"
	"github.com/spf13/cobra"
)

type blockMergeOption struct {
	sourceFile  string
	targetFile  string
	ignoreAttrs []string
	setAttrs    map[string]string
}

func newBlockMergeCmd() *cobra.Command {
	opt := &blockMergeOption{}

	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge module blocks from a source file into a target file",
		Long: `Intelligently syncs blocks (e.g., modules) between two HCL files.
Updates existing blocks and safely injects new blocks while respecting environment-specific variables.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runBlockMerge(cmd, opt)
		},
	}

	cmd.Flags().StringVarP(&opt.sourceFile, "source", "s", "", "Source HCL file (required)")
	cmd.Flags().StringVarP(&opt.targetFile, "target", "t", "", "Target HCL file (required)")
	cmd.Flags().StringSliceVarP(&opt.ignoreAttrs, "ignore-attr", "i", []string{}, "Attributes to ignore during merge (e.g., env)")
	cmd.Flags().StringToStringVarP(&opt.setAttrs, "set-attr", "a", map[string]string{}, "Attributes to force-set on new blocks (e.g., env=prod)")

	_ = cmd.MarkFlagRequired("source")
	_ = cmd.MarkFlagRequired("target")

	return cmd
}

func runBlockMerge(cmd *cobra.Command, opt *blockMergeOption) error {
	srcBytes, err := os.ReadFile(opt.sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	src, diags := hclwrite.ParseConfig(srcBytes, opt.sourceFile, hcl.InitialPos)
	if diags.HasErrors() {
		return fmt.Errorf("source HCL parsing errors: %s", diags.Error())
	}

	tgtBytes, err := os.ReadFile(opt.targetFile)
	if err != nil {
		return fmt.Errorf("failed to read target file: %w", err)
	}
	tgt, diags := hclwrite.ParseConfig(tgtBytes, opt.targetFile, hcl.InitialPos)
	if diags.HasErrors() {
		return fmt.Errorf("target HCL parsing errors: %s", diags.Error())
	}

	merge.Blocks(src, tgt, "module", opt.ignoreAttrs, opt.setAttrs)

	formattedBytes := hclwrite.Format(tgt.Bytes())

	// Mode 0644 is standard for text files: rw-r--r--
	// nolint:gosec // 0644 is standard and expected for Terraform text files
	if err := os.WriteFile(opt.targetFile, formattedBytes, 0644); err != nil {
		return fmt.Errorf("failed to write target file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Successfully merged '%s' into '%s'\n", opt.sourceFile, opt.targetFile)
	return nil
}
