package examples

import (
	"fmt"
	"log"

	"github.com/cockscomb/cel2sql"
	"github.com/cockscomb/cel2sql/generic"
	"github.com/cockscomb/cel2sql/sqltypes"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

var employeeSchema = generic.Schema{
	{
		Name: "name",
		Type: generic.StringFieldType,
	},
	{
		Name: "hired_at",
		Type: generic.TimestampFieldType,
	},
}

func ExampleSimple() {
	env, _ := cel.NewEnv(
		cel.CustomTypeProvider(generic.NewTypeProvider(map[string]generic.Schema{
			"Employee": employeeSchema,
		})),
		sqltypes.SQLTypeDeclarations,
		cel.Declarations(
			decls.NewVar("employee", decls.NewObjectType("Employee")),
		),
	)

	// Convert CEL to SQL
	ast, iss := env.Compile(`employee.name == "John Doe" && employee.hired_at >= current_timestamp() - duration("24h")`)
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	sqlCondition, err := cel2sql.Convert(ast)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(sqlCondition)
	// Output: `employee`.`name` = "John Doe" AND `employee`.`hired_at` >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 24 HOUR)
}
