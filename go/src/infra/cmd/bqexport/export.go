// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"infra/libs/bqschema/tabledef"

	"go.chromium.org/luci/common/errors"

	"github.com/golang/protobuf/proto"

	"golang.org/x/net/context"
)

func sanitizeComment(v string) string {
	prevSpace := false
	v = strings.Map(func(r rune) rune {
		switch {
		case unicode.IsSpace(r):
			if prevSpace {
				return -1
			}
			prevSpace = true
			return ' '
		case unicode.IsPrint(r):
			prevSpace = false
			return r
		default:
			return -1
		}
	}, v)
	return strings.TrimSpace(v)
}

var structTemplate = template.Must(template.New("").
	Funcs(template.FuncMap{
		"sanitizeComment": sanitizeComment,
	}).
	Parse(`
// THIS FILE IS AUTOGENERATED. DO NOT MODIFY.

package {{.Package}}

import pb "infra/libs/bqschema/tabledef"
{{range .Imports -}}
import {{.Alias}} "{{.Path}}"
{{end -}}

// {{.TableDefName}} is the TableDef for the
// "{{.DatasetName}}" dataset's "{{.TableID}}" table.
var {{.TableDefName}} = &pb.TableDef{
	Dataset: pb.TableDef_{{.Dataset}},
	TableId: "{{.TableID}}",
}
{{range .Structs -}}

// {{sanitizeComment .Comment}}
type {{.Name}} struct {
{{range .Fields -}}
	{{if .Description -}}
	// {{sanitizeComment .Description}}
	{{end -}}
	{{.Name}} {{.Type}} ` + "`bigquery:\"{{.FieldName}}\"`" + `

{{end -}}
}
{{end -}}
`))

// LoadTableDef loads a TableDef text protobuf.
func LoadTableDef(path string) (*tabledef.TableDef, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tdef tabledef.TableDef
	return &tdef, proto.UnmarshalText(string(content), &tdef)
}

// Export exports the specified TableDef as a set of Go structs for the
// specific package. The primary exported struct will be named "structName",
// and the associated TableDef will be named "<structName>Table".
//
// The resulting generated Go source will be written to "out".
func Export(ctx context.Context, td *tabledef.TableDef, packageName, structName, out string) error {
	tableDefName := structName + "Table"

	type packageImport struct {
		Alias string
		Path  string
	}

	type params struct {
		Package      string
		Imports      []*packageImport
		TableDefName string
		DatasetName  string
		Dataset      string
		TableID      string
		Structs      []*structDef
	}

	a := analyzer{
		structBase: structName,
	}
	a.ensureStruct("", fmt.Sprintf("the schema for %q.", tableDefName), td.Fields)

	p := params{
		Package:      packageName,
		Imports:      make([]*packageImport, 0, len(a.packages)),
		TableDefName: tableDefName,
		DatasetName:  td.Dataset.ID(),
		Dataset:      td.Dataset.String(),
		TableID:      td.TableId,
		Structs:      a.structDefs,
	}

	// Add sorted package imports.
	for alias, path := range a.packages {
		p.Imports = append(p.Imports, &packageImport{alias, path})
	}
	sort.Slice(p.Imports, func(i, j int) bool { return p.Imports[i].Alias < p.Imports[j].Alias })

	// Generate the template.
	fd, err := os.Create(out)
	if err != nil {
		return errors.Annotate(err, "could not open output file %q", out).Err()
	}
	err = structTemplate.Execute(fd, &p)
	if err != nil {
		_ = fd.Close()
		return errors.Annotate(err, "could not generate output file").Err()
	}
	if err := fd.Close(); err != nil {
		return errors.Annotate(err, "could not close output file").Err()
	}

	cmd := exec.CommandContext(ctx, "gofmt", "-s", "-w", out)
	if err := cmd.Run(); err != nil {
		return errors.Annotate(err, "could not format output file").Err()
	}

	return nil
}

type structDef struct {
	Comment string
	Name    string
	Fields  []*fieldEntry
}

type fieldEntry struct {
	Name        string
	Type        string
	FieldName   string
	Description string
}

type analyzer struct {
	structBase string
	packages   map[string]string

	structs    map[string]*structDef
	structDefs []*structDef
}

func (a *analyzer) ensurePackage(name, path string) {
	if _, ok := a.packages[name]; ok {
		return
	}
	if a.packages == nil {
		a.packages = make(map[string]string)
	}
	a.packages[name] = path
}

func (a *analyzer) ensureStruct(fieldName, comment string, schema []*tabledef.FieldSchema) string {
	// Create "key", a deterministic rendering of "schema".
	keyBaser := tabledef.FieldSchema{
		Schema: schema,
	}
	schemaBytes, err := proto.Marshal(&keyBaser)
	if err != nil {
		panic(err)
	}
	key := string(schemaBytes)

	// Is an identical struct already defined?
	if sd := a.structs[key]; sd != nil {
		return sd.Name
	}

	name := a.structBase
	if fieldName != "" {
		name = fmt.Sprintf("%s_%s", name, toCamelCase(fieldName))
		if comment == "" {
			comment = fmt.Sprintf("a record for the %q field.", fieldName)
		}
	}

	if comment != "" {
		comment = fmt.Sprintf("%s is %s", name, comment)
	}

	// Define a new struct type.
	sd := structDef{
		Name:    name,
		Fields:  make([]*fieldEntry, len(schema)),
		Comment: comment,
	}
	for i, field := range schema {
		sd.Fields[i] = a.makeFieldEntry(field)
	}
	if a.structs == nil {
		a.structs = make(map[string]*structDef)
	}
	a.structs[key] = &sd
	a.structDefs = append(a.structDefs, &sd)
	return sd.Name
}

func (a *analyzer) makeFieldEntry(fs *tabledef.FieldSchema) *fieldEntry {
	var typ string
	switch fs.Type {
	case tabledef.Type_STRING:
		typ = "string"
	case tabledef.Type_BYTES:
		typ = "[]byte"
	case tabledef.Type_INTEGER:
		typ = "int64"
	case tabledef.Type_FLOAT:
		typ = "float64"
	case tabledef.Type_BOOLEAN:
		typ = "bool"
	case tabledef.Type_TIMESTAMP:
		a.ensurePackage("time", "time")
		typ = "time.Time"
	case tabledef.Type_RECORD:
		typ = "*" + a.ensureStruct(fs.Name, "", fs.Schema)
	case tabledef.Type_DATE:
		a.ensurePackage("civil", "cloud.google.com/go/civil")
		typ = "civil.Date"
	case tabledef.Type_TIME:
		a.ensurePackage("civil", "cloud.google.com/go/civil")
		typ = "civil.Time"
	case tabledef.Type_DATETIME:
		a.ensurePackage("civil", "cloud.google.com/go/civil")
		typ = "civil.DateTime"
	}
	if fs.IsRepeated {
		typ = "[]" + typ
	}

	return &fieldEntry{
		Name:        toCamelCase(fs.Name),
		Type:        typ,
		FieldName:   fs.Name,
		Description: fs.Description,
	}
}
