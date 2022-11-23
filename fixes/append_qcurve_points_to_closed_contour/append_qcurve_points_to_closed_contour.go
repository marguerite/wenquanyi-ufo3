package main

import (
	"encoding/xml"
	"fmt"
	"github.com/marguerite/wenq/ufo3"
  "github.com/marguerite/go-stdlib/dir"
  "github.com/marguerite/wenq/glyphutils"
	"io/ioutil"
	"math"
	"strconv"
  "strings"
  "os"
  "path/filepath"
)

// findContoursWithoutType find closed Contours that contains all type-less points,
func findContoursWithoutType(glyph ufo3.Glyph) []int {
	var nums []int
	for i, contour := range glyph.Outline.Contours {
		var found bool
		for _, point := range contour.Points {
			if len(point.Type) != 0 {
				found = true
			}
		}
		if !found {
			nums = append(nums, i)
		}
	}
	return nums
}

func high(n, n1 float64) float64 {
	if n > n1 {
		return n
	}
	return n1
}

func low(n, n1 float64) float64 {
	if n > n1 {
		return n1
	}
	return n
}

func round(s, s1 string) string {
	n, _ := strconv.ParseFloat(s, 64)
	n1, _ := strconv.ParseFloat(s1, 64)
	l := low(n, n1)
	h := high(n, n1)
	return fmt.Sprintf("%d", int(math.Floor(l+(h-l)/2+0.5)))
}

// python fontTools/pens/pointPen.py will append a qcurve point without X and Y if a contour has no on-curve point
// we fix such contours by calculating and appending a qcurve point
func main() {
	cwd, _ := os.Getwd()
	directories, _ := dir.Glob(filepath.Dir(filepath.Dir(cwd)) + "/WenQuanYiZenHei*.ufo3")

	for _, v := range directories {
		files, _ := dir.Ls(filepath.Join(v, "glyphs"), true, true)

		for _, f := range files {
			if !strings.HasSuffix(f, ".glif") {
				continue
			}
			glyph := ufo3.NewGlyphFromFile(f)

			nums := findContoursWithoutType(glyph)

			for i := 0; i < len(nums); i++ {
        fmt.Printf("fixing %s\n", glyphutils.CodepointFromGlifFileName(f))
				contour := glyph.Outline.Contours[nums[i]]
				for j := len(contour.Points) - 1; j > 0; j-- {
					var prev ufo3.Point
					if j == len(contour.Points)-1 {
						if contour.Points[j].X == contour.Points[j-1].X || contour.Points[j].Y == contour.Points[j-1].Y {
							prev = contour.Points[j-1]
						} else {
							prev = contour.Points[0]
						}
						if contour.Points[j].X == prev.X {
							y := round(contour.Points[j].Y, prev.Y)
							p := ufo3.Point{xml.Name{}, prev.X, y, "qcurve", "yes"}
							glyph.AppendPoint(nums[i], j, p, true)
						}
						if contour.Points[j].Y == prev.Y {
							x := round(contour.Points[j].X, prev.X)
							p := ufo3.Point{xml.Name{}, x, prev.Y, "qcurve", "yes"}
							glyph.AppendPoint(nums[i], j, p, true)
						}
					} else {
						prev = contour.Points[j-1]
						if contour.Points[j].X == prev.X {
							y := round(contour.Points[j].Y, prev.Y)
							p := ufo3.Point{xml.Name{}, prev.X, y, "qcurve", "yes"}
							glyph.AppendPoint(nums[i], j, p, false)
						}
						if contour.Points[j].Y == prev.Y {
							x := round(contour.Points[j].X, prev.X)
							p := ufo3.Point{xml.Name{}, x, prev.Y, "qcurve", "yes"}
							glyph.AppendPoint(nums[i], j, p, false)
						}
					}
				}
			}

			ioutil.WriteFile(f, glyph.Bytes(), 0644)
		}
	}
}
