package clientgen

import (
	"bytes"
	"fmt"
	"go/types"

	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"golang.org/x/xerrors"
)

type Source struct {
	schema          *ast.Schema
	queryDocument   *ast.QueryDocument
	sourceGenerator *SourceGenerator
}

func NewSource(schema *ast.Schema, queryDocument *ast.QueryDocument, sourceGenerator *SourceGenerator) *Source {
	return &Source{
		schema:          schema,
		queryDocument:   queryDocument,
		sourceGenerator: sourceGenerator,
	}
}

type Fragment struct {
	Name string
	Type types.Type
}

func (s *Source) Fragments() ([]*Fragment, error) {
	fragments := make([]*Fragment, 0, len(s.queryDocument.Fragments))
	for _, fragment := range s.queryDocument.Fragments {
		responseFields := s.sourceGenerator.NewResponseFields(fragment.SelectionSet)
		if s.sourceGenerator.cfg.Models.Exists(fragment.Name) {
			return nil, xerrors.New(fmt.Sprintf("%s is duplicated", fragment.Name))
		}

		fragment := &Fragment{
			Name: fragment.Name,
			Type: responseFields.StructType(),
		}

		fragments = append(fragments, fragment)
	}

	for _, fragment := range fragments {
		name := fragment.Name
		s.sourceGenerator.cfg.Models.Add(
			name,
			fmt.Sprintf("%s.%s", s.sourceGenerator.client.Pkg(), templates.ToGo(name)),
		)
	}

	return fragments, nil
}

type Operation struct {
	Name                string
	ResponseStructName  string
	Operation           string
	Args                []*Argument
	VariableDefinitions ast.VariableDefinitionList
}

func NewOperation(operation *ast.OperationDefinition, queryDocument *ast.QueryDocument, args []*Argument) *Operation {
	return &Operation{
		Name:                operation.Name,
		ResponseStructName:  getResponseStructName(operation),
		Operation:           queryString(queryDocument),
		Args:                args,
		VariableDefinitions: operation.VariableDefinitions,
	}
}

func (s *Source) Operations(queryDocuments []*ast.QueryDocument) []*Operation {
	operations := make([]*Operation, 0, len(s.queryDocument.Operations))

	queryDocumentsMap := queryDocumentMapByOperationName(queryDocuments)
	operationArgsMap := s.operationArgsMapByOperationName()
	for _, operation := range s.queryDocument.Operations {
		queryDocument := queryDocumentsMap[operation.Name]
		args := operationArgsMap[operation.Name]
		operations = append(operations, NewOperation(
			operation,
			queryDocument,
			args,
		))
	}

	return operations
}

func (s *Source) operationArgsMapByOperationName() map[string][]*Argument {
	operationArgsMap := make(map[string][]*Argument)
	for _, operation := range s.queryDocument.Operations {
		operationArgsMap[operation.Name] = s.sourceGenerator.OperationArguments(operation.VariableDefinitions)
	}

	return operationArgsMap
}

func queryDocumentMapByOperationName(queryDocuments []*ast.QueryDocument) map[string]*ast.QueryDocument {
	queryDocumentMap := make(map[string]*ast.QueryDocument)
	for _, queryDocument := range queryDocuments {
		operation := queryDocument.Operations[0]
		queryDocumentMap[operation.Name] = queryDocument
	}

	return queryDocumentMap
}

func queryString(queryDocument *ast.QueryDocument) string {
	var buf bytes.Buffer
	astFormatter := formatter.NewFormatter(&buf)
	astFormatter.FormatQueryDocument(queryDocument)

	return buf.String()
}

type OperationResponse struct {
	Name string
	Type types.Type
}

func (s *Source) OperationResponses() ([]*OperationResponse, error) {
	operationResponse := make([]*OperationResponse, 0, len(s.queryDocument.Operations))
	for _, operation := range s.queryDocument.Operations {
		responseFields := s.sourceGenerator.NewResponseFields(operation.SelectionSet)
		name := getResponseStructName(operation)
		if s.sourceGenerator.cfg.Models.Exists(name) {
			return nil, xerrors.New(fmt.Sprintf("%s is duplicated", name))
		}
		operationResponse = append(operationResponse, &OperationResponse{
			Name: name,
			Type: responseFields.StructType(),
		})
	}

	for _, operationResponse := range operationResponse {
		name := operationResponse.Name
		s.sourceGenerator.cfg.Models.Add(
			name,
			fmt.Sprintf("%s.%s", s.sourceGenerator.client.Pkg(), templates.ToGo(name)),
		)
	}

	return operationResponse, nil
}

type Query struct {
	Name string
	Type types.Type
}

func (s *Source) Query() (*Query, error) {
	fields, err := s.sourceGenerator.NewResponseFieldsByDefinition(s.schema.Query)
	if err != nil {
		return nil, xerrors.Errorf("generate failed for query struct type : %w", err)
	}

	return &Query{
		Name: s.schema.Query.Name,
		Type: fields.StructType(),
	}, nil
}

type Mutation struct {
	Name string
	Type types.Type
}

func (s *Source) Mutation() (*Mutation, error) {
	fields, err := s.sourceGenerator.NewResponseFieldsByDefinition(s.schema.Mutation)
	if err != nil {
		return nil, xerrors.Errorf("generate failed for mutation struct type : %w", err)
	}

	return &Mutation{
		Name: s.schema.Mutation.Name,
		Type: fields.StructType(),
	}, nil
}

func getResponseStructName(operation *ast.OperationDefinition) string {
	if operation.Operation == ast.Mutation {
		return fmt.Sprintf("%sPayload", operation.Name)
	}

	return operation.Name
}
