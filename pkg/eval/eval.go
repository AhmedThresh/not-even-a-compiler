package eval

import (
	"github.com/AhmedThresh/not-even-a-compiler/pkg/ast"
	"github.com/AhmedThresh/not-even-a-compiler/pkg/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(right, node.Operator)

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.InfixExpression:
		right := Eval(node.Right)
		left := Eval(node.Left)
		return evalInfixExpression(right, left, node.Operator)

	case *ast.IfExpression:
		return evalIfExpression(node.Condition, node.Consequence, node.Alternative)

	case *ast.BlockStatement:
		return evalStatements(node.Statements)

	default:
		return NULL
	}
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}

func evalPrefixExpression(right object.Object, operator string) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalMinuxOperator(right)
	default:
		return NULL
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

func evalMinuxOperator(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return NULL
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

	return NULL
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
		return NULL
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

	return NULL
}

func evalIfExpression(condition ast.Expression, consequence *ast.BlockStatement, alternative *ast.BlockStatement) object.Object {
	conditionValue := Eval(condition)
	if conditionValue == FALSE || conditionValue == NULL {
		if alternative != nil {
			return Eval(alternative)
		}
		return NULL
	}
	return Eval(consequence)
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
