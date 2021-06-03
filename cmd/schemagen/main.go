package main

import (
	"github.com/calmera/schemagen"
	"github.com/riferrei/srclient"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	schemagen.GenerationRequest
	Registry string `yaml:"registry" valid:"required"`
}

func main() {
	var (
		cfgFile  string
		cfg      Config
		username string
		password string
	)

	rootCmd := &cobra.Command{
		Use: "schemagen",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cfgFile, &cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client := srclient.CreateSchemaRegistryClient(cfg.Registry)

			if username != "" && password != "" {
				client.SetCredentials(username, password)
			}

			gen := schemagen.NewGenerator(cfg.Registry, client)

			return gen.Generate(cfg.GenerationRequest)
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default .schemagen.yaml in the current directory")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "the username for the schema registry")
	rootCmd.PersistentFlags().StringVar(&username, "password", "", "the password for the schema registry")

	if err := bindFlags(rootCmd.PersistentFlags(), viper.GetViper()); err != nil {
		log.Panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
	}
}

func bindFlags(fs *pflag.FlagSet, v *viper.Viper) error {
	var err error

	fs.VisitAll(func(f *pflag.Flag) {
		err = v.BindPFlag(strings.Replace(f.Name, "-", "_", -1), f)
		if err != nil {
			return
		}
	})

	return err
}

func initConfig(cfgFile string, cfg *Config) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".schemagen")
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(cfg)
}
