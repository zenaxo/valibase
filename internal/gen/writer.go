package gen

import "strings"

// tsw is a tiny indent-aware TypeScript writer used to generate readable TS output
// without excessive manual string concatenation.
type tsw struct {
	b      strings.Builder
	indent int
}

func (w *tsw) W(s string) { w.b.WriteString(s) }

func (w *tsw) WL(line string) {
	if line == "" {
		w.b.WriteString("\n")
		return
	}
	w.b.WriteString(strings.Repeat("\t", w.indent))
	w.b.WriteString(line)
	w.b.WriteString("\n")
}

func (w *tsw) Indent()        { w.indent++ }
func (w *tsw) Dedent()        { w.indent-- }
func (w *tsw) String() string { return w.b.String() }
