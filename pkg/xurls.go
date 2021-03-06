// Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

// Package xurls extracts urls from plain text using regular expressions.
package xurls

import (
	"bytes"
	"regexp"
)

//go:generate go run cmd/generate/tldsgen/main.go
//go:generate go run cmd/generate/schemesgen/main.go

const (
	letter    = `\p{L}`
	mark      = `\p{M}`
	number    = `\p{N}`
	iriChar   = letter + mark + number
	currency  = `\p{Sc}`
	otherSymb = `\p{So}`
	endChar   = iriChar + `/\-+_&~*%=#` + currency + otherSymb
	otherPunc = `\p{Po}`
	midChar   = endChar + `|` + otherPunc
	wellParen = `\([` + midChar + `]*(\([` + midChar + `]*\)[` + midChar + `]*)*\)`
	wellBrack = `\[[` + midChar + `]*(\[[` + midChar + `]*\][` + midChar + `]*)*\]`
	wellBrace = `\{[` + midChar + `]*(\{[` + midChar + `]*\}[` + midChar + `]*)*\}`
	wellAll   = wellParen + `|` + wellBrack + `|` + wellBrace
	pathCont  = `([` + midChar + `]*(` + wellAll + `|[` + endChar + `])+)+`

	iri      = `[` + iriChar + `]([` + iriChar + `\-]*[` + iriChar + `])?`
	domain   = `(` + iri + `\.)+`
	octet    = `(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])`
	ipv4Addr = `\b` + octet + `\.` + octet + `\.` + octet + `\.` + octet + `\b`
	ipv6Addr = `([0-9a-fA-F]{1,4}:([0-9a-fA-F]{1,4}:([0-9a-fA-F]{1,4}:([0-9a-fA-F]{1,4}:([0-9a-fA-F]{1,4}:[0-9a-fA-F]{0,4}|:[0-9a-fA-F]{1,4})?|(:[0-9a-fA-F]{1,4}){0,2})|(:[0-9a-fA-F]{1,4}){0,3})|(:[0-9a-fA-F]{1,4}){0,4})|:(:[0-9a-fA-F]{1,4}){0,5})((:[0-9a-fA-F]{1,4}){2}|:(25[0-5]|(2[0-4]|1[0-9]|[1-9])?[0-9])(\.(25[0-5]|(2[0-4]|1[0-9]|[1-9])?[0-9])){3})|(([0-9a-fA-F]{1,4}:){1,6}|:):[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){7}:`
	ipAddr   = `(` + ipv4Addr + `|` + ipv6Addr + `)`
	port     = `(:[0-9]*)?`

	// ref. https://github.com/mingrammer/commonregex/blob/master/commonregex.go#L7-L35
	universalLinkPattern = `(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`
	linkPattern          = `(?:(?:https?:\/\/)?(?:[a-z0-9.\-]+|www|[a-z0-9.\-])[.](?:[^\s()<>]+|\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\))+(?:\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:\'".,<>?]))`
	emailPattern         = `(?i)([A-Za-z0-9!#$%&'*+\/=?^_{|.}~-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`
	ipv4Pattern          = `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
	ipv6Pattern          = `(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`
	ipPattern            = ipv4Pattern + `|` + ipv6Pattern
	macAddressPattern    = `(([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2}))`
	notKnownPortPattern  = `6[0-5]{2}[0-3][0-5]|[1-5][\d]{4}|[2-9][\d]{3}|1[1-9][\d]{2}|10[3-9][\d]|102[4-9]`
)

// AnyScheme can be passed to StrictMatchingScheme to match any possibly
// valid scheme.
var AnyScheme = `([a-zA-Z][a-zA-Z.\-+]*://|` + anyOf(SchemesNoAuthority...) + `:)`

// SchemesNoAuthority is a sorted list of some well-known url schemes that are
// followed by ":" instead of "://".
var SchemesNoAuthority = []string{
	`bitcoin`, // Bitcoin
	`file`,    // Files
	`magnet`,  // Torrent magnets
	`mailto`,  // Mail
	`sms`,     // SMS
	`tel`,     // Telephone
	`xmpp`,    // XMPP
}

func anyOf(strs ...string) string {
	var b bytes.Buffer
	b.WriteByte('(')
	for i, s := range strs {
		if i != 0 {
			b.WriteByte('|')
		}
		b.WriteString(regexp.QuoteMeta(s))
	}
	b.WriteByte(')')
	return b.String()
}

func strictExp() string {
	schemes := `(` + anyOf(Schemes...) + `://|` + anyOf(SchemesNoAuthority...) + `:)`
	return `(?i)` + schemes + `(?-i)` + pathCont
}

func relaxedExp() string {
	site := domain + `(?i)` + anyOf(append(TLDs, PseudoTLDs...)...) + `(?-i)`
	hostName := `(` + site + `|` + ipAddr + `)`
	webURL := hostName + port + `(/|/` + pathCont + `?|\b|$)`
	return strictExp() + `|` + webURL
}

func Relaxed() *regexp.Regexp {
	re := regexp.MustCompile(relaxedExp())
	re.Longest()
	return re
}

func Strict() *regexp.Regexp {
	re := regexp.MustCompile(strictExp())
	re.Longest()
	return re
}

// StrictMatchingScheme produces a regexp that matches urls like Strict but
// whose scheme matches the given regular expression.
func StrictMatchingScheme(exp string) (*regexp.Regexp, error) {
	strictMatching := `(?i)(` + exp + `)(?-i)` + pathCont
	re, err := regexp.Compile(strictMatching)
	if err != nil {
		return nil, err
	}
	re.Longest()
	return re, nil
}
