// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package restful

import (
	"bytes"
	"regexp"
	"strings"
)

// PathExpression holds a compiled path expression (RegExp) needed to match against
// Http request paths and to extract path parameter values.
type pathExpression struct {
	LiteralCount int // the number of literal characters (means those not resulting from template variable substitution)
	VarCount     int // the number of named parameters (enclosed by {}) in the path
	Matcher      *regexp.Regexp
	Source       string // Path as defined by the RouteBuilder
	tokens       []string
}

// NewPathExpression creates a PathExpression from the input URL path.
// Returns an error if the path is invalid.
func newPathExpression(path string) (*pathExpression, error) {
	expression, literalCount, varCount, tokens := templateToRegularExpression(path)
	compiled, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}
	return &pathExpression{literalCount, varCount, compiled, expression, tokens}, nil
}

// http://jsr311.java.net/nonav/releases/1.1/spec/spec3.html#x3-370003.7.3
func templateToRegularExpression(template string) (expression string, literalCount int, varCount int, tokens []string) {
	var buffer bytes.Buffer
	buffer.WriteString("^")
	//tokens = strings.Split(template, "/")
	tokens = tokenizePath(template)
	for _, each := range tokens {
		if each == "" {
			continue
		}
		buffer.WriteString("/")
		if strings.HasPrefix(each, "{") {
			// ignore var spec
			varCount += 1
			buffer.WriteString("([^/]+?)")
		} else {
			literalCount += len(each)
			encoded := each // TODO URI encode
			buffer.WriteString(regexp.QuoteMeta(encoded))
		}
	}
	return strings.TrimRight(buffer.String(), "/") + "(/.*)?$", literalCount, varCount, tokens
}
