// Package ast exposes AST elements used by River.
//
// The various interfaces exposed by ast are all closed; only types within this
// package can satisfy an AST interface.
package ast

import (
	"fmt"

	"github.com/grafana/agent/pkg/river/token"
)

// Node represents any node in the AST.
type Node interface {
	astNode()
}

// Stmt is a type of statement wthin the body of a file or block.
type Stmt interface {
	Node
	astStmt()
}

// Expr is an expression within the AST.
type Expr interface {
	Node
	astExpr()
}

// File is a parsed file.
type File struct {
	Name     string         // Filename provided to parser
	Body     Body           // Content of File
	Comments []CommentGroup // List of all comments in the File
}

// Body is a list of statements.
type Body []Stmt

// A CommentGroup represents a sequence of comments that are not separated by
// any empty lines or other non-comment tokens.
type CommentGroup []*Comment

// A Comment represents a single line or block comment.
//
// The Text field contains the comment text without any carriage returns (\r)
// that may have been present in the source. Since carriage returns get
// removed, EndPos will not be accurate for any comment which contained
// carriage returns.
type Comment struct {
	Start token.Pos // Starting position of comment
	// Text of the comment. Text will not contain '\n' for line comments.
	Text string
}

// AttributeStmt is a key-value pair being set in a Body or BlockStmt.
type AttributeStmt struct {
	Name  *IdentifierExpr
	Value Expr
}

// BlockStmt declares a block.
type BlockStmt struct {
	Name    []string
	NamePos token.Pos
	Label   string
	Body    Body

	LCurly, RCurly token.Pos
}

// IdentifierExpr refers to a named value.
type IdentifierExpr struct {
	Name    string
	NamePos token.Pos
}

// LiteralExpr is a constant value of a specific token kind.
type LiteralExpr struct {
	Kind     token.Token
	ValuePos token.Pos

	// Value holds the unparsed literal value. For example, if Kind ==
	// token.STRING, then Value would be wrapped in the original quotes (e.g.,
	// `"foobar"`).
	Value string
}

// ArrayExpr is an array of values.
type ArrayExpr struct {
	Elements       []Expr
	LBrack, RBrack token.Pos
}

// ObjectExpr declares an object of key-value pairs.
type ObjectExpr struct {
	Fields         []*ObjectField
	LCurly, RCurly token.Pos
}

// ObjectField defines an individual key-value pair within an object.
// ObjectField does not implement Node.
type ObjectField struct {
	Name   *IdentifierExpr
	Quoted bool // True if the name was wrapped in quotes
	Value  Expr
}

// AccessExpr accesses a field in an object value by name.
type AccessExpr struct {
	Value Expr
	Name  *IdentifierExpr
}

// IndexExpr accesses an index in an array value.
type IndexExpr struct {
	Value, Index   Expr
	LBrack, RBrack token.Pos
}

// CallExpr invokes a function value with a set of arguments.
type CallExpr struct {
	Value Expr
	Args  []Expr

	LParen, RParen token.Pos
}

// UnaryExpr performs a unary operation on a single value.
type UnaryExpr struct {
	Kind    token.Token
	KindPos token.Pos
	Value   Expr
}

// BinaryExpr performs a binary operation against two values.
type BinaryExpr struct {
	Kind        token.Token
	KindPos     token.Pos
	Left, Right Expr
}

// ParenExpr represents an expression wrapped in parenthesis.
type ParenExpr struct {
	Inner          Expr
	LParen, RParen token.Pos
}

// Type assertions

var (
	_ Node = (*File)(nil)
	_ Node = (*Body)(nil)
	_ Node = (*AttributeStmt)(nil)
	_ Node = (*BlockStmt)(nil)
	_ Node = (*IdentifierExpr)(nil)
	_ Node = (*LiteralExpr)(nil)
	_ Node = (*ArrayExpr)(nil)
	_ Node = (*ObjectExpr)(nil)
	_ Node = (*AccessExpr)(nil)
	_ Node = (*IndexExpr)(nil)
	_ Node = (*CallExpr)(nil)
	_ Node = (*UnaryExpr)(nil)
	_ Node = (*BinaryExpr)(nil)
	_ Node = (*ParenExpr)(nil)

	_ Stmt = (*AttributeStmt)(nil)
	_ Stmt = (*BlockStmt)(nil)

	_ Expr = (*IdentifierExpr)(nil)
	_ Expr = (*LiteralExpr)(nil)
	_ Expr = (*ArrayExpr)(nil)
	_ Expr = (*ObjectExpr)(nil)
	_ Expr = (*AccessExpr)(nil)
	_ Expr = (*IndexExpr)(nil)
	_ Expr = (*CallExpr)(nil)
	_ Expr = (*UnaryExpr)(nil)
	_ Expr = (*BinaryExpr)(nil)
	_ Expr = (*ParenExpr)(nil)
)

