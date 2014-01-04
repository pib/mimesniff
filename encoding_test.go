package mimesniff

import (
	"code.google.com/p/go.text/encoding"
	"testing"
)

func TestDetermineEncoding(t *testing.T) {
	head := `

<?xml version="1.0" encoding="UTF-8"?>

<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
	"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">

<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">

	<head>

		<meta http-equiv="content-type" content="text/html; charset=utf-8">

		<title>
			
			⚃ Rolling Dice with JS and Unicode ⚂ | Probably Programming
			
		</title>`

	mimeType := ParseMimeType("text/html")
	if enc := DetermineEncoding(mimeType, []byte(head)); enc != encoding.Nop {
		t.Error("Expected", encoding.Nop, "got", enc)
	}
}
