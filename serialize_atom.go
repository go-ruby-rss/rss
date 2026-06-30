// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

// String serializes the Atom feed to MRI-compatible XML.
//
// MRI declares the default Atom namespace plus the Dublin Core namespace on
// the root <feed> (xmlns, then xmlns:dc), and emits child elements in this
// fixed order: author, category, contributor, generator, icon, id, link,
// logo, rights, subtitle, title, updated, then the entries. Atom dates use
// W3CDTF (RFC3339).
func (f *AtomFeed) String() string {
	root := newNode("feed")
	root.setAttr("xmlns", uriAtom)
	root.setAttr("xmlns:dc", uriDC)
	for _, a := range f.Authors {
		root.add(atomPersonNode("author", a))
	}
	for _, cat := range f.Categories {
		root.add(atomCategoryNode(cat))
	}
	addLeaf(root, "generator", f.Generator, false)
	addLeaf(root, "id", f.ID, false)
	for _, l := range f.Links {
		root.add(atomLinkNode(l))
	}
	addLeaf(root, "rights", f.Rights, false)
	addLeaf(root, "subtitle", f.Subtitle, false)
	addLeaf(root, "title", f.Title, false)
	if f.Updated != nil {
		root.add(leaf("updated", formatW3CDTF(*f.Updated)))
	}
	for _, e := range f.Entries {
		root.add(e.node())
	}
	return prolog(f.parsed) + root.render("")
}

func (e *AtomEntry) node() *node {
	n := newNode("entry")
	for _, a := range e.Authors {
		n.add(atomPersonNode("author", a))
	}
	for _, cat := range e.Categories {
		n.add(atomCategoryNode(cat))
	}
	addLeaf(n, "content", e.Content, false)
	addLeaf(n, "id", e.ID, false)
	for _, l := range e.Links {
		n.add(atomLinkNode(l))
	}
	if e.Published != nil {
		n.add(leaf("published", formatW3CDTF(*e.Published)))
	}
	addLeaf(n, "rights", e.Rights, false)
	addLeaf(n, "summary", e.Summary, false)
	addLeaf(n, "title", e.Title, false)
	if e.Updated != nil {
		n.add(leaf("updated", formatW3CDTF(*e.Updated)))
	}
	return n
}

func atomPersonNode(tag string, p *AtomPerson) *node {
	n := newNode(tag)
	addLeaf(n, "name", p.Name, false)
	addLeaf(n, "uri", p.URI, false)
	addLeaf(n, "email", p.Email, false)
	return n
}

func atomLinkNode(l *AtomLink) *node {
	n := newNode("link")
	if l.Href != "" {
		n.setAttr("href", l.Href)
	}
	if l.Rel != "" {
		n.setAttr("rel", l.Rel)
	}
	if l.Type != "" {
		n.setAttr("type", l.Type)
	}
	if l.Title != "" {
		n.setAttr("title", l.Title)
	}
	return n
}

func atomCategoryNode(c *AtomCategory) *node {
	n := newNode("category")
	if c.Term != "" {
		n.setAttr("term", c.Term)
	}
	if c.Scheme != "" {
		n.setAttr("scheme", c.Scheme)
	}
	if c.Label != "" {
		n.setAttr("label", c.Label)
	}
	return n
}
