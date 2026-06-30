// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import "time"

// Feed is the common interface returned by Parser.Parse. The concrete type is
// *Rss (RSS 0.9x/2.0), *RDF (RSS 1.0) or *AtomFeed (Atom). FeedType reports
// the dialect, mirroring RSS feed#feed_type.
type Feed interface {
	// FeedType is "rss", "rss" (for 1.0 too, via RDF) or "atom" — see the
	// concrete types for the exact value. Use the type switch for dispatch.
	FeedType() string
	// String serializes the feed to MRI-compatible XML.
	String() string
}

// ---- RSS 0.9x / 2.0 -------------------------------------------------------

// Rss is the <rss> document of RSS 0.91/0.92/2.0. It mirrors RSS::Rss.
type Rss struct {
	Version string // e.g. "2.0"
	Channel *Channel

	// parsed records that this feed came from Parse, which makes String emit
	// encoding="UTF-8" in the XML prolog, matching MRI: its parser defaults
	// the document encoding to UTF-8, so reserialized feeds carry it, whereas
	// directly-constructed feeds do not.
	parsed bool
}

// Channel mirrors RSS::Rss::Channel.
type Channel struct {
	Title             string
	Link              string
	Description       string
	Language          string
	Copyright         string
	ManagingEditor    string
	WebMaster         string
	PubDate           *time.Time // <pubDate>, RFC822
	LastBuildDate     *time.Time // <lastBuildDate>, RFC822
	Generator         string
	Docs              string
	TTL               string
	Image             *Image
	Categories        []string
	DCDate            *time.Time // dc:date, W3CDTF
	SyUpdatePeriod    string     // sy:updatePeriod
	SyUpdateFrequency string     // sy:updateFrequency
	SyUpdateBase      string     // sy:updateBase
	Items             []*Item
}

// Image mirrors RSS::Rss::Channel::Image.
type Image struct {
	URL    string
	Title  string
	Link   string
	Width  string
	Height string
}

// Item mirrors RSS::Rss::Channel::Item.
type Item struct {
	Title          string
	Link           string
	Description    string
	Author         string
	Comments       string
	PubDate        *time.Time // <pubDate>, RFC822
	Guid           *Guid
	Categories     []string
	DCDate         *time.Time // dc:date, W3CDTF
	DCCreator      string     // dc:creator
	DCSubject      string     // dc:subject
	ContentEncoded string     // content:encoded
}

// Guid mirrors RSS::Rss::Channel::Item::Guid. IsPermaLink is a tri-state in
// MRI (nil/true/false); HasPermaLink reports whether the attribute is set.
type Guid struct {
	Content      string
	IsPermaLink  bool
	HasPermaLink bool
}

func (r *Rss) FeedType() string { return "rss" }

// ---- RSS 1.0 (RDF) --------------------------------------------------------

// RDF is the <rdf:RDF> document of RSS 1.0. It mirrors RSS::RDF.
type RDF struct {
	Channel   *RDFChannel
	Image     *RDFImage
	Items     []*RDFItem
	Textinput *RDFTextinput

	parsed bool // see Rss.parsed
}

// RDFChannel mirrors RSS::RDF::Channel. About is the rdf:about attribute.
type RDFChannel struct {
	About             string
	Title             string
	Link              string
	Description       string
	DCDate            *time.Time
	DCCreator         string
	SyUpdatePeriod    string
	SyUpdateFrequency string
	SyUpdateBase      string
	ImageResource     string // <image rdf:resource="..."/>, when an image is present
	// ItemResources is the <items><rdf:Seq><rdf:li resource="..."/> sequence.
	ItemResources     []string
	TextinputResource string
}

// RDFImage mirrors RSS::RDF::Image.
type RDFImage struct {
	About string
	Title string
	URL   string
	Link  string
}

// RDFItem mirrors RSS::RDF::Item.
type RDFItem struct {
	About          string
	Title          string
	Link           string
	Description    string
	DCDate         *time.Time
	DCCreator      string
	DCSubject      string
	ContentEncoded string
}

// RDFTextinput mirrors RSS::RDF::Textinput.
type RDFTextinput struct {
	About       string
	Title       string
	Description string
	Name        string
	Link        string
}

func (r *RDF) FeedType() string { return "rss" }

// ---- Atom -----------------------------------------------------------------

// AtomFeed is the <feed> document of Atom (RFC 4287). It mirrors
// RSS::Atom::Feed.
type AtomFeed struct {
	ID         string
	Title      string
	Subtitle   string
	Updated    *time.Time // W3CDTF
	Rights     string
	Generator  string
	Authors    []*AtomPerson
	Links      []*AtomLink
	Categories []*AtomCategory
	Entries    []*AtomEntry

	parsed bool // see Rss.parsed
}

// AtomEntry mirrors RSS::Atom::Feed::Entry.
type AtomEntry struct {
	ID         string
	Title      string
	Summary    string
	Content    string
	Updated    *time.Time
	Published  *time.Time
	Rights     string
	Authors    []*AtomPerson
	Links      []*AtomLink
	Categories []*AtomCategory
}

// AtomPerson mirrors an Atom person construct (author/contributor).
type AtomPerson struct {
	Name  string
	URI   string
	Email string
}

// AtomLink mirrors RSS::Atom::*::Link.
type AtomLink struct {
	Href  string
	Rel   string
	Type  string
	Title string
}

// AtomCategory mirrors RSS::Atom::*::Category.
type AtomCategory struct {
	Term   string
	Scheme string
	Label  string
}

func (f *AtomFeed) FeedType() string { return "atom" }
