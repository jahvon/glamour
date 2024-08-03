package ansi

import (
	"io"

	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/termenv"
)

// MarginWriter is a Writer that applies indentation and padding around
// whatever you write to it.
type MarginWriter struct {
	indentation, margin  uint
	indentPos, marginPos uint
	indentToken          string

	profile termenv.Profile
	rules   StylePrimitive

	w  io.Writer
	pw *padding.Writer
	iw *indent.Writer
}

// NewMarginWriter returns a new MarginWriter.
func NewMarginWriter(ctx RenderContext, w io.Writer, rules StyleBlock) *MarginWriter {
	mw := &MarginWriter{w: w}
	bs := ctx.blockStack

	if rules.Indent != nil {
		mw.indentation = *rules.Indent
		mw.indentToken = " "
		if rules.IndentToken != nil {
			mw.indentToken = *rules.IndentToken
		}
	}
	if rules.Margin != nil {
		mw.margin = *rules.Margin
	}

	mw.pw = padding.NewWriterPipe(w, bs.Width(ctx), func(wr io.Writer) {
		renderText(w, ctx.options.ColorProfile, rules.StylePrimitive, " ")
	})

	mw.iw = indent.NewWriterPipe(mw.pw, mw.indentation+mw.margin, mw.indentFunc)
	return mw
}

func (w *MarginWriter) Write(b []byte) (int, error) {
	return w.iw.Write(b)
}

// indentFunc is called when writing each the margin and indentation tokens.
// The margin is written first, using an empty space character as the token.
// The indentation is written next, using the token specified in the rules.
func (mw *MarginWriter) indentFunc(w io.Writer) {
	ic := " "
	switch {
	case mw.margin == 0 && mw.indentation == 0:
		return
	case mw.margin >= 1 && mw.indentation == 0:
		break
	case mw.margin >= 1 && mw.marginPos < mw.margin:
		mw.marginPos++
	case mw.indentation >= 1 && mw.indentPos < mw.indentation:
		mw.indentPos++
		ic = mw.indentToken
		if mw.indentPos == mw.indentation {
			mw.marginPos = 0
			mw.indentPos = 0
		}
	}
	renderText(w, mw.profile, mw.rules, ic)
}
