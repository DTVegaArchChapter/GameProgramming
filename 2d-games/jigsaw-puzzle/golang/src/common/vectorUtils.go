package common

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Point struct {
	X, Y float64
}

func FillPathWithColor(screen *ebiten.Image, path *vector.Path, x, y float32, c color.Color) {
	whiteImage := ebiten.NewImage(3, 3)
	whiteImage.Fill(c) // Fill with white color

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage := whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	var vertices []ebiten.Vertex
	var indices []uint16
	op := &ebiten.DrawTrianglesOptions{
		FillRule:  ebiten.FillRuleEvenOdd,
		AntiAlias: true,
	}

	vertices, indices = path.AppendVerticesAndIndicesForFilling(nil, nil)

	for iv := range vertices {
		v := &vertices[iv]
		v.DstX += x
		v.DstY += y
	}

	screen.DrawTriangles(vertices, indices, whiteSubImage, op)
}

func StrokePathWithColor(screen *ebiten.Image, path *vector.Path, x, y, strokeWidth float32, c color.Color) {
	whiteImage := ebiten.NewImage(3, 3)
	whiteImage.Fill(c) // Fill with white color

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage := whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	var vertices []ebiten.Vertex
	var indices []uint16
	op := &ebiten.DrawTrianglesOptions{
		FillRule:  ebiten.FillRuleEvenOdd,
		AntiAlias: true,
	}

	vertices, indices = path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{Width: strokeWidth})

	for iv := range vertices {
		v := &vertices[iv]
		v.DstX += x
		v.DstY += y
	}

	screen.DrawTriangles(vertices, indices, whiteSubImage, op)
}

func GetPathVertices(vPath *vector.Path) []Point {
	var vertices []ebiten.Vertex
	var indices []uint16

	// Triangulate the path (but we only need the vertices)
	vertices, _ = vPath.AppendVerticesAndIndicesForFilling(vertices, indices)

	// Extract unique vertices (approximation)
	result := make([]Point, 0, len(vertices))
	for _, v := range vertices {
		result = append(result, Point{X: float64(v.DstX), Y: float64(v.DstY)})
	}

	return result
}

func IsPointInPath(p Point, vPath *vector.Path) bool {
	if vPath == nil {
		return false // No path to check against
	}

	// Get the vertices of the path
	pathVertices := GetPathVertices(vPath)

	// Use a point-in-polygon algorithm to check if the point is inside the path
	return IsPointInPathVertices(p, pathVertices)
}

func IsPointInPathVertices(p Point, pathVertices []Point) bool {
	if len(pathVertices) < 3 {
		return false // Not a valid polygon
	}

	inside := false
	n := len(pathVertices)

	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		vi := pathVertices[i]
		vj := pathVertices[j]

		// Ray-casting algorithm
		if ((vi.Y > p.Y) != (vj.Y > p.Y)) &&
			(p.X < (vj.X-vi.X)*(p.Y-vi.Y)/(vj.Y-vi.Y+1e-10)+vi.X) {
			inside = !inside
		}
	}

	return inside
}
