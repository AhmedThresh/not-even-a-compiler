package eval

import (
	"fmt"

	"github.com/AhmedThresh/not-even-a-compiler/pkg/ast"
	"github.com/AhmedThresh/not-even-a-compiler/pkg/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		env.Store(node.Name.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(right, node.Operator)

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.InfixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		return evalInfixExpression(right, left, node.Operator)

	case *ast.IfExpression:
		return evalIfExpression(node.Condition, node.Consequence, node.Alternative, env)

	case *ast.BlockStatement:
		return evalBlockStatements(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.Identifier:
		val, ok := env.Get(node.Value)
		if !ok {
			return newError("identifier not found: " + node.Value)
		}
		return val

	default:
		return NULL
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result.(*object.ReturnValue).Value
		}

		if result != nil && result.Type() == object.ERROR_OBJ {
			return result.(*object.Error)
		}

	}

	return result
}

func evalPrefixExpression(right object.Object, operator string) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalMinusOperator(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperator(right object.Object) object.Object {
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

func evalMinusOperator(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(right object.Object, left object.Object, operator string) object.Object {
	if right.Type() == object.INTEGER && left.Type() == object.INTEGER {
		return evalIntegerInfixOperation(right, left, operator)
	}

	if right.Type() == object.BOOLEAN && left.Type() == object.BOOLEAN {
		return evalBooleanInfixOperation(right, left, operator)
	}

	if left.Type() != right.Type() {
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalBooleanInfixOperation(right object.Object, left object.Object, operator string) object.Object {
	rightVal := right.(*object.Boolean).Value
	leftVal := left.(*object.Boolean).Value
	switch operator {
	case "==":
		if rightVal == leftVal {
			return TRUE
		}
		return FALSE
	case "!=":
		if rightVal != leftVal {
			return TRUE
		}
		return FALSE
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixOperation(right object.Object, left object.Object, operator string) object.Object {
	rightVal := right.(*object.Integer).Value
	leftVal := left.(*object.Integer).Value
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
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalIfExpression(condition ast.Expression, consequence *ast.BlockStatement, alternative *ast.BlockStatement, env *object.Environment) object.Object {
	conditionValue := Eval(condition, env)
	if isError(conditionValue) {
		return conditionValue
	}

	if conditionValue == FALSE || conditionValue == NULL {
		if alternative != nil {
			return Eval(alternative, env)
		}
		return NULL
	}
	return Eval(consequence, env)
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
			return result
		}
	}

	return result
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

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
