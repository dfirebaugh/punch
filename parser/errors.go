package parser

import (
	"fmt"
	"strings"

	"github.com/dfirebaugh/punch/token"
	"github.com/sirupsen/logrus"
)

var showFileName bool = true

func init() {
	logrus.SetLevel(logrus.ErrorLevel)
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("%s:\tno prefix parse function for %s found", p.curToken.Position, t)
	p.error(msg)
}

func (p *Parser) errorf(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	if !showFileName {
		return fmt.Errorf("[%d:%d]: %s", p.curToken.Position.Line, p.curToken.Position.Column, message)
	}
	return fmt.Errorf("%s:[%d:%d]: %s", p.curToken.Position.Filename, p.curToken.Position.Line, p.curToken.Position.Column, message)
}

func (p *Parser) error(msg ...string) error {
	if !showFileName {
		return fmt.Errorf("[%d:%d]: %s", p.curToken.Position.Line, p.curToken.Position.Column, strings.Join(msg, " "))
	}
	return fmt.Errorf("%s:[%d:%d]: %s", p.curToken.Position.Filename, p.curToken.Position.Line, p.curToken.Position.Column, strings.Join(msg, " "))
}

func (p *Parser) debug(msg ...string) {
	if !showFileName {
		logrus.Debugf("[%d:%d]: %s", p.curToken.Position.Line, p.curToken.Position.Column, strings.Join(msg, " "))
		return
	}
	logrus.Debugf("%s:[%d:%d]: %s", p.curToken.Position.Filename, p.curToken.Position.Line, p.curToken.Position.Column, strings.Join(msg, " "))
}

func (p *Parser) trace(msg ...string) {
	if !showFileName {
		logrus.Tracef("[%d:%d]: %s", p.curToken.Position.Line, p.curToken.Position.Column, strings.Join(msg, " "))
		return
	}
	logrus.Tracef("%s:[%d:%d]: %s", p.curToken.Position.Filename, p.curToken.Position.Line, p.curToken.Position.Column, strings.Join(msg, " "))
}

func (p *Parser) Errors() []string {
	return p.errors
}
