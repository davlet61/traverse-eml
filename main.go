package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/mail"
	"os"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

var charsetDecoders = map[string]encoding.Encoding{
	"windows-1252": charmap.Windows1252,
	"windows-1257": charmap.Windows1257,
	"iso-8859-1":   charmap.ISO8859_1,
	"iso-8859-2":   charmap.ISO8859_2,
	"iso-8859-3":   charmap.ISO8859_3,
	"iso-8859-4":   charmap.ISO8859_4,
	"iso-8859-5":   charmap.ISO8859_5,
	"iso-8859-6":   charmap.ISO8859_6,
	"iso-8859-7":   charmap.ISO8859_7,
	"iso-8859-8":   charmap.ISO8859_8,
	"iso-8859-9":   charmap.ISO8859_9,
	"iso-8859-10":  charmap.ISO8859_10,
	"iso-8859-13":  charmap.ISO8859_13,
	"iso-8859-14":  charmap.ISO8859_14,
	"iso-8859-15":  charmap.ISO8859_15,
	"iso-8859-16":  charmap.ISO8859_16,
	"koi8-r":       charmap.KOI8R,
	"koi8-u":       charmap.KOI8U,
}

func main() {
	dir, err := os.ReadDir("./emails")
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		msg, err := os.ReadFile("emails/" + entry.Name())
		if err != nil {
			log.Fatal(err)
		}

		r := strings.NewReader(string(msg))
		m, err := mail.ReadMessage(r)
		if err != nil {
			log.Fatal(err)
		}

		header := m.Header
		from := header.Get("From")
		decodedFrom, err := decodeMIMEHeader(from)
		if err != nil {
			log.Fatal(err)
		}

		emailAddr, err := parseEmailAddress(decodedFrom)
		if err != nil {
			log.Fatal(err)
		}

		senderDomain := strings.Split(emailAddr, "@")[1]

		if senderDomain == "fotball.no" {
			fmt.Printf("Email from %s found in %s\n", emailAddr, entry.Name())
		}
	}
}

func decodeMIMEHeader(header string) (string, error) {
	decoder := new(mime.WordDecoder)
	decoder.CharsetReader = charsetReader

	decoded, err := decoder.DecodeHeader(header)
	if err != nil {
		return "", err
	}
	return decoded, nil
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	decoder, found := charsetDecoders[strings.ToLower(charset)]
	if found {
		return decoder.NewDecoder().Reader(input), nil
	}
	return nil, fmt.Errorf("unhandled charset %q", charset)
}

func parseEmailAddress(address string) (string, error) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "", err
	}
	return addr.Address, nil
}
