package query

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 4

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleroot
	ruleselectStmt
	ruledescribeStmt
	ruledescribeAllStmt
	ruledescribeSingleStmt
	rulerangeClause
	ruleoptionalPredicateClause
	ruleexpressionList
	ruleexpression_1
	ruleexpression_2
	ruleexpression_3
	ruleexpression_function
	ruleexpression_metric
	rulegroupByClause
	rulepredicateClause
	rulepredicate_1
	rulepredicate_2
	rulepredicate_3
	ruletagMatcher
	ruleliteralString
	ruleliteralList
	ruleliteralListString
	ruletagName
	ruleCOLUMN_NAME
	ruleMETRIC_NAME
	ruleTAG_NAME
	ruleTIMESTAMP
	ruleIDENTIFIER
	ruleID_SEGMENT
	ruleID_START
	ruleID_CONT
	ruleKEYWORD
	ruleOP_ADD
	ruleOP_SUB
	ruleOP_MULT
	ruleOP_DIV
	ruleOP_AND
	ruleOP_OR
	ruleOP_NOT
	ruleSTRING
	ruleCHAR
	ruleESCAPE_CLASS
	ruleNUMBER
	ruleNUMBER_NATURAL
	ruleNUMBER_FRACTION
	ruleNUMBER_INTEGER
	ruleNUMBER_EXP
	rulePAREN_OPEN
	rulePAREN_CLOSE
	ruleCOMMA
	rule_
	rule__
	ruleSPACE
	ruleAction0
	ruleAction1
	rulePegText
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19
	ruleAction20
	ruleAction21
	ruleAction22
	ruleAction23
	ruleAction24
	ruleAction25
	ruleAction26
	ruleAction27
	ruleAction28
	ruleAction29
	ruleAction30

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"root",
	"selectStmt",
	"describeStmt",
	"describeAllStmt",
	"describeSingleStmt",
	"rangeClause",
	"optionalPredicateClause",
	"expressionList",
	"expression_1",
	"expression_2",
	"expression_3",
	"expression_function",
	"expression_metric",
	"groupByClause",
	"predicateClause",
	"predicate_1",
	"predicate_2",
	"predicate_3",
	"tagMatcher",
	"literalString",
	"literalList",
	"literalListString",
	"tagName",
	"COLUMN_NAME",
	"METRIC_NAME",
	"TAG_NAME",
	"TIMESTAMP",
	"IDENTIFIER",
	"ID_SEGMENT",
	"ID_START",
	"ID_CONT",
	"KEYWORD",
	"OP_ADD",
	"OP_SUB",
	"OP_MULT",
	"OP_DIV",
	"OP_AND",
	"OP_OR",
	"OP_NOT",
	"STRING",
	"CHAR",
	"ESCAPE_CLASS",
	"NUMBER",
	"NUMBER_NATURAL",
	"NUMBER_FRACTION",
	"NUMBER_INTEGER",
	"NUMBER_EXP",
	"PAREN_OPEN",
	"PAREN_CLOSE",
	"COMMA",
	"_",
	"__",
	"SPACE",
	"Action0",
	"Action1",
	"PegText",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
	"Action20",
	"Action21",
	"Action22",
	"Action23",
	"Action24",
	"Action25",
	"Action26",
	"Action27",
	"Action28",
	"Action29",
	"Action30",

	"Pre_",
	"_In_",
	"_Suf",
}

type tokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule pegRule, begin, end, next, depth int)
	Expand(index int) tokenTree
	Tokens() <-chan token32
	AST() *node32
	Error() []token32
	trim(length int)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (ast *node32) Print(buffer string) {
	ast.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token16 struct {
	pegRule
	begin, end, next int16
}

func (t *token16) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token16) isParentOf(u token16) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token16) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
}

func (t *token16) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens16 struct {
	tree    []token16
	ordered [][]token16
}

func (t *tokens16) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens16) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens16) Order() [][]token16 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int16, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token16, len(depths)), make([]token16, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = int16(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state16 struct {
	token16
	depths []int16
	leaf   bool
}

func (t *tokens16) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens16) PreOrder() (<-chan state16, [][]token16) {
	s, ordered := make(chan state16, 6), t.Order()
	go func() {
		var states [8]state16
		for i, _ := range states {
			states[i].depths = make([]int16, len(ordered))
		}
		depths, state, depth := make([]int16, len(ordered)), 0, 1
		write := func(t token16, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, int16(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token16 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token16{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token16{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token16{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens16) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens16) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens16) Add(rule pegRule, begin, end, depth, index int) {
	t.tree[index] = token16{pegRule: rule, begin: int16(begin), end: int16(end), next: int16(depth)}
}

func (t *tokens16) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens16) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next int32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = int32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i, _ := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, int32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth, index int) {
	t.tree[index] = token32{pegRule: rule, begin: int32(begin), end: int32(end), next: int32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

func (t *tokens16) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		for i, v := range tree {
			expanded[i] = v.getToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}

func (t *tokens32) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	return nil
}

type Parser struct {

	// temporary variables
	// ===================

	// stack of nodes used during the AST traversal.
	// a non-empty stack at the finish implies a programming error.
	nodeStack []Node

	// user errors accumulated during the AST traversal.
	// a non-empty list at the finish time means an invalid query is provided.
	errors []SyntaxError

	// programming errors accumulated during the AST traversal.
	// a non-empty list at the finish time implies a programming error.
	assertions []error

	// final result
	command Command

	Buffer string
	buffer []rune
	rules  [86]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	tokenTree
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer string, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer[0:] {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p *Parser
}

func (e *parseError) Error() string {
	tokens, error := e.p.tokenTree.Error(), "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.Buffer, positions)
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf("parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n",
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			/*strconv.Quote(*/ e.p.Buffer[begin:end] /*)*/)
	}

	return error
}

func (p *Parser) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *Parser) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *Parser) Execute() {
	buffer, begin, end := p.Buffer, 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)

		case ruleAction0:

			p.makeSelect()

		case ruleAction1:
			p.makeDescribeAll()
		case ruleAction2:
			p.addStringLiteral(unescapeLiteral(buffer[begin:end]))
		case ruleAction3:
			p.makeDescribe()
		case ruleAction4:
			p.addNullPredicate()
		case ruleAction5:
			p.addExpressionList()
		case ruleAction6:
			p.appendExpression()
		case ruleAction7:
			p.appendExpression()
		case ruleAction8:
			p.addOperatorLiteral("*")
		case ruleAction9:
			p.addOperatorLiteral("-")
		case ruleAction10:
			p.addOperatorFunction()
		case ruleAction11:
			p.addOperatorLiteral("*")
		case ruleAction12:
			p.addOperatorLiteral("*")
		case ruleAction13:
			p.addOperatorFunction()
		case ruleAction14:
			p.addNumberNode(buffer[begin:end])
		case ruleAction15:

			p.addStringLiteral(unescapeLiteral(buffer[begin:end]))

		case ruleAction16:

			p.addFunctionInvocation()

		case ruleAction17:

			p.addStringLiteral(unescapeLiteral(buffer[begin:end]))

		case ruleAction18:
			p.addNullPredicate()
		case ruleAction19:

			p.addMetricExpression()

		case ruleAction20:
			p.addAndPredicate()
		case ruleAction21:
			p.addOrPredicate()
		case ruleAction22:
			p.addNotPredicate()
		case ruleAction23:

			p.addLiteralMatcher()

		case ruleAction24:

			p.addLiteralMatcher()
			p.addNotPredicate()

		case ruleAction25:

			p.addRegexMatcher()

		case ruleAction26:

			p.addListMatcher()

		case ruleAction27:

			p.addStringLiteral(unescapeLiteral(buffer[begin:end]))

		case ruleAction28:
			p.addLiteralListNode()
		case ruleAction29:

			p.appendLiteral(unescapeLiteral(buffer[begin:end]))

		case ruleAction30:
			p.addTagLiteral(unescapeLiteral(buffer[begin:end]))

		}
	}
	_, _, _ = buffer, begin, end
}

