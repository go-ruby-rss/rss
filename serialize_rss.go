// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

// String serializes the RSS 2.0/0.9x document to MRI-compatible XML.
//
// MRI declares the bundled-module namespaces on the root <rss> element
// unconditionally (content, dc, itunes, trackback), in that fixed order,
// regardless of whether the modules are used; this method reproduces that.
func (r *Rss) String() string {
	root := newNode("rss")
	version := r.Version
	if version == "" {
		version = "2.0"
	}
	root.setAttr("version", version)
	root.setAttr("xmlns:content", uriContent)
	root.setAttr("xmlns:dc", uriDC)
	root.setAttr("xmlns:itunes", uriItunes)
	root.setAttr("xmlns:trackback", uriTrackbck)
	if r.Channel != nil {
		root.add(r.Channel.node())
	}
	return prolog(r.parsed) + root.render("")
}

func (c *Channel) node() *node {
	ch := newNode("channel")
	addLeaf(ch, "title", c.Title, true)
	addLeaf(ch, "link", c.Link, true)
	addLeaf(ch, "description", c.Description, true)
	addLeaf(ch, "language", c.Language, false)
	addLeaf(ch, "copyright", c.Copyright, false)
	addLeaf(ch, "managingEditor", c.ManagingEditor, false)
	addLeaf(ch, "webMaster", c.WebMaster, false)
	if c.PubDate != nil {
		ch.add(leaf("pubDate", formatRFC822(*c.PubDate)))
	}
	if c.LastBuildDate != nil {
		ch.add(leaf("lastBuildDate", formatRFC822(*c.LastBuildDate)))
	}
	for _, cat := range c.Categories {
		ch.add(leaf("category", cat))
	}
	addLeaf(ch, "generator", c.Generator, false)
	addLeaf(ch, "docs", c.Docs, false)
	addLeaf(ch, "ttl", c.TTL, false)
	if c.Image != nil {
		ch.add(c.Image.node())
	}
	// Bundled syndication module emitted on the channel.
	if c.SyUpdatePeriod != "" {
		ch.add(leaf("sy:updatePeriod", c.SyUpdatePeriod))
	}
	if c.SyUpdateFrequency != "" {
		ch.add(leaf("sy:updateFrequency", c.SyUpdateFrequency))
	}
	if c.SyUpdateBase != "" {
		ch.add(leaf("sy:updateBase", c.SyUpdateBase))
	}
	for _, it := range c.Items {
		ch.add(it.node())
	}
	if c.DCDate != nil {
		ch.add(leaf("dc:date", formatW3CDTF(*c.DCDate)))
	}
	return ch
}

func (im *Image) node() *node {
	n := newNode("image")
	addLeaf(n, "url", im.URL, true)
	addLeaf(n, "title", im.Title, true)
	addLeaf(n, "link", im.Link, true)
	addLeaf(n, "width", im.Width, false)
	addLeaf(n, "height", im.Height, false)
	return n
}

func (it *Item) node() *node {
	n := newNode("item")
	addLeaf(n, "title", it.Title, false)
	addLeaf(n, "link", it.Link, false)
	addLeaf(n, "description", it.Description, false)
	addLeaf(n, "author", it.Author, false)
	for _, cat := range it.Categories {
		n.add(leaf("category", cat))
	}
	addLeaf(n, "comments", it.Comments, false)
	if it.PubDate != nil {
		n.add(leaf("pubDate", formatRFC822(*it.PubDate)))
	}
	if it.Guid != nil {
		n.add(it.Guid.node())
	}
	if it.ContentEncoded != "" {
		n.add(leaf("content:encoded", it.ContentEncoded))
	}
	addLeaf(n, "dc:creator", it.DCCreator, false)
	addLeaf(n, "dc:subject", it.DCSubject, false)
	if it.DCDate != nil {
		n.add(leaf("dc:date", formatW3CDTF(*it.DCDate)))
	}
	return n
}

func (g *Guid) node() *node {
	n := leaf("guid", g.Content)
	if g.HasPermaLink {
		if g.IsPermaLink {
			n.setAttr("isPermaLink", "true")
		} else {
			n.setAttr("isPermaLink", "false")
		}
	}
	return n
}

// addLeaf appends <name>value</name>. When required is false, an empty value
// is skipped (MRI omits unset optional elements). A required element is always
// emitted even when empty, matching MRI's behavior for the mandatory fields.
func addLeaf(parent *node, name, value string, required bool) {
	if value == "" && !required {
		return
	}
	parent.add(leaf(name, value))
}
