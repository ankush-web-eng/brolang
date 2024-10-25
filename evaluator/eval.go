package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ankush-web-eng/brolang/ast"
	"github.com/ankush-web-eng/brolang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.IfExpression:
		return evalIfExpression(node)

	case *ast.WhileExpression:
		return evalWhileExpression(node)

	case *ast.ForExpression:
		return evalForExpression(node)

	case *ast.LetStatement:
		val := Eval(node.Value)
		if isError(val) {
			return val
		}
		return val

	case *ast.Identifier:
		if node.Value == "suna_bhai" {
			return evalInput()
		}
		return &object.Error{Message: "identifier not found: " + node.Value}

	case *ast.CallExpression:
		if node.Function.(*ast.Identifier).Value == "bol_bhai" {
			for _, arg := range node.Arguments {
				val := Eval(arg)
				if isError(val) {
					return val
				}
				fmt.Println(val.Inspect())
			}
			return NULL
		}
		return &object.Error{Message: "unknown function: " + node.Function.(*ast.Identifier).Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		index := Eval(node.Index)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}

	return NULL
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
		if errObj, ok := result.(*object.Error); ok {
			return errObj
		}
	}

	return result
}

func evalExpressions(exps []ast.Expression) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

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
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func evalWhileExpression(we *ast.WhileExpression) object.Object {
	condition := Eval(we.Condition)
	if isError(condition) {
		return condition
	}

	var result object.Object = NULL
	for isTruthy(condition) {
		result = Eval(we.Body)
		if isError(result) {
			return result
		}
		condition = Eval(we.Condition)
		if isError(condition) {
			return condition
		}
	}
	return result
}

func evalForExpression(fe *ast.ForExpression) object.Object {
	if fe.Init != nil {
		initVal := Eval(fe.Init)
		if isError(initVal) {
			return initVal
		}
	}

	var result object.Object = NULL
	for {
		condition := Eval(fe.Condition)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(fe.Body)
		if isError(result) {
			return result
		}

		if fe.Update != nil {
			updateVal := Eval(fe.Update)
			if isError(updateVal) {
				return updateVal
			}
		}
	}
	return result
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return &object.Error{Message: "aukat me reh le, aukat me, array ke bahar mat jaa!!"}
	}

	return arrayObject.Elements[idx]
}

func evalInput() object.Object {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Try to parse as integer
	if intVal, err := strconv.ParseInt(input, 10, 64); err == nil {
		return &object.Integer{Value: intVal}
	}

	// Try to parse as boolean
	if input == "sach" {
		return TRUE
	} else if input == "jhuth" {
		return FALSE
	}

	// Return as string if not integer or boolean
	return &object.String{Value: input}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// Environment related functions
var environment = make(map[string]object.Object)

func GetEnv(name string) (object.Object, bool) {
	obj, ok := environment[name]
	return obj, ok
}

func SetEnv(name string, val object.Object) object.Object {
	environment[name] = val
	return val
}

// Helper functions for array operations
func evalArrayIndexAssignment(array *object.Array, index int64, value object.Object) object.Object {
	if index < 0 || index >= int64(len(array.Elements)) {
		return &object.Error{Message: "aukat me reh le, aukat me, array ke bahar mat jaa!!"}
	}
	array.Elements[index] = value
	return value
}

// Custom error messages
func semicolonError() *object.Error {
	return &object.Error{Message: "bhadwe kya coder banega tu, semicolon bhool gaya!!!"}
}

func arrayBoundsError() *object.Error {
	return &object.Error{Message: "aukat me reh le, aukat me, array ke bahar mat jaa!!"}
}

// Helper function for type checking
func typeCheck(obj object.Object, expected ...object.ObjectType) bool {
	for _, t := range expected {
		if obj.Type() == t {
			return true
		}
	}
	return false
}

// Helper function for array operations
func coerceToInteger(obj object.Object) (int64, bool) {
	switch obj := obj.(type) {
	case *object.Integer:
		return obj.Value, true
	case *object.String:
		if i, err := strconv.ParseInt(obj.Value, 10, 64); err == nil {
			return i, true
		}
	}
	return 0, false
}

// Helper function for string operations
func coerceToString(obj object.Object) string {
	switch obj := obj.(type) {
	case *object.String:
		return obj.Value
	case *object.Integer:
		return strconv.FormatInt(obj.Value, 10)
	case *object.Boolean:
		return strconv.FormatBool(obj.Value)
	case *object.Array:
		elements := make([]string, len(obj.Elements))
		for i, elem := range obj.Elements {
			elements[i] = coerceToString(elem)
		}
		return "[" + strings.Join(elements, ", ") + "]"
	case *object.Null:
		return "null"
	default:
		return obj.Inspect()
	}
}

// Helper function for arithmetic operations
func performArithmetic(left, right object.Object, op string) object.Object {
	leftVal, leftOk := coerceToInteger(left)
	rightVal, rightOk := coerceToInteger(right)

	if !leftOk || !rightOk {
		return newError("cannot perform arithmetic on %s and %s", left.Type(), right.Type())
	}

	switch op {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &object.Integer{Value: leftVal / rightVal}
	default:
		return newError("unknown operator: %s", op)
	}
}

// Helper function for comparison operations
func performComparison(left, right object.Object, op string) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return compareIntegers(left.(*object.Integer).Value, right.(*object.Integer).Value, op)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return compareStrings(left.(*object.String).Value, right.(*object.String).Value, op)
	default:
		return newError("cannot compare %s and %s", left.Type(), right.Type())
	}
}

func compareIntegers(left, right int64, op string) object.Object {
	switch op {
	case "<":
		return nativeBoolToBooleanObject(left < right)
	case ">":
		return nativeBoolToBooleanObject(left > right)
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s", op)
	}
}

func compareStrings(left, right string, op string) object.Object {
	switch op {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s", op)
	}
}
