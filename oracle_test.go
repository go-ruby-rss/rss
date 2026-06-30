// Copyright (c) the go-ruby-rss/rss authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rss

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

// rubyBin locates a `ruby` that can `require 'rss'`. The oracle tests skip
// themselves when it is absent or the rss gem is unavailable (the qemu
// cross-arch lanes and the Windows lane), so the deterministic suite alone
// drives the 100% gate there.
func rubyBin(t *testing.T) string {
	t.Helper()
	path, err := exec.LookPath("ruby")
	if err != nil {
		t.Skip("ruby not on PATH; skipping MRI oracle")
	}
	if err := exec.Command(path, "-rrss", "-e", "").Run(); err != nil {
		t.Skip("ruby cannot require 'rss'; skipping MRI oracle")
	}
	return path
}

// mriRoundTrip parses src with MRI's RSS::Parser (validation off, to accept the
// minimal fixtures) and returns its to_s. $stdout.binmode keeps Windows text
// mode from polluting the bytes (the go-ruby-erb lesson). do_validation=false
// matches what our parser does — it does not enforce MRI's strict schema.
func mriRoundTrip(t *testing.T, bin, src string) string {
	t.Helper()
	script := "$stdout.binmode\n" +
		"require 'rss'\n" +
		"feed = RSS::Parser.parse($stdin.read, false)\n" +
		"print feed.to_s\n"
	cmd := exec.Command(bin, "-e", script)
	cmd.Stdin = strings.NewReader(src)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ruby error: %v\noutput:\n%s", err, out)
	}
	return string(out)
}

// TestOracleRoundTrip is the headline cross-check: a real feed in each dialect,
// parsed and reserialized by MRI and by this library, must be byte-for-byte
// identical. It covers RSS 2.0, RSS 1.0 (RDF) and Atom, including the RFC822
// and W3CDTF date formats and the bundled dc/content/syndication modules.
func TestOracleRoundTrip(t *testing.T) {
	bin := rubyBin(t)
	cases := []struct {
		name string
		src  string
	}{
		{"rss2", oracleRss2Src},
		{"rdf", oracleRdfSrc},
		{"atom", oracleAtomSrc},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mri := mriRoundTrip(t, bin, c.src)
			f, err := Parse(c.src)
			if err != nil {
				t.Fatalf("our Parse failed: %v", err)
			}
			got := f.String()
			if got != mri {
				t.Errorf("round-trip mismatch with MRI:\n--- ours ---\n%s\n--- MRI ---\n%s", got, mri)
			}
		})
	}
}

// TestOracleDateFormats checks our date formatters against MRI's Time#rfc822
// and Time#w3cdtf directly.
func TestOracleDateFormats(t *testing.T) {
	bin := rubyBin(t)
	script := "$stdout.binmode\nrequire 'time'\nrequire 'rss'\n" +
		"t = Time.at(1700000000).getlocal('+01:00')\n" +
		"u = Time.at(1700000000).utc\n" +
		"puts t.rfc822\nputs u.w3cdtf\nputs t.w3cdtf\n"
	cmd := exec.Command(bin, "-e", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ruby error: %v\n%s", err, out)
	}
	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %q", out)
	}
	// 1700000000 == 2023-11-14T22:13:20Z == 23:13:20+01:00.
	tLoc := ptime(2023, 11, 14, 23, 13, 20, loc1)
	tUTC := ptime(2023, 11, 14, 22, 13, 20, time.UTC)
	if got := formatRFC822(*tLoc); got != lines[0] {
		t.Errorf("rfc822 = %q, MRI = %q", got, lines[0])
	}
	if got := formatW3CDTF(*tUTC); got != lines[1] {
		t.Errorf("w3cdtf(utc) = %q, MRI = %q", got, lines[1])
	}
	if got := formatW3CDTF(*tLoc); got != lines[2] {
		t.Errorf("w3cdtf(+01) = %q, MRI = %q", got, lines[2])
	}
}

const oracleRss2Src = `<?xml version="1.0"?>
<rss version="2.0">
  <channel>
    <title>News</title>
    <link>http://news.example/</link>
    <description>Daily news &amp; views</description>
    <language>en-us</language>
    <pubDate>Mon, 06 Sep 2021 10:00:00 GMT</pubDate>
    <item>
      <title>Story &lt;1&gt;</title>
      <link>http://news.example/1</link>
      <description>Body</description>
      <pubDate>Sun, 05 Sep 2021 09:30:00 +0200</pubDate>
      <guid isPermaLink="false">tag:news,1</guid>
    </item>
  </channel>
</rss>`

const oracleRdfSrc = `<?xml version="1.0" encoding="UTF-8"?>
<rdf:RDF xmlns="http://purl.org/rss/1.0/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:dc="http://purl.org/dc/elements/1.1/">
  <channel rdf:about="http://example.com/">
    <title>RDF News</title>
    <link>http://example.com/</link>
    <description>desc</description>
    <items>
      <rdf:Seq>
        <rdf:li resource="http://example.com/1"/>
      </rdf:Seq>
    </items>
    <dc:date>2023-11-14T23:13:20+01:00</dc:date>
  </channel>
  <item rdf:about="http://example.com/1">
    <title>One</title>
    <link>http://example.com/1</link>
    <description>body</description>
    <dc:creator>Jane</dc:creator>
  </item>
</rdf:RDF>`

const oracleAtomSrc = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <id>urn:uuid:feed</id>
  <title>Atom News</title>
  <updated>2023-11-14T22:13:20Z</updated>
  <author>
    <name>Jane</name>
  </author>
  <link href="http://example.com/" rel="alternate"/>
  <entry>
    <id>urn:uuid:e1</id>
    <title>Entry</title>
    <updated>2023-11-14T22:13:20Z</updated>
    <content>Hello</content>
  </entry>
</feed>`
