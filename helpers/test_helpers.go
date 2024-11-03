package helper

import (
	"testing"

	"github.com/ankush-web-eng/brolang/ast"
	"github.com/ankush-web-eng/brolang/object"
)

func TestLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "bhai_sun" {
		t.Errorf("s.TokenLiteral not 'bhai_sun'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	return true
}

func TestIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}
