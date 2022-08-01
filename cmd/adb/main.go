package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/meanguy/automato/internal/vm"
)

type Args struct {
	Debug bool
}

func Execute(args *Args, cmd *cobra.Command, _ []string) error {
	opts := []vm.VMOption{}
	if args.Debug {
		opts = append(opts, vm.EnableDebug())
	}

	runtime := vm.NewVM(opts...)
	if err := runtime.Interpret("5*5+7-2"); err != nil {
		return err
	}

	return nil
}

func main() {
	args := Args{}
	cmd := &cobra.Command{
		Use:   "adb",
		Short: "automation debugger",
		RunE: func(cmd *cobra.Command, cargs []string) error {
			return Execute(&args, cmd, cargs)
		},
	}

	cmd.PersistentFlags().BoolVar(&args.Debug, "debug", false, "enable debug tracing")

	if err := cmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())

		os.Exit(1)
	}
}
