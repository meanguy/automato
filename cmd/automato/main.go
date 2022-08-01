package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/meanguy/automato/internal/mem"
	"github.com/meanguy/automato/internal/vm"
)

type Args struct {
	Debug bool
}

func Execute(opts *Args, cmd *cobra.Command, args []string) error {
	vmOpts := []vm.VMOption{}
	if opts.Debug {
		vmOpts = append(vmOpts, vm.EnableDebug())
	}
	automato := vm.NewVM(vmOpts...)

	if len(args) == 0 {
		return REPL(automato)
	}

	return automato.InterpretChunk(&mem.Chunk{})
}

func REPL(vm *vm.VM) error {
	scan := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprintf(os.Stderr, "> ")

		if !scan.Scan() {
			return scan.Err()
		}

		if err := vm.Interpret(scan.Text()); err != nil {
			return err
		}
	}
}

func main() {
	opts := Args{}
	cmd := &cobra.Command{
		Use:   "adb",
		Short: "automation programming language, based on Crafting Interpreters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Execute(&opts, cmd, args)
		},
	}

	cmd.PersistentFlags().BoolVar(&opts.Debug, "debug", false, "enable debug tracing")

	if err := cmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())

		os.Exit(1)
	}
}
