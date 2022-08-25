package object

import (
	"fmt"
	"hek/ast"
	"hek/token"
)

func Eval(node ast.Node, envs *Env) Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n.Statements, envs)
	case *ast.ExpressionStatement:
		return Eval(n.Expression, envs)
	case *ast.IntegerLiteral:
		return &Integer{Value: n.Value}
	case *ast.BoolExpression:
		return boolObject(n.Value)
	case *ast.PrefixExpression:
		res := Eval(n.Right, envs)
		return evalPrefix(n.Token.Type, res)
	case *ast.InfixExpression:
		left := Eval(n.Left, envs)
		right := Eval(n.Right, envs)
		return evalInfixExpression(n.Token.Type, left, right)
	case *ast.IFExpression:
		return evalIF(n, envs)
	case *ast.BlockStatement:
		return evalStatement(n, envs)
	case *ast.ReturnStatement:
		return evalReturn(n, envs)
	case *ast.LetStatement:
		res := Eval(n.Value, envs)
		if isError(res) {
			return res
		}
		envs.Set(n.Name.Value, res)
	case *ast.Identifier:
		return envs.Get(n.Value)
	case *ast.FunExpression:
		return evalFun(n, envs)
	case *ast.CallExpression:
		return evalCall(n, envs)
	case *ast.StringExpression:
		return &String{Value: n.Value}
	case *ast.ArrayExpression:
		return evalArray(n)
	case *ast.IndexExpression:
		return evalIndex(n, envs)
	case *ast.HashExpression:
		return evalHash(n)
	default:
		return newError("未知语法")
	}
	return NULL_
}
func evalProgram(arr []ast.Statement, envs *Env) Object {
	var result Object
	for _, statement := range arr {
		result = Eval(statement, envs)

		if v, ok := result.(*Return); ok {
			return v.Value
		}
	}
	return result
}
func boolObject(b bool) Object {
	if b {
		return TRUE
	}
	return FALSE
}
func evalPrefix(types token.Type, object Object) Object {
	switch types {
	case token.BANG:
		return evalPrefixBangExpression(object)
	case token.MINUS:
		return evalPrefixMinusExpression(object.(*Integer))
	}

	return NULL_
}
func evalPrefixBangExpression(object Object) Object {
	switch object {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL_:
		return TRUE
	default:
		return FALSE
	}
}
func evalPrefixMinusExpression(val *Integer) Object {
	return &Integer{Value: -val.Value}
}
func evalInfixExpression(types token.Type, left Object, right Object) Object {
	if object := infixTypes(left, right); object.Type() == ERROR {
		return object
	}
	if left.Type() == STRING {
		if types == token.PLUS {
			return &String{Value: fmt.Sprintf("%s%s", left.(*String).Value, right.(*String).Value)}
		} else if types == token.EQ {
			return boolObject(left.(*String).Value == right.(*String).Value)
		} else if types == token.NotEq {
			return boolObject(left.(*String).Value != right.(*String).Value)
		}
		return newError("string 不支持该操作 " + types.ToString())
	}
	switch types {
	case token.MINUS:
		return &Integer{Value: left.(*Integer).Value - right.(*Integer).Value}
	case token.PLUS:
		return &Integer{Value: left.(*Integer).Value + right.(*Integer).Value}
	case token.SLASH:
		return &Integer{Value: left.(*Integer).Value / right.(*Integer).Value}
	case token.ASTERISK:
		return &Integer{Value: left.(*Integer).Value * right.(*Integer).Value}
	case token.LT:
		return boolObject(left.(*Integer).Value < right.(*Integer).Value)
	case token.GT:
		return boolObject(left.(*Integer).Value > right.(*Integer).Value)
	case token.EQ:
		return boolObject(left.(*Integer).Value == right.(*Integer).Value)
	case token.NotEq:
		return boolObject(left.(*Integer).Value != right.(*Integer).Value)
	}
	return NULL_
}
func evalIF(if_ *ast.IFExpression, envs *Env) Object {
	condition := Eval(if_.Condition, envs)
	if isError(condition) {
		return condition
	}
	if isTrue(condition) {
		return Eval(if_.Consequence, envs)
	} else {
		if if_.Alternative != nil {
			return Eval(if_.Alternative, envs)
		}
	}
	return nil
}
func isTrue(object Object) bool {
	switch object {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL_:
		return false
	default:
		return true
	}
}
func evalStatement(stmt *ast.BlockStatement, envs *Env) Object {
	var result Object
	for _, statement := range stmt.Statements {
		result = Eval(statement, envs)
		if isError(result) {
			return result
		}
		if result.Type() == RETURN {
			return result
		}
	}
	return result
}
func evalReturn(ret *ast.ReturnStatement, envs *Env) Object {
	object := &Return{Value: NULL_}
	result := Eval(ret.Value, envs)
	object.Value = result
	return object
}
func infixTypes(left Object, right Object) Object {
	if left.Type() != right.Type() {
		return newError(fmt.Sprintf("%s and %s type atypism", left.Type().String(), right.Type().String()))
	}
	return NULL_
}
func newError(msg string) *Error {
	return &Error{Msg: msg}
}
func isError(object Object) bool {
	return object.Type() == ERROR
}
func evalFun(f *ast.FunExpression, envs *Env) Object {
	funObject := &Fun{
		Params: f.Params,
		Block:  f.Block,
		Env:    envs,
	}
	if f.Name != nil {
		envs.Set(f.Name.Value, funObject)
		return NULL_
	}
	return funObject
}
func evalCall(c *ast.CallExpression, envs *Env) Object {
	fun := Eval(c.Fun, envs)
	if isError(fun) {
		return fun
	}
	params := evalCallParamExpression(c.Params, envs)
	//if ifun, ok := fun.(*BuiltFun); ok {
	//	return ifun.Fun_(params...)
	//}
	return applyFun(fun.(*Fun), params)
}
func evalCallParamExpression(arr []ast.Expression, envs *Env) []Object {
	var arr_ []Object
	for _, expression := range arr {
		result := Eval(expression, envs)
		if isError(result) {
			return []Object{result}
		}
		arr_ = append(arr_, result)
	}
	return arr_
}
func applyFun(f *Fun, params []Object) Object {
	env := newFunEvn(params, f)
	if len(params) != len(f.Params) {
		return newError("参数数量不一致")
	}
	result := Eval(f.Block, env)
	return unwrapRet(result)
}
func newFunEvn(params []Object, fp *Fun) *Env {
	e := NewEnv(fp.Env)
	for index, name := range fp.Params {
		e.Set(name.Value, params[index])
	}
	return e
}
func unwrapRet(object Object) Object {
	if v, ok := object.(*Return); ok {
		return v.Value
	}
	return object
}
func evalArray(array *ast.ArrayExpression) Object {
	object := &Array{}
	for _, expression := range array.Value {
		object.Value = append(object.Value, Eval(expression, nil))
	}
	return object
}
func evalIndex(index *ast.IndexExpression, envs *Env) Object {
	arr := Eval(index.Left, envs)

	if arr.Type() == NULL {
		return NULL_
	}
	if array, ok := arr.(*Array); ok {
		i := Eval(index.Index, envs)
		if i == NULL_ || isError(i) {
			return i
		}
		index_ := i.(*Integer).Value
		if int(index_) >= len(array.Value) {
			return NULL_
		}
		return array.Value[index_]
	} else if hash, ok := arr.(*Hash); ok {
		i := Eval(index.Index, envs)
		val, ok := hash.Value[i]
		if !ok {
			return NULL_
		}
		return val
	}
	return NULL_
}
func evalHash(hash *ast.HashExpression) Object {
	h := &Hash{map[Object]Object{}}
	for index, val := range hash.Value {
		key := Eval(index, nil)
		if key == NULL_ {
			return newError("hash key err")
		}
		value := Eval(val, nil)
		if key == NULL_ {
			return newError("hash value err")
		}
		h.Value[key] = value
	}
	return h
}
