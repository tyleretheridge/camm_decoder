package cmd

import (
	"bytes"
	"camm_extractor/internal/containers"
	"camm_extractor/internal/decoder"
	"camm_extractor/internal/writers"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return decodeCallback(args[0])
	},
}

func decodeCallback(binaryFile string) error {
	// Binary Stream
	binaryData, err := os.ReadFile(binaryFile)
	if err != nil {
		fmt.Printf("error opening binary file: %s", err)

	}
	buf := bytes.NewBuffer(binaryData)
	if err != nil {
		fmt.Printf("error reading binary stream: %s", err)
	}
	// Init Packet Container
	cammStream := containers.NewBaseContainer()
	// Decode binary into container
	err = decoder.DecodeCAMMData(buf, cammStream)
	if err != nil {
		fmt.Printf("error decoding binary stream: %s", err)
	}
	// Create a writer
	w := writers.TerminalWriter{}
	// Write
	err = w.Write(cammStream)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mycli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mycli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mycli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
