// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import "github.com/go-ruby-rexml/rexml"

// parseAtom parses an Atom <feed> root into *AtomFeed.
func parseAtom(root *rexml.Element) (*AtomFeed, error) {
	f := &AtomFeed{parsed: true}
	for _, c := range root.ChildElements() {
		switch localName(c.QName()) {
		case "id":
			f.ID = c.Text()
		case "title":
			f.Title = c.Text()
		case "subtitle":
			f.Subtitle = c.Text()
		case "rights":
			f.Rights = c.Text()
		case "generator":
			f.Generator = c.Text()
		case "updated":
			t, err := parseW3CDTF(c.Text())
			if err != nil {
				return nil, err
			}
			f.Updated = &t
		case "author":
			f.Authors = append(f.Authors, parseAtomPerson(c))
		case "link":
			f.Links = append(f.Links, parseAtomLink(c))
		case "category":
			f.Categories = append(f.Categories, parseAtomCategory(c))
		case "entry":
			e, err := parseAtomEntry(c)
			if err != nil {
				return nil, err
			}
			f.Entries = append(f.Entries, e)
		}
	}
	return f, nil
}

func parseAtomPerson(c *rexml.Element) *AtomPerson {
	p := &AtomPerson{}
	if v, ok := childText(c, "name"); ok {
		p.Name = v
	}
	if v, ok := childText(c, "uri"); ok {
		p.URI = v
	}
	if v, ok := childText(c, "email"); ok {
		p.Email = v
	}
	return p
}

func parseAtomLink(c *rexml.Element) *AtomLink {
	l := &AtomLink{}
	if v, ok := c.Attr("href"); ok {
		l.Href = v
	}
	if v, ok := c.Attr("rel"); ok {
		l.Rel = v
	}
	if v, ok := c.Attr("type"); ok {
		l.Type = v
	}
	if v, ok := c.Attr("title"); ok {
		l.Title = v
	}
	return l
}

func parseAtomCategory(c *rexml.Element) *AtomCategory {
	cat := &AtomCategory{}
	if v, ok := c.Attr("term"); ok {
		cat.Term = v
	}
	if v, ok := c.Attr("scheme"); ok {
		cat.Scheme = v
	}
	if v, ok := c.Attr("label"); ok {
		cat.Label = v
	}
	return cat
}

func parseAtomEntry(c *rexml.Element) (*AtomEntry, error) {
	e := &AtomEntry{}
	for _, x := range c.ChildElements() {
		switch localName(x.QName()) {
		case "id":
			e.ID = x.Text()
		case "title":
			e.Title = x.Text()
		case "summary":
			e.Summary = x.Text()
		case "content":
			e.Content = x.Text()
		case "rights":
			e.Rights = x.Text()
		case "updated":
			t, err := parseW3CDTF(x.Text())
			if err != nil {
				return nil, err
			}
			e.Updated = &t
		case "published":
			t, err := parseW3CDTF(x.Text())
			if err != nil {
				return nil, err
			}
			e.Published = &t
		case "author":
			e.Authors = append(e.Authors, parseAtomPerson(x))
		case "link":
			e.Links = append(e.Links, parseAtomLink(x))
		case "category":
			e.Categories = append(e.Categories, parseAtomCategory(x))
		}
	}
	return e, nil
}
