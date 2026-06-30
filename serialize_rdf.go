// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

// String serializes the RSS 1.0 RDF document to MRI-compatible XML.
//
// MRI declares, on the root <rdf:RDF>, the default RSS 1.0 namespace and the
// rdf: namespace first, then the bundled-module namespaces (content, dc,
// image, slash, sy, taxo, trackback) in that fixed order, regardless of use.
// Children appear as <channel>, then <image>, then the <item> sequence, then
// <textinput>. The channel carries an <items><rdf:Seq><rdf:li resource=…/>
// table of item URIs.
func (r *RDF) String() string {
	root := newNode("rdf:RDF")
	root.setAttr("xmlns", uriRSS10)
	root.setAttr("xmlns:rdf", uriRDF)
	root.setAttr("xmlns:content", uriContent)
	root.setAttr("xmlns:dc", uriDC)
	root.setAttr("xmlns:image", uriImage)
	root.setAttr("xmlns:slash", uriSlash)
	root.setAttr("xmlns:sy", uriSy)
	root.setAttr("xmlns:taxo", uriTaxo)
	root.setAttr("xmlns:trackback", uriTrackbck)
	if r.Channel != nil {
		root.add(r.Channel.node())
	}
	if r.Image != nil {
		root.add(r.Image.node())
	}
	for _, it := range r.Items {
		root.add(it.node())
	}
	if r.Textinput != nil {
		root.add(r.Textinput.node())
	}
	return prolog(r.parsed) + root.render("")
}

func (c *RDFChannel) node() *node {
	n := newNode("channel")
	if c.About != "" {
		n.setAttr("rdf:about", c.About)
	}
	addLeaf(n, "title", c.Title, true)
	addLeaf(n, "link", c.Link, true)
	addLeaf(n, "description", c.Description, true)
	if c.ImageResource != "" {
		img := newNode("image")
		img.setAttr("rdf:resource", c.ImageResource)
		n.add(img)
	}
	if len(c.ItemResources) > 0 {
		items := newNode("items")
		seq := newNode("rdf:Seq")
		for _, res := range c.ItemResources {
			li := newNode("rdf:li")
			li.setAttr("resource", res)
			seq.add(li)
		}
		items.add(seq)
		n.add(items)
	}
	if c.TextinputResource != "" {
		ti := newNode("textinput")
		ti.setAttr("rdf:resource", c.TextinputResource)
		n.add(ti)
	}
	addLeaf(n, "dc:creator", c.DCCreator, false)
	if c.SyUpdatePeriod != "" {
		n.add(leaf("sy:updatePeriod", c.SyUpdatePeriod))
	}
	if c.SyUpdateFrequency != "" {
		n.add(leaf("sy:updateFrequency", c.SyUpdateFrequency))
	}
	if c.SyUpdateBase != "" {
		n.add(leaf("sy:updateBase", c.SyUpdateBase))
	}
	if c.DCDate != nil {
		n.add(leaf("dc:date", formatW3CDTF(*c.DCDate)))
	}
	return n
}

func (im *RDFImage) node() *node {
	n := newNode("image")
	if im.About != "" {
		n.setAttr("rdf:about", im.About)
	}
	addLeaf(n, "title", im.Title, true)
	addLeaf(n, "url", im.URL, true)
	addLeaf(n, "link", im.Link, true)
	return n
}

func (it *RDFItem) node() *node {
	n := newNode("item")
	if it.About != "" {
		n.setAttr("rdf:about", it.About)
	}
	addLeaf(n, "title", it.Title, true)
	addLeaf(n, "link", it.Link, true)
	addLeaf(n, "description", it.Description, false)
	addLeaf(n, "dc:creator", it.DCCreator, false)
	addLeaf(n, "dc:subject", it.DCSubject, false)
	if it.ContentEncoded != "" {
		n.add(leaf("content:encoded", it.ContentEncoded))
	}
	if it.DCDate != nil {
		n.add(leaf("dc:date", formatW3CDTF(*it.DCDate)))
	}
	return n
}

func (ti *RDFTextinput) node() *node {
	n := newNode("textinput")
	if ti.About != "" {
		n.setAttr("rdf:about", ti.About)
	}
	addLeaf(n, "title", ti.Title, true)
	addLeaf(n, "description", ti.Description, true)
	addLeaf(n, "name", ti.Name, true)
	addLeaf(n, "link", ti.Link, true)
	return n
}
