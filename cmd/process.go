package cmd

import (
	"github.com/gkwa/fullmaine/core"
	"github.com/spf13/cobra"
)

var startNum int

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process files in testdata directory",
	Long:  `Recurse into testdata directory, find and parse filenames, and create new ones with incremented numbers.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		logger.Info("Running process command", "startNum", startNum)
		fileProcessor := core.NewFileProcessor()
		err := fileProcessor.ProcessFiles("testdata", startNum, logger)
		if err != nil {
			logger.Error(err, "Failed to process files")
		}
	},
}

func init() {
	rootCmd.AddCommand(processCmd)
	processCmd.Flags().IntVarP(&startNum, "number", "n", 100, "Starting number for file naming")
}
