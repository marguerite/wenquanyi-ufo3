package ufo3

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// Glyph represent an unicode character in a font
type Glyph struct {
	XMLName xml.Name `xml:"glyph"`
	Name    string   `xml:"name,attr"`
	Format  string   `xml:"format,attr"`
	Advance Advance  `xml:"advance"`
	Unicode Unicode  `xml:"unicode"`
	Outline Outline  `xml:"outline"`
}

// Advance ?
type Advance struct {
	XMLName xml.Name `xml:"advance"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
}

// Unicode the unicode string
type Unicode struct {
	XMLName xml.Name `xml:"unicode"`
	Hex     string   `xml:"hex,attr"`
}

// Outline the glyph outline
type Outline struct {
	XMLName    xml.Name    `xml:"outline"`
	Contours   []Contour   `xml:"contour"`
	Components []Component `xml:"component"`
}

type Component struct {
	Base    string `xml:"base,attr"`
	Xscale  string `xml:"xScale,attr"`
	Yscale  string `xml:"yScale,attr"`
	XYscale string `xml:"xyScale,attr"`
	YXscale string `xml:"yxScale,attr"`
	Xoffset string `xml:"xOffset,attr"`
	Yoffset string `xml:"yOffset,attr"`
}

// Contour the self-closed "shape" in a glyph
type Contour struct {
	XMLName xml.Name `xml:"contour"`
	Points  []Point  `xml:"point"`
}

// Point a (x,y) point in a shape
type Point struct {
	XMLName xml.Name `xml:"point"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Type    string   `xml:"type,attr"`
	Smooth  string   `xml:"smooth,attr"`
}

// NewGlyphFromFile open glyph in a .glif file
func NewGlyphFromFile(file string) (glyph Glyph) {
	f, err := os.Open(file)
	if err != nil {
		f.Close()
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		f.Close()
		panic(err)
	}
	xml.Unmarshal(b, &glyph)
	f.Close()
	return glyph
}

// Bytes turn the glyph to bytes again
func (g Glyph) Bytes() []byte {
	b, err := xml.MarshalIndent(g, "", "  ")
	if err != nil {
		panic(err)
	}
	// xml header
	text := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
	// replace <string></string> to <string/>
	scanner := bufio.NewScanner(bytes.NewReader(b))
	re := regexp.MustCompile(`^.*(><\/\w+>)$`)
	// ignore empty
	re1 := regexp.MustCompile(`\s\w+=""`)
	re2 := regexp.MustCompile(`<[^=]+/>`)
	for scanner.Scan() {
		line := scanner.Text()
		if re1.MatchString(line) {
			m := re1.FindAllStringSubmatch(line, -1)
			for _, v := range m {
				line = strings.ReplaceAll(line, v[0], "")
			}
		}
		if re.MatchString(line) {
			m := re.FindStringSubmatch(line)
			line = strings.Replace(line, m[1], "/>", 1)
		}
		if re2.MatchString(line) {
			continue
		}
		text += line + "\n"
	}
	return []byte(text)
}

// DeletePoint delete idxPoint Point in the idxContour shape from glyph
func (g *Glyph) DeletePoint(idxContour, idxPoint int) {
	cv := reflect.Indirect(reflect.ValueOf(g)).FieldByName("Outline").FieldByName("Contours").Index(idxContour).FieldByName("Points")
	s := reflect.MakeSlice(reflect.SliceOf(cv.Index(0).Type()), 0, 0)
	for i := 0; i < cv.Len(); i++ {
		if i == idxPoint {
			continue
		}
		s = reflect.Append(s, cv.Index(i))
	}
	cv.Set(s)
}

// FindPointByX find a point with by its X value
func (c Contour) FindPointByX(x string, options ...string) (idx int, point Point) {
	var typ, smooth string
	if len(options) > 0 {
		typ = options[0]
	}
	if len(options) > 1 {
		smooth = options[1]
	}

	for i, p := range c.Points {
		if p.X == x {
			if len(typ) > 0 {
				if p.Type == typ {
					if len(smooth) > 0 {
						if p.Smooth == smooth {
							return i, p
						}
					} else {
						return i, p
					}
				}
			} else {
				return i, p
			}
		}
	}
	return idx, point
}

// IsNil if this point is uninitialized
func (p Point) IsNil() bool {
	for _, v := range []string{p.XMLName.Space, p.XMLName.Local, p.X, p.Y, p.Type, p.Smooth} {
		if len(v) > 0 {
			return false
		}
	}
	return true
}
