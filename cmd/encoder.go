// SPDX-FileCopyrightText: The Go Authors. All rights reserved.
// SPDX-FileCopyrightText: 2021 Kalle Fagerberg
//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"net/url"
	"strings"

	"github.com/fatih/color"
	"github.com/jilleJr/urlencode/pkg/flagtype"
)

// The code in this file has been taken from the source code of the `net/url`
// Go package, v1.17.1.

var escapedColor = color.New(color.FgMagenta)
var unescapedColor = color.New(color.FgRed)

const upperHex = "0123456789ABCDEF"

// isHex has been copied from
// https://cs.opensource.google/go/go/+/refs/tags/go1.17.1:src/net/url/url.go;l=47-57
func isHex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

// unHex has been copied from
// https://cs.opensource.google/go/go/+/refs/tags/go1.17.1:src/net/url/url.go;l=59-69
func unHex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

// shouldEscape has been copied from
// https://cs.opensource.google/go/go/+/refs/tags/go1.17.1:src/net/url/url.go;l=100-175
func shouldEscape(c byte, mode flagtype.Encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}

	if mode == flagtype.EncodeHost || mode == flagtype.EncodeZone {
		// §3.2.2 Host allows
		//	sub-delims = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		// as part of reg-name.
		// We add : because we include :port as part of host.
		// We add [ ] because we include [ipv6]:port as part of host.
		// We add < > because they're the only characters left that
		// we could possibly allow, and Parse will reject them if we
		// escape them (because hosts can't use %-encoding for
		// ASCII bytes).
		switch c {
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '[', ']', '<', '>', '"':
			return false
		}
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch mode {
		case flagtype.EncodePath: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments. This package
			// only manipulates the path as a whole, so we allow those
			// last three as well. That leaves only ? to escape.
			return c == '?'

		case flagtype.EncodePathSegment: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments.
			return c == '/' || c == ';' || c == ',' || c == '?'

		case flagtype.EncodeUserPassword: // §3.2.1
			// The RFC allows ';', ':', '&', '=', '+', '$', and ',' in
			// userinfo, so we must escape only '@', '/', and '?'.
			// The parsing of userinfo treats ':' as special so we must escape
			// that too.
			return c == '@' || c == '/' || c == '?' || c == ':'

		case flagtype.EncodeQueryComponent: // §3.4
			// The RFC reserves (so we must escape) everything.
			return true

		case flagtype.EncodeFragment: // §4.1
			// The RFC text is silent but the grammar allows
			// everything, so escape nothing.
			return false
		}
	}

	if mode == flagtype.EncodeFragment {
		// RFC 3986 §2.2 allows not escaping sub-delims. A subset of sub-delims are
		// included in reserved from RFC 2396 §2.2. The remaining sub-delims do not
		// need to be escaped. To minimize potential breakage, we apply two restrictions:
		// (1) we always escape sub-delims outside of the fragment, and (2) we always
		// escape single quote to avoid breaking callers that had previously assumed that
		// single quotes would be escaped. See issue #19917.
		switch c {
		case '!', '(', ')', '*':
			return false
		}
	}

	// Everything else must be escaped.
	return true
}

// unescape has been copied and modified from
// https://cs.opensource.google/go/go/+/refs/tags/go1.17.1:src/net/url/url.go;l=199-270
func unescape(s string, mode flagtype.Encoding) (string, error) {
	// Count %, check that they're well-formed.
	n := 0
	hasPlus := false
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			n++
			if i+2 >= len(s) || !isHex(s[i+1]) || !isHex(s[i+2]) {
				s = s[i:]
				if len(s) > 3 {
					s = s[:3]
				}
				return "", url.EscapeError(s)
			}
			// Per https://tools.ietf.org/html/rfc3986#page-21
			// in the host component %-encoding can only be used
			// for non-ASCII bytes.
			// But https://tools.ietf.org/html/rfc6874#section-2
			// introduces %25 being allowed to escape a percent sign
			// in IPv6 scoped-address literals. Yay.
			if mode == flagtype.EncodeHost && unHex(s[i+1]) < 8 && s[i:i+3] != "%25" {
				return "", url.EscapeError(s[i : i+3])
			}
			if mode == flagtype.EncodeZone {
				// RFC 6874 says basically "anything goes" for zone identifiers
				// and that even non-ASCII can be redundantly escaped,
				// but it seems prudent to restrict %-escaped bytes here to those
				// that are valid host name bytes in their unescaped form.
				// That is, you can use escaping in the zone identifier but not
				// to introduce bytes you couldn't just write directly.
				// But Windows puts spaces here! Yay.
				v := unHex(s[i+1])<<4 | unHex(s[i+2])
				if s[i:i+3] != "%25" && v != ' ' && shouldEscape(v, flagtype.EncodeHost) {
					return "", url.EscapeError(s[i : i+3])
				}
			}
			i += 3
		case '+':
			hasPlus = mode == flagtype.EncodeQueryComponent
			i++
		default:
			if (mode == flagtype.EncodeHost || mode == flagtype.EncodeZone) && s[i] < 0x80 && shouldEscape(s[i], mode) {
				return "", url.InvalidHostError(s[i : i+1])
			}
			i++
		}
	}

	if n == 0 && !hasPlus {
		return s, nil
	}

	var t strings.Builder
	t.Grow(len(s) - 2*n)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '%':
			var strByte [1]byte
			strByte[0] = unHex(s[i+1])<<4 | unHex(s[i+2])
			unescapedColor.Fprint(&t, string(strByte[:]))
			i += 2
		case '+':
			if mode == flagtype.EncodeQueryComponent {
				unescapedColor.Fprint(&t, " ")
			} else {
				t.WriteByte('+')
			}
		default:
			t.WriteByte(s[i])
		}
	}
	return t.String(), nil
}

// unescape has been copied and modified from
// https://cs.opensource.google/go/go/+/refs/tags/go1.17.1:src/net/url/url.go;l=284-338
func escape(s string, mode flagtype.Encoding) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == flagtype.EncodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	var sb strings.Builder

	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == flagtype.EncodeQueryComponent:
			escapedColor.Fprint(&sb, "+")
		case shouldEscape(c, mode):
			var strByte [3]byte
			strByte[0] = '%'
			strByte[1] = upperHex[c>>4]
			strByte[2] = upperHex[c&15]
			escapedColor.Fprint(&sb, string(strByte[:]))
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}
