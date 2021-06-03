package schemagen

import (
	"fmt"
	"github.com/actgardner/gogen-avro/v8/generator"
	"github.com/actgardner/gogen-avro/v8/generator/flat"
	"github.com/actgardner/gogen-avro/v8/parser"
	"github.com/actgardner/gogen-avro/v8/resolver"
	"github.com/riferrei/srclient"
	"log"
	"os"
	"strconv"
	"strings"
)

type GenerationRequest struct {
	Package   string            `yaml:"package"`
	TargetDir string            `yaml:"targetDir"`
	Schemas   []SchemaReference `yaml:"schemas"`
}

type SchemaReference struct {
	Subject string `yaml:"subject" valid:"-"`
	Version string `yaml:"version" valid:"-"`
}

func NewGenerator(registryId string, client srclient.ISchemaRegistryClient) *Generator {
	return &Generator{
		registryUrl: registryId,
		client:      client,
	}
}

type Generator struct {
	registryUrl string
	client      srclient.ISchemaRegistryClient
}

func (g *Generator) Generate(req GenerationRequest) error {
	log.Printf("Generating to package %q in directory %q\n", req.Package, req.TargetDir)
	for _, s := range req.Schemas {
		log.Printf(" - %s@%s", s.Subject, s.Version)
	}

	schemas, err := g.fetchSchemas(req.Schemas...)
	if err != nil {
		return err
	}

	// -- actually generate the data
	return g.compile(req.Package, req.TargetDir, schemas...)
}

func (g *Generator) fetchSchemas(refs ...SchemaReference) ([]*Schema, error) {
	schemas := make([]*Schema, len(refs))

	for idx, ref := range refs {
		s, err := g.fetchSchema(ref.Subject, ref.Version)
		if err != nil {
			return nil, err
		}

		schemas[idx] = s
	}

	return schemas, nil
}

func (g *Generator) fetchSchema(subject string, version string) (*Schema, error) {
	// -- make it explicit we change an empty string to the latest version
	if version == "" {
		version = "latest"
	}

	if version == "latest" {
		s, err := g.client.GetLatestSchemaWithArbitrarySubject(subject)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve the latest schema for %s: %v", subject, err)
		}

		return &Schema{
			Schema:  []byte(s.Schema()),
			Subject: subject,
			Version: s.Version(),
			ID:      s.ID(),
		}, nil
	} else {
		v, err := strconv.Atoi(version)
		if err != nil {
			return nil, fmt.Errorf("version %q is not valid: %v", version, err)
		}

		s, err := g.client.GetSchemaByVersionWithArbitrarySubject(subject, v)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve schema version %d for %s: %v", v, subject, err)
		}

		return &Schema{
			Schema:  []byte(s.Schema()),
			Subject: subject,
			Version: s.Version(),
			ID:      s.ID(),
		}, nil
	}
}

func (g *Generator) compile(gopkg string, targetDir string, schemas ...*Schema) error {
	pkg := generator.NewPackage(gopkg, g.generateFileHeader(schemas...))
	namespace := parser.NewNamespace(true)
	gen := flat.NewFlatPackageGenerator(pkg, false)

	for _, s := range schemas {
		if _, err := namespace.TypeForSchema(s.Schema); err != nil {
			return fmt.Errorf("error decoding schema: %v", err)
		}
	}

	for _, def := range namespace.Roots {
		if err := resolver.ResolveDefinition(def, namespace.Definitions); err != nil {
			return fmt.Errorf("error resolving definition for type %q: %v", def.Name(), err)
		}
	}

	for _, def := range namespace.Roots {
		if err := gen.Add(def); err != nil {
			return fmt.Errorf("error generating code for schema: %v", err)
		}
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0777); err != nil {
			return fmt.Errorf("unable to create output directory: %v", err)
		}
	}

	if err := pkg.WriteFiles(targetDir); err != nil {
		return fmt.Errorf("error writing source files to directory %q: %v", targetDir, err)
	}

	return nil
}

func (g *Generator) generateFileHeader(schemas ...*Schema) string {
	const fileComment = `/*
 *              !!! THIS FILE HAS BEEN GENERATED !!!
 *
 * Please do not modify since future regenerations will just overwrite your changes.
 * This file is based on the following avro schemas from %s:
%s
 */`
	var sourceBlock []string
	for _, s := range schemas {
		sourceBlock = append(sourceBlock, fmt.Sprintf(" *    - %s@%d", s.Subject, s.Version))
	}

	return fmt.Sprintf(fileComment, g.registryUrl, strings.Join(sourceBlock, "\n"))
}
