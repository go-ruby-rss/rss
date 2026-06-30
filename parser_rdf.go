// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import "github.com/go-ruby-rexml/rexml"

// parseRDF parses an <rdf:RDF> root (RSS 1.0) into *RDF.
func parseRDF(root *rexml.Element) (*RDF, error) {
	r := &RDF{parsed: true}
	for _, c := range root.ChildElements() {
		switch localName(c.QName()) {
		case "channel":
			ch, err := parseRDFChannel(c)
			if err != nil {
				return nil, err
			}
			r.Channel = ch
		case "image":
			r.Image = parseRDFImage(c)
		case "item":
			it, err := parseRDFItem(c)
			if err != nil {
				return nil, err
			}
			r.Items = append(r.Items, it)
		case "textinput":
			r.Textinput = parseRDFTextinput(c)
		}
	}
	return r, nil
}

func rdfAbout(e *rexml.Element) string {
	if v, ok := e.Attr("rdf:about"); ok {
		return v
	}
	return ""
}

func rdfResource(e *rexml.Element) string {
	if v, ok := e.Attr("rdf:resource"); ok {
		return v
	}
	return ""
}

func parseRDFChannel(c *rexml.Element) (*RDFChannel, error) {
	ch := &RDFChannel{About: rdfAbout(c)}
	for _, e := range c.ChildElements() {
		switch e.QName() {
		case "title":
			ch.Title = e.Text()
		case "link":
			ch.Link = e.Text()
		case "description":
			ch.Description = e.Text()
		case "image":
			ch.ImageResource = rdfResource(e)
		case "textinput":
			ch.TextinputResource = rdfResource(e)
		case "dc:creator":
			ch.DCCreator = e.Text()
		case "dc:date":
			t, err := parseW3CDTF(e.Text())
			if err != nil {
				return nil, err
			}
			ch.DCDate = &t
		case "sy:updatePeriod":
			ch.SyUpdatePeriod = e.Text()
		case "sy:updateFrequency":
			ch.SyUpdateFrequency = e.Text()
		case "sy:updateBase":
			ch.SyUpdateBase = e.Text()
		case "items":
			ch.ItemResources = parseRDFSeq(e)
		}
	}
	return ch, nil
}

// parseRDFSeq extracts the resource URIs of <items><rdf:Seq><rdf:li .../>.
func parseRDFSeq(items *rexml.Element) []string {
	var out []string
	seq := childElem(items, "Seq")
	if seq == nil {
		return out
	}
	for _, li := range seq.ChildElements() {
		if localName(li.QName()) == "li" {
			if v, ok := li.Attr("resource"); ok {
				out = append(out, v)
			} else if v, ok := li.Attr("rdf:resource"); ok {
				out = append(out, v)
			}
		}
	}
	return out
}

func parseRDFImage(c *rexml.Element) *RDFImage {
	im := &RDFImage{About: rdfAbout(c)}
	if v, ok := childText(c, "title"); ok {
		im.Title = v
	}
	if v, ok := childText(c, "url"); ok {
		im.URL = v
	}
	if v, ok := childText(c, "link"); ok {
		im.Link = v
	}
	return im
}

func parseRDFItem(c *rexml.Element) (*RDFItem, error) {
	it := &RDFItem{About: rdfAbout(c)}
	for _, e := range c.ChildElements() {
		switch e.QName() {
		case "title":
			it.Title = e.Text()
		case "link":
			it.Link = e.Text()
		case "description":
			it.Description = e.Text()
		case "dc:creator":
			it.DCCreator = e.Text()
		case "dc:subject":
			it.DCSubject = e.Text()
		case "content:encoded":
			it.ContentEncoded = e.Text()
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

func parseRDFTextinput(c *rexml.Element) *RDFTextinput {
	ti := &RDFTextinput{About: rdfAbout(c)}
	if v, ok := childText(c, "title"); ok {
		ti.Title = v
	}
	if v, ok := childText(c, "description"); ok {
		ti.Description = v
	}
	if v, ok := childText(c, "name"); ok {
		ti.Name = v
	}
	if v, ok := childText(c, "link"); ok {
		ti.Link = v
	}
	return ti
}
