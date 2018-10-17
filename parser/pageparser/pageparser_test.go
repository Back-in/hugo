// Copyright 2018 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pageparser

import (
	"testing"
)

type shortCodeLexerTest struct {
	name  string
	input string
	items []Item
}

var (
	tstEOF       = Item{tEOF, 0, ""}
	tstLeftNoMD  = Item{tLeftDelimScNoMarkup, 0, "{{<"}
	tstRightNoMD = Item{tRightDelimScNoMarkup, 0, ">}}"}
	tstLeftMD    = Item{tLeftDelimScWithMarkup, 0, "{{%"}
	tstRightMD   = Item{tRightDelimScWithMarkup, 0, "%}}"}
	tstSCClose   = Item{tScClose, 0, "/"}
	tstSC1       = Item{tScName, 0, "sc1"}
	tstSC2       = Item{tScName, 0, "sc2"}
	tstSC3       = Item{tScName, 0, "sc3"}
	tstSCSlash   = Item{tScName, 0, "sc/sub"}
	tstParam1    = Item{tScParam, 0, "param1"}
	tstParam2    = Item{tScParam, 0, "param2"}
	tstVal       = Item{tScParamVal, 0, "Hello World"}
)

var shortCodeLexerTests = []shortCodeLexerTest{
	{"empty", "", []Item{tstEOF}},
	{"spaces", " \t\n", []Item{{tText, 0, " \t\n"}, tstEOF}},
	{"text", `to be or not`, []Item{{tText, 0, "to be or not"}, tstEOF}},
	{"no markup", `{{< sc1 >}}`, []Item{tstLeftNoMD, tstSC1, tstRightNoMD, tstEOF}},
	{"with EOL", "{{< sc1 \n >}}", []Item{tstLeftNoMD, tstSC1, tstRightNoMD, tstEOF}},

	{"forward slash inside name", `{{< sc/sub >}}`, []Item{tstLeftNoMD, tstSCSlash, tstRightNoMD, tstEOF}},

	{"simple with markup", `{{% sc1 %}}`, []Item{tstLeftMD, tstSC1, tstRightMD, tstEOF}},
	{"with spaces", `{{<     sc1     >}}`, []Item{tstLeftNoMD, tstSC1, tstRightNoMD, tstEOF}},
	{"mismatched rightDelim", `{{< sc1 %}}`, []Item{tstLeftNoMD, tstSC1,
		{tError, 0, "unrecognized character in shortcode action: U+0025 '%'. Note: Parameters with non-alphanumeric args must be quoted"}}},
	{"inner, markup", `{{% sc1 %}} inner {{% /sc1 %}}`, []Item{
		tstLeftMD,
		tstSC1,
		tstRightMD,
		{tText, 0, " inner "},
		tstLeftMD,
		tstSCClose,
		tstSC1,
		tstRightMD,
		tstEOF,
	}},
	{"close, but no open", `{{< /sc1 >}}`, []Item{
		tstLeftNoMD, {tError, 0, "got closing shortcode, but none is open"}}},
	{"close wrong", `{{< sc1 >}}{{< /another >}}`, []Item{
		tstLeftNoMD, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose,
		{tError, 0, "closing tag for shortcode 'another' does not match start tag"}}},
	{"close, but no open, more", `{{< sc1 >}}{{< /sc1 >}}{{< /another >}}`, []Item{
		tstLeftNoMD, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose,
		{tError, 0, "closing tag for shortcode 'another' does not match start tag"}}},
	{"close with extra keyword", `{{< sc1 >}}{{< /sc1 keyword>}}`, []Item{
		tstLeftNoMD, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose, tstSC1,
		{tError, 0, "unclosed shortcode"}}},
	{"Youtube id", `{{< sc1 -ziL-Q_456igdO-4 >}}`, []Item{
		tstLeftNoMD, tstSC1, {tScParam, 0, "-ziL-Q_456igdO-4"}, tstRightNoMD, tstEOF}},
	{"non-alphanumerics param quoted", `{{< sc1 "-ziL-.%QigdO-4" >}}`, []Item{
		tstLeftNoMD, tstSC1, {tScParam, 0, "-ziL-.%QigdO-4"}, tstRightNoMD, tstEOF}},

	{"two params", `{{< sc1 param1   param2 >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstParam2, tstRightNoMD, tstEOF}},
	// issue #934
	{"self-closing", `{{< sc1 />}}`, []Item{
		tstLeftNoMD, tstSC1, tstSCClose, tstRightNoMD, tstEOF}},
	// Issue 2498
	{"multiple self-closing", `{{< sc1 />}}{{< sc1 />}}`, []Item{
		tstLeftNoMD, tstSC1, tstSCClose, tstRightNoMD,
		tstLeftNoMD, tstSC1, tstSCClose, tstRightNoMD, tstEOF}},
	{"self-closing with param", `{{< sc1 param1 />}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstSCClose, tstRightNoMD, tstEOF}},
	{"multiple self-closing with param", `{{< sc1 param1 />}}{{< sc1 param1 />}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstSCClose, tstRightNoMD,
		tstLeftNoMD, tstSC1, tstParam1, tstSCClose, tstRightNoMD, tstEOF}},
	{"multiple different self-closing with param", `{{< sc1 param1 />}}{{< sc2 param1 />}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstSCClose, tstRightNoMD,
		tstLeftNoMD, tstSC2, tstParam1, tstSCClose, tstRightNoMD, tstEOF}},
	{"nested simple", `{{< sc1 >}}{{< sc2 >}}{{< /sc1 >}}`, []Item{
		tstLeftNoMD, tstSC1, tstRightNoMD,
		tstLeftNoMD, tstSC2, tstRightNoMD,
		tstLeftNoMD, tstSCClose, tstSC1, tstRightNoMD, tstEOF}},
	{"nested complex", `{{< sc1 >}}ab{{% sc2 param1 %}}cd{{< sc3 >}}ef{{< /sc3 >}}gh{{% /sc2 %}}ij{{< /sc1 >}}kl`, []Item{
		tstLeftNoMD, tstSC1, tstRightNoMD,
		{tText, 0, "ab"},
		tstLeftMD, tstSC2, tstParam1, tstRightMD,
		{tText, 0, "cd"},
		tstLeftNoMD, tstSC3, tstRightNoMD,
		{tText, 0, "ef"},
		tstLeftNoMD, tstSCClose, tstSC3, tstRightNoMD,
		{tText, 0, "gh"},
		tstLeftMD, tstSCClose, tstSC2, tstRightMD,
		{tText, 0, "ij"},
		tstLeftNoMD, tstSCClose, tstSC1, tstRightNoMD,
		{tText, 0, "kl"}, tstEOF,
	}},

	{"two quoted params", `{{< sc1 "param nr. 1" "param nr. 2" >}}`, []Item{
		tstLeftNoMD, tstSC1, {tScParam, 0, "param nr. 1"}, {tScParam, 0, "param nr. 2"}, tstRightNoMD, tstEOF}},
	{"two named params", `{{< sc1 param1="Hello World" param2="p2Val">}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstVal, tstParam2, {tScParamVal, 0, "p2Val"}, tstRightNoMD, tstEOF}},
	{"escaped quotes", `{{< sc1 param1=\"Hello World\"  >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstVal, tstRightNoMD, tstEOF}},
	{"escaped quotes, positional param", `{{< sc1 \"param1\"  >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstRightNoMD, tstEOF}},
	{"escaped quotes inside escaped quotes", `{{< sc1 param1=\"Hello \"escaped\" World\"  >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1,
		{tScParamVal, 0, `Hello `}, {tError, 0, `got positional parameter 'escaped'. Cannot mix named and positional parameters`}}},
	{"escaped quotes inside nonescaped quotes",
		`{{< sc1 param1="Hello \"escaped\" World"  >}}`, []Item{
			tstLeftNoMD, tstSC1, tstParam1, {tScParamVal, 0, `Hello "escaped" World`}, tstRightNoMD, tstEOF}},
	{"escaped quotes inside nonescaped quotes in positional param",
		`{{< sc1 "Hello \"escaped\" World"  >}}`, []Item{
			tstLeftNoMD, tstSC1, {tScParam, 0, `Hello "escaped" World`}, tstRightNoMD, tstEOF}},
	{"unterminated quote", `{{< sc1 param2="Hello World>}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam2, {tError, 0, "unterminated quoted string in shortcode parameter-argument: 'Hello World>}}'"}}},
	{"one named param, one not", `{{< sc1 param1="Hello World" p2 >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstVal,
		{tError, 0, "got positional parameter 'p2'. Cannot mix named and positional parameters"}}},
	{"one named param, one quoted positional param", `{{< sc1 param1="Hello World" "And Universe" >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstVal,
		{tError, 0, "got quoted positional parameter. Cannot mix named and positional parameters"}}},
	{"one quoted positional param, one named param", `{{< sc1 "param1" param2="And Universe" >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1,
		{tError, 0, "got named parameter 'param2'. Cannot mix named and positional parameters"}}},
	{"ono positional param, one not", `{{< sc1 param1 param2="Hello World">}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1,
		{tError, 0, "got named parameter 'param2'. Cannot mix named and positional parameters"}}},
	{"commented out", `{{</* sc1 */>}}`, []Item{
		{tText, 0, "{{<"}, {tText, 0, " sc1 "}, {tText, 0, ">}}"}, tstEOF}},
	{"commented out, with asterisk inside", `{{</* sc1 "**/*.pdf" */>}}`, []Item{
		{tText, 0, "{{<"}, {tText, 0, " sc1 \"**/*.pdf\" "}, {tText, 0, ">}}"}, tstEOF}},
	{"commented out, missing close", `{{</* sc1 >}}`, []Item{
		{tError, 0, "comment must be closed"}}},
	{"commented out, misplaced close", `{{</* sc1 >}}*/`, []Item{
		{tError, 0, "comment must be closed"}}},
}

func TestShortcodeLexer(t *testing.T) {
	t.Parallel()
	for i, test := range shortCodeLexerTests {
		items := collect(&test)
		if !equal(items, test.items) {
			t.Errorf("[%d] %s: got\n\t%v\nexpected\n\t%v", i, test.name, items, test.items)
		}
	}
}

func BenchmarkShortcodeLexer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, test := range shortCodeLexerTests {
			items := collect(&test)
			if !equal(items, test.items) {
				b.Errorf("%s: got\n\t%v\nexpected\n\t%v", test.name, items, test.items)
			}
		}
	}
}

func collect(t *shortCodeLexerTest) (items []Item) {
	l := newPageLexer(t.name, t.input, 0).run()
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == tEOF || item.typ == tError {
			break
		}
	}
	return
}

// no positional checking, for now ...
func equal(i1, i2 []Item) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].Val != i2[k].Val {
			return false
		}
	}
	return true
}