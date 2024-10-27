package evaluator

import (
	"fmt"

	"github.com/ankush-web-eng/brolang/ast"
	"github.com/ankush-web-eng/brolang/object"
)

// Eval evaluates the given AST node in the specified environment.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.LetStatement:
		fmt.Printf("Evaluating LetStatement: Name: %s, Value: %+v\n", node.Name.Value, node.Value)
		return evalLetStatement(node, env)
	case *ast.ForExpression:
		return evalForExpression(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.PrintStatement:
		return evalPrintStatement(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	default:
		return newError("unknown node type: %T", node)
	}
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	if program == nil {
		return newError("program is nil")
	}

	fmt.Println("Evaluating program with statements:", program.Statements)

	var result object.Object
	for _, stmt := range program.Statements {
		if stmt == nil {
			return newError("statement is nil")
		}

		fmt.Printf("Evaluating statement: %+v\n", stmt) // Log the statement
		result = Eval(stmt, env)
		if result.Type() == object.ERROR_OBJ {
			return result
		}
	}
	return result
}

func evalLetStatement(ls *ast.LetStatement, env *object.Environment) object.Object {
	value := Eval(ls.Value, env) // Ensure that Eval for the value is not nil
	if value.Type() == object.ERROR_OBJ {
		return value
	}
	env.Set(ls.Name.Value, value) // Ensure 'env' is not nil and 'Set' is implemented correctly
	return nil                    // Or appropriate return value
}

// evalForExpression evaluates a for expression.
func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	if fe.Init != nil {
		initVal := Eval(fe.Init, env)
		if isError(initVal) {
			return initVal
		}
	}

	for {
		condition := Eval(fe.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result := Eval(fe.Body, env)
		if isError(result) {
			return result
		}

		if fe.Update != nil {
			updateVal := Eval(fe.Update, env)
			if isError(updateVal) {
				return updateVal
			}
		}
	}
	return NULL
}

// evalBlockStatement evaluates a block of statements.
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	newEnv := object.NewEnvironment() // Create a new environment for block
	var result object.Object
	for _, stmt := range block.Statements {
		result = Eval(stmt, newEnv)
		if isError(result) {
			return result
		}
	}
	return result
}

// evalPrintStatement evaluates a print statement.
func evalPrintStatement(ps *ast.PrintStatement, env *object.Environment) object.Object {
	value := Eval(ps.Expression, env)
	if isError(value) {
		return value
	}
	fmt.Println(value.Inspect()) // Print the value
	return NULL
}

// evalIdentifier evaluates an identifier to find its value in the environment.
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return newError("identifier not found: %s", node.Value)
}

// Helper functions for error handling and truthiness.
func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJ
}

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

// NULL represents a null object.
var NULL = &object.Null{}
