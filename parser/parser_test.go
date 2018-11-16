package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 5;
let hoge = 1111111111111111;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}

	if len(program.Statements) != 3 {
		t.Errorf("len(program.Statements) = %v, want %v", len(program.Statements), 3)
	}

	tests := []struct {
		expect string
	}{
		{"x"},
		{"y"},
		{"hoge"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testLetStatement(t, stmt, tt.expect)
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral = %v, want %v", s.TokenLiteral(), "let")
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s isn't *ast.LetStatement. got=%T", s)
		return
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value = %v, want %v", letStmt.Name.Value, name)
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral = %v, want %v", letStmt.Name.TokenLiteral(), name)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 12345;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}

	if len(program.Statements) != 3 {
		t.Errorf("len(program.Statements) = %v, want %v", len(program.Statements), 3)
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt isn't *ast.ReturnStatement, got=%T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("letStmt.TokenLiteral = %v, want %v", returnStmt.TokenLiteral(), "return")
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) = %v, want %v", len(program.Statements), 1)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] isn't *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp isn't *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value is %v, want %v", ident.Value, "foobar")
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral is %v, want %v", ident.TokenLiteral(), "foobar")
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) = %v, want %v", len(program.Statements), 1)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] isn't *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	integer, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp isn't *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if integer.Value != 5 {
		t.Errorf("integer.Value is %v, want %v", integer.Value, "foobar")
	}
	if integer.TokenLiteral() != "5" {
		t.Errorf("integer.TokenLiteral is %v, want %v", integer.TokenLiteral(), "foobar")
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements) = %v, want %v", len(program.Statements), 1)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] isn't *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp isn't *ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is %v. want %v", exp.Operator, tt.operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}

	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value is %v, want %v", integ.Value, value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral() is %v, want %v", integ.TokenLiteral(), value)
		return false

	}

	return true
}
