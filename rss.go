// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

// Package rss is a pure-Go (CGO=0) reimplementation of Ruby's standard
// library RSS module (gem rss 0.3.2, shipped with MRI 4.0.5).
//
// It parses and generates the three feed dialects MRI's RSS supports:
//
//   - RSS 2.0 (and the historical 0.9x family), rooted at <rss>;
//   - RSS 1.0, the RDF dialect rooted at <rdf:RDF>;
//   - Atom (RFC 4287), rooted at <feed>.
//
// Parser.Parse auto-detects the dialect from the document's root element,
// exactly like RSS::Parser.parse. The generated XML is byte-for-byte
// compatible with MRI's to_s for the round-trip (parse a feed, re-serialize
// it) on the supported element set, including MRI's 2-space indentation, its
// always-on namespace declarations on the root element, and its date
// formats (RFC822 for RSS <pubDate>/<lastBuildDate>, W3CDTF/RFC3339 for
// <dc:date> and every Atom date).
//
// The XML layer (tokenizing on parse, attribute/text escaping) is shared with
// github.com/go-ruby-rexml/rexml, mirroring how MRI's RSS sits on REXML.
//
// Bundled modules: Dublin Core (dc:), content (content:encoded) and
// syndication (sy:) — the common modules MRI bundles and emits by default.
// See the package documentation in the repository README for the exact
// element coverage boundary.
package rss

// Namespace URIs, matching RSS::* constants in MRI.
const (
	uriRSS10    = "http://purl.org/rss/1.0/"
	uriRDF      = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	uriContent  = "http://purl.org/rss/1.0/modules/content/"
	uriDC       = "http://purl.org/dc/elements/1.1/"
	uriSy       = "http://purl.org/rss/1.0/modules/syndication/"
	uriItunes   = "http://www.itunes.com/dtds/podcast-1.0.dtd"
	uriTrackbck = "http://madskills.com/public/xml/rss/module/trackback/"
	uriImage    = "http://purl.org/rss/1.0/modules/image/"
	uriSlash    = "http://purl.org/rss/1.0/modules/slash/"
	uriTaxo     = "http://purl.org/rss/1.0/modules/taxonomy/"
	uriAtom     = "http://www.w3.org/2005/Atom"
)

// indent is RSS::Element::INDENT — two spaces.
const indent = "  "
