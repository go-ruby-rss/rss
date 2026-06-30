// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import (
	"strings"
	"testing"
	"time"
)

const goldenAtom = `<?xml version="1.0"?>
<feed xmlns="http://www.w3.org/2005/Atom"
  xmlns:dc="http://purl.org/dc/elements/1.1/">
  <author>
    <name>Jane</name>
  </author>
  <id>urn:uuid:1</id>
  <link href="http://example.com/"/>
  <title>Atom Example</title>
  <updated>2023-11-14T23:13:20+01:00</updated>
</feed>`

func sampleAtom() *AtomFeed {
	return &AtomFeed{
		ID: "urn:uuid:1", Title: "Atom Example",
		Updated: ptime(2023, 11, 14, 23, 13, 20, loc1),
		Authors: []*AtomPerson{{Name: "Jane"}},
		Links:   []*AtomLink{{Href: "http://example.com/"}},
	}
}

func TestAtomGenerate(t *testing.T) {
	if got := sampleAtom().String(); got != goldenAtom {
		t.Errorf("Atom generate mismatch:\n got:\n%s\nwant:\n%s", got, goldenAtom)
	}
	if ft := sampleAtom().FeedType(); ft != "atom" {
		t.Errorf("FeedType = %q, want atom", ft)
	}
}

func TestAtomFull(t *testing.T) {
	f := &AtomFeed{
		ID: "i", Title: "t", Subtitle: "sub", Rights: "(c)", Generator: "g",
		Updated:    ptime(2023, 1, 1, 0, 0, 0, time.UTC),
		Authors:    []*AtomPerson{{Name: "n", URI: "u", Email: "e"}},
		Links:      []*AtomLink{{Href: "h", Rel: "alternate", Type: "text/html", Title: "lt"}},
		Categories: []*AtomCategory{{Term: "term", Scheme: "sch", Label: "lab"}},
		Entries: []*AtomEntry{{
			ID: "ei", Title: "et", Summary: "es", Content: "ec", Rights: "er",
			Updated:    ptime(2023, 1, 1, 0, 0, 0, time.UTC),
			Published:  ptime(2023, 1, 2, 0, 0, 0, loc1),
			Authors:    []*AtomPerson{{Name: "ea"}},
			Links:      []*AtomLink{{Href: "eh"}},
			Categories: []*AtomCategory{{Term: "ect"}},
		}},
	}
	out := f.String()
	for _, want := range []string{
		"<subtitle>sub</subtitle>", "<rights>(c)</rights>", "<generator>g</generator>",
		"<uri>u</uri>", "<email>e</email>",
		`<link href="h"`, `rel="alternate"`, `type="text/html"`, `title="lt"`,
		`<category term="term"`, `scheme="sch"`, `label="lab"`,
		"<entry>", "<content>ec</content>", "<summary>es</summary>",
		"<published>2023-01-02T00:00:00+01:00</published>",
		"<updated>2023-01-01T00:00:00Z</updated>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("full Atom missing %q in:\n%s", want, out)
		}
	}
}

func TestAtomRoundTrip(t *testing.T) {
	f, err := Parse(goldenAtom)
	if err != nil {
		t.Fatal(err)
	}
	af, ok := f.(*AtomFeed)
	if !ok {
		t.Fatalf("Parse returned %T, want *AtomFeed", f)
	}
	if af.ID != "urn:uuid:1" || af.Authors[0].Name != "Jane" || af.Links[0].Href != "http://example.com/" {
		t.Errorf("Atom parse wrong: %+v", af)
	}
	if !strings.HasPrefix(af.String(), `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("parsed Atom should carry encoding prolog")
	}
}

func TestAtomParseAll(t *testing.T) {
	src := `<?xml version="1.0"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <id>i</id><title>t</title><subtitle>sub</subtitle>
  <rights>r</rights><generator>g</generator>
  <updated>2023-11-14T22:13:20Z</updated>
  <author><name>n</name><uri>u</uri><email>e</email></author>
  <link href="h" rel="self" type="text/html" title="lt"/>
  <category term="ct" scheme="cs" label="cl"/>
  <entry>
    <id>ei</id><title>et</title><summary>es</summary><content>ec</content>
    <rights>er</rights>
    <updated>2023-11-14T22:13:20Z</updated>
    <published>2023-11-14T22:13:20Z</published>
    <author><name>ea</name></author>
    <link href="eh"/>
    <category term="ect"/>
  </entry>
</feed>`
	f, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	af := f.(*AtomFeed)
	if af.Subtitle != "sub" || af.Rights != "r" || af.Generator != "g" || af.Updated == nil ||
		af.Authors[0].URI != "u" || af.Authors[0].Email != "e" ||
		af.Links[0].Rel != "self" || af.Links[0].Type != "text/html" || af.Links[0].Title != "lt" ||
		af.Categories[0].Scheme != "cs" || af.Categories[0].Label != "cl" {
		t.Errorf("feed parse incomplete: %+v", af)
	}
	e := af.Entries[0]
	if e.ID != "ei" || e.Summary != "es" || e.Content != "ec" || e.Rights != "er" ||
		e.Updated == nil || e.Published == nil || e.Authors[0].Name != "ea" ||
		e.Links[0].Href != "eh" || e.Categories[0].Term != "ect" {
		t.Errorf("entry parse incomplete: %+v", e)
	}
}

func TestAtomParseErrors(t *testing.T) {
	cases := []string{
		`<feed xmlns="http://www.w3.org/2005/Atom"><updated>bad</updated></feed>`,
		`<feed xmlns="http://www.w3.org/2005/Atom"><entry><updated>bad</updated></entry></feed>`,
		`<feed xmlns="http://www.w3.org/2005/Atom"><entry><published>bad</published></entry></feed>`,
	}
	for i, src := range cases {
		if _, err := Parse(src); err == nil {
			t.Errorf("case %d: expected error", i)
		}
	}
}

// ---- dispatch ------------------------------------------------------------

func TestParseDispatchErrors(t *testing.T) {
	if _, err := Parse(`<<<bad`); err == nil {
		t.Error("malformed XML should error")
	}
	if _, err := Parse(`<html></html>`); err == nil {
		t.Error("unknown root should error")
	}
	if _, err := Parse(``); err == nil {
		t.Error("empty document should error")
	}
}
