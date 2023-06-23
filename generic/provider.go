package generic

import (
	"strings"

	"github.com/cockscomb/cel2sql/sqltypes"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type FieldType string

const (
	StringFieldType    FieldType = "string"
	BytesFieldType     FieldType = "bytes"
	BooleanFieldType   FieldType = "boolean"
	IntegerFieldType   FieldType = "integer"
	FloatFieldType     FieldType = "float"
	TimestampFieldType FieldType = "timestamp"
	RecordFieldType    FieldType = "record"
	DateFieldType      FieldType = "date"
	TimeFieldType      FieldType = "time"
	DateTimeFieldType  FieldType = "datetime"
)

type Schema []*FieldSchema

type FieldSchema struct {
	Name string
	Type FieldType
	// Array type
	Repeated bool
	// TODO: support for nested record types
	Schema Schema
}

type typeProvider struct {
	schemas map[string]Schema
}

func NewTypeProvider(schemas map[string]Schema) *typeProvider {
	return &typeProvider{schemas: schemas}
}

func (p *typeProvider) EnumValue(enumName string) ref.Val {
	return types.NewErr("unknown enum name '%s'", enumName)
}

func (p *typeProvider) FindIdent(identName string) (ref.Val, bool) {
	return nil, false
}

func (p *typeProvider) findSchema(typeName string) (Schema, bool) {
	typeNames := strings.Split(typeName, ".")
	schema, found := p.schemas[typeNames[0]]
	if !found {
		return nil, false
	}
	for _, tn := range typeNames[1:] {
		var s *Schema
		for _, fieldSchema := range schema {
			if fieldSchema.Name == tn {
				s = &fieldSchema.Schema
				break
			}
		}
		if s == nil {
			return nil, false
		}
		schema = *s
	}
	return schema, true
}

func (p *typeProvider) FindType(typeName string) (*exprpb.Type, bool) {
	_, found := p.findSchema(typeName)
	if !found {
		return nil, false
	}
	return decls.NewTypeType(decls.NewObjectType(typeName)), true
}

func (p *typeProvider) FindFieldType(messageType string, fieldName string) (*ref.FieldType, bool) {
	schema, found := p.findSchema(messageType)
	if !found {
		return nil, false
	}
	var field *FieldSchema
	for _, fieldSchema := range schema {
		if fieldSchema.Name == fieldName {
			field = fieldSchema
			break
		}
	}
	if field == nil {
		return nil, false
	}
	var typ *exprpb.Type
	switch field.Type {
	case StringFieldType:
		typ = decls.String
	case BytesFieldType:
		typ = decls.Bytes
	case BooleanFieldType:
		typ = decls.Bool
	case IntegerFieldType:
		typ = decls.Int
	case FloatFieldType:
		typ = decls.Double
	case TimestampFieldType:
		typ = decls.Timestamp
	case RecordFieldType:
		typ = decls.NewObjectType(strings.Join([]string{messageType, fieldName}, "."))
	case DateFieldType:
		typ = sqltypes.Date
	case TimeFieldType:
		typ = sqltypes.Time
	case DateTimeFieldType:
		typ = sqltypes.DateTime
	}
	if field.Repeated {
		typ = decls.NewListType(typ)
	}
	return &ref.FieldType{
		Type: typ,
	}, true
}

func (p *typeProvider) NewValue(typeName string, fields map[string]ref.Val) ref.Val {
	return types.NewErr("unknown type '%s'", typeName)
}

var _ ref.TypeProvider = new(typeProvider)
