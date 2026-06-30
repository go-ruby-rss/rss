// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import (
	"fmt"
	"strings"
	"time"
)

// escapeHTML replicates Ruby's CGI.escapeHTML, which MRI's RSS uses for both
// element text content and attribute values. It escapes, in this order, the
// five characters & < > " ' to their named/numeric entities.
func escapeHTML(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&#39;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// monthsShort and daysShort drive the RFC822 formatter. MRI's Time#rfc822
// always emits English abbreviations regardless of locale.
var monthsShort = [...]string{
	"Jan", "Feb", "Mar", "Apr", "May", "Jun",
	"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
}
var daysShort = [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// formatRFC822 renders t the way Ruby's Time#rfc822 does, e.g.
// "Tue, 14 Nov 2023 23:13:20 +0100". The offset uses ±HHMM; a UTC time
// (zero offset) renders as "-0000" only if the location is UTC, matching
// MRI, which prints "-0000" for an explicit UTC time.
func formatRFC822(t time.Time) string {
	_, off := t.Zone()
	zone := formatOffsetCompact(off, t)
	return fmt.Sprintf("%s, %02d %s %04d %02d:%02d:%02d %s",
		daysShort[int(t.Weekday())],
		t.Day(), monthsShort[int(t.Month())-1], t.Year(),
		t.Hour(), t.Minute(), t.Second(), zone)
}

// formatOffsetCompact renders the zone for RFC822: "+HHMM"/"-HHMM", with a
// UTC location rendered as "-0000" (MRI's Time#utc.rfc822 behavior).
func formatOffsetCompact(off int, t time.Time) string {
	if off == 0 && t.Location() == time.UTC {
		return "-0000"
	}
	sign := "+"
	if off < 0 {
		sign = "-"
		off = -off
	}
	h := off / 3600
	m := (off % 3600) / 60
	return fmt.Sprintf("%s%02d%02d", sign, h, m)
}

// formatW3CDTF renders t as Ruby's Time#w3cdtf does (an ISO8601 / RFC3339
// profile). Fractional seconds are emitted only when non-zero, trimmed to the
// significant digits. A UTC location renders the zone as "Z"; any other
// location renders "±HH:MM" — matching MRI, where only an explicit utc Time
// yields "Z".
func formatW3CDTF(t time.Time) string {
	base := fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	if ns := t.Nanosecond(); ns != 0 {
		frac := fmt.Sprintf("%09d", ns)
		frac = strings.TrimRight(frac, "0")
		base += "." + frac
	}
	return base + formatOffsetColon(t)
}

// formatOffsetColon renders the zone as "Z" (UTC location) or "±HH:MM".
func formatOffsetColon(t time.Time) string {
	_, off := t.Zone()
	if off == 0 && t.Location() == time.UTC {
		return "Z"
	}
	sign := "+"
	if off < 0 {
		sign = "-"
		off = -off
	}
	h := off / 3600
	m := (off % 3600) / 60
	return fmt.Sprintf("%s%02d:%02d", sign, h, m)
}
