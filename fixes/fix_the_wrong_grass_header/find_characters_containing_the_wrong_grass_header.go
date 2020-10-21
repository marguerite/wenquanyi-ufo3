package main

import (
	"bufio"
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fmt"

	"github.com/golang/freetype"
	"github.com/marguerite/go-stdlib/dir"
	"github.com/marguerite/wenq/fixes/utils"
	"golang.org/x/image/font"
)

const (
	fontFile = "/usr/share/fonts/truetype/wqy-zenhei.ttc"
)

func main() {
	cwd, _ := os.Getwd()
	directories, _ := dir.Glob(filepath.Dir(filepath.Dir(cwd)) + "/WenQuanYiZenHei*.ufo3")

	//for _, v := range directories {
	files, _ := dir.Ls(filepath.Join(directories[0], "glyphs"), true, true)

	var chars string

	for _, f := range files {
		glyph, err := os.Open(f)
		if err != nil {
			glyph.Close()
			continue
		}
		b, err := ioutil.ReadAll(glyph)
		if err != nil {
			glyph.Close()
			continue
		}
		scanner := bufio.NewScanner(bytes.NewReader(b))
		var found bool
		for scanner.Scan() {
			if found {
				if strings.Contains(scanner.Text(), "</contour>") {
					chars += utils.Codepoint(f)
				}
				break
			}
			if strings.Contains(scanner.Text(), "<point x=\"4\"") {
				found = true
			}
		}
		glyph.Close()
	}

	imgfile, _ := os.Create(filepath.Join(cwd, "result.png"))
	defer imgfile.Close()

	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		panic(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
		return
	}

	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
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
	for _, s := range utils.SplitSubN(chars, 20) {
		fmt.Println(s)
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(12 * 1.5)
	}

	b := bufio.NewWriter(imgfile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	//}
}
