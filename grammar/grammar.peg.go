package grammar

//go:generate peg grammar.peg

import (
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleSciName
	ruleUninomial
	ruleAuthor
	ruleYear
	ruleCapWord
	ruleWord
	rulenum
	rulehASCII
	rulelASCII
	rule_
)

var rul3s = [...]string{
	"Unknown",
	"SciName",
	"Uninomial",
	"Author",
	"Year",
	"CapWord",
	"Word",
	"num",
	"hASCII",
	"lASCII",
	"_",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
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
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type GNParser struct {
	Buffer string
	buffer []rune
	rules  [11]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *GNParser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *GNParser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
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
	p   *GNParser
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *GNParser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *GNParser) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func (p *GNParser) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
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
		/* 0 SciName <- <(Uninomial !.)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[ruleUninomial]() {
					goto l0
				}
				{
					position2, tokenIndex2 := position, tokenIndex
					if !matchDot() {
						goto l2
					}
					goto l0
				l2:
					position, tokenIndex = position2, tokenIndex2
				}
				add(ruleSciName, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Uninomial <- <(CapWord (_ Author)?)> */
		func() bool {
			position3, tokenIndex3 := position, tokenIndex
			{
				position4 := position
				if !_rules[ruleCapWord]() {
					goto l3
				}
				{
					position5, tokenIndex5 := position, tokenIndex
					if !_rules[rule_]() {
						goto l5
					}
					if !_rules[ruleAuthor]() {
						goto l5
					}
					goto l6
				l5:
					position, tokenIndex = position5, tokenIndex5
				}
			l6:
				add(ruleUninomial, position4)
			}
			return true
		l3:
			position, tokenIndex = position3, tokenIndex3
			return false
		},
		/* 2 Author <- <(CapWord (_? Year)?)> */
		func() bool {
			position7, tokenIndex7 := position, tokenIndex
			{
				position8 := position
				if !_rules[ruleCapWord]() {
					goto l7
				}
				{
					position9, tokenIndex9 := position, tokenIndex
					{
						position11, tokenIndex11 := position, tokenIndex
						if !_rules[rule_]() {
							goto l11
						}
						goto l12
					l11:
						position, tokenIndex = position11, tokenIndex11
					}
				l12:
					if !_rules[ruleYear]() {
						goto l9
					}
					goto l10
				l9:
					position, tokenIndex = position9, tokenIndex9
				}
			l10:
				add(ruleAuthor, position8)
			}
			return true
		l7:
			position, tokenIndex = position7, tokenIndex7
			return false
		},
		/* 3 Year <- <(('1' / '2') ('0' / '7' / '8' / '9') num num)> */
		func() bool {
			position13, tokenIndex13 := position, tokenIndex
			{
				position14 := position
				{
					position15, tokenIndex15 := position, tokenIndex
					if buffer[position] != rune('1') {
						goto l16
					}
					position++
					goto l15
				l16:
					position, tokenIndex = position15, tokenIndex15
					if buffer[position] != rune('2') {
						goto l13
					}
					position++
				}
			l15:
				{
					position17, tokenIndex17 := position, tokenIndex
					if buffer[position] != rune('0') {
						goto l18
					}
					position++
					goto l17
				l18:
					position, tokenIndex = position17, tokenIndex17
					if buffer[position] != rune('7') {
						goto l19
					}
					position++
					goto l17
				l19:
					position, tokenIndex = position17, tokenIndex17
					if buffer[position] != rune('8') {
						goto l20
					}
					position++
					goto l17
				l20:
					position, tokenIndex = position17, tokenIndex17
					if buffer[position] != rune('9') {
						goto l13
					}
					position++
				}
			l17:
				if !_rules[rulenum]() {
					goto l13
				}
				if !_rules[rulenum]() {
					goto l13
				}
				add(ruleYear, position14)
			}
			return true
		l13:
			position, tokenIndex = position13, tokenIndex13
			return false
		},
		/* 4 CapWord <- <(hASCII Word)> */
		func() bool {
			position21, tokenIndex21 := position, tokenIndex
			{
				position22 := position
				if !_rules[rulehASCII]() {
					goto l21
				}
				if !_rules[ruleWord]() {
					goto l21
				}
				add(ruleCapWord, position22)
			}
			return true
		l21:
			position, tokenIndex = position21, tokenIndex21
			return false
		},
		/* 5 Word <- <lASCII+> */
		func() bool {
			position23, tokenIndex23 := position, tokenIndex
			{
				position24 := position
				if !_rules[rulelASCII]() {
					goto l23
				}
			l25:
				{
					position26, tokenIndex26 := position, tokenIndex
					if !_rules[rulelASCII]() {
						goto l26
					}
					goto l25
				l26:
					position, tokenIndex = position26, tokenIndex26
				}
				add(ruleWord, position24)
			}
			return true
		l23:
			position, tokenIndex = position23, tokenIndex23
			return false
		},
		/* 6 num <- <[0-9]> */
		func() bool {
			position27, tokenIndex27 := position, tokenIndex
			{
				position28 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l27
				}
				position++
				add(rulenum, position28)
			}
			return true
		l27:
			position, tokenIndex = position27, tokenIndex27
			return false
		},
		/* 7 hASCII <- <[A-Z]> */
		func() bool {
			position29, tokenIndex29 := position, tokenIndex
			{
				position30 := position
				if c := buffer[position]; c < rune('A') || c > rune('Z') {
					goto l29
				}
				position++
				add(rulehASCII, position30)
			}
			return true
		l29:
			position, tokenIndex = position29, tokenIndex29
			return false
		},
		/* 8 lASCII <- <[a-z]> */
		func() bool {
			position31, tokenIndex31 := position, tokenIndex
			{
				position32 := position
				if c := buffer[position]; c < rune('a') || c > rune('z') {
					goto l31
				}
				position++
				add(rulelASCII, position32)
			}
			return true
		l31:
			position, tokenIndex = position31, tokenIndex31
			return false
		},
		/* 9 _ <- <' '+> */
		func() bool {
			position33, tokenIndex33 := position, tokenIndex
			{
				position34 := position
				if buffer[position] != rune(' ') {
					goto l33
				}
				position++
			l35:
				{
					position36, tokenIndex36 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l36
					}
					position++
					goto l35
				l36:
					position, tokenIndex = position36, tokenIndex36
				}
				add(rule_, position34)
			}
			return true
		l33:
			position, tokenIndex = position33, tokenIndex33
			return false
		},
	}
	p.rules = _rules
}
