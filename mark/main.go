// Command mark writes the go-geometry mark, svg/color/go-geometry.svg:
//
//	cd mark && go run . > ../svg/color/go-geometry.svg
//
// go-geometry is computational geometry: robust predicates, convex hulls,
// Delaunay triangulations, Voronoi diagrams, k-d trees and meshes. The mark is
// the object all of them start from: a point set, its convex hull, and the
// triangulation that fans out of an interior point.
//
// Nothing is placed by hand. The program lists the hull vertices as an angle
// and a radius around the centre of the square, adds one interior point, joins
// the hull in order and the interior point to every hull vertex, and puts a
// dot on each point. Move a point in the table and the mark redraws itself.
//
// The rest is the family grammar every go-* mark shares: a 256-unit rounded
// square (rx 56) filled with a diagonal two-stop gradient, a white line glyph
// at stroke-width 13 with round caps and joins. logos/brandgen renders the
// format matrix from the file this writes; the kit workflow runs this again
// and fails if a byte of svg/color differs.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	light = "#5CCBE0" // the gradient's light stop
	dark  = "#0B6B87" // its dark stop

	stroke = 13.0 // the family's line width
	dotR   = 11.0 // radius of a point
)

// point is a position on the 256-unit square.
type point struct{ x, y float64 }

// polar is a hull vertex: its angle in degrees from straight up, clockwise,
// and its distance from the centre of the square.
type polar struct{ deg, r float64 }

func (p polar) at() point {
	a := p.deg * math.Pi / 180
	return point{128 + p.r*math.Sin(a), 128 - p.r*math.Cos(a)}
}

// hull lists the vertices in order around the hull. The angles are uneven on
// purpose: a regular pentagon reads as a symbol, an uneven one as a point set.
var hull = []polar{
	{0, 84},
	{70, 80},
	{140, 86},
	{216, 78},
	{292, 82},
}

// interior is the point the triangulation fans out of.
var interior = point{134, 134}

// num prints a coordinate with at most one decimal and no trailing zero.
func num(v float64) string {
	s := fmt.Sprintf("%.1f", math.Round(v*10)/10)
	s = strings.TrimSuffix(s, ".0")
	if s == "-0" {
		s = "0"
	}
	return s
}

func main() {
	var vs []point
	for _, h := range hull {
		vs = append(vs, h.at())
	}

	var lines []string

	// The hull, closed.
	var d []string
	for i, v := range vs {
		cmd := "L"
		if i == 0 {
			cmd = "M"
		}
		d = append(d, cmd+num(v.x)+" "+num(v.y))
	}
	lines = append(lines, fmt.Sprintf(`<path d="%sZ"/>`, strings.Join(d, "")))

	// The triangulation: the interior point joined to every hull vertex.
	for _, v := range vs {
		lines = append(lines, fmt.Sprintf(`<path d="M%s %sL%s %s"/>`,
			num(interior.x), num(interior.y), num(v.x), num(v.y)))
	}

	// The points: the hull vertices and the interior point.
	var dots []string
	for _, v := range append(vs, interior) {
		dots = append(dots, fmt.Sprintf(`<circle cx="%s" cy="%s" r="%s"/>`, num(v.x), num(v.y), num(dotR)))
	}

	fmt.Fprintf(os.Stdout, `<svg viewBox="0 0 256 256" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="go-geometry">
  <defs><linearGradient id="gmgrad" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="%s"/><stop offset="1" stop-color="%s"/></linearGradient></defs>
  <rect width="256" height="256" rx="56" fill="url(#gmgrad)"/>
  <g fill="none" stroke="#FFFFFF" stroke-width="%s" stroke-linecap="round" stroke-linejoin="round">
    %s
  </g>
  <g fill="#FFFFFF">
    %s
  </g>
</svg>
`, light, dark, num(stroke), strings.Join(lines, "\n    "), strings.Join(dots, "\n    "))
}
