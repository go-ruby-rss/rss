// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import "testing"

// TestRenderEmptyNode covers the branch where an element has no text, no
// attributes and only empty children: it renders to the empty string, matching
// MRI's "return ” if attrs.empty?".
func TestRenderEmptyNode(t *testing.T) {
	n := newNode("wrap")
	n.add(newNode("empty")) // child with no text/attrs/children → renders ""
	if got := n.render(""); got != "" {
		t.Errorf("empty node should render to empty string, got %q", got)
	}
}

// TestRenderEmptyWithAttrs covers the self-closing branch: a childless element
// that carries an attribute renders as "<name .../>".
func TestRenderEmptyWithAttrs(t *testing.T) {
	n := newNode("link")
	n.setAttr("href", "h")
	if got := n.render(""); got != `<link href="h"/>` {
		t.Errorf("got %q", got)
	}
}

// TestRdfResourceMissing covers rdfResource's fallback when no rdf:resource
// attribute is present (returns "").
func TestRdfResourceMissing(t *testing.T) {
	src := `<rdf:RDF xmlns:rdf="x"><channel rdf:about="a"><title>t</title><link>l</link><description>d</description><image/></channel></rdf:RDF>`
	f, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if r := f.(*RDF); r.Channel.ImageResource != "" {
		t.Errorf("ImageResource = %q, want empty", r.Channel.ImageResource)
	}
}
