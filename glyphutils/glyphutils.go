package glyphutils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
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

// CodepointFromGlifFileName produce the unicode codepoint from .glif filename
func CodepointFromGlifFileName(f string) string {
	codepoint := strings.TrimPrefix(filepath.Base(f), "uni")
	codepoint = strings.TrimSuffix(codepoint, ".glif")
	codepoint = strings.ReplaceAll(codepoint, "_", "")
	s, _ := u2s(strings.ToLower(codepoint))
	return s
}

// SplitStringByLength split string by length
func SplitStringByLength(s string, n int) []string {
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

// GenImageWithFont generate `text` rendered by font `file` to `img`
func GenImageWithFont(file, img, text string) {
	if len(img) == 0 {
		cwd, _ := os.Getwd()
		img = filepath.Join(cwd, text+".png")
	}

	f, err := os.Create(img)
	if err != nil {
		f.Close()
		panic(err)
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		f.Close()
		panic(err)
	}
	ft, err := freetype.ParseFont(b)
	if err != nil {
		f.Close()
		panic(err)
	}

	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(ft)
	c.SetFontSize(12)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	c.SetHinting(font.HintingFull)
	// Draw the guidelines.
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}
	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(12)>>6))
	for _, s := range SplitStringByLength(text, 20) {
		fmt.Println(s)
		_, err = c.DrawString(s, pt)
		if err != nil {
			f.Close()
			panic(err)
		}
		pt.Y += c.PointToFixed(12 * 1.5)
	}

	wt := bufio.NewWriter(f)
	err = png.Encode(wt, rgba)
	if err != nil {
		f.Close()
		panic(err)
	}
	err = wt.Flush()
	if err != nil {
		f.Close()
		panic(err)
	}
	f.Close()
}
