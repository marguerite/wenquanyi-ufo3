package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/marguerite/go-stdlib/dir"
	"github.com/marguerite/wenq/fixes/utils"
)

func main() {
	cwd, _ := os.Getwd()
	directories, _ := dir.Glob(filepath.Dir(filepath.Dir(cwd)) + "/WenQuanYiZenHei*.ufo3")

	for _, v := range directories {
		files, _ := dir.Ls(filepath.Join(v, "glyphs"), true, true)

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
			var n int
			var text []string
			for scanner.Scan() {
				text = append(text, scanner.Text())
				if found {
					if n == 0 && strings.Contains(scanner.Text(), "</contour>") {
						fmt.Printf("fixing %s\n", utils.Codepoint(f))
						text = text[:len(text)-2]
						text = append(text, scanner.Text())
					}
					n++
				}
				if strings.Contains(scanner.Text(), "<point x=\"4\"") {
					found = true
				}
			}
			glyph.Close()

			ioutil.WriteFile(f, []byte(strings.Join(text, "\n")+"\n"), 0644)
		}
	}
}
