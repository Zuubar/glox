package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func writeStringLn(file *os.File, str string) {
	_, err := file.WriteString(str + "\n")
	if err != nil {
		panic(err)
	}
}

func generateVisitor(file *os.File, baseName string, types []string) {
	writeStringLn(file, "package parser")
	writeStringLn(file, "")
	writeStringLn(file, fmt.Sprintf("type Visitor%s interface {", baseName))

	for _, t := range types {
		split := strings.Split(t, ":")
		typeName := strings.TrimSpace(split[0])

		fName := "Visit" + typeName + baseName
		fArgs := fmt.Sprintf("%s", typeName+baseName)
		writeStringLn(file, fmt.Sprintf("\t%s(%s) (any, error)", fName, fArgs))
	}
	writeStringLn(file, "}")
	writeStringLn(file, "")
}

func generateTypes(file *os.File, baseName string, types []string) {
	visitorInterfaceName := "Visitor" + baseName
	writeStringLn(file, fmt.Sprintf("type %s interface {", baseName))
	writeStringLn(file, fmt.Sprintf("\tAccept(visitor %s) (any, error)", visitorInterfaceName))
	writeStringLn(file, "}")
	writeStringLn(file, "")

	for _, t := range types {
		split := strings.Split(t, ":")
		name := strings.TrimSpace(split[0]) + baseName
		writeStringLn(file, fmt.Sprintf("type %s struct {", name))

		for _, member := range strings.Split(split[1], ",") {
			writeStringLn(file, "\t"+member)
		}
		writeStringLn(file, "}")
		writeStringLn(file, "")

		receiver := strings.ToLower(string(name[0]))
		writeStringLn(file, fmt.Sprintf("func (%s %s) Accept(visitor %s) (any, error) {", receiver, name, visitorInterfaceName))
		writeStringLn(file, fmt.Sprintf("\t return visitor.Visit%s(%s)", name, receiver))
		writeStringLn(file, "}")
		writeStringLn(file, "")
	}
}

func defineAst(outputDir string, baseName string, types []string) {
	file, err := os.Create(outputDir + "/" + strings.ToLower(baseName) + ".go")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	generateVisitor(file, baseName, types)
	generateTypes(file, baseName, types)
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Println("Output directory required.")
		return
	}
	outputDir := os.Args[1]

	defineAst(outputDir, "Expr", []string{
		"Ternary  	: Condition Expr, Left Expr, Right Expr",
		"Assignment : Name scanner.Token, Value Expr",
		"Logical	: Left Expr, Operator scanner.Token, Right Expr",
		"Set		: Object Expr, Name scanner.Token, Value Expr",
		"Super		: Keyword scanner.Token, Method scanner.Token",
		"Binary		: Left Expr, Operator scanner.Token, Right Expr",
		"Grouping	: Expr Expr",
		"Literal	: Value any",
		"Unary		: Operator scanner.Token, Right Expr",
		"Get		: Object Expr, Name scanner.Token",
		"Call		: Callee Expr, Parenthesis scanner.Token, Arguments []Expr",
		"Lambda		: Parenthesis scanner.Token, Parameters []scanner.Token, Body []Stmt",
		"This 		: Keyword scanner.Token",
		"Variable 	: Name scanner.Token",
	})

	defineAst(outputDir, "Stmt", []string{
		"Expression : Expression Expr",
		"Print      : Expression Expr",
		"Var 		: Name scanner.Token, Initializer Expr",
		"Class 		: Name scanner.Token, Superclass VariableExpr, Methods []FunctionStmt, StaticMethods []FunctionStmt",
		"Function 	: Name scanner.Token, Parameters []scanner.Token, Body []Stmt",
		"Block 		: Declarations []Stmt",
		"If 		: Expression Expr, ThenBranch Stmt, ElseBranch Stmt",
		"While 		: Condition Expr, Body Stmt",
		"For 		: Initializer Stmt, Condition Expr, Increment Stmt, Body Stmt",
		"Break 		: Keyword scanner.Token",
		"Continue 	: Keyword scanner.Token",
		"Return 	: Keyword scanner.Token, Expr Expr",
	})
}
