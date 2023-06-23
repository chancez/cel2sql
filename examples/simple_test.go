package examples

import (
	"fmt"
	"log"

	"github.com/cockscomb/cel2sql"
	"github.com/cockscomb/cel2sql/generic"
	"github.com/cockscomb/cel2sql/sqltypes"
	"github.com/google/cel-go/cel"
)

var employeeSchema = &generic.Schema{
	TableName:    "employees",
	VariableName: "employee",
	ObjectType:   "Employee",
	Fields: []*generic.FieldSchema{
		{
			Name: "name",
			Type: generic.StringFieldType,
		},
		{
			Name: "hired_at",
			Type: generic.TimestampFieldType,
		},
	},
}

func ExampleSimple() {
	sqlTypeProvider := generic.NewTypeProvider([]*generic.Schema{
		employeeSchema,
	})
	env, _ := cel.NewEnv(
		cel.CustomTypeProvider(sqlTypeProvider),
		sqltypes.SQLTypeDeclarations,
		cel.Declarations(sqlTypeProvider.Declarations()...),
	)

	// Convert CEL to SQL
	ast, iss := env.Compile(`employee.name == "John Doe" && employee.hired_at >= current_timestamp() - duration("24h")`)
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	sqlCondition, err := cel2sql.ConvertWithMapper(ast, sqlTypeProvider)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(sqlCondition)
	// Output: `employees`.`name` = 'John Doe' AND `employees`.`hired_at` >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 24 HOUR)
}
