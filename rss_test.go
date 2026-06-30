// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import (
	"strings"
	"testing"
	"time"
)

// loc1 is the +01:00 fixed zone used across the golden vectors.
var loc1 = time.FixedZone("", 3600)

func ptime(y, mo, d, h, mi, s int, loc *time.Location) *time.Time {
	t := time.Date(y, time.Month(mo), d, h, mi, s, 0, loc)
	return &t
}

// ---- RSS 2.0 generation (golden, MRI-verified) ---------------------------

const goldenRss2 = `<?xml version="1.0"?>
<rss version="2.0"
  xmlns:content="http://purl.org/rss/1.0/modules/content/"
  xmlns:dc="http://purl.org/dc/elements/1.1/"
  xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd"
  xmlns:trackback="http://madskills.com/public/xml/rss/module/trackback/">
  <channel>
    <title>Example Feed</title>
    <link>http://example.com/</link>
    <description>An example.</description>
    <language>en-us</language>
    <pubDate>Tue, 14 Nov 2023 23:13:20 +0100</pubDate>
    <item>
      <title>Item One</title>
      <link>http://example.com/1</link>
      <description>First item.</description>
      <pubDate>Tue, 14 Nov 2023 23:13:20 +0100</pubDate>
      <guid isPermaLink="true">http://example.com/1</guid>
    </item>
  </channel>
</rss>`

func sampleRss2() *Rss {
	return &Rss{
		Version: "2.0",
		Channel: &Channel{
			Title:       "Example Feed",
			Link:        "http://example.com/",
			Description: "An example.",
			Language:    "en-us",
			PubDate:     ptime(2023, 11, 14, 23, 13, 20, loc1),
			Items: []*Item{{
				Title:       "Item One",
				Link:        "http://example.com/1",
				Description: "First item.",
				PubDate:     ptime(2023, 11, 14, 23, 13, 20, loc1),
				Guid:        &Guid{Content: "http://example.com/1", IsPermaLink: true, HasPermaLink: true},
			}},
		},
	}
}

func TestRss2Generate(t *testing.T) {
	if got := sampleRss2().String(); got != goldenRss2 {
		t.Errorf("RSS 2.0 generate mismatch:\n got:\n%s\nwant:\n%s", got, goldenRss2)
	}
	if ft := sampleRss2().FeedType(); ft != "rss" {
		t.Errorf("FeedType = %q, want rss", ft)
	}
}

func TestRss2DefaultVersion(t *testing.T) {
	r := &Rss{Channel: &Channel{Title: "t", Link: "l", Description: "d"}}
	if !strings.Contains(r.String(), `version="2.0"`) {
		t.Error("empty Version should default to 2.0")
	}
}

func TestRss2Full(t *testing.T) {
	// Exercises every optional channel/item element and the image, category,
	// syndication, dc and content branches.
	r := &Rss{
		Version: "2.0",
		Channel: &Channel{
			Title: "t", Link: "l", Description: "d",
			Language: "en", Copyright: "(c)", ManagingEditor: "me",
			WebMaster: "wm", Generator: "g", Docs: "docs", TTL: "60",
			PubDate:        ptime(2023, 1, 1, 0, 0, 0, time.UTC),
			LastBuildDate:  ptime(2023, 1, 2, 0, 0, 0, loc1),
			Categories:     []string{"a", "b"},
			Image:          &Image{URL: "u", Title: "it", Link: "il", Width: "88", Height: "31"},
			DCDate:         ptime(2023, 1, 1, 0, 0, 0, loc1),
			SyUpdatePeriod: "hourly", SyUpdateFrequency: "1", SyUpdateBase: "2000-01-01T00:00+00:00",
			Items: []*Item{{
				Title: "it", Link: "il", Description: "id",
				Author: "au", Comments: "co",
				Categories:     []string{"c1"},
				PubDate:        ptime(2023, 1, 1, 0, 0, 0, loc1),
				Guid:           &Guid{Content: "g", IsPermaLink: false, HasPermaLink: true},
				ContentEncoded: "<b>hi</b>", DCCreator: "dc", DCSubject: "ds",
				DCDate: ptime(2023, 1, 1, 0, 0, 0, loc1),
			}},
		},
	}
	out := r.String()
	for _, want := range []string{
		"<copyright>(c)</copyright>", "<managingEditor>me</managingEditor>",
		"<webMaster>wm</webMaster>", "<category>a</category>", "<category>b</category>",
		"<image>", "<sy:updatePeriod>hourly</sy:updatePeriod>",
		"<content:encoded>&lt;b&gt;hi&lt;/b&gt;</content:encoded>",
		"<dc:creator>dc</dc:creator>", "<dc:subject>ds</dc:subject>",
		`<guid isPermaLink="false">g</guid>`,
		"<lastBuildDate>", "<dc:date>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("full RSS2 missing %q in:\n%s", want, out)
		}
	}
}