func (p *Parser) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != end_symbol {
		p.buffer = append(p.buffer, end_symbol)
	}

	var tree tokenTree = &tokens16{tree: make([]token16, math.MaxInt16)}
	position, depth, tokenIndex, buffer, _rules := 0, 0, 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokenTree = tree
		if matches {
			p.tokenTree.trim(tokenIndex)
			return nil
		}
		return &parseError{p}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin int) {
		if t := tree.Expand(tokenIndex); t != nil {
			tree = t
		}
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
	}

	matchDot := func() bool {
		if buffer[position] != end_symbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 root <- <((selectStmt / describeStmt) !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					{
						position4 := position
						depth++
						{
							position5, tokenIndex5, depth5 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l6
							}
							position++
							goto l5
						l6:
							position, tokenIndex, depth = position5, tokenIndex5, depth5
							if buffer[position] != rune('S') {
								goto l3
							}
							position++
						}
					l5:
						{
							position7, tokenIndex7, depth7 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l8
							}
							position++
							goto l7
						l8:
							position, tokenIndex, depth = position7, tokenIndex7, depth7
							if buffer[position] != rune('E') {
								goto l3
							}
							position++
						}
					l7:
						{
							position9, tokenIndex9, depth9 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l10
							}
							position++
							goto l9
						l10:
							position, tokenIndex, depth = position9, tokenIndex9, depth9
							if buffer[position] != rune('L') {
								goto l3
							}
							position++
						}
					l9:
						{
							position11, tokenIndex11, depth11 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l12
							}
							position++
							goto l11
						l12:
							position, tokenIndex, depth = position11, tokenIndex11, depth11
							if buffer[position] != rune('E') {
								goto l3
							}
							position++
						}
					l11:
						{
							position13, tokenIndex13, depth13 := position, tokenIndex, depth
							if buffer[position] != rune('c') {
								goto l14
							}
							position++
							goto l13
						l14:
							position, tokenIndex, depth = position13, tokenIndex13, depth13
							if buffer[position] != rune('C') {
								goto l3
							}
							position++
						}
					l13:
						{
							position15, tokenIndex15, depth15 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l16
							}
							position++
							goto l15
						l16:
							position, tokenIndex, depth = position15, tokenIndex15, depth15
							if buffer[position] != rune('T') {
								goto l3
							}
							position++
						}
					l15:
						if !_rules[rule__]() {
							goto l3
						}
						if !_rules[ruleexpressionList]() {
							goto l3
						}
						if !_rules[ruleoptionalPredicateClause]() {
							goto l3
						}
						{
							position17 := position
							depth++
							if !_rules[rule_]() {
								goto l3
							}
							{
								position18, tokenIndex18, depth18 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l19
								}
								position++
								goto l18
							l19:
								position, tokenIndex, depth = position18, tokenIndex18, depth18
								if buffer[position] != rune('F') {
									goto l3
								}
								position++
							}
						l18:
							{
								position20, tokenIndex20, depth20 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l21
								}
								position++
								goto l20
							l21:
								position, tokenIndex, depth = position20, tokenIndex20, depth20
								if buffer[position] != rune('R') {
									goto l3
								}
								position++
							}
						l20:
							{
								position22, tokenIndex22, depth22 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l23
								}
								position++
								goto l22
							l23:
								position, tokenIndex, depth = position22, tokenIndex22, depth22
								if buffer[position] != rune('O') {
									goto l3
								}
								position++
							}
						l22:
							{
								position24, tokenIndex24, depth24 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l25
								}
								position++
								goto l24
							l25:
								position, tokenIndex, depth = position24, tokenIndex24, depth24
								if buffer[position] != rune('M') {
									goto l3
								}
								position++
							}
						l24:
							if !_rules[rule__]() {
								goto l3
							}
							if !_rules[ruleTIMESTAMP]() {
								goto l3
							}
							if !_rules[rule__]() {
								goto l3
							}
							{
								position26, tokenIndex26, depth26 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l27
								}
								position++
								goto l26
							l27:
								position, tokenIndex, depth = position26, tokenIndex26, depth26
								if buffer[position] != rune('T') {
									goto l3
								}
								position++
							}
						l26:
							{
								position28, tokenIndex28, depth28 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l29
								}
								position++
								goto l28
							l29:
								position, tokenIndex, depth = position28, tokenIndex28, depth28
								if buffer[position] != rune('O') {
									goto l3
								}
								position++
							}
						l28:
							if !_rules[rule__]() {
								goto l3
							}
							if !_rules[ruleTIMESTAMP]() {
								goto l3
							}
							depth--
							add(rulerangeClause, position17)
						}
						{
							add(ruleAction0, position)
						}
						depth--
						add(ruleselectStmt, position4)
					}
					goto l2
				l3:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					{
						position31 := position
						depth++
						{
							position32, tokenIndex32, depth32 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l33
							}
							position++
							goto l32
						l33:
							position, tokenIndex, depth = position32, tokenIndex32, depth32
							if buffer[position] != rune('D') {
								goto l0
							}
							position++
						}
					l32:
						{
							position34, tokenIndex34, depth34 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l35
							}
							position++
							goto l34
						l35:
							position, tokenIndex, depth = position34, tokenIndex34, depth34
							if buffer[position] != rune('E') {
								goto l0
							}
							position++
						}
					l34:
						{
							position36, tokenIndex36, depth36 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l37
							}
							position++
							goto l36
						l37:
							position, tokenIndex, depth = position36, tokenIndex36, depth36
							if buffer[position] != rune('S') {
								goto l0
							}
							position++
						}
					l36:
						{
							position38, tokenIndex38, depth38 := position, tokenIndex, depth
							if buffer[position] != rune('c') {
								goto l39
							}
							position++
							goto l38
						l39:
							position, tokenIndex, depth = position38, tokenIndex38, depth38
							if buffer[position] != rune('C') {
								goto l0
							}
							position++
						}
					l38:
						{
							position40, tokenIndex40, depth40 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l41
							}
							position++
							goto l40
						l41:
							position, tokenIndex, depth = position40, tokenIndex40, depth40
							if buffer[position] != rune('R') {
								goto l0
							}
							position++
						}
					l40:
						{
							position42, tokenIndex42, depth42 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l43
							}
							position++
							goto l42
						l43:
							position, tokenIndex, depth = position42, tokenIndex42, depth42
							if buffer[position] != rune('I') {
								goto l0
							}
							position++
						}
					l42:
						{
							position44, tokenIndex44, depth44 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l45
							}
							position++
							goto l44
						l45:
							position, tokenIndex, depth = position44, tokenIndex44, depth44
							if buffer[position] != rune('B') {
								goto l0
							}
							position++
						}
					l44:
						{
							position46, tokenIndex46, depth46 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l47
							}
							position++
							goto l46
						l47:
							position, tokenIndex, depth = position46, tokenIndex46, depth46
							if buffer[position] != rune('E') {
								goto l0
							}
							position++
						}
					l46:
						if !_rules[rule__]() {
							goto l0
						}
						{
							position48, tokenIndex48, depth48 := position, tokenIndex, depth
							{
								position50 := position
								depth++
								{
									position51, tokenIndex51, depth51 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l52
									}
									position++
									goto l51
								l52:
									position, tokenIndex, depth = position51, tokenIndex51, depth51
									if buffer[position] != rune('A') {
										goto l49
									}
									position++
								}
							l51:
								{
									position53, tokenIndex53, depth53 := position, tokenIndex, depth
									if buffer[position] != rune('l') {
										goto l54
									}
									position++
									goto l53
								l54:
									position, tokenIndex, depth = position53, tokenIndex53, depth53
									if buffer[position] != rune('L') {
										goto l49
									}
									position++
								}
							l53:
								{
									position55, tokenIndex55, depth55 := position, tokenIndex, depth
									if buffer[position] != rune('l') {
										goto l56
									}
									position++
									goto l55
								l56:
									position, tokenIndex, depth = position55, tokenIndex55, depth55
									if buffer[position] != rune('L') {
										goto l49
									}
									position++
								}
							l55:
								{
									add(ruleAction1, position)
								}
								depth--
								add(ruledescribeAllStmt, position50)
							}
							goto l48
						l49:
							position, tokenIndex, depth = position48, tokenIndex48, depth48
							{
								position58 := position
								depth++
								{
									position59 := position
									depth++
									{
										position60 := position
										depth++
										if !_rules[ruleIDENTIFIER]() {
											goto l0
										}
										depth--
										add(ruleMETRIC_NAME, position60)
									}
									depth--
									add(rulePegText, position59)
								}
								{
									add(ruleAction2, position)
								}
								if !_rules[ruleoptionalPredicateClause]() {
									goto l0
								}
								{
									add(ruleAction3, position)
								}
								depth--
								add(ruledescribeSingleStmt, position58)
							}
						}
					l48:
						depth--
						add(ruledescribeStmt, position31)
					}
				}
			l2:
				{
					position63, tokenIndex63, depth63 := position, tokenIndex, depth
					if !matchDot() {
						goto l63
					}
					goto l0
				l63:
					position, tokenIndex, depth = position63, tokenIndex63, depth63
				}
				depth--
				add(ruleroot, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 selectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') __ expressionList optionalPredicateClause rangeClause Action0)> */
		nil,
		/* 2 describeStmt <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') __ (describeAllStmt / describeSingleStmt))> */
		nil,
		/* 3 describeAllStmt <- <(('a' / 'A') ('l' / 'L') ('l' / 'L') Action1)> */
		nil,
		/* 4 describeSingleStmt <- <(<METRIC_NAME> Action2 optionalPredicateClause Action3)> */
		nil,
		/* 5 rangeClause <- <(_ (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) __ TIMESTAMP __ (('t' / 'T') ('o' / 'O')) __ TIMESTAMP)> */
		nil,
		/* 6 optionalPredicateClause <- <((__ predicateClause) / Action4)> */
		func() bool {
			{
				position70 := position
				depth++
				{
					position71, tokenIndex71, depth71 := position, tokenIndex, depth
					if !_rules[rule__]() {
						goto l72
					}
					{
						position73 := position
						depth++
						{
							position74, tokenIndex74, depth74 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l75
							}
							position++
							goto l74
						l75:
							position, tokenIndex, depth = position74, tokenIndex74, depth74
							if buffer[position] != rune('W') {
								goto l72
							}
							position++
						}
					l74:
						{
							position76, tokenIndex76, depth76 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l77
							}
							position++
							goto l76
						l77:
							position, tokenIndex, depth = position76, tokenIndex76, depth76
							if buffer[position] != rune('H') {
								goto l72
							}
							position++
						}
					l76:
						{
							position78, tokenIndex78, depth78 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l79
							}
							position++
							goto l78
						l79:
							position, tokenIndex, depth = position78, tokenIndex78, depth78
							if buffer[position] != rune('E') {
								goto l72
							}
							position++
						}
					l78:
						{
							position80, tokenIndex80, depth80 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l81
							}
							position++
							goto l80
						l81:
							position, tokenIndex, depth = position80, tokenIndex80, depth80
							if buffer[position] != rune('R') {
								goto l72
							}
							position++
						}
					l80:
						{
							position82, tokenIndex82, depth82 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l83
							}
							position++
							goto l82
						l83:
							position, tokenIndex, depth = position82, tokenIndex82, depth82
							if buffer[position] != rune('E') {
								goto l72
							}
							position++
						}
					l82:
						if !_rules[rule__]() {
							goto l72
						}
						if !_rules[rulepredicate_1]() {
							goto l72
						}
						depth--
						add(rulepredicateClause, position73)
					}
					goto l71
				l72:
					position, tokenIndex, depth = position71, tokenIndex71, depth71
					{
						add(ruleAction4, position)
					}
				}
			l71:
				depth--
				add(ruleoptionalPredicateClause, position70)
			}
			return true
		},
		/* 7 expressionList <- <(Action5 expression_1 Action6 (COMMA expression_1 Action7)*)> */
		func() bool {
			position85, tokenIndex85, depth85 := position, tokenIndex, depth
			{
				position86 := position
				depth++
				{
					add(ruleAction5, position)
				}
				if !_rules[ruleexpression_1]() {
					goto l85
				}
				{
					add(ruleAction6, position)
				}
			l89:
				{
					position90, tokenIndex90, depth90 := position, tokenIndex, depth
					if !_rules[ruleCOMMA]() {
						goto l90
					}
					if !_rules[ruleexpression_1]() {
						goto l90
					}
					{
						add(ruleAction7, position)
					}
					goto l89
				l90:
					position, tokenIndex, depth = position90, tokenIndex90, depth90
				}
				depth--
				add(ruleexpressionList, position86)
			}
			return true
		l85:
			position, tokenIndex, depth = position85, tokenIndex85, depth85
			return false
		},
		/* 8 expression_1 <- <(expression_2 (((OP_ADD Action8) / (OP_SUB Action9)) expression_2 Action10)*)> */
		func() bool {
			position92, tokenIndex92, depth92 := position, tokenIndex, depth
			{
				position93 := position
				depth++
				if !_rules[ruleexpression_2]() {
					goto l92
				}
			l94:
				{
					position95, tokenIndex95, depth95 := position, tokenIndex, depth
					{
						position96, tokenIndex96, depth96 := position, tokenIndex, depth
						{
							position98 := position
							depth++
							if !_rules[rule_]() {
								goto l97
							}
							if buffer[position] != rune('+') {
								goto l97
							}
							position++
							if !_rules[rule_]() {
								goto l97
							}
							depth--
							add(ruleOP_ADD, position98)
						}
						{
							add(ruleAction8, position)
						}
						goto l96
					l97:
						position, tokenIndex, depth = position96, tokenIndex96, depth96
						{
							position100 := position
							depth++
							if !_rules[rule_]() {
								goto l95
							}
							if buffer[position] != rune('-') {
								goto l95
							}
							position++
							if !_rules[rule_]() {
								goto l95
							}
							depth--
							add(ruleOP_SUB, position100)
						}
						{
							add(ruleAction9, position)
						}
					}
				l96:
					if !_rules[ruleexpression_2]() {
						goto l95
					}
					{
						add(ruleAction10, position)
					}
					goto l94
				l95:
					position, tokenIndex, depth = position95, tokenIndex95, depth95
				}
				depth--
				add(ruleexpression_1, position93)
			}
			return true
		l92:
			position, tokenIndex, depth = position92, tokenIndex92, depth92
			return false
		},
		/* 9 expression_2 <- <(expression_3 (((OP_DIV Action11) / (OP_MULT Action12)) expression_3 Action13)*)> */
		func() bool {
			position103, tokenIndex103, depth103 := position, tokenIndex, depth
			{
				position104 := position
				depth++
				if !_rules[ruleexpression_3]() {
					goto l103
				}
			l105:
				{
					position106, tokenIndex106, depth106 := position, tokenIndex, depth
					{
						position107, tokenIndex107, depth107 := position, tokenIndex, depth
						{
							position109 := position
							depth++
							if !_rules[rule_]() {
								goto l108
							}
							if buffer[position] != rune('/') {
								goto l108
							}
							position++
							if !_rules[rule_]() {
								goto l108
							}
							depth--
							add(ruleOP_DIV, position109)
						}
						{
							add(ruleAction11, position)
						}
						goto l107
					l108:
						position, tokenIndex, depth = position107, tokenIndex107, depth107
						{
							position111 := position
							depth++
							if !_rules[rule_]() {
								goto l106
							}
							if buffer[position] != rune('*') {
								goto l106
							}
							position++
							if !_rules[rule_]() {
								goto l106
							}
							depth--
							add(ruleOP_MULT, position111)
						}
						{
							add(ruleAction12, position)
						}
					}
				l107:
					if !_rules[ruleexpression_3]() {
						goto l106
					}
					{
						add(ruleAction13, position)
					}
					goto l105
				l106:
					position, tokenIndex, depth = position106, tokenIndex106, depth106
				}
				depth--
				add(ruleexpression_2, position104)
			}
			return true
		l103:
			position, tokenIndex, depth = position103, tokenIndex103, depth103
			return false
		},
		/* 10 expression_3 <- <(expression_function / ((&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') (<NUMBER> Action14)) | (&('\t' | '\n' | ' ' | '(') (PAREN_OPEN expression_1 PAREN_CLOSE)) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | '`' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') expression_metric)))> */
		func() bool {
			position114, tokenIndex114, depth114 := position, tokenIndex, depth
			{
				position115 := position
				depth++
				{
					position116, tokenIndex116, depth116 := position, tokenIndex, depth
					{
						position118 := position
						depth++
						{
							position119 := position
							depth++
							if !_rules[ruleIDENTIFIER]() {
								goto l117
							}
							depth--
							add(rulePegText, position119)
						}
						{
							add(ruleAction15, position)
						}
						if !_rules[rulePAREN_OPEN]() {
							goto l117
						}
						if !_rules[ruleexpressionList]() {
							goto l117
						}
						{
							position121, tokenIndex121, depth121 := position, tokenIndex, depth
							if !_rules[rule__]() {
								goto l121
							}
							{
								position123 := position
								depth++
								{
									position124, tokenIndex124, depth124 := position, tokenIndex, depth
									if buffer[position] != rune('g') {
										goto l125
									}
									position++
									goto l124
								l125:
									position, tokenIndex, depth = position124, tokenIndex124, depth124
									if buffer[position] != rune('G') {
										goto l121
									}
									position++
								}
							l124:
								{
									position126, tokenIndex126, depth126 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l127
									}
									position++
									goto l126
								l127:
									position, tokenIndex, depth = position126, tokenIndex126, depth126
									if buffer[position] != rune('R') {
										goto l121
									}
									position++
								}
							l126:
								{
									position128, tokenIndex128, depth128 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l129
									}
									position++
									goto l128
								l129:
									position, tokenIndex, depth = position128, tokenIndex128, depth128
									if buffer[position] != rune('O') {
										goto l121
									}
									position++
								}
							l128:
								{
									position130, tokenIndex130, depth130 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l131
									}
									position++
									goto l130
								l131:
									position, tokenIndex, depth = position130, tokenIndex130, depth130
									if buffer[position] != rune('U') {
										goto l121
									}
									position++
								}
							l130:
								{
									position132, tokenIndex132, depth132 := position, tokenIndex, depth
									if buffer[position] != rune('p') {
										goto l133
									}
									position++
									goto l132
								l133:
									position, tokenIndex, depth = position132, tokenIndex132, depth132
									if buffer[position] != rune('P') {
										goto l121
									}
									position++
								}
							l132:
								if !_rules[rule__]() {
									goto l121
								}
								{
									position134, tokenIndex134, depth134 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l135
									}
									position++
									goto l134
								l135:
									position, tokenIndex, depth = position134, tokenIndex134, depth134
									if buffer[position] != rune('B') {
										goto l121
									}
									position++
								}
							l134:
								{
									position136, tokenIndex136, depth136 := position, tokenIndex, depth
									if buffer[position] != rune('y') {
										goto l137
									}
									position++
									goto l136
								l137:
									position, tokenIndex, depth = position136, tokenIndex136, depth136
									if buffer[position] != rune('Y') {
										goto l121
									}
									position++
								}
							l136:
								if !_rules[rule__]() {
									goto l121
								}
								if !_rules[ruleCOLUMN_NAME]() {
									goto l121
								}
							l138:
								{
									position139, tokenIndex139, depth139 := position, tokenIndex, depth
									if !_rules[ruleCOMMA]() {
										goto l139
									}
									if !_rules[ruleCOLUMN_NAME]() {
										goto l139
									}
									goto l138
								l139:
									position, tokenIndex, depth = position139, tokenIndex139, depth139
								}
								depth--
								add(rulegroupByClause, position123)
							}
							goto l122
						l121:
							position, tokenIndex, depth = position121, tokenIndex121, depth121
						}
					l122:
						if !_rules[rulePAREN_CLOSE]() {
							goto l117
						}
						{
							add(ruleAction16, position)
						}
						depth--
						add(ruleexpression_function, position118)
					}
					goto l116
				l117:
					position, tokenIndex, depth = position116, tokenIndex116, depth116
					{
						switch buffer[position] {
						case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							{
								position142 := position
								depth++
								{
									position143 := position
									depth++
									{
										position144 := position
										depth++
										{
											position145, tokenIndex145, depth145 := position, tokenIndex, depth
											if buffer[position] != rune('-') {
												goto l145
											}
											position++
											goto l146
										l145:
											position, tokenIndex, depth = position145, tokenIndex145, depth145
										}
									l146:
										if !_rules[ruleNUMBER_NATURAL]() {
											goto l114
										}
										depth--
										add(ruleNUMBER_INTEGER, position144)
									}
									{
										position147, tokenIndex147, depth147 := position, tokenIndex, depth
										{
											position149 := position
											depth++
											if buffer[position] != rune('.') {
												goto l147
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l147
											}
											position++
										l150:
											{
												position151, tokenIndex151, depth151 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l151
												}
												position++
												goto l150
											l151:
												position, tokenIndex, depth = position151, tokenIndex151, depth151
											}
											depth--
											add(ruleNUMBER_FRACTION, position149)
										}
										goto l148
									l147:
										position, tokenIndex, depth = position147, tokenIndex147, depth147
									}
								l148:
									{
										position152, tokenIndex152, depth152 := position, tokenIndex, depth
										{
											position154 := position
											depth++
											{
												position155, tokenIndex155, depth155 := position, tokenIndex, depth
												if buffer[position] != rune('e') {
													goto l156
												}
												position++
												goto l155
											l156:
												position, tokenIndex, depth = position155, tokenIndex155, depth155
												if buffer[position] != rune('E') {
													goto l152
												}
												position++
											}
										l155:
											{
												position157, tokenIndex157, depth157 := position, tokenIndex, depth
												{
													position159, tokenIndex159, depth159 := position, tokenIndex, depth
													if buffer[position] != rune('+') {
														goto l160
													}
													position++
													goto l159
												l160:
													position, tokenIndex, depth = position159, tokenIndex159, depth159
													if buffer[position] != rune('-') {
														goto l157
													}
													position++
												}
											l159:
												goto l158
											l157:
												position, tokenIndex, depth = position157, tokenIndex157, depth157
											}
										l158:
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l152
											}
											position++
										l161:
											{
												position162, tokenIndex162, depth162 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l162
												}
												position++
												goto l161
											l162:
												position, tokenIndex, depth = position162, tokenIndex162, depth162
											}
											depth--
											add(ruleNUMBER_EXP, position154)
										}
										goto l153
									l152:
										position, tokenIndex, depth = position152, tokenIndex152, depth152
									}
								l153:
									depth--
									add(ruleNUMBER, position143)
								}
								depth--
								add(rulePegText, position142)
							}
							{
								add(ruleAction14, position)
							}
							break
						case '\t', '\n', ' ', '(':
							if !_rules[rulePAREN_OPEN]() {
								goto l114
							}
							if !_rules[ruleexpression_1]() {
								goto l114
							}
							if !_rules[rulePAREN_CLOSE]() {
								goto l114
							}
							break
						default:
							{
								position164 := position
								depth++
								{
									position165 := position
									depth++
									if !_rules[ruleIDENTIFIER]() {
										goto l114
									}
									depth--
									add(rulePegText, position165)
								}
								{
									add(ruleAction17, position)
								}
								{
									position167, tokenIndex167, depth167 := position, tokenIndex, depth
									{
										position169, tokenIndex169, depth169 := position, tokenIndex, depth
										if buffer[position] != rune('[') {
											goto l170
										}
										position++
										if !_rules[rule_]() {
											goto l170
										}
										if !_rules[rulepredicate_1]() {
											goto l170
										}
										if !_rules[rule_]() {
											goto l170
										}
										if buffer[position] != rune(']') {
											goto l170
										}
										position++
										goto l169
									l170:
										position, tokenIndex, depth = position169, tokenIndex169, depth169
										{
											add(ruleAction18, position)
										}
									}
								l169:
									goto l168

									position, tokenIndex, depth = position167, tokenIndex167, depth167
								}
							l168:
								{
									add(ruleAction19, position)
								}
								depth--
								add(ruleexpression_metric, position164)
							}
							break
						}
					}

				}
			l116:
				depth--
				add(ruleexpression_3, position115)
			}
			return true
		l114:
			position, tokenIndex, depth = position114, tokenIndex114, depth114
			return false
		},
		/* 11 expression_function <- <(<IDENTIFIER> Action15 PAREN_OPEN expressionList (__ groupByClause)? PAREN_CLOSE Action16)> */
		nil,
		/* 12 expression_metric <- <(<IDENTIFIER> Action17 (('[' _ predicate_1 _ ']') / Action18)? Action19)> */
		nil,
		/* 13 groupByClause <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') __ (('b' / 'B') ('y' / 'Y')) __ COLUMN_NAME (COMMA COLUMN_NAME)*)> */
		nil,
		/* 14 predicateClause <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') __ predicate_1)> */
		nil,
		/* 15 predicate_1 <- <((predicate_2 OP_AND predicate_1 Action20) / predicate_2 / )> */
		func() bool {
			{
				position178 := position
				depth++
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					if !_rules[rulepredicate_2]() {
						goto l180
					}
					{
						position181 := position
						depth++
						if !_rules[rule__]() {
							goto l180
						}
						{
							position182, tokenIndex182, depth182 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l183
							}
							position++
							goto l182
						l183:
							position, tokenIndex, depth = position182, tokenIndex182, depth182
							if buffer[position] != rune('A') {
								goto l180
							}
							position++
						}
					l182:
						{
							position184, tokenIndex184, depth184 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l185
							}
							position++
							goto l184
						l185:
							position, tokenIndex, depth = position184, tokenIndex184, depth184
							if buffer[position] != rune('N') {
								goto l180
							}
							position++
						}
					l184:
						{
							position186, tokenIndex186, depth186 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l187
							}
							position++
							goto l186
						l187:
							position, tokenIndex, depth = position186, tokenIndex186, depth186
							if buffer[position] != rune('D') {
								goto l180
							}
							position++
						}
					l186:
						if !_rules[rule__]() {
							goto l180
						}
						depth--
						add(ruleOP_AND, position181)
					}
					if !_rules[rulepredicate_1]() {
						goto l180
					}
					{
						add(ruleAction20, position)
					}
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if !_rules[rulepredicate_2]() {
						goto l189
					}
					goto l179
				l189:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
				}
			l179:
				depth--
				add(rulepredicate_1, position178)
			}
			return true
		},
		/* 16 predicate_2 <- <((predicate_3 OP_OR predicate_2 Action21) / predicate_3)> */
		func() bool {
			position190, tokenIndex190, depth190 := position, tokenIndex, depth
			{
				position191 := position
				depth++
				{
					position192, tokenIndex192, depth192 := position, tokenIndex, depth
					if !_rules[rulepredicate_3]() {
						goto l193
					}
					{
						position194 := position
						depth++
						if !_rules[rule__]() {
							goto l193
						}
						{
							position195, tokenIndex195, depth195 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l196
							}
							position++
							goto l195
						l196:
							position, tokenIndex, depth = position195, tokenIndex195, depth195
							if buffer[position] != rune('O') {
								goto l193
							}
							position++
						}
					l195:
						{
							position197, tokenIndex197, depth197 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l198
							}
							position++
							goto l197
						l198:
							position, tokenIndex, depth = position197, tokenIndex197, depth197
							if buffer[position] != rune('R') {
								goto l193
							}
							position++
						}
					l197:
						if !_rules[rule__]() {
							goto l193
						}
						depth--
						add(ruleOP_OR, position194)
					}
					if !_rules[rulepredicate_2]() {
						goto l193
					}
					{
						add(ruleAction21, position)
					}
					goto l192
				l193:
					position, tokenIndex, depth = position192, tokenIndex192, depth192
					if !_rules[rulepredicate_3]() {
						goto l190
					}
				}
			l192:
				depth--
				add(rulepredicate_2, position191)
			}
			return true
		l190:
			position, tokenIndex, depth = position190, tokenIndex190, depth190
			return false
		},
		/* 17 predicate_3 <- <((OP_NOT predicate_3 Action22) / (PAREN_OPEN predicate_1 PAREN_CLOSE) / tagMatcher)> */
		func() bool {
			position200, tokenIndex200, depth200 := position, tokenIndex, depth
			{
				position201 := position
				depth++
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					{
						position204 := position
						depth++
						{
							position205, tokenIndex205, depth205 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l206
							}
							position++
							goto l205
						l206:
							position, tokenIndex, depth = position205, tokenIndex205, depth205
							if buffer[position] != rune('N') {
								goto l203
							}
							position++
						}
					l205:
						{
							position207, tokenIndex207, depth207 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l208
							}
							position++
							goto l207
						l208:
							position, tokenIndex, depth = position207, tokenIndex207, depth207
							if buffer[position] != rune('O') {
								goto l203
							}
							position++
						}
					l207:
						{
							position209, tokenIndex209, depth209 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l210
							}
							position++
							goto l209
						l210:
							position, tokenIndex, depth = position209, tokenIndex209, depth209
							if buffer[position] != rune('T') {
								goto l203
							}
							position++
						}
					l209:
						if !_rules[rule__]() {
							goto l203
						}
						depth--
						add(ruleOP_NOT, position204)
					}
					if !_rules[rulepredicate_3]() {
						goto l203
					}
					{
						add(ruleAction22, position)
					}
					goto l202
				l203:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					if !_rules[rulePAREN_OPEN]() {
						goto l212
					}
					if !_rules[rulepredicate_1]() {
						goto l212
					}
					if !_rules[rulePAREN_CLOSE]() {
						goto l212
					}
					goto l202
				l212:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					{
						position213 := position
						depth++
						{
							position214, tokenIndex214, depth214 := position, tokenIndex, depth
							if !_rules[ruletagName]() {
								goto l215
							}
							if !_rules[rule_]() {
								goto l215
							}
							if buffer[position] != rune('=') {
								goto l215
							}
							position++
							if !_rules[rule_]() {
								goto l215
							}
							if !_rules[ruleliteralString]() {
								goto l215
							}
							{
								add(ruleAction23, position)
							}
							goto l214
						l215:
							position, tokenIndex, depth = position214, tokenIndex214, depth214
							if !_rules[ruletagName]() {
								goto l217
							}
							if !_rules[rule_]() {
								goto l217
							}
							if buffer[position] != rune('!') {
								goto l217
							}
							position++
							if buffer[position] != rune('=') {
								goto l217
							}
							position++
							if !_rules[rule_]() {
								goto l217
							}
							if !_rules[ruleliteralString]() {
								goto l217
							}
							{
								add(ruleAction24, position)
							}
							goto l214
						l217:
							position, tokenIndex, depth = position214, tokenIndex214, depth214
							if !_rules[ruletagName]() {
								goto l219
							}
							if !_rules[rule__]() {
								goto l219
							}
							{
								position220, tokenIndex220, depth220 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l221
								}
								position++
								goto l220
							l221:
								position, tokenIndex, depth = position220, tokenIndex220, depth220
								if buffer[position] != rune('M') {
									goto l219
								}
								position++
							}
						l220:
							{
								position222, tokenIndex222, depth222 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l223
								}
								position++
								goto l222
							l223:
								position, tokenIndex, depth = position222, tokenIndex222, depth222
								if buffer[position] != rune('A') {
									goto l219
								}
								position++
							}
						l222:
							{
								position224, tokenIndex224, depth224 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l225
								}
								position++
								goto l224
							l225:
								position, tokenIndex, depth = position224, tokenIndex224, depth224
								if buffer[position] != rune('T') {
									goto l219
								}
								position++
							}
						l224:
							{
								position226, tokenIndex226, depth226 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l227
								}
								position++
								goto l226
							l227:
								position, tokenIndex, depth = position226, tokenIndex226, depth226
								if buffer[position] != rune('C') {
									goto l219
								}
								position++
							}
						l226:
							{
								position228, tokenIndex228, depth228 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l229
								}
								position++
								goto l228
							l229:
								position, tokenIndex, depth = position228, tokenIndex228, depth228
								if buffer[position] != rune('H') {
									goto l219
								}
								position++
							}
						l228:
							{
								position230, tokenIndex230, depth230 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l231
								}
								position++
								goto l230
							l231:
								position, tokenIndex, depth = position230, tokenIndex230, depth230
								if buffer[position] != rune('E') {
									goto l219
								}
								position++
							}
						l230:
							{
								position232, tokenIndex232, depth232 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l233
								}
								position++
								goto l232
							l233:
								position, tokenIndex, depth = position232, tokenIndex232, depth232
								if buffer[position] != rune('S') {
									goto l219
								}
								position++
							}
						l232:
							if !_rules[rule__]() {
								goto l219
							}
							if !_rules[ruleliteralString]() {
								goto l219
							}
							{
								add(ruleAction25, position)
							}
							goto l214
						l219:
							position, tokenIndex, depth = position214, tokenIndex214, depth214
							if !_rules[ruletagName]() {
								goto l200
							}
							if !_rules[rule__]() {
								goto l200
							}
							{
								position235, tokenIndex235, depth235 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l236
								}
								position++
								goto l235
							l236:
								position, tokenIndex, depth = position235, tokenIndex235, depth235
								if buffer[position] != rune('I') {
									goto l200
								}
								position++
							}
						l235:
							{
								position237, tokenIndex237, depth237 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l238
								}
								position++
								goto l237
							l238:
								position, tokenIndex, depth = position237, tokenIndex237, depth237
								if buffer[position] != rune('N') {
									goto l200
								}
								position++
							}
						l237:
							if !_rules[rule__]() {
								goto l200
							}
							{
								position239 := position
								depth++
								{
									add(ruleAction28, position)
								}
								if !_rules[rulePAREN_OPEN]() {
									goto l200
								}
								if !_rules[ruleliteralListString]() {
									goto l200
								}
							l241:
								{
									position242, tokenIndex242, depth242 := position, tokenIndex, depth
									if !_rules[ruleCOMMA]() {
										goto l242
									}
									if !_rules[ruleliteralListString]() {
										goto l242
									}
									goto l241
								l242:
									position, tokenIndex, depth = position242, tokenIndex242, depth242
								}
								if !_rules[rulePAREN_CLOSE]() {
									goto l200
								}
								depth--
								add(ruleliteralList, position239)
							}
							{
								add(ruleAction26, position)
							}
						}
					l214:
						depth--
						add(ruletagMatcher, position213)
					}
				}
			l202:
				depth--
				add(rulepredicate_3, position201)
			}
			return true
		l200:
			position, tokenIndex, depth = position200, tokenIndex200, depth200
			return false
		},
		/* 18 tagMatcher <- <((tagName _ '=' _ literalString Action23) / (tagName _ ('!' '=') _ literalString Action24) / (tagName __ (('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S')) __ literalString Action25) / (tagName __ (('i' / 'I') ('n' / 'N')) __ literalList Action26))> */
		nil,
		/* 19 literalString <- <(<STRING> Action27)> */
		func() bool {
			position245, tokenIndex245, depth245 := position, tokenIndex, depth
			{
				position246 := position
				depth++
				{
					position247 := position
					depth++
					if !_rules[ruleSTRING]() {
						goto l245
					}
					depth--
					add(rulePegText, position247)
				}
				{
					add(ruleAction27, position)
				}
				depth--
				add(ruleliteralString, position246)
			}
			return true
		l245:
			position, tokenIndex, depth = position245, tokenIndex245, depth245
			return false
		},
		/* 20 literalList <- <(Action28 PAREN_OPEN literalListString (COMMA literalListString)* PAREN_CLOSE)> */
		nil,
		/* 21 literalListString <- <(STRING Action29)> */
		func() bool {
			position250, tokenIndex250, depth250 := position, tokenIndex, depth
			{
				position251 := position
				depth++
				if !_rules[ruleSTRING]() {
					goto l250
				}
				{
					add(ruleAction29, position)
				}
				depth--
				add(ruleliteralListString, position251)
			}
			return true
		l250:
			position, tokenIndex, depth = position250, tokenIndex250, depth250
			return false
		},
		/* 22 tagName <- <(<TAG_NAME> Action30)> */
		func() bool {
			position253, tokenIndex253, depth253 := position, tokenIndex, depth
			{
				position254 := position
				depth++
				{
					position255 := position
					depth++
					{
						position256 := position
						depth++
						if !_rules[ruleIDENTIFIER]() {
							goto l253
						}
						depth--
						add(ruleTAG_NAME, position256)
					}
					depth--
					add(rulePegText, position255)
				}
				{
					add(ruleAction30, position)
				}
				depth--
				add(ruletagName, position254)
			}
			return true
		l253:
			position, tokenIndex, depth = position253, tokenIndex253, depth253
			return false
		},
		/* 23 COLUMN_NAME <- <IDENTIFIER> */
		func() bool {
			position258, tokenIndex258, depth258 := position, tokenIndex, depth
			{
				position259 := position
				depth++
				if !_rules[ruleIDENTIFIER]() {
					goto l258
				}
				depth--
				add(ruleCOLUMN_NAME, position259)
			}
			return true
		l258:
			position, tokenIndex, depth = position258, tokenIndex258, depth258
			return false
		},
		/* 24 METRIC_NAME <- <IDENTIFIER> */
		nil,
		/* 25 TAG_NAME <- <IDENTIFIER> */
		nil,
		/* 26 TIMESTAMP <- <(NUMBER_NATURAL / STRING)> */
		func() bool {
			position262, tokenIndex262, depth262 := position, tokenIndex, depth
			{
				position263 := position
				depth++
				{
					position264, tokenIndex264, depth264 := position, tokenIndex, depth
					if !_rules[ruleNUMBER_NATURAL]() {
						goto l265
					}
					goto l264
				l265:
					position, tokenIndex, depth = position264, tokenIndex264, depth264
					if !_rules[ruleSTRING]() {
						goto l262
					}
				}
			l264:
				depth--
				add(ruleTIMESTAMP, position263)
			}
			return true
		l262:
			position, tokenIndex, depth = position262, tokenIndex262, depth262
			return false
		},
		/* 27 IDENTIFIER <- <(('`' CHAR* '`') / (!KEYWORD ID_SEGMENT ('.' ID_SEGMENT)*))> */
		func() bool {
			position266, tokenIndex266, depth266 := position, tokenIndex, depth
			{
				position267 := position
				depth++
				{
					position268, tokenIndex268, depth268 := position, tokenIndex, depth
					if buffer[position] != rune('`') {
						goto l269
					}
					position++
				l270:
					{
						position271, tokenIndex271, depth271 := position, tokenIndex, depth
						if !_rules[ruleCHAR]() {
							goto l271
						}
						goto l270
					l271:
						position, tokenIndex, depth = position271, tokenIndex271, depth271
					}
					if buffer[position] != rune('`') {
						goto l269
					}
					position++
					goto l268
				l269:
					position, tokenIndex, depth = position268, tokenIndex268, depth268
					{
						position272, tokenIndex272, depth272 := position, tokenIndex, depth
						{
							position273 := position
							depth++
							{
								position274, tokenIndex274, depth274 := position, tokenIndex, depth
								{
									position276, tokenIndex276, depth276 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l277
									}
									position++
									goto l276
								l277:
									position, tokenIndex, depth = position276, tokenIndex276, depth276
									if buffer[position] != rune('A') {
										goto l275
									}
									position++
								}
							l276:
								{
									position278, tokenIndex278, depth278 := position, tokenIndex, depth
									if buffer[position] != rune('l') {
										goto l279
									}
									position++
									goto l278
								l279:
									position, tokenIndex, depth = position278, tokenIndex278, depth278
									if buffer[position] != rune('L') {
										goto l275
									}
									position++
								}
							l278:
								{
									position280, tokenIndex280, depth280 := position, tokenIndex, depth
									if buffer[position] != rune('l') {
										goto l281
									}
									position++
									goto l280
								l281:
									position, tokenIndex, depth = position280, tokenIndex280, depth280
									if buffer[position] != rune('L') {
										goto l275
									}
									position++
								}
							l280:
								goto l274
							l275:
								position, tokenIndex, depth = position274, tokenIndex274, depth274
								{
									position283, tokenIndex283, depth283 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l284
									}
									position++
									goto l283
								l284:
									position, tokenIndex, depth = position283, tokenIndex283, depth283
									if buffer[position] != rune('A') {
										goto l282
									}
									position++
								}
							l283:
								{
									position285, tokenIndex285, depth285 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l286
									}
									position++
									goto l285
								l286:
									position, tokenIndex, depth = position285, tokenIndex285, depth285
									if buffer[position] != rune('N') {
										goto l282
									}
									position++
								}
							l285:
								{
									position287, tokenIndex287, depth287 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l288
									}
									position++
									goto l287
								l288:
									position, tokenIndex, depth = position287, tokenIndex287, depth287
									if buffer[position] != rune('D') {
										goto l282
									}
									position++
								}
							l287:
								goto l274
							l282:
								position, tokenIndex, depth = position274, tokenIndex274, depth274
								{
									switch buffer[position] {
									case 'W', 'w':
										{
											position290, tokenIndex290, depth290 := position, tokenIndex, depth
											if buffer[position] != rune('w') {
												goto l291
											}
											position++
											goto l290
										l291:
											position, tokenIndex, depth = position290, tokenIndex290, depth290
											if buffer[position] != rune('W') {
												goto l272
											}
											position++
										}
									l290:
										{
											position292, tokenIndex292, depth292 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l293
											}
											position++
											goto l292
										l293:
											position, tokenIndex, depth = position292, tokenIndex292, depth292
											if buffer[position] != rune('H') {
												goto l272
											}
											position++
										}
									l292:
										{
											position294, tokenIndex294, depth294 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l295
											}
											position++
											goto l294
										l295:
											position, tokenIndex, depth = position294, tokenIndex294, depth294
											if buffer[position] != rune('E') {
												goto l272
											}
											position++
										}
									l294:
										{
											position296, tokenIndex296, depth296 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l297
											}
											position++
											goto l296
										l297:
											position, tokenIndex, depth = position296, tokenIndex296, depth296
											if buffer[position] != rune('R') {
												goto l272
											}
											position++
										}
									l296:
										{
											position298, tokenIndex298, depth298 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l299
											}
											position++
											goto l298
										l299:
											position, tokenIndex, depth = position298, tokenIndex298, depth298
											if buffer[position] != rune('E') {
												goto l272
											}
											position++
										}
									l298:
										break
									case 'T', 't':
										{
											position300, tokenIndex300, depth300 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l301
											}
											position++
											goto l300
										l301:
											position, tokenIndex, depth = position300, tokenIndex300, depth300
											if buffer[position] != rune('T') {
												goto l272
											}
											position++
										}
									l300:
										{
											position302, tokenIndex302, depth302 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l303
											}
											position++
											goto l302
										l303:
											position, tokenIndex, depth = position302, tokenIndex302, depth302
											if buffer[position] != rune('O') {
												goto l272
											}
											position++
										}
									l302:
										break
									case 'S', 's':
										{
											position304, tokenIndex304, depth304 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l305
											}
											position++
											goto l304
										l305:
											position, tokenIndex, depth = position304, tokenIndex304, depth304
											if buffer[position] != rune('S') {
												goto l272
											}
											position++
										}
									l304:
										{
											position306, tokenIndex306, depth306 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l307
											}
											position++
											goto l306
										l307:
											position, tokenIndex, depth = position306, tokenIndex306, depth306
											if buffer[position] != rune('E') {
												goto l272
											}
											position++
										}
									l306:
										{
											position308, tokenIndex308, depth308 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l309
											}
											position++
											goto l308
										l309:
											position, tokenIndex, depth = position308, tokenIndex308, depth308
											if buffer[position] != rune('L') {
												goto l272
											}
											position++
										}
									l308:
										{
											position310, tokenIndex310, depth310 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l311
											}
											position++
											goto l310
										l311:
											position, tokenIndex, depth = position310, tokenIndex310, depth310
											if buffer[position] != rune('E') {
												goto l272
											}
											position++
										}
									l310:
										{
											position312, tokenIndex312, depth312 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l313
											}
											position++
											goto l312
										l313:
											position, tokenIndex, depth = position312, tokenIndex312, depth312
											if buffer[position] != rune('C') {
												goto l272
											}
											position++
										}
									l312:
										{
											position314, tokenIndex314, depth314 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l315
											}
											position++
											goto l314
										l315:
											position, tokenIndex, depth = position314, tokenIndex314, depth314
											if buffer[position] != rune('T') {
												goto l272
											}
											position++
										}
									l314:
										break
									case 'O', 'o':
										{
											position316, tokenIndex316, depth316 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l317
											}
											position++
											goto l316
										l317:
											position, tokenIndex, depth = position316, tokenIndex316, depth316
											if buffer[position] != rune('O') {
												goto l272
											}
											position++
										}
									l316:
										{
											position318, tokenIndex318, depth318 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l319
											}
											position++
											goto l318
										l319:
											position, tokenIndex, depth = position318, tokenIndex318, depth318
											if buffer[position] != rune('R') {
												goto l272
											}
											position++
										}
									l318:
										break
									case 'N', 'n':
										{
											position320, tokenIndex320, depth320 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l321
											}
											position++
											goto l320
										l321:
											position, tokenIndex, depth = position320, tokenIndex320, depth320
											if buffer[position] != rune('N') {
												goto l272
											}
											position++
										}
									l320:
										{
											position322, tokenIndex322, depth322 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l323
											}
											position++
											goto l322
										l323:
											position, tokenIndex, depth = position322, tokenIndex322, depth322
											if buffer[position] != rune('O') {
												goto l272
											}
											position++
										}
									l322:
										{
											position324, tokenIndex324, depth324 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l325
											}
											position++
											goto l324
										l325:
											position, tokenIndex, depth = position324, tokenIndex324, depth324
											if buffer[position] != rune('T') {
												goto l272
											}
											position++
										}
									l324:
										break
									case 'M', 'm':
										{
											position326, tokenIndex326, depth326 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l327
											}
											position++
											goto l326
										l327:
											position, tokenIndex, depth = position326, tokenIndex326, depth326
											if buffer[position] != rune('M') {
												goto l272
											}
											position++
										}
									l326:
										{
											position328, tokenIndex328, depth328 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l329
											}
											position++
											goto l328
										l329:
											position, tokenIndex, depth = position328, tokenIndex328, depth328
											if buffer[position] != rune('A') {
												goto l272
											}
											position++
										}
									l328:
										{
											position330, tokenIndex330, depth330 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l331
											}
											position++
											goto l330
										l331:
											position, tokenIndex, depth = position330, tokenIndex330, depth330
											if buffer[position] != rune('T') {
												goto l272
											}
											position++
										}
									l330:
										{
											position332, tokenIndex332, depth332 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l333
											}
											position++
											goto l332
										l333:
											position, tokenIndex, depth = position332, tokenIndex332, depth332
											if buffer[position] != rune('C') {
												goto l272
											}
											position++
										}
									l332:
										{
											position334, tokenIndex334, depth334 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l335
											}
											position++
											goto l334
										l335:
											position, tokenIndex, depth = position334, tokenIndex334, depth334
											if buffer[position] != rune('H') {
												goto l272
											}
											position++
										}
									l334:
										{
											position336, tokenIndex336, depth336 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l337
											}
											position++
											goto l336
										l337:
											position, tokenIndex, depth = position336, tokenIndex336, depth336
											if buffer[position] != rune('E') {
												goto l272
											}
											position++
										}
									l336:
										{
											position338, tokenIndex338, depth338 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l339
											}
											position++
											goto l338
										l339:
											position, tokenIndex, depth = position338, tokenIndex338, depth338
											if buffer[position] != rune('S') {
												goto l272
											}
											position++
										}
									l338:
										break
									case 'I', 'i':
										{
											position340, tokenIndex340, depth340 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l341
											}
											position++
											goto l340
										l341:
											position, tokenIndex, depth = position340, tokenIndex340, depth340
											if buffer[position] != rune('I') {
												goto l272
											}
											position++
										}
									l340:
										{
											position342, tokenIndex342, depth342 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l343
											}
											position++
											goto l342
										l343:
											position, tokenIndex, depth = position342, tokenIndex342, depth342
											if buffer[position] != rune('N') {
												goto l272
											}
											position++
										}
									l342:
										break
									case 'G', 'g':
										{
											position344, tokenIndex344, depth344 := position, tokenIndex, depth
											if buffer[position] != rune('g') {
												goto l345
											}
											position++
											goto l344
										l345:
											position, tokenIndex, depth = position344, tokenIndex344, depth344
											if buffer[position] != rune('G') {
												goto l272
											}
											position++
										}
									l344:
										{
											position346, tokenIndex346, depth346 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l347
											}
											position++
											goto l346
										l347:
											position, tokenIndex, depth = position346, tokenIndex346, depth346
											if buffer[position] != rune('R') {
												goto l272
											}
											position++
										}
									l346:
										{
											position348, tokenIndex348, depth348 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l349
											}
											position++
											goto l348
										l349:
											position, tokenIndex, depth = position348, tokenIndex348, depth348
											if buffer[position] != rune('O') {
												goto l272
											}
											position++
										}
									l348:
										{
											position350, tokenIndex350, depth350 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l351
											}
											position++
											goto l350
										l351:
											position, tokenIndex, depth = position350, tokenIndex350, depth350
											if buffer[position] != rune('U') {
												goto l272
											}
											position++
										}
									l350:
										{
											position352, tokenIndex352, depth352 := position, tokenIndex, depth
											if buffer[position] != rune('p') {
												goto l353
											}
											position++
											goto l352
										l353:
											position, tokenIndex, depth = position352, tokenIndex352, depth352
											if buffer[position] != rune('P') {
												goto l272
											}
											position++
										}
									l352:
										break
									case 'F', 'f':
										{
											position354, tokenIndex354, depth354 := position, tokenIndex, depth
											if buffer[position] != rune('f') {
												goto l355
											}
											position++
											goto l354
										l355:
											position, tokenIndex, depth = position354, tokenIndex354, depth354
											if buffer[position] != rune('F') {
												goto l272
											}
											position++
										}
									l354:
										{
											position356, tokenIndex356, depth356 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l357
											}
											position++
											goto l356
										l357:
											position, tokenIndex, depth = position356, tokenIndex356, depth356
											if buffer[position] != rune('R') {
												goto l272
											}
											position++
										}
									l356:
										{
											position358, tokenIndex358, depth358 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l359
											}
											position++
											goto l358
										l359:
											position, tokenIndex, depth = position358, tokenIndex358, depth358
											if buffer[position] != rune('O') {
												goto l272
											}
											position++
										}
									l358:
										{
											position360, tokenIndex360, depth360 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l361
											}
											position++
											goto l360
										l361:
											position, tokenIndex, depth = position360, tokenIndex360, depth360
											if buffer[position] != rune('M') {
												goto l272
											}
											position++
										}
									l360:
										break
									case 'D', 'd':
										{
											position362, tokenIndex362, depth362 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l363
											}
											position++
											goto l362
										l363:
											position, tokenIndex, depth = position362, tokenIndex362, depth362
											if buffer[position] != rune('D') {
												goto l272
											}
											position++
										}
									l362:
										{
											position364, tokenIndex364, depth364 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l365
											}
											position++
											goto l364
										l365:
											position, tokenIndex, depth = position364, tokenIndex364, depth364
											if buffer[position] != rune('E') {
												goto l272
											}
											position++
										}
									l364:
										{
											position366, tokenIndex366, depth366 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l367
											}
											position++
											goto l366
										l367:
											position, tokenIndex, depth = position366, tokenIndex366, depth366
											if buffer[position] != rune('S') {
												goto l272
											}
											position++
										}
									l366:
										{
											position368, tokenIndex368, depth368 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l369
											}
											position++
											goto l368
										l369:
											position, tokenIndex, depth = position368, tokenIndex368, depth368
											if buffer[position] != rune('C') {
												goto l272
											}
											position++
										}
									l368:
										{
											position370, tokenIndex370, depth370 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l371
											}
											position++
											goto l370
										l371:
											position, tokenIndex, depth = position370, tokenIndex370, depth370
											if buffer[position] != rune('R') {
												goto l272
											}
											position++
										}
									l370:
										{
											position372, tokenIndex372, depth372 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l373
											}
											position++
											goto l372
										l373:
											position, tokenIndex, depth = position372, tokenIndex372, depth372
											if buffer[position] != rune('I') {
												goto l272
											}
											position++
										}
									l372:
										{
											position374, tokenIndex374, depth374 := position, tokenIndex, depth
											if buffer[position] != rune('b') {
												goto l375
											}
											position++
											goto l374
										l375:
											position, tokenIndex, depth = position374, tokenIndex374, depth374
											if buffer[position] != rune('B') {
												goto l272
											}
											position++
										}
									l374:
										{
											position376, tokenIndex376, depth376 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l377
											}
											position++
											goto l376
										l377:
											position, tokenIndex, depth = position376, tokenIndex376, depth376
											if buffer[position] != rune('E') {
												goto l272
											}
											position++
										}
									l376:
										break
									case 'B', 'b':
										{
											position378, tokenIndex378, depth378 := position, tokenIndex, depth
											if buffer[position] != rune('b') {
												goto l379
											}
											position++
											goto l378
										l379:
											position, tokenIndex, depth = position378, tokenIndex378, depth378
											if buffer[position] != rune('B') {
												goto l272
											}
											position++
										}
									l378:
										{
											position380, tokenIndex380, depth380 := position, tokenIndex, depth
											if buffer[position] != rune('y') {
												goto l381
											}
											position++
											goto l380
										l381:
											position, tokenIndex, depth = position380, tokenIndex380, depth380
											if buffer[position] != rune('Y') {
												goto l272
											}
											position++
										}
									l380:
										break
									default:
										{
											position382, tokenIndex382, depth382 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l383
											}
											position++
											goto l382
										l383:
											position, tokenIndex, depth = position382, tokenIndex382, depth382
											if buffer[position] != rune('A') {
												goto l272
											}
											position++
										}
									l382:
										{
											position384, tokenIndex384, depth384 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l385
											}
											position++
											goto l384
										l385:
											position, tokenIndex, depth = position384, tokenIndex384, depth384
											if buffer[position] != rune('S') {
												goto l272
											}
											position++
										}
									l384:
										break
									}
								}

							}
						l274:
							depth--
							add(ruleKEYWORD, position273)
						}
						goto l266
					l272:
						position, tokenIndex, depth = position272, tokenIndex272, depth272
					}
					if !_rules[ruleID_SEGMENT]() {
						goto l266
					}
				l386:
					{
						position387, tokenIndex387, depth387 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l387
						}
						position++
						if !_rules[ruleID_SEGMENT]() {
							goto l387
						}
						goto l386
					l387:
						position, tokenIndex, depth = position387, tokenIndex387, depth387
					}
				}
			l268:
				depth--
				add(ruleIDENTIFIER, position267)
			}
			return true
		l266:
			position, tokenIndex, depth = position266, tokenIndex266, depth266
			return false
		},
		/* 28 ID_SEGMENT <- <(ID_START ID_CONT*)> */
		func() bool {
			position388, tokenIndex388, depth388 := position, tokenIndex, depth
			{
				position389 := position
				depth++
				if !_rules[ruleID_START]() {
					goto l388
				}
			l390:
				{
					position391, tokenIndex391, depth391 := position, tokenIndex, depth
					{
						position392 := position
						depth++
						{
							position393, tokenIndex393, depth393 := position, tokenIndex, depth
							if !_rules[ruleID_START]() {
								goto l394
							}
							goto l393
						l394:
							position, tokenIndex, depth = position393, tokenIndex393, depth393
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l391
							}
							position++
						}
					l393:
						depth--
						add(ruleID_CONT, position392)
					}
					goto l390
				l391:
					position, tokenIndex, depth = position391, tokenIndex391, depth391
				}
				depth--
				add(ruleID_SEGMENT, position389)
			}
			return true
		l388:
			position, tokenIndex, depth = position388, tokenIndex388, depth388
			return false
		},
		/* 29 ID_START <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position395, tokenIndex395, depth395 := position, tokenIndex, depth
			{
				position396 := position
				depth++
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l395
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l395
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l395
						}
						position++
						break
					}
				}

				depth--
				add(ruleID_START, position396)
			}
			return true
		l395:
			position, tokenIndex, depth = position395, tokenIndex395, depth395
			return false
		},
		/* 30 ID_CONT <- <(ID_START / [0-9])> */
		nil,
		/* 31 KEYWORD <- <((('a' / 'A') ('l' / 'L') ('l' / 'L')) / (('a' / 'A') ('n' / 'N') ('d' / 'D')) / ((&('W' | 'w') (('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E'))) | (&('T' | 't') (('t' / 'T') ('o' / 'O'))) | (&('S' | 's') (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T'))) | (&('O' | 'o') (('o' / 'O') ('r' / 'R'))) | (&('N' | 'n') (('n' / 'N') ('o' / 'O') ('t' / 'T'))) | (&('M' | 'm') (('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S'))) | (&('I' | 'i') (('i' / 'I') ('n' / 'N'))) | (&('G' | 'g') (('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P'))) | (&('F' | 'f') (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M'))) | (&('D' | 'd') (('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E'))) | (&('B' | 'b') (('b' / 'B') ('y' / 'Y'))) | (&('A' | 'a') (('a' / 'A') ('s' / 'S')))))> */
		nil,
		/* 32 OP_ADD <- <(_ '+' _)> */
		nil,
		/* 33 OP_SUB <- <(_ '-' _)> */
		nil,
		/* 34 OP_MULT <- <(_ '*' _)> */
		nil,
		/* 35 OP_DIV <- <(_ '/' _)> */
		nil,
		/* 36 OP_AND <- <(__ (('a' / 'A') ('n' / 'N') ('d' / 'D')) __)> */
		nil,
		/* 37 OP_OR <- <(__ (('o' / 'O') ('r' / 'R')) __)> */
		nil,
		/* 38 OP_NOT <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') __)> */
		nil,
		/* 39 STRING <- <(('\'' CHAR* '\'') / ('"' CHAR* '"'))> */
		func() bool {
			position407, tokenIndex407, depth407 := position, tokenIndex, depth
			{
				position408 := position
				depth++
				{
					position409, tokenIndex409, depth409 := position, tokenIndex, depth
					if buffer[position] != rune('\'') {
						goto l410
					}
					position++
				l411:
					{
						position412, tokenIndex412, depth412 := position, tokenIndex, depth
						if !_rules[ruleCHAR]() {
							goto l412
						}
						goto l411
					l412:
						position, tokenIndex, depth = position412, tokenIndex412, depth412
					}
					if buffer[position] != rune('\'') {
						goto l410
					}
					position++
					goto l409
				l410:
					position, tokenIndex, depth = position409, tokenIndex409, depth409
					if buffer[position] != rune('"') {
						goto l407
					}
					position++
				l413:
					{
						position414, tokenIndex414, depth414 := position, tokenIndex, depth
						if !_rules[ruleCHAR]() {
							goto l414
						}
						goto l413
					l414:
						position, tokenIndex, depth = position414, tokenIndex414, depth414
					}
					if buffer[position] != rune('"') {
						goto l407
					}
					position++
				}
			l409:
				depth--
				add(ruleSTRING, position408)
			}
			return true
		l407:
			position, tokenIndex, depth = position407, tokenIndex407, depth407
			return false
		},
		/* 40 CHAR <- <(('\\' ESCAPE_CLASS) / (!ESCAPE_CLASS .))> */
		func() bool {
			position415, tokenIndex415, depth415 := position, tokenIndex, depth
			{
				position416 := position
				depth++
				{
					position417, tokenIndex417, depth417 := position, tokenIndex, depth
					if buffer[position] != rune('\\') {
						goto l418
					}
					position++
					if !_rules[ruleESCAPE_CLASS]() {
						goto l418
					}
					goto l417
				l418:
					position, tokenIndex, depth = position417, tokenIndex417, depth417
					{
						position419, tokenIndex419, depth419 := position, tokenIndex, depth
						if !_rules[ruleESCAPE_CLASS]() {
							goto l419
						}
						goto l415
					l419:
						position, tokenIndex, depth = position419, tokenIndex419, depth419
					}
					if !matchDot() {
						goto l415
					}
				}
			l417:
				depth--
				add(ruleCHAR, position416)
			}
			return true
		l415:
			position, tokenIndex, depth = position415, tokenIndex415, depth415
			return false
		},
		/* 41 ESCAPE_CLASS <- <((&('\\') '\\') | (&('"') '"') | (&('`') '`') | (&('\'') '\''))> */
		func() bool {
			position420, tokenIndex420, depth420 := position, tokenIndex, depth
			{
				position421 := position
				depth++
				{
					switch buffer[position] {
					case '\\':
						if buffer[position] != rune('\\') {
							goto l420
						}
						position++
						break
					case '"':
						if buffer[position] != rune('"') {
							goto l420
						}
						position++
						break
					case '`':
						if buffer[position] != rune('`') {
							goto l420
						}
						position++
						break
					default:
						if buffer[position] != rune('\'') {
							goto l420
						}
						position++
						break
					}
				}

				depth--
				add(ruleESCAPE_CLASS, position421)
			}
			return true
		l420:
			position, tokenIndex, depth = position420, tokenIndex420, depth420
			return false
		},
		/* 42 NUMBER <- <(NUMBER_INTEGER NUMBER_FRACTION? NUMBER_EXP?)> */
		nil,
		/* 43 NUMBER_NATURAL <- <('0' / ([1-9] [0-9]*))> */
		func() bool {
			position424, tokenIndex424, depth424 := position, tokenIndex, depth
			{
				position425 := position
				depth++
				{
					position426, tokenIndex426, depth426 := position, tokenIndex, depth
					if buffer[position] != rune('0') {
						goto l427
					}
					position++
					goto l426
				l427:
					position, tokenIndex, depth = position426, tokenIndex426, depth426
					if c := buffer[position]; c < rune('1') || c > rune('9') {
						goto l424
					}
					position++
				l428:
					{
						position429, tokenIndex429, depth429 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l429
						}
						position++
						goto l428
					l429:
						position, tokenIndex, depth = position429, tokenIndex429, depth429
					}
				}
			l426:
				depth--
				add(ruleNUMBER_NATURAL, position425)
			}
			return true
		l424:
			position, tokenIndex, depth = position424, tokenIndex424, depth424
			return false
		},
		/* 44 NUMBER_FRACTION <- <('.' [0-9]+)> */
		nil,
		/* 45 NUMBER_INTEGER <- <('-'? NUMBER_NATURAL)> */
		nil,
		/* 46 NUMBER_EXP <- <(('e' / 'E') ('+' / '-')? [0-9]+)> */
		nil,
		/* 47 PAREN_OPEN <- <(_ '(' _)> */
		func() bool {
			position433, tokenIndex433, depth433 := position, tokenIndex, depth
			{
				position434 := position
				depth++
				if !_rules[rule_]() {
					goto l433
				}
				if buffer[position] != rune('(') {
					goto l433
				}
				position++
				if !_rules[rule_]() {
					goto l433
				}
				depth--
				add(rulePAREN_OPEN, position434)
			}
			return true
		l433:
			position, tokenIndex, depth = position433, tokenIndex433, depth433
			return false
		},
		/* 48 PAREN_CLOSE <- <(_ ')' _)> */
		func() bool {
			position435, tokenIndex435, depth435 := position, tokenIndex, depth
			{
				position436 := position
				depth++
				if !_rules[rule_]() {
					goto l435
				}
				if buffer[position] != rune(')') {
					goto l435
				}
				position++
				if !_rules[rule_]() {
					goto l435
				}
				depth--
				add(rulePAREN_CLOSE, position436)
			}
			return true
		l435:
			position, tokenIndex, depth = position435, tokenIndex435, depth435
			return false
		},
		/* 49 COMMA <- <(_ ',' _)> */
		func() bool {
			position437, tokenIndex437, depth437 := position, tokenIndex, depth
			{
				position438 := position
				depth++
				if !_rules[rule_]() {
					goto l437
				}
				if buffer[position] != rune(',') {
					goto l437
				}
				position++
				if !_rules[rule_]() {
					goto l437
				}
				depth--
				add(ruleCOMMA, position438)
			}
			return true
		l437:
			position, tokenIndex, depth = position437, tokenIndex437, depth437
			return false
		},
		/* 50 _ <- <SPACE*> */
		func() bool {
			{
				position440 := position
				depth++
			l441:
				{
					position442, tokenIndex442, depth442 := position, tokenIndex, depth
					if !_rules[ruleSPACE]() {
						goto l442
					}
					goto l441
				l442:
					position, tokenIndex, depth = position442, tokenIndex442, depth442
				}
				depth--
				add(rule_, position440)
			}
			return true
		},
		/* 51 __ <- <SPACE+> */
		func() bool {
			position443, tokenIndex443, depth443 := position, tokenIndex, depth
			{
				position444 := position
				depth++
				if !_rules[ruleSPACE]() {
					goto l443
				}
			l445:
				{
					position446, tokenIndex446, depth446 := position, tokenIndex, depth
					if !_rules[ruleSPACE]() {
						goto l446
					}
					goto l445
				l446:
					position, tokenIndex, depth = position446, tokenIndex446, depth446
				}
				depth--
				add(rule__, position444)
			}
			return true
		l443:
			position, tokenIndex, depth = position443, tokenIndex443, depth443
			return false
		},
		/* 52 SPACE <- <((&('\t') '\t') | (&('\n') '\n') | (&(' ') ' '))> */
		func() bool {
			position447, tokenIndex447, depth447 := position, tokenIndex, depth
			{
				position448 := position
				depth++
				{
					switch buffer[position] {
					case '\t':
						if buffer[position] != rune('\t') {
							goto l447
						}
						position++
						break
					case '\n':
						if buffer[position] != rune('\n') {
							goto l447
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l447
						}
						position++
						break
					}
				}

				depth--
				add(ruleSPACE, position448)
			}
			return true
		l447:
			position, tokenIndex, depth = position447, tokenIndex447, depth447
			return false
		},
		/* 54 Action0 <- <{
		   p.makeSelect()
		 }> */
		nil,
		/* 55 Action1 <- <{ p.makeDescribeAll() }> */
		nil,
		nil,
		/* 57 Action2 <- <{ p.addStringLiteral(unescapeLiteral(buffer[begin:end])) }> */
		nil,
		/* 58 Action3 <- <{ p.makeDescribe() }> */
		nil,
		/* 59 Action4 <- <{ p.addNullPredicate() }> */
		nil,
		/* 60 Action5 <- <{ p.addExpressionList() }> */
		nil,
		/* 61 Action6 <- <{ p.appendExpression() }> */
		nil,
		/* 62 Action7 <- <{ p.appendExpression() }> */
		nil,
		/* 63 Action8 <- <{ p.addOperatorLiteral("*") }> */
		nil,
		/* 64 Action9 <- <{ p.addOperatorLiteral("-") }> */
		nil,
		/* 65 Action10 <- <{ p.addOperatorFunction() }> */
		nil,
		/* 66 Action11 <- <{ p.addOperatorLiteral("*") }> */
		nil,
		/* 67 Action12 <- <{ p.addOperatorLiteral("*") }> */
		nil,
		/* 68 Action13 <- <{ p.addOperatorFunction() }> */
		nil,
		/* 69 Action14 <- <{ p.addNumberNode(buffer[begin:end]) }> */
		nil,
		/* 70 Action15 <- <{
		   p.addStringLiteral(unescapeLiteral(buffer[begin:end]))
		 }> */
		nil,
		/* 71 Action16 <- <{
		   p.addFunctionInvocation()
		 }> */
		nil,
		/* 72 Action17 <- <{
		   p.addStringLiteral(unescapeLiteral(buffer[begin:end]))
		 }> */
		nil,
		/* 73 Action18 <- <{ p.addNullPredicate() }> */
		nil,
		/* 74 Action19 <- <{
		   p.addMetricExpression()
		 }> */
		nil,
		/* 75 Action20 <- <{ p.addAndPredicate() }> */
		nil,
		/* 76 Action21 <- <{ p.addOrPredicate() }> */
		nil,
		/* 77 Action22 <- <{ p.addNotPredicate() }> */
		nil,
		/* 78 Action23 <- <{
		   p.addLiteralMatcher()
		 }> */
		nil,
		/* 79 Action24 <- <{
		   p.addLiteralMatcher()
		   p.addNotPredicate()
		 }> */
		nil,
		/* 80 Action25 <- <{
		   p.addRegexMatcher()
		 }> */
		nil,
		/* 81 Action26 <- <{
		   p.addListMatcher()
		 }> */
		nil,
		/* 82 Action27 <- <{
		  p.addStringLiteral(unescapeLiteral(buffer[begin:end]))
		}> */
		nil,
		/* 83 Action28 <- <{ p.addLiteralListNode() }> */
		nil,
		/* 84 Action29 <- <{
		  p.appendLiteral(unescapeLiteral(buffer[begin:end]))
		}> */
		nil,
		/* 85 Action30 <- <{ p.addTagLiteral(unescapeLiteral(buffer[begin:end])) }> */
		nil,
	}
	p.rules = _rules
}
