package main

import (
	"github.com/calmera/schemagen"
	"github.com/riferrei/srclient"
	"log"
	"os"
)

const (
	EnvTargetDir        = "SGEN_TARGET_DIR"
	EnvRegistry         = "SGEN_REGISTRY"
	EnvRegistryUser     = "SGEN_REGISTRY_USERNAME"
	EnvRegistryPassword = "SGEN_REGISTRY_PASSWORD"
)

func main() {
	targetDir := os.Getenv(EnvTargetDir)
	if targetDir == "" {
		log.Panic("no target directory specified")
	}

	registry := os.Getenv(EnvRegistry)
	if registry == "" {
		log.Panic("no registry url specified")
	}

	rUser := os.Getenv(EnvRegistryUser)
	rPwd := os.Getenv(EnvRegistryPassword)

	client := srclient.CreateSchemaRegistryClient(registry)

	if rUser != "" && rPwd != "" {
		client.SetCredentials(rUser, rPwd)
	}

	gen := schemagen.NewGenerator(registry, client)

	req := schemagen.GenerationRequest{
		Package:   "mypkg",
		TargetDir: targetDir,
		Schemas: []schemagen.SchemaReference{
			{Subject: "country", Version: "latest"},
			{Subject: "city", Version: "latest"},
		},
	}

	if err := gen.Generate(req); err != nil {
		log.Panic(err)
	}
}
