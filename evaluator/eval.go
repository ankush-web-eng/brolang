package evaluator

import (
	"fmt"
	"log"

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
	case *ast.CallExpression:
		if node.Function.TokenLiteral() == "bol_bhai" {
			args := evalExpressions(node.Arguments, env)
			if len(args) == 1 && isError(args[0]) {
				return args[0]
			}
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		}
		return newError("unknown function: %s", node.Function.TokenLiteral())
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
			log.Printf("Warning: nil statement encountered")
			continue
		}

		fmt.Printf("Evaluating statement: %+v\n", stmt)
		result = Eval(stmt, env)

		if result == nil {
			log.Printf("Warning: statement evaluation returned nil")
			continue
		}

		if result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	if result == nil {
		return NULL
	}
	return result
}

func evalLetStatement(ls *ast.LetStatement, env *object.Environment) object.Object {
	if ls == nil || ls.Value == nil {
		return newError("invalid let statement")
	}

	value := Eval(ls.Value, env)
	if value == nil || value.Type() == object.ERROR_OBJ {
		return value
	}

	fmt.Printf("Setting variable: %s = %v\n", ls.Name.Value, value.Inspect())
	env.Set(ls.Name.Value, value)
	return value // Return the value instead of nil
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

func evalPrintStatement(ps *ast.PrintStatement, env *object.Environment) object.Object {
	if ps == nil || ps.Expression == nil {
		return newError("invalid print statement")
	}

	value := Eval(ps.Expression, env)
	if value == nil {
		return newError("cannot print nil value")
	}

	if value.Type() == object.ERROR_OBJ {
		return value
	}

	fmt.Printf("Output: %s\n", value.Inspect())
	return value
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// evalIdentifier evaluates an identifier to find its value in the environment.
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if node == nil {
		return newError("nil identifier")
	}

	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return newError("identifier not found: %s", node.Value)
}

// Helper functions for error handling and truthiness.
func isError(obj object.Object) bool {
	if obj == nil {
		return false
	}
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
