package main

import (
	"bytes"
	"testing"
)

// decodeSent decodes what the client's last encodeHuffman produced, the way
// the server reads it: decodedChatLength letters from the encoded bytes.
func decodeSent(c *client) []byte {
	out := make([]byte, c.decodedChatLength)
	decodeHuffman(c.encodedChatBuffer[:c.encodedChatLength], out, 0, 0, c.decodedChatLength)
	return out
}

// A letter that takes two bytes in UTF-8 is sent as its one Latin-1 byte,
// with nothing after it (#21).
func TestEncodeHuffmanSendsLatin1Letters(t *testing.T) {
	c := &client{}
	msg := "Ångel says hi"
	c.encodeHuffman(msg)

	want := append([]byte{0xC5}, []byte("ngel says hi")...)
	if c.decodedChatLength != len(want) {
		t.Fatalf("length: want %d letters, got %d", len(want), c.decodedChatLength)
	}
	if got := decodeSent(c); !bytes.Equal(got, want) {
		t.Fatalf("decoded: want %q, got %q", want, got)
	}
	if got := chatLength(msg); got != len(want) {
		t.Fatalf("chatLength: want %d, got %d", len(want), got)
	}
}

// A short message after a longer one carries nothing left over from it (#21).
func TestEncodeHuffmanAfterLongerMessageCarriesNothingStale(t *testing.T) {
	c := &client{}
	c.encodeHuffman("a much longer message from before")
	c.encodeHuffman("Åh")

	fresh := &client{}
	fresh.encodeHuffman("Åh")

	if got, want := decodeSent(c), []byte{0xC5, 'h'}; !bytes.Equal(got, want) {
		t.Fatalf("decoded: want %q, got %q", want, got)
	}
	if !bytes.Equal(c.encodedChatBuffer[:c.encodedChatLength], fresh.encodedChatBuffer[:fresh.encodedChatLength]) {
		t.Fatalf("encoded bytes differ from a fresh client's")
	}
}
