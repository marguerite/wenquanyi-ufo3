package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/marguerite/go-stdlib/dir"
	"github.com/marguerite/wenq/glyphutils"
	"github.com/marguerite/wenq/ufo3"
)

func main() {
	cwd, _ := os.Getwd()
	directories, _ := dir.Glob(filepath.Dir(filepath.Dir(cwd)) + "/WenQuanYiZenHei*.ufo3")

	//for _, v := range directories {
	files, _ := dir.Ls(filepath.Join(directories[0], "glyphs"), true, true)

	var chars string

	for _, f := range files {
		if !strings.HasSuffix(f, ".glif") {
			continue
		}
		glyph := ufo3.NewGlyphFromFile(f)

		for _, v := range glyph.Outline.Contours {
			i, p := v.FindPointByX("49")
			if p.IsNil() {
				continue
			}
			j, p1 := v.FindPointByX("243")
			k, p2 := v.FindPointByX("253", "qcurve", "yes")
			if p1.IsNil() || p2.IsNil() {
				continue
			}
			if i-j == 1 && k-i == 1 {
				chars += glyphutils.CodepointFromGlifFileName(f)
			}
		}
	}

	glyphutils.GenImageWithFont("/usr/share/fonts/truetype/wqy-zenhei.ttc", filepath.Join(cwd, "result.png"), chars)

	//}
}
