package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	root := &cobra.Command{
		Use:           "isc-rag",
		Short:         "Tools for the ISC-RAG research project",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.AddCommand(newBootstrapCommand())
	root.AddCommand(newMigrateCommand())
	root.AddCommand(newIngestCommand())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
