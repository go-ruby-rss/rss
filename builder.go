// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import "strings"

// attr is one XML attribute, kept in insertion order so output matches MRI,
// which preserves the declaration order of attributes.
type attr struct {
	name  string
	value string
}

// node is the minimal element model the serializer walks. It mirrors how
// RSS::Element#tag builds XML directly as strings (MRI does not route output
// through REXML): a node is either a leaf carrying text content, or a branch
// carrying child nodes.
type node struct {
	name     string
	attrs    []attr
	text     string // leaf text content (already unescaped)
	hasText  bool   // distinguishes empty-string text from no text
	children []*node
}

func newNode(name string) *node { return &node{name: name} }

func (n *node) setAttr(name, value string) {
	n.attrs = append(n.attrs, attr{name, value})
}

func (n *node) setText(s string) { n.text, n.hasText = s, true }

func (n *node) add(c *node) {
	if c != nil {
		n.children = append(n.children, c)
	}
}

// leaf builds a text-carrying element <name>text</name>.
func leaf(name, text string) *node {
	n := newNode(name)
	n.setText(text)
	return n
}

// render serializes the node tree the way RSS::Element#tag does, at the given
// indentation. It returns the empty string for an element that has neither
// content, attributes, nor children — matching MRI's "return ” if attrs.empty?".
func (n *node) render(ind string) string {
	next := ind + indent
	var sb strings.Builder
	sb.WriteString(ind)
	sb.WriteByte('<')
	sb.WriteString(n.name)
	if len(n.attrs) > 0 {
		sb.WriteByte(' ')
		parts := make([]string, len(n.attrs))
		for i, a := range n.attrs {
			parts[i] = escapeHTML(a.name) + `="` + escapeHTML(a.value) + `"`
		}
		sb.WriteString(strings.Join(parts, "\n"+next))
	}
	start := sb.String()

	if n.hasText {
		// String content: <start>text</name>, content HTML-escaped.
		return start + ">" + escapeHTML(n.text) + "</" + n.name + ">"
	}

	// Branch content: collect non-empty child renders.
	var kids []string
	for _, c := range n.children {
		s := c.render(next)
		if s != "" {
			kids = append(kids, s)
		}
	}
	if len(kids) == 0 {
		if len(n.attrs) == 0 {
			return ""
		}
		return start + "/>"
	}
	return start + ">\n" + strings.Join(kids, "\n") + "\n" + ind + "</" + n.name + ">"
}

// xmlDecl is the prolog a directly-constructed feed emits: version only, no
// encoding (matching RSS::Rss#to_s when @encoding is nil).
const xmlDecl = `<?xml version="1.0"?>` + "\n"

// xmlDeclEncoded is the prolog a parsed feed emits: MRI's parser defaults the
// document encoding to UTF-8, so reserialized feeds carry encoding="UTF-8".
const xmlDeclEncoded = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"

// prolog selects the XML declaration based on whether the feed was parsed.
func prolog(parsed bool) string {
	if parsed {
		return xmlDeclEncoded
	}
	return xmlDecl
}