func (n *File) astNode()           {}
func (n Body) astNode()            {}
func (n CommentGroup) astNode()    {}
func (n *Comment) astNode()        {}
func (n *AttributeStmt) astNode()  {}
func (n *BlockStmt) astNode()      {}
func (n *IdentifierExpr) astNode() {}
func (n *LiteralExpr) astNode()    {}
func (n *ArrayExpr) astNode()      {}
func (n *ObjectExpr) astNode()     {}
func (n *AccessExpr) astNode()     {}
func (n *IndexExpr) astNode()      {}
func (n *CallExpr) astNode()       {}
func (n *UnaryExpr) astNode()      {}
func (n *BinaryExpr) astNode()     {}
func (n *ParenExpr) astNode()      {}

func (n *AttributeStmt) astStmt() {}
func (n *BlockStmt) astStmt()     {}

func (n *IdentifierExpr) astExpr() {}
func (n *LiteralExpr) astExpr()    {}
func (n *ArrayExpr) astExpr()      {}
func (n *ObjectExpr) astExpr()     {}
func (n *AccessExpr) astExpr()     {}
func (n *IndexExpr) astExpr()      {}
func (n *CallExpr) astExpr()       {}
func (n *UnaryExpr) astExpr()      {}
func (n *BinaryExpr) astExpr()     {}
func (n *ParenExpr) astExpr()      {}

// StartPos returns the position of the first character belonging to a Node.
func StartPos(n Node) token.Pos {
	if n == nil {
		return token.NoPos
	}
	switch n := n.(type) {
	case *File:
		return StartPos(n.Body)
	case Body:
		if len(n) == 0 {
			return token.NoPos
		}
		return StartPos(n[0])
	case CommentGroup:
		if len(n) == 0 {
			return token.NoPos
		}
		return StartPos(n[0])
	case *Comment:
		return n.Start
	case *AttributeStmt:
		return StartPos(n.Name)
	case *BlockStmt:
		return n.NamePos
	case *IdentifierExpr:
		return n.NamePos
	case *LiteralExpr:
		return n.ValuePos
	case *ArrayExpr:
		return n.LBrack
	case *ObjectExpr:
		return n.LCurly
	case *AccessExpr:
		return StartPos(n.Value)
	case *IndexExpr:
		return StartPos(n.Value)
	case *CallExpr:
		return StartPos(n.Value)
	case *UnaryExpr:
		return n.KindPos
	case *BinaryExpr:
		return StartPos(n.Left)
	case *ParenExpr:
		return n.LParen
	default:
		panic(fmt.Sprintf("Unhandled Node type %T", n))
	}
}

// EndPos returns the position of the final character in a Node.
func EndPos(n Node) token.Pos {
	if n == nil {
		return token.NoPos
	}
	switch n := n.(type) {
	case *File:
		return EndPos(n.Body)
	case Body:
		if len(n) == 0 {
			return token.NoPos
		}
		return EndPos(n[len(n)-1])
	case CommentGroup:
		if len(n) == 0 {
			return token.NoPos
		}
		return EndPos(n[len(n)-1])
	case *Comment:
		return n.Start.Add(len(n.Text) - 1)
	case *AttributeStmt:
		return EndPos(n.Name)
	case *BlockStmt:
		return n.RCurly
	case *IdentifierExpr:
		return n.NamePos.Add(len(n.Name) - 1)
	case *LiteralExpr:
		return n.ValuePos.Add(len(n.Value) - 1)
	case *ArrayExpr:
		return n.RBrack
	case *ObjectExpr:
		return n.RCurly
	case *AccessExpr:
		return EndPos(n.Name)
	case *IndexExpr:
		return n.RBrack
	case *CallExpr:
		return n.RParen
	case *UnaryExpr:
		return EndPos(n.Value)
	case *BinaryExpr:
		return EndPos(n.Right)
	case *ParenExpr:
		return n.RParen
	default:
		panic(fmt.Sprintf("Unhandled Node type %T", n))
	}
}
