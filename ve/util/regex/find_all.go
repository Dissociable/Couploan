package regex

import (
	"github.com/Dissociable/IPFA/ipfc/services/util"
	"regexp"
)

// MatchGroup returns a slice of strings holding the text of the
// leftmost match of the regular expression in s and the matches, if any, of
// its subexpressions, as defined by the 'Submatch' description in the
// package comment.
// A return value of nil indicates no match.
func (c *Regex) MatchGroup(regEx *regexp.Regexp, s string, group int) *string {
	r := MatchGroup(regEx, s, group)
	if c.options.ManipulateStringResult != nil {
		r = c.options.ManipulateStringResult(r)
	}
	return r
}

// MatchGroup returns a slice of strings holding the text of the
// leftmost match of the regular expression in s and the matches, if any, of
// its subexpressions, as defined by the 'Submatch' description in the
// package comment.
// A return value of nil indicates no match.
func MatchGroup(regEx *regexp.Regexp, s string, group int) *string {
	r := regEx.FindStringSubmatch(s)
	if r == nil {
		return nil
	}
	if group >= len(r) {
		return nil
	}
	return util.StringP(r[group])
}

// MatchGroups returns a slice of strings holding the text of the
// leftmost match of the regular expression in s and the matches, if any, of
// its subexpressions, as defined by the 'Submatch' description in the
// package comment.
// A return value of nil indicates no match.
func (c *Regex) MatchGroups(regEx *regexp.Regexp, s string, group ...int) []string {
	r := MatchGroups(regEx, s, group...)
	if c.options.ManipulateStringResult != nil {
		for i := 0; i < len(r); i++ {
			manipulated := c.options.ManipulateStringResult(&r[i])
			if manipulated == nil {
				continue
			}
			r[i] = *manipulated
		}
	}
	return r
}

// MatchGroups returns a slice of strings holding the text of the
// leftmost match of the regular expression in s and the matches, if any, of
// its subexpressions, as defined by the 'Submatch' description in the
// package comment.
// A return value of nil indicates no match.
func MatchGroups(regEx *regexp.Regexp, s string, group ...int) []string {
	r := regEx.FindStringSubmatch(s)
	if r == nil {
		return nil
	}
	if len(group) == 0 {
		group = []int{0}
	}
	var res []string
	for _, g := range group {
		if g >= len(r) {
			continue
		}
		res = append(res, r[g])
	}
	return res
}

// MatchAllGroup returns a slice of strings holding the text of the
// leftmost match of the regular expression in s and the matches, if any, of
// its subexpressions, as defined by the 'Submatch' description in the
// package comment.
// A return value of nil indicates no match.
func (c *Regex) MatchAllGroup(regEx *regexp.Regexp, s string, n int, group int) []string {
	r := MatchAllGroup(regEx, s, n, group)
	var rNew []string
	if c.options.ManipulateStringResult != nil {
		for i := 0; i < len(r); i++ {
			manipulated := c.options.ManipulateStringResult(&r[i])
			if manipulated == nil {
				continue
			}
			rNew = append(rNew, *manipulated)
		}
	}
	return rNew
}

// MatchAllGroup is the 'All' version of [Regex.MatchGroup]; it
// returns a slice of all successive matches of the expression, as defined by
// the 'All' description in the package comment.
// A return value of nil indicates no match.
func MatchAllGroup(regEx *regexp.Regexp, s string, n int, group int) []string {
	r := regEx.FindAllStringSubmatch(s, n)
	if r == nil {
		return nil
	}
	var result []string
	for _, v := range r {
		if group >= len(v) {
			continue
		}
		result = append(result, v[group])
	}
	return result
}

// MatchAllGroups returns a slice of strings holding the text of the
// leftmost match of the regular expression in s and the matches, if any, of
// its subexpressions, as defined by the 'Submatch' description in the
// package comment.
// A return value of nil indicates no match.
func (c *Regex) MatchAllGroups(regEx *regexp.Regexp, s string, n int, group ...int) [][]string {
	r := MatchAllGroups(regEx, s, n, group...)
	var rNew [][]string
	if c.options.ManipulateStringResult != nil {
		for i := 0; i < len(r); i++ {
			var rNewInner []string
			for j := 0; j < len(r[i]); j++ {
				manipulated := c.options.ManipulateStringResult(&r[i][j])
				if manipulated == nil {
					continue
				}
				rNewInner = append(rNewInner, *manipulated)
			}
			if len(rNewInner) > 0 {
				rNew = append(rNew, rNewInner)
			}
		}
	}
	return r
}

// MatchAllGroups is the 'All' version of [Regex.MatchGroups]; it
// returns a slice of all successive matches of the expression, as defined by
// the 'All' description in the package comment.
// A return value of nil indicates no match.
func MatchAllGroups(regEx *regexp.Regexp, s string, n int, group ...int) [][]string {
	r := regEx.FindAllStringSubmatch(s, n)
	if r == nil {
		return nil
	}
	var result [][]string
	for _, v := range r {
		var innerResult []string
		for _, g := range group {
			if g >= len(v) {
				continue
			}
			innerResult = append(innerResult, v[g])
		}
		if len(innerResult) > 0 {
			result = append(result, innerResult)
		}
	}
	return result
}