func TestRss2GuidNoPermaLink(t *testing.T) {
	r := &Rss{Channel: &Channel{Title: "t", Link: "l", Description: "d",
		Items: []*Item{{Title: "i", Guid: &Guid{Content: "g"}}}}}
	if !strings.Contains(r.String(), "<guid>g</guid>") {
		t.Error("guid without HasPermaLink should omit the attribute")
	}
}

// ---- RSS 2.0 parsing -----------------------------------------------------

func TestRss2RoundTrip(t *testing.T) {
	// A parsed feed reserializes with the encoding prolog.
	src := goldenRss2
	f, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	r, ok := f.(*Rss)
	if !ok {
		t.Fatalf("Parse returned %T, want *Rss", f)
	}
	if r.Channel.Title != "Example Feed" || r.Channel.Items[0].Guid.Content != "http://example.com/1" {
		t.Errorf("parsed accessors wrong: %+v", r.Channel)
	}
	if !r.Channel.Items[0].Guid.IsPermaLink {
		t.Error("isPermaLink should be true")
	}
	out := r.String()
	if !strings.HasPrefix(out, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Errorf("parsed feed should carry encoding prolog, got:\n%s", out)
	}
}

func TestRss2ParseAllFields(t *testing.T) {
	src := `<?xml version="1.0"?>
<rss version="2.0">
  <channel>
    <title>t</title><link>l</link><description>d</description>
    <language>en</language><copyright>c</copyright>
    <managingEditor>me</managingEditor><webMaster>wm</webMaster>
    <generator>g</generator><docs>do</docs><ttl>60</ttl>
    <category>cat</category>
    <pubDate>Mon, 06 Sep 2021 10:00:00 GMT</pubDate>
    <lastBuildDate>Tue, 07 Sep 2021 10:00:00 +0200</lastBuildDate>
    <dc:date xmlns:dc="http://purl.org/dc/elements/1.1/">2023-11-14T23:13:20+01:00</dc:date>
    <sy:updatePeriod xmlns:sy="x">hourly</sy:updatePeriod>
    <sy:updateFrequency xmlns:sy="x">1</sy:updateFrequency>
    <sy:updateBase xmlns:sy="x">b</sy:updateBase>
    <image><url>u</url><title>it</title><link>il</link><width>88</width><height>31</height></image>
    <item>
      <title>it</title><link>il</link><description>id</description>
      <author>au</author><comments>co</comments><category>c1</category>
      <pubDate>Sun, 05 Sep 2021 09:30:00 +0200</pubDate>
      <guid isPermaLink="false">g</guid>
      <content:encoded xmlns:content="x">body</content:encoded>
      <dc:creator xmlns:dc="x">dc</dc:creator>
      <dc:subject xmlns:dc="x">ds</dc:subject>
      <dc:date xmlns:dc="x">2023-11-14T23:13:20+01:00</dc:date>
    </item>
  </channel>
</rss>`
	f, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	r := f.(*Rss)
	ch := r.Channel
	if ch.Copyright != "c" || ch.ManagingEditor != "me" || ch.WebMaster != "wm" ||
		ch.Generator != "g" || ch.Docs != "do" || ch.TTL != "60" ||
		len(ch.Categories) != 1 || ch.Image.Width != "88" || ch.DCDate == nil ||
		ch.SyUpdatePeriod != "hourly" || ch.SyUpdateFrequency != "1" || ch.SyUpdateBase != "b" {
		t.Errorf("channel fields not all parsed: %+v", ch)
	}
	it := ch.Items[0]
	if it.Author != "au" || it.Comments != "co" || len(it.Categories) != 1 ||
		it.ContentEncoded != "body" || it.DCCreator != "dc" || it.DCSubject != "ds" ||
		it.DCDate == nil || it.PubDate == nil || it.Guid.IsPermaLink {
		t.Errorf("item fields not all parsed: %+v", it)
	}
	if y := ch.PubDate.Year(); y != 2021 {
		t.Errorf("pubDate year = %d", y)
	}
}

func TestRss2ParseErrors(t *testing.T) {
	cases := map[string]string{
		"malformed xml":      `<rss>`,
		"no channel":         `<rss version="2.0"></rss>`,
		"bad channel date":   `<rss version="2.0"><channel><pubDate>nonsense</pubDate></channel></rss>`,
		"bad lastbuild":      `<rss version="2.0"><channel><lastBuildDate>nope</lastBuildDate></channel></rss>`,
		"bad channel dcdate": `<rss version="2.0"><channel><dc:date xmlns:dc="x">nope</dc:date></channel></rss>`,
		"bad item pubdate":   `<rss version="2.0"><channel><item><pubDate>x</pubDate></item></channel></rss>`,
		"bad item dcdate":    `<rss version="2.0"><channel><item><dc:date xmlns:dc="x">x</dc:date></item></channel></rss>`,
	}
	for name, src := range cases {
		if _, err := Parse(src); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
}
