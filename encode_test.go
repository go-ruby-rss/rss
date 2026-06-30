// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import (
	"testing"
	"time"
)

func TestEscapeHTML(t *testing.T) {
	got := escapeHTML(`<a> & "b" 'c`)
	want := `&lt;a&gt; &amp; &quot;b&quot; &#39;c`
	if got != want {
		t.Errorf("escapeHTML = %q, want %q", got, want)
	}
	if escapeHTML("plain") != "plain" {
		t.Error("plain text should be unchanged")
	}
}

func TestFormatRFC822(t *testing.T) {
	cases := []struct {
		t    time.Time
		want string
	}{
		{time.Date(2023, 11, 14, 23, 13, 20, 0, time.FixedZone("", 3600)), "Tue, 14 Nov 2023 23:13:20 +0100"},
		{time.Date(2023, 11, 14, 22, 13, 20, 0, time.UTC), "Tue, 14 Nov 2023 22:13:20 -0000"},
		{time.Date(2023, 11, 14, 9, 30, 0, 0, time.FixedZone("", -5*3600)), "Tue, 14 Nov 2023 09:30:00 -0500"},
	}
	for _, c := range cases {
		if got := formatRFC822(c.t); got != c.want {
			t.Errorf("formatRFC822 = %q, want %q", got, c.want)
		}
	}
}

func TestFormatW3CDTF(t *testing.T) {
	cases := []struct {
		t    time.Time
		want string
	}{
		{time.Date(2023, 11, 14, 22, 13, 20, 0, time.UTC), "2023-11-14T22:13:20Z"},
		{time.Date(2023, 11, 14, 23, 13, 20, 0, time.FixedZone("", 3600)), "2023-11-14T23:13:20+01:00"},
		{time.Date(2023, 11, 14, 9, 30, 0, 0, time.FixedZone("", -5*3600)), "2023-11-14T09:30:00-05:00"},
		{time.Date(2023, 11, 14, 22, 13, 20, 500000000, time.UTC), "2023-11-14T22:13:20.5Z"},
	}
	for _, c := range cases {
		if got := formatW3CDTF(c.t); got != c.want {
			t.Errorf("formatW3CDTF(%v) = %q, want %q", c.t, got, c.want)
		}
	}
}

func TestParseRFC822(t *testing.T) {
	cases := []struct {
		in       string
		wantYear int
		wantOff  int // seconds
	}{
		{"Tue, 14 Nov 2023 23:13:20 +0100", 2023, 3600},
		{"14 Nov 2023 23:13:20 +0100", 2023, 3600}, // no weekday
		{"Mon, 06 Sep 2021 10:00:00 GMT", 2021, 0},
		{"Sun, 05 Sep 2021 09:30:00 -0200", 2021, -7200},
		{"01 Jan 99 00:00:00 GMT", 1999, 0},        // two-digit year >= 70
		{"01 Jan 23 00:00:00 GMT", 2023, 0},        // two-digit year < 70
		{"01 Jan 2023 00:00 EST", 2023, -5 * 3600}, // no seconds, named zone
	}
	for _, c := range cases {
		got, err := parseRFC822(c.in)
		if err != nil {
			t.Fatalf("parseRFC822(%q): %v", c.in, err)
		}
		if got.Year() != c.wantYear {
			t.Errorf("parseRFC822(%q) year = %d, want %d", c.in, got.Year(), c.wantYear)
		}
		if _, off := got.Zone(); off != c.wantOff {
			t.Errorf("parseRFC822(%q) offset = %d, want %d", c.in, off, c.wantOff)
		}
	}
}

func TestParseRFC822Errors(t *testing.T) {
	bad := []string{
		"", "x", "too few", "Tue, 14 Nov 2023 23:13:20",
		"14 Bad 2023 23:13:20 +0100",
		"xx Nov 2023 23:13:20 +0100",
		"14 Nov xxxx 23:13:20 +0100",
		"14 Nov 2023 99 +0100",         // bad time (one part)
		"14 Nov 2023 xx:13:20 +0100",   // bad hour
		"14 Nov 2023 23:xx:20 +0100",   // bad minute
		"14 Nov 2023 23:13:xx +0100",   // bad second
		"14 Nov 2023 23:13:20 BADZONE", // bad zone
		"14 Nov 2023 23:13:20 +xx00",   // bad numeric zone hour
		"14 Nov 2023 23:13:20 +00xx",   // bad numeric zone minute
	}
	for _, s := range bad {
		if _, err := parseRFC822(s); err == nil {
			t.Errorf("parseRFC822(%q) should fail", s)
		}
	}
}

func TestParseW3CDTF(t *testing.T) {
	cases := []struct {
		in       string
		wantUTC  bool
		wantYear int
	}{
		{"2023-11-14T22:13:20Z", true, 2023},
		{"2023-11-14T23:13:20+01:00", false, 2023},
		{"2023-11-14T23:13:20.5+01:00", false, 2023},
		{"2023-11-14T23:13+01:00", false, 2023}, // no seconds
		{"2023-11-14", false, 2023},             // date only
	}
	for _, c := range cases {
		got, err := parseW3CDTF(c.in)
		if err != nil {
			t.Fatalf("parseW3CDTF(%q): %v", c.in, err)
		}
		if got.Year() != c.wantYear {
			t.Errorf("parseW3CDTF(%q) year = %d", c.in, got.Year())
		}
		if c.wantUTC && got.Location() != time.UTC {
			t.Errorf("parseW3CDTF(%q) should be UTC", c.in)
		}
	}
	if _, err := parseW3CDTF("not a date"); err == nil {
		t.Error("parseW3CDTF should reject garbage")
	}
}

func TestZParsesAsUTC(t *testing.T) {
	// A "Z" RFC822 named zone maps to UTC.
	tm, err := parseRFC822("01 Jan 2023 00:00:00 Z")
	if err != nil {
		t.Fatal(err)
	}
	if _, off := tm.Zone(); off != 0 {
		t.Errorf("Z zone offset = %d, want 0", off)
	}
}
