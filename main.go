package main

import (
	"fmt"
	"os"

	"github.com/minamijoyo/hcledit/cmd"
	"github.com/spf13/viper"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		if viper.GetBool("debug") {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		os.Exit(1)
	}
}
