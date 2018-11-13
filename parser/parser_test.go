package parser

import (
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
	} {
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

	for _, stmt:= range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt isn't *ast.ReturnStatement, got=%T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("letStmt.TokenLiteral = %v, want %v", returnStmt.TokenLiteral(), "return")
		}
	}
}
