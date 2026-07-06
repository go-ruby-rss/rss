<p align="center"><img src="https://raw.githubusercontent.com/go-ruby-rss/brand/main/social/go-ruby-rss-rss.png" alt="go-ruby-rss/rss" width="720"></p>

# rss — go-ruby-rss

[![Docs](https://img.shields.io/badge/docs-mkdocs--material-DC2626)](https://go-ruby-rss.github.io/docs/)
[![License](https://img.shields.io/badge/license-BSD--3--Clause-blue)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.26.4%2B-00ADD8)](https://go.dev/dl/)
[![Coverage](https://img.shields.io/badge/coverage-100%25-1a7f37)](#tests--coverage)

**A pure-Go (no cgo) reimplementation of Ruby's standard-library
[`RSS`](https://docs.ruby-lang.org/en/master/RSS.html) module** (gem `rss`
0.3.2, shipped with MRI 4.0.5) — it **parses** and **generates** the three feed
dialects Ruby's `RSS` supports, **without any Ruby runtime**:

- **RSS 2.0** (and the historical 0.9x family), rooted at `<rss>`;
- **RSS 1.0**, the RDF dialect rooted at `<rdf:RDF>`;
- **Atom** (RFC 4287), rooted at `<feed>`.

`Parse` auto-detects the dialect from the document's root element, exactly like
`RSS::Parser.parse`, and the generated XML is **byte-for-byte identical to MRI's
`to_s`** on the supported element set — including MRI's two-space indentation,
its always-on namespace declarations on the root element, and its date formats
(RFC822 for RSS `<pubDate>`/`<lastBuildDate>`, W3CDTF/RFC3339 for `<dc:date>`
and every Atom date).

It is built on [go-ruby-rexml](https://github.com/go-ruby-rexml/rexml) for the
XML layer (the same way MRI's `RSS` sits on REXML), and is the feed backend for
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) — but a
**standalone, reusable** module, a sibling of
[go-ruby-regexp](https://github.com/go-ruby-regexp/regexp),
[go-ruby-erb](https://github.com/go-ruby-erb/erb) and
[go-ruby-marshal](https://github.com/go-ruby-marshal/marshal).

> **What it is — and isn't.** Parsing and generating the feed XML for the Ruby
> value model (the element tree, typed accessors, and date formats) is fully
> deterministic and needs **no interpreter**, so it lives here as pure Go. It
> hands back a small, explicit value model (`*Rss`, `*RDF`, `*AtomFeed`, …) the
> host maps to and from its own objects.

## Features

- **Parse** any of the three dialects with `Parse(xml)`, which returns a `Feed`
  whose concrete type (`*Rss` / `*RDF` / `*AtomFeed`) is selected by the root
  element. Typed accessors mirror MRI's: `feed.Channel.Title`,
  `item.Guid.IsPermaLink`, `entry.Updated`, and so on. Dates parse to
  `*time.Time` preserving their original offset (RFC822 named zones such as
  `GMT`/`EST` included).
- **Generate** spec-valid XML with `(&Rss{…}).String()`, `(&RDF{…}).String()`,
  `(&AtomFeed{…}).String()`, matching MRI's element order, indentation,
  namespace block and date formatting.
- **Bundled modules** MRI emits by default: **Dublin Core** (`dc:date`,
  `dc:creator`, `dc:subject`), **content** (`content:encoded`) and
  **syndication** (`sy:updatePeriod`/`updateFrequency`/`updateBase`).
- **Date formats** faithful to MRI's `Time#rfc822` and `Time#w3cdtf`: RFC822
  for RSS pub/build dates, W3CDTF (RFC3339, `Z` for UTC, trimmed fractional
  seconds) for Dublin Core and Atom dates.

CGO-free, **100% test coverage**, `gofmt` + `go vet` clean, and green across the
six 64-bit Go targets (amd64, arm64, riscv64, loong64, ppc64le, s390x) on Linux,
macOS and Windows.

## Install

```sh
go get github.com/go-ruby-rss/rss
```

## Usage

### Parse (auto-detect dialect)

```go
feed, err := rss.Parse(xml)
if err != nil { /* ... */ }

switch f := feed.(type) {
case *rss.Rss: // RSS 2.0 / 0.9x
    fmt.Println(f.Channel.Title, f.Channel.Items[0].PubDate)
case *rss.RDF: // RSS 1.0
    fmt.Println(f.Channel.About, f.Items[0].DCDate)
case *rss.AtomFeed: // Atom
    fmt.Println(f.Title, f.Entries[0].Updated)
}
```

### Generate RSS 2.0

```go
t := time.Unix(1700000000, 0)
r := &rss.Rss{Version: "2.0", Channel: &rss.Channel{
    Title: "Example Feed", Link: "http://example.com/",
    Description: "An example.", PubDate: &t,
    Items: []*rss.Item{{
        Title: "Item One", Link: "http://example.com/1",
        Guid: &rss.Guid{Content: "http://example.com/1", IsPermaLink: true, HasPermaLink: true},
    }},
}}
fmt.Print(r.String())
```

This produces the same bytes as `RSS::Rss.new("2.0")` … `.to_s` in MRI 4.0.5,
namespace block and all.

### Generate Atom

```go
f := &rss.AtomFeed{
    ID: "urn:uuid:1", Title: "Atom Example", Updated: &t,
    Authors: []*rss.AtomPerson{{Name: "Jane"}},
    Links:   []*rss.AtomLink{{Href: "http://example.com/"}},
}
fmt.Print(f.String())
```

## Coverage boundary

This library targets the **common, default-emitted** element set MRI's `RSS`
exposes, which round-trips real feeds byte-for-byte. Specifically supported:

| Dialect  | Supported elements |
| -------- | ------------------ |
| RSS 2.0  | `<channel>`: title, link, description, language, copyright, managingEditor, webMaster, pubDate, lastBuildDate, generator, docs, ttl, category, image. `<item>`: title, link, description, author, comments, category, pubDate, guid (+ isPermaLink). |
| RSS 1.0  | `<channel>` (rdf:about), title, link, description, image/textinput resources, `<items><rdf:Seq>`. `<image>`, `<item>`, `<textinput>` resources. |
| Atom     | `<feed>`/`<entry>`: id, title, subtitle, rights, generator, updated, published, author/contributor (name/uri/email), link (href/rel/type/title), category (term/scheme/label), content, summary. |
| Modules  | **Dublin Core** `dc:date` / `dc:creator` / `dc:subject`; **content** `content:encoded`; **syndication** `sy:updatePeriod` / `sy:updateFrequency` / `sy:updateBase`. |

**Outside the boundary** (parsed leniently — unknown elements are ignored rather
than rejected): MRI's strict schema *validation* (e.g. raising `MissingTagError`
for an RSS 1.0 channel without `<items>`), the high-level `RSS::Maker` builder
DSL, the iTunes/taxonomy/trackback/slash/image modules (their namespaces are
still declared on the root to match MRI's output), and per-element `xml:lang` /
`xml:base` attributes. The namespace declarations MRI always emits on the root
element are reproduced verbatim so generation matches byte-for-byte.

## How rexml is used

The **parse** path tokenizes the document with
[`rexml.ParseDocument`](https://github.com/go-ruby-rexml/rexml) and walks the
resulting element tree (`ChildElements`, `Text`, `Attr`, `QName`) to populate
the typed model. The **generate** path mirrors MRI, which builds its XML as
strings directly (it does *not* route `to_s` through REXML): a small internal
node builder replicates `RSS::Element#tag` — two-space indent, attribute
wrapping, `CGI.escapeHTML` text/attribute escaping, and the empty-element rules
— so the bytes match exactly.

## Tests & coverage

```sh
go test -cover ./...
```

The suite is **deterministic and Ruby-free** (golden vectors captured from MRI
4.0.5), keeping coverage at **100%** on every platform — including the
cross-arch qemu lanes and the Windows lane where no `ruby` is present. A
differential **MRI oracle** (`oracle_test.go`) additionally runs the real `ruby
-rrss` binary where available: it parses each sample feed through MRI and through
this library and asserts the reserialized XML is byte-identical, and checks the
date formatters against `Time#rfc822` / `Time#w3cdtf`. The oracle self-skips when
`ruby` is absent.

## License

BSD-3-Clause — see [LICENSE](LICENSE). Copyright (c) the go-ruby-rss/rss authors.

## WebAssembly

Being pure Go (CGO=0), this library also compiles to **WebAssembly** — both
`GOOS=js GOARCH=wasm` (browser / Node.js) and `GOOS=wasip1 GOARCH=wasm` (WASI).
CI builds both targets on every push, alongside the six 64-bit native/qemu arches.

```sh
GOOS=js     GOARCH=wasm go build ./...   # browser / Node
GOOS=wasip1 GOARCH=wasm go build ./...   # WASI (wasmtime, wasmer, wasmedge, …)
```
