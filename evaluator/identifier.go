package evaluator

import (
	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalIdentifier(node *ast.Identifier, s *scope.Scope) object.Object {
	if val, ok := s.Get(node.Value); ok {
		return val
	}
	if builtin, ok := e.builtins[node.Value]; ok {
		return builtin
	}
	return newError("name error: %s is not defined", node.Value)
}
