package money

import (
	"html/template"
	"strconv"
	"strings"
)

// Formatter stores Money formatting information
type Formatter struct {
	Fraction int
	Decimal  string
	Thousand string
	Grapheme string
	Template string
}

// NewFormatter creates new Formatter instance
func NewFormatter(fraction int, decimal, thousand, grapheme, template string) *Formatter {
	return &Formatter{
		Fraction: fraction,
		Decimal:  decimal,
		Thousand: thousand,
		Grapheme: grapheme,
		Template: template,
	}
}

// Format returns string of formatted integer using given currency template
func (f *Formatter) Format(amount int64) string {
	// Work with absolute amount value
	sa := strconv.FormatInt(f.abs(amount), 10)

	if len(sa) <= f.Fraction {
		sa = strings.Repeat("0", f.Fraction-len(sa)+1) + sa
	}

	if f.Thousand != "" {
		for i := len(sa) - f.Fraction - 3; i > 0; i -= 3 {
			sa = sa[:i] + f.Thousand + sa[i:]
		}
	}

	if f.Fraction > 0 {
		sa = sa[:len(sa)-f.Fraction] + f.Decimal + sa[len(sa)-f.Fraction:]
	}
	sa = strings.Replace(f.Template, "1", sa, 1)
	sa = strings.Replace(sa, "$", f.Grapheme, 1)

	// Add minus sign for negative amount
	if amount < 0 {
		sa = "-" + sa
	}

	return sa
}

// abs return absolute value of given integer
func (f Formatter) abs(amount int64) int64 {
	if amount < 0 {
		return -amount
	}

	return amount
}

func (f *Formatter) FormatAccounting(amount int64) string {
	// Work with absolute amount value
	sa := strconv.FormatInt(f.abs(amount), 10)

	if len(sa) <= f.Fraction {
		sa = strings.Repeat("0", f.Fraction-len(sa)+1) + sa
	}

	if f.Thousand != "" {
		for i := len(sa) - f.Fraction - 3; i > 0; i -= 3 {
			sa = sa[:i] + f.Thousand + sa[i:]
		}
	}

	if f.Fraction > 0 {
		sa = sa[:len(sa)-f.Fraction] + f.Decimal + sa[len(sa)-f.Fraction:]
	}

	if amount < 0 {
		sa = "( " + sa + " )"
	}

	return sa
}

func (f *Formatter) FormatDIV(amount int64, positive string, negative string) template.HTML {
	// Work with absolute amount value
	sa := strconv.FormatInt(f.abs(amount), 10)

	if len(sa) <= f.Fraction {
		sa = strings.Repeat("0", f.Fraction-len(sa)+1) + sa
	}

	if f.Thousand != "" {
		for i := len(sa) - f.Fraction - 3; i > 0; i -= 3 {
			sa = sa[:i] + f.Thousand + sa[i:]
		}
	}

	if f.Fraction > 0 {
		sa = sa[:len(sa)-f.Fraction] + f.Decimal + sa[len(sa)-f.Fraction:]
	}

	sa = strings.Replace(f.Template, "1", sa, 1)
	sa = strings.Replace(sa, "$", f.Grapheme, 1)

	if amount < 0 {
		sa = "<div class='" + negative + "'>" + sa + "</div>"
	} else {
		sa = "<div class='" + positive + "'>" + sa + "</div>"
	}

	return template.HTML(sa)
}

func (f *Formatter) FormatAccountingDIV(amount int64, positive string, negative string) template.HTML {
	// Work with absolute amount value
	sa := strconv.FormatInt(f.abs(amount), 10)

	if len(sa) <= f.Fraction {
		sa = strings.Repeat("0", f.Fraction-len(sa)+1) + sa
	}

	if f.Thousand != "" {
		for i := len(sa) - f.Fraction - 3; i > 0; i -= 3 {
			sa = sa[:i] + f.Thousand + sa[i:]
		}
	}

	if f.Fraction > 0 {
		sa = sa[:len(sa)-f.Fraction] + f.Decimal + sa[len(sa)-f.Fraction:]
	}

	if amount < 0 {
		sa = "<div class='" + negative + "'>( " + sa + " )</div>"
	} else {
		sa = "<div class='" + positive + "'>" + sa + "</div>"
	}

	return template.HTML(sa)
}
