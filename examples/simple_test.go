package examples

import (
	"fmt"
	"log"

	"github.com/cockscomb/cel2sql"
	"github.com/cockscomb/cel2sql/generic"
	"github.com/cockscomb/cel2sql/sqltypes"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
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
	idents := []*exprpb.Decl{decls.NewVar("employee", decls.NewObjectType("Employee"))}
	sqlTypeProvider := generic.NewTypeProvider(map[string]generic.Schema{
		"Employee": employeeSchema,
	}, idents)
	env, _ := cel.NewEnv(
		cel.CustomTypeProvider(sqlTypeProvider),
		sqltypes.SQLTypeDeclarations,
		cel.Declarations(idents...),
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
	// Output: `employee`.`name` = "John Doe" AND `employee`.`hired_at` >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 24 HOUR)
}
