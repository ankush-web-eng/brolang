package evaluator

import (
	"fmt"

	"github.com/ankush-web-eng/brolang/ast"
	"github.com/ankush-web-eng/brolang/object"
)

// Eval evaluates the given AST node in the specified environment.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.Program:
		return evalProgram(node, env)
	case *ast.LetStatement:
		return evalLetStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ForExpression:
		return evalForExpression(node, env)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.BreakStatement:
		return evalBreakStatement()
	case *ast.ContinueStatement:
		return evalContinueStatement()

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
		return newError("Ye konsa function h ??: %s", node.Function.TokenLiteral())

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}

	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.AssignStatement:
		return evalAssignStatement(node, env)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	default:
		return newError("unknown node type: %T", node)
	}
}

// evalProgram evaluates a program by evaluating each statement in order.
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	if program == nil {
		return newError("Kuch likh to sahi be!!")
	}

	var result object.Object
	for _, stmt := range program.Statements {
		result = Eval(stmt, env)
		if result != nil && result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	if result == nil {
		return NULL
	}
	return result
}

// evaluates a let statement by evaluating the value and setting it in the environment.
func evalLetStatement(ls *ast.LetStatement, env *object.Environment) object.Object {
	if ls == nil || ls.Value == nil {
		return newError("invalid let statement")
	}

	value := Eval(ls.Value, env)
	if value == nil || value.Type() == object.ERROR_OBJ {
		return value
	}

	env.Set(ls.Name.Value, value)
	return value
}

// evaluates an assign statement by evaluating the value and setting it in the environment.
func evalAssignStatement(stmt *ast.AssignStatement, env *object.Environment) object.Object {
	// First evaluate the value to be assigned
	val := Eval(stmt.Value, env)
	if isError(val) {
		return val
	}

	// Try to set directly in current environment first
	if _, ok := env.Get(stmt.Name.Value); ok {
		return env.Set(stmt.Name.Value, val)
	}

	// If not found in current, explicitly check outer environments
	if env.Outer != nil {
		current := env.Outer
		for current != nil {
			if _, ok := current.Get(stmt.Name.Value); ok {
				return current.Set(stmt.Name.Value, val)
			}
			current = current.Outer
		}
	}

	// If variable doesn't exist anywhere in chain, create it in current environment
	return env.Set(stmt.Name.Value, val)
}

// evaluates an if expression by evaluating the condition and then the consequence or alternative.
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		result := Eval(ie.Consequence, env)
		// Propagate break/continue through if statements
		if result != nil && (result.Type() == "BREAK" || result.Type() == "CONTINUE") {
			return result
		}
		return result
	}

	// Handle else-if blocks
	for _, elseIf := range ie.ElseIf {
		condition := Eval(elseIf.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			result := Eval(elseIf.Consequence, env)
			// Propagate break/continue through else-if statements
			if result != nil && (result.Type() == "BREAK" || result.Type() == "CONTINUE") {
				return result
			}
			return result
		}
	}

	if ie.Alternative != nil {
		result := Eval(ie.Alternative, env)
		// Propagate break/continue through else statements
		if result != nil && (result.Type() == "BREAK" || result.Type() == "CONTINUE") {
			return result
		}
		return result
	}

	return NULL
}

