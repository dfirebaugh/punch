package ast

import (
	"bytes"
	"encoding/json"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Files []*File
}

func (p *Program) TokenLiteral() string {
	if len(p.Files) > 0 && len(p.Files[0].Statements) > 0 {
		return p.Files[0].Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, f := range p.Files {
		out.WriteString(f.String())
		out.WriteString("\n")
	}
	return out.String()
}

func (p *Program) JSON() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (p *Program) JSONPretty() (string, error) {
	b, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

type File struct {
	Filename    string
	PackageName string
	Imports     []string
	Statements  []Statement
}

func (f *File) statementNode() {}

func (f *File) TokenLiteral() string {
	if len(f.Statements) > 0 {
		return f.Statements[0].TokenLiteral()
	}
	return ""
}

func (f *File) String() string {
	var out bytes.Buffer
	out.WriteString("package " + f.PackageName + "\n\n")
	if len(f.Imports) > 0 {
		out.WriteString("import (\n")
		for _, imp := range f.Imports {
			out.WriteString("\t\"" + imp + "\"\n")
		}
		out.WriteString(")\n\n")
	}
	for _, s := range f.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	return out.String()
}

func (f *File) JSON() (string, error) {
	b, err := json.Marshal(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (f *File) JSONPretty() (string, error) {
	b, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
