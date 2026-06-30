// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parseRFC822 parses an RSS <pubDate>-style date. MRI accepts both numeric
// offsets ("+0100") and the classic named zones ("GMT", "UT", "EST", …),
// returning a Time that preserves the offset (named zones map to their fixed
// offsets). The day-of-week prefix is optional.
func parseRFC822(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	// Drop an optional leading weekday, e.g. "Tue, ".
	if i := strings.IndexByte(s, ','); i >= 0 && i <= 3 {
		s = strings.TrimSpace(s[i+1:])
	}
	fields := strings.Fields(s)
	if len(fields) < 5 {
		return time.Time{}, fmt.Errorf("rss: not an RFC822 date: %q", s)
	}
	day, err := strconv.Atoi(fields[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("rss: bad day in %q", s)
	}
	mon, ok := monthIndex(fields[1])
	if !ok {
		return time.Time{}, fmt.Errorf("rss: bad month in %q", s)
	}
	year, err := strconv.Atoi(fields[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("rss: bad year in %q", s)
	}
	if year < 100 { // two-digit years, RFC822 legacy
		if year < 70 {
			year += 2000
		} else {
			year += 1900
		}
	}
	hh, mm, ss, err := parseClock(fields[3])
	if err != nil {
		return time.Time{}, err
	}
	loc, err := parseZone(fields[4])
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(year, time.Month(mon+1), day, hh, mm, ss, 0, loc), nil
}

func parseClock(s string) (hh, mm, ss int, err error) {
	parts := strings.Split(s, ":")
	if len(parts) < 2 {
		return 0, 0, 0, fmt.Errorf("rss: bad time %q", s)
	}
	if hh, err = strconv.Atoi(parts[0]); err != nil {
		return 0, 0, 0, fmt.Errorf("rss: bad hour %q", s)
	}
	if mm, err = strconv.Atoi(parts[1]); err != nil {
		return 0, 0, 0, fmt.Errorf("rss: bad minute %q", s)
	}
	if len(parts) >= 3 {
		if ss, err = strconv.Atoi(parts[2]); err != nil {
			return 0, 0, 0, fmt.Errorf("rss: bad second %q", s)
		}
	}
	return hh, mm, ss, nil
}

// namedZones maps the RFC822 obsolete zone names to their offsets in seconds.
// MRI maps these to fixed-offset Times (e.g. "GMT" → +0000), not to a UTC
// location, so a reserialized RFC822 date prints "+0000" rather than "-0000".
var namedZones = map[string]int{
	"UT": 0, "GMT": 0, "Z": 0,
	"EST": -5 * 3600, "EDT": -4 * 3600,
	"CST": -6 * 3600, "CDT": -5 * 3600,
	"MST": -7 * 3600, "MDT": -6 * 3600,
	"PST": -8 * 3600, "PDT": -7 * 3600,
}

func parseZone(s string) (*time.Location, error) {
	if off, ok := namedZones[s]; ok {
		return time.FixedZone(s, off), nil
	}
	if len(s) == 5 && (s[0] == '+' || s[0] == '-') {
		h, err1 := strconv.Atoi(s[1:3])
		m, err2 := strconv.Atoi(s[3:5])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("rss: bad zone %q", s)
		}
		off := (h*60 + m) * 60
		if s[0] == '-' {
			off = -off
		}
		return time.FixedZone(s, off), nil
	}
	return nil, fmt.Errorf("rss: bad zone %q", s)
}

func monthIndex(s string) (int, bool) {
	for i, m := range monthsShort {
		if strings.EqualFold(m, s) {
			return i, true
		}
	}
	return 0, false
}

// parseW3CDTF parses an ISO8601/RFC3339 date as MRI's Time.w3cdtf accepts it.
// A trailing "Z" yields a UTC-located Time; an explicit offset yields a fixed
// zone preserving that offset.
func parseW3CDTF(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	// Try the common layouts with fractional seconds first.
	layouts := []string{
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04Z07:00",
		"2006-01-02",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			// Normalize a "+00:00" Z-suffix into UTC so re-serialization
			// matches MRI: a "Z" input round-trips as "Z".
			if strings.HasSuffix(s, "Z") {
				return t.UTC(), nil
			}
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("rss: not a W3CDTF date: %q", s)
}