// evalBlockStatement evaluates a block of statements.
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	// newEnv := object.NewEnvironment() // Create a new environment for block
	var result object.Object
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
		// if result != nil {
		// 	// Add output to environment's OutputBuilder
		// 	if result.Type() != object.ERROR_OBJ {
		// 		env.OutputBuilder.WriteString(result.Inspect())
		// 		env.OutputBuilder.WriteString("\n")
		// 	}
		// }

		if result != nil {
			rt := result.Type()
			if rt == "BREAK" || rt == "CONTINUE" || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

// the print statement is evaluated by evaluating the expression and printing its value.
func evalPrintStatement(ps *ast.PrintStatement, env *object.Environment) object.Object {
	if ps == nil || ps.Expression == nil {
		return newError("invalid print statement")
	}

	value := Eval(ps.Expression, env)
	if value == nil {
		return newError("kya coder banega re tu!! Print karana bhi nahi seekha!!")
	}

	if value.Type() == object.ERROR_OBJ {
		return value
	}

	// Keep adding the ouput to the Output' Builder string -- can be handled better
	env.OutputBuilder.WriteString(value.Inspect())
	env.OutputBuilder.WriteString("\n")

	return value
}

// -------All about loops and blocked scopes-------

// To prevent infinite loops and server loads
const maxIterations = 10000

type LoopControlFlow int

const (
	LoopNormal LoopControlFlow = iota
	LoopBreak
	LoopContinue
)

func evalBreakStatement() object.Object {
	return &object.BreakControl{}
}

func evalContinueStatement() object.Object {
	return &object.ContinueControl{}
}

// evaluate for-expressions (loops)
func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	loopEnv := object.NewEnclosedEnvironment(env)
	iterations := 0

	if fe.Init != nil {
		initResult := Eval(fe.Init, loopEnv)
		if isError(initResult) {
			return initResult
		}
	}

	var result object.Object = NULL

	for {
		iterations++
		if iterations > maxIterations {
			return newError("Mere server ka bill tera BAAP bharega ? Itni badi loop chala raha h!!")
		}

		if fe.Condition != nil {
			condition := Eval(fe.Condition, loopEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		result = Eval(fe.Body, loopEnv)

		// Append loopEnv's output to env's output after each iteration
		env.OutputBuilder.WriteString(loopEnv.OutputBuilder.String())
		loopEnv.OutputBuilder.Reset()

		if isError(result) {
			return result
		}

		// Handle break/continue
		switch result.(type) {
		case *object.BreakControl:
			return NULL
		case *object.ContinueControl:
			if fe.Update != nil {
				updateResult := Eval(fe.Update, loopEnv)
				if isError(updateResult) {
					return updateResult
				}
			}
			continue
		}

		if fe.Update != nil {
			updateResult := Eval(fe.Update, loopEnv)
			if isError(updateResult) {
				return updateResult
			}
		}
	}

	return result
}

// evaluate while-expressions (loops)
func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	loopEnv := object.NewEnclosedEnvironment(env)
	var result object.Object = NULL
	iterations := 0

	for {
		iterations++
		if iterations > maxIterations {
			return newError("Mere server ka bill tera BAAP bharega ? Itni badi loop chala raha h!!")
		}

		condition := Eval(we.Condition, loopEnv)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(we.Body, loopEnv)

		// Append loopEnv's output to env's output after each iteration
		env.OutputBuilder.WriteString(loopEnv.OutputBuilder.String())
		loopEnv.OutputBuilder.Reset()

		if isError(result) {
			return result
		}

		// Handle break/continue
		switch result.(type) {
		case *object.BreakControl:
			return NULL
		case *object.ContinueControl:
			continue
		}
	}

	return result
}

// -------About Arrays and Indexing-------

// evaluates an array literal by evaluating each element.
func evalArrayLiteral(node *ast.ArrayLiteral, env *object.Environment) object.Object {
	elements := evalExpressions(node.Elements, env)

	// Check for errors in elements
	for _, el := range elements {
		if isError(el) {
			return el
		}
	}

	// Check type consistency
	if len(elements) > 0 {
		firstType := elements[0].Type()
		for _, el := range elements[1:] {
			if el.Type() != firstType {
				return newError("Girgit mat ban, datatype mat badle array ke elements ka. %s ko %s se saath mix mat kar!!",
					firstType, el.Type())
			}
		}
	}

	return &object.Array{Elements: elements}
}

// evaluates an index expression by evaluating the left and index expressions.
func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("Kya coder banega re tu!! Sabse basic data structure bhi nahi aata tujhe!!: %s", left.Type())
	}
}

// evaluates an array index expression by checking if the index is within bounds.
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx, ok := index.(*object.Integer)
	if !ok {
		return newError("Beta tum se nahi ho payega, jao arrays padh ke aao striver sir se! Integer daal be,S %s", index.Type())
	}

	if idx.Value < 0 || idx.Value >= int64(len(arrayObject.Elements)) {
		return newError("Aukaat m rehle aukaat m, %d index pe kuch nahi hai! Bahar mat jaa array se!!", idx.Value)
	}

	return arrayObject.Elements[idx.Value]
}

// evaluates a list of expressions.
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
	return newError("Abe hosh me rehle! %s kaha likha h tune bataiyo zara...", node.Value)
}

// -------Infix Expressions-------

// evaluates an infix expression by evaluating the left and right expressions.
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	default:
		return newError("Bete %s, '%s', aur %s ka sambandh nahi ban sakta!!", left.Type(), operator, right.Type())
	}
}

// evaluates an integer infix expression.
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return &object.Boolean{Value: leftVal < rightVal}
	case ">":
		return &object.Boolean{Value: leftVal > rightVal}
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	case "<=":
		return &object.Boolean{Value: leftVal <= rightVal}
	case ">=":
		return &object.Boolean{Value: leftVal >= rightVal}
	default:
		return newError("Ye konsa operator h!?!?: %s", operator)
	}
}

// -------Helper functions for error handling and truthiness-------

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
	case *object.Integer:
		return obj.Value != 0
	case *object.Null:
		return false
	default:
		return true
	}
}

// NULL represents a null object.
var NULL = &object.Null{}
