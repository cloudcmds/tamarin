package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalListLiteral(
	ctx context.Context,
	node *ast.ListLiteral,
	s *scope.Scope,
) object.Object {
	elements := e.evalExpressions(ctx, node.Items, s)
	if len(elements) == 1 && isError(elements[0]) {
		return elements[0]
	}
	return &object.List{Items: elements}
}
