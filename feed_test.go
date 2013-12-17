package mimesniff

import "testing"

type skipTest struct {
	start  string
	prefix string
	suffix string
	end    string
	ret    bool
}

var skipTests = []skipTest{
	{"Hello there", "H", "o", " there", true},
	{"Hello there", "e", "e", "Hello there", false},
	{"<yodawg><ihearyoulike>", "<", ">", "<ihearyoulike>", true},
}

func TestMaybeSkip(t *testing.T) {
	for _, test := range skipTests {
		s := test.start
		ret := maybeSkip(&s, test.prefix, test.suffix)
		if ret != test.ret || s != test.end {
			t.Error("For", test, "expected", test.ret, test.end, "got", ret, s)
		}
	}
}

type feedTest struct {
	text     string
	realType string
}

var feedTests = []feedTest{
	{"Not at all html", "text/html"},
	{"<?php Not at all html, still, ?php>", "text/html"},
	{"<rss>", "application/rss+xml"},
	{"\n <!--not the best idea, but whatever --><rss>", "application/rss+xml"},
	{"<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n<rss version=\"2.0\">", "application/rss+xml"},
	{`<?xml version="1.0"?>
     <rdf:RDF 
       xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
       xmlns="http://purl.org/rss/1.0/"
     >`, "application/rss+xml"},
	{`<?xml version="1.0" encoding="utf-8"?>
      <feed xmlns="http://www.w3.org/2005/Atom">`, "application/atom+xml"},
	{`<?xml version="1.0" encoding="iso-8859-1"?>
      <rss version="2.0" xmlns:atom="http://purl.org/atom/ns#">`, "application/rss+xml"},
	{`<?xml version=”1.0” encoding=”utf-8”?>
      <feed xmlns=”http://www.w3.org/2005/Atom”
        xml:base=”http://example.org/”
        xml:lang=”en“>`, "application/atom+xml"},
}

func TestDistinguishFeed(t *testing.T) {
	for _, test := range feedTests {
		realType := distinguishFeed("text/html", test.text)

		if realType != test.realType {
			t.Error("For", test, "expected", test.realType, "got", realType)
		}
	}
}
