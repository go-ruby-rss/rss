// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import (
	"fmt"
	"strings"

	"github.com/go-ruby-rexml/rexml"
)

// Parse parses an RSS/Atom document and auto-detects the dialect from the root
// element, mirroring RSS::Parser.parse:
//
//   - <rss>      → *Rss      (RSS 0.9x / 2.0)
//   - <rdf:RDF>  → *RDF      (RSS 1.0)
//   - <feed>     → *AtomFeed (Atom)
//
// It returns an error for malformed XML or an unrecognized root element.
func Parse(xml string) (Feed, error) {
	doc, err := rexml.ParseDocument(xml)
	if err != nil {
		return nil, fmt.Errorf("rss: %w", err)
	}
	// rexml.ParseDocument guarantees a non-nil root on success (it errors on a
	// rootless document), so doc.Root() is always set here.
	root := doc.Root()
	switch localName(root.QName()) {
	case "rss":
		return parseRss(root)
	case "RDF":
		return parseRDF(root)
	case "feed":
		return parseAtom(root)
	default:
		return nil, fmt.Errorf("rss: unknown feed root <%s>", root.QName())
	}
}

// localName strips a namespace prefix, e.g. "rdf:RDF" → "RDF".
func localName(qname string) string {
	if i := strings.IndexByte(qname, ':'); i >= 0 {
		return qname[i+1:]
	}
	return qname
}

// childText returns the trimmed text of the first child element whose local
// name matches, and whether it was present.
func childText(e *rexml.Element, name string) (string, bool) {
	if c := childElem(e, name); c != nil {
		return c.Text(), true
	}
	return "", false
}

// childElem returns the first child element with the given local name.
func childElem(e *rexml.Element, name string) *rexml.Element {
	for _, c := range e.ChildElements() {
		if localName(c.QName()) == name {
			return c
		}
	}
	return nil
}

// parseRss parses an <rss> root into *Rss.
func parseRss(root *rexml.Element) (*Rss, error) {
	r := &Rss{parsed: true}
	if v, ok := root.Attr("version"); ok {
		r.Version = v
	}
	chElem := childElem(root, "channel")
	if chElem == nil {
		return nil, fmt.Errorf("rss: <rss> without <channel>")
	}
	ch := &Channel{}
	r.Channel = ch
	for _, c := range chElem.ChildElements() {
		if err := assignChannelChild(ch, c); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func assignChannelChild(ch *Channel, c *rexml.Element) error {
	switch c.QName() {
	case "title":
		ch.Title = c.Text()
	case "link":
		ch.Link = c.Text()
	case "description":
		ch.Description = c.Text()
	case "language":
		ch.Language = c.Text()
	case "copyright":
		ch.Copyright = c.Text()
	case "managingEditor":
		ch.ManagingEditor = c.Text()
	case "webMaster":
		ch.WebMaster = c.Text()
	case "generator":
		ch.Generator = c.Text()
	case "docs":
		ch.Docs = c.Text()
	case "ttl":
		ch.TTL = c.Text()
	case "category":
		ch.Categories = append(ch.Categories, c.Text())
	case "pubDate":
		t, err := parseRFC822(c.Text())
		if err != nil {
			return err
		}
		ch.PubDate = &t
	case "lastBuildDate":
		t, err := parseRFC822(c.Text())
		if err != nil {
			return err
		}
		ch.LastBuildDate = &t
	case "dc:date":
		t, err := parseW3CDTF(c.Text())
		if err != nil {
			return err
		}
		ch.DCDate = &t
	case "sy:updatePeriod":
		ch.SyUpdatePeriod = c.Text()
	case "sy:updateFrequency":
		ch.SyUpdateFrequency = c.Text()
	case "sy:updateBase":
		ch.SyUpdateBase = c.Text()
	case "image":
		ch.Image = parseImage(c)
	case "item":
		it, err := parseItem(c)
		if err != nil {
			return err
		}
		ch.Items = append(ch.Items, it)
	}
	return nil
}

func parseImage(c *rexml.Element) *Image {
	im := &Image{}
	if v, ok := childText(c, "url"); ok {
		im.URL = v
	}
	if v, ok := childText(c, "title"); ok {
		im.Title = v
	}
	if v, ok := childText(c, "link"); ok {
		im.Link = v
	}
	if v, ok := childText(c, "width"); ok {
		im.Width = v
	}
	if v, ok := childText(c, "height"); ok {
		im.Height = v
	}
	return im
}

func parseItem(c *rexml.Element) (*Item, error) {
	it := &Item{}
	for _, e := range c.ChildElements() {
		switch e.QName() {
		case "title":
			it.Title = e.Text()
		case "link":
			it.Link = e.Text()
		case "description":
			it.Description = e.Text()
		case "author":
			it.Author = e.Text()
		case "comments":
			it.Comments = e.Text()
		case "category":
			it.Categories = append(it.Categories, e.Text())
		case "pubDate":
			t, err := parseRFC822(e.Text())
			if err != nil {
				return nil, err
			}
			it.PubDate = &t
		case "guid":
			g := &Guid{Content: e.Text()}
			if v, ok := e.Attr("isPermaLink"); ok {
				g.HasPermaLink = true
				g.IsPermaLink = v == "true"
			}
			it.Guid = g
		case "content:encoded":
			it.ContentEncoded = e.Text()
		case "dc:creator":
			it.DCCreator = e.Text()
		case "dc:subject":
			it.DCSubject = e.Text()
		case "dc:date":
			t, err := parseW3CDTF(e.Text())
			if err != nil {
				return nil, err
			}
			it.DCDate = &t
		}
	}
	return it, nil
}
