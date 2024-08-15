package util

import "github.com/microcosm-cc/bluemonday"

var HtmlTagStripper = bluemonday.StripTagsPolicy()
