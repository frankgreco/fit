package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	config = struct {
		file string
	}{}

	rootCmd = &cobra.Command{
		Use: "fit",
	}

	decodeCmd = &cobra.Command{
		Use:  "decode",
		RunE: decode,
	}
)

func decode(cmd *cobra.Command, args []string) error {
	file, err := Decode(config.file)
	if err != nil {
		return err
	}

	data, err := json.Marshal(file)
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func init() {
	decodeCmd.Flags().StringVarP(&config.file, "file", "f", "", "location of .fit file")
	decodeCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(decodeCmd)
}
