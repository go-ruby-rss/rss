// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import (
	"strings"
	"testing"
	"time"
)

const goldenRdf = `<?xml version="1.0"?>
<rdf:RDF xmlns="http://purl.org/rss/1.0/"
  xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
  xmlns:content="http://purl.org/rss/1.0/modules/content/"
  xmlns:dc="http://purl.org/dc/elements/1.1/"
  xmlns:image="http://purl.org/rss/1.0/modules/image/"
  xmlns:slash="http://purl.org/rss/1.0/modules/slash/"
  xmlns:sy="http://purl.org/rss/1.0/modules/syndication/"
  xmlns:taxo="http://purl.org/rss/1.0/modules/taxonomy/"
  xmlns:trackback="http://madskills.com/public/xml/rss/module/trackback/">
  <channel rdf:about="http://example.com/">
    <title>Example 1.0</title>
    <link>http://example.com/</link>
    <description>RDF feed.</description>
    <items>
      <rdf:Seq>
        <rdf:li resource="http://example.com/1"/>
      </rdf:Seq>
    </items>
    <dc:date>2023-11-14T23:13:20+01:00</dc:date>
  </channel>
  <item rdf:about="http://example.com/1">
    <title>Item One</title>
    <link>http://example.com/1</link>
    <description>First.</description>
  </item>
</rdf:RDF>`

func sampleRdf() *RDF {
	return &RDF{
		Channel: &RDFChannel{
			About: "http://example.com/", Title: "Example 1.0",
			Link: "http://example.com/", Description: "RDF feed.",
			DCDate:        ptime(2023, 11, 14, 23, 13, 20, loc1),
			ItemResources: []string{"http://example.com/1"},
		},
		Items: []*RDFItem{{
			About: "http://example.com/1", Title: "Item One",
			Link: "http://example.com/1", Description: "First.",
		}},
	}
}

func TestRdfGenerate(t *testing.T) {
	if got := sampleRdf().String(); got != goldenRdf {
		t.Errorf("RDF generate mismatch:\n got:\n%s\nwant:\n%s", got, goldenRdf)
	}
	if ft := sampleRdf().FeedType(); ft != "rss" {
		t.Errorf("FeedType = %q, want rss", ft)
	}
}

func TestRdfFull(t *testing.T) {
	r := &RDF{
		Channel: &RDFChannel{
			About: "a", Title: "t", Link: "l", Description: "d",
			ImageResource: "img", TextinputResource: "ti",
			ItemResources: []string{"i1"},
			DCCreator:     "dc", DCDate: ptime(2023, 1, 1, 0, 0, 0, time.UTC),
			SyUpdatePeriod: "hourly", SyUpdateFrequency: "1", SyUpdateBase: "b",
		},
		Image: &RDFImage{About: "ia", Title: "it", URL: "u", Link: "il"},
		Items: []*RDFItem{{
			About: "i1", Title: "t", Link: "l", Description: "d",
			DCCreator: "c", DCSubject: "s", ContentEncoded: "body",
			DCDate: ptime(2023, 1, 1, 0, 0, 0, loc1),
		}},
		Textinput: &RDFTextinput{About: "ta", Title: "tt", Description: "td", Name: "tn", Link: "tl"},
	}
	out := r.String()
	for _, want := range []string{
		`<image rdf:resource="img"/>`, `<textinput rdf:resource="ti"/>`,
		"<dc:creator>dc</dc:creator>", "<sy:updatePeriod>hourly</sy:updatePeriod>",
		`<image rdf:about="ia">`, "<content:encoded>body</content:encoded>",
		`<textinput rdf:about="ta">`, "<dc:date>2023-01-01T00:00:00Z</dc:date>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("full RDF missing %q in:\n%s", want, out)
		}
	}
}

func TestRdfRoundTrip(t *testing.T) {
	f, err := Parse(goldenRdf)
	if err != nil {
		t.Fatal(err)
	}
	r, ok := f.(*RDF)
	if !ok {
		t.Fatalf("Parse returned %T, want *RDF", f)
	}
	if r.Channel.About != "http://example.com/" || r.Channel.DCDate == nil ||
		len(r.Channel.ItemResources) != 1 || r.Items[0].Title != "Item One" {
		t.Errorf("RDF parse wrong: %+v", r.Channel)
	}
}

func TestRdfParseAll(t *testing.T) {
	src := `<?xml version="1.0"?>
<rdf:RDF xmlns="http://purl.org/rss/1.0/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:dc="x" xmlns:content="x" xmlns:sy="x">
  <channel rdf:about="a">
    <title>t</title><link>l</link><description>d</description>
    <image rdf:resource="img"/><textinput rdf:resource="ti"/>
    <dc:creator>dc</dc:creator>
    <sy:updatePeriod>hourly</sy:updatePeriod>
    <sy:updateFrequency>1</sy:updateFrequency>
    <sy:updateBase>b</sy:updateBase>
    <dc:date>2023-11-14T23:13:20+01:00</dc:date>
    <items><rdf:Seq><rdf:li resource="i1"/><rdf:li rdf:resource="i2"/></rdf:Seq></items>
  </channel>
  <image rdf:about="ia"><title>it</title><url>u</url><link>il</link></image>
  <item rdf:about="i1">
    <title>t</title><link>l</link><description>d</description>
    <dc:creator>c</dc:creator><dc:subject>s</dc:subject>
    <content:encoded>body</content:encoded>
    <dc:date>2023-11-14T23:13:20+01:00</dc:date>
  </item>
  <textinput rdf:about="ta"><title>tt</title><description>td</description><name>tn</name><link>tl</link></textinput>
</rdf:RDF>`
	f, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	r := f.(*RDF)
	if r.Channel.ImageResource != "img" || r.Channel.TextinputResource != "ti" ||
		r.Channel.DCCreator != "dc" || r.Channel.SyUpdatePeriod != "hourly" ||
		r.Channel.SyUpdateFrequency != "1" || r.Channel.SyUpdateBase != "b" ||
		len(r.Channel.ItemResources) != 2 || r.Channel.DCDate == nil {
		t.Errorf("RDF channel parse incomplete: %+v", r.Channel)
	}
	if r.Image.URL != "u" || r.Image.About != "ia" {
		t.Errorf("RDF image parse: %+v", r.Image)
	}
	it := r.Items[0]
	if it.DCCreator != "c" || it.DCSubject != "s" || it.ContentEncoded != "body" || it.DCDate == nil {
		t.Errorf("RDF item parse: %+v", it)
	}
	if r.Textinput.Name != "tn" || r.Textinput.Link != "tl" {
		t.Errorf("RDF textinput parse: %+v", r.Textinput)
	}
}

func TestRdfParseErrors(t *testing.T) {
	cases := []string{
		`<rdf:RDF xmlns:rdf="x"><channel><dc:date xmlns:dc="x">bad</dc:date></channel></rdf:RDF>`,
		`<rdf:RDF xmlns:rdf="x"><item rdf:about="i"><dc:date xmlns:dc="x">bad</dc:date></item></rdf:RDF>`,
	}
	for i, src := range cases {
		if _, err := Parse(src); err == nil {
			t.Errorf("case %d: expected error", i)
		}
	}
}

func TestRdfEmptySeq(t *testing.T) {
	// An <items> without a <rdf:Seq> yields no resources (parseRDFSeq nil path).
	src := `<rdf:RDF xmlns:rdf="x"><channel rdf:about="a"><title>t</title><link>l</link><description>d</description><items></items></channel></rdf:RDF>`
	f, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.(*RDF).Channel.ItemResources) != 0 {
		t.Error("expected no item resources")
	}
}
