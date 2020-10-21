package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"path/filepath"
	"strings"
)

func u2s(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, "", -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

// Codepoint produce the unicode codepoint from .glif filename
func Codepoint(f string) string {
	codepoint := strings.TrimPrefix(filepath.Base(f), "uni")
	codepoint = strings.TrimSuffix(codepoint, ".glif")
	codepoint = strings.ReplaceAll(codepoint, "_", "")
	s, _ := u2s(strings.ToLower(codepoint))
	return s
}

// SplitSubN split string by length
func SplitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}
