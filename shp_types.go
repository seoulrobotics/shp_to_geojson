package shp_to_geojson

type FileHeader struct {
	FileCode        uint32
	Unused          [20]byte
	FileLength      uint32
	Version         uint32
	ShapeType       uint32
	BoundingBoxInfo [8]float64
}

type RecordHeader struct {
	RecordNumber uint32
	/* contentLength measured in 16Bit words */
	ContentLength uint32
}

type Shape struct {
	ShapeType uint32
	Points    []Point
}

type Point struct {
	/* both in little */
	X float64
	Y float64
}

type PointM struct {
	X float64
	Y float64
	M float64
}

type PointZ struct {
	X float64
	Y float64
	Z float64
	M float64
}

type MultiPoint struct {
	Box       [4]float64
	NumPoints uint32
	Points    []Point
}

type MultiPointM struct {
	Box       [4]float64
	NumPoints uint32
	Points    []Point
	Mmin      float64
	Mmax      float64
	Marray    []float64
}

type MultiPointZ struct {
	Box       [4]float64
	NumPoints uint32
	Points    []Point
	Zmin      float64
	Zmax      float64
	Zarray    []float64
	Mmin      float64
	Mmax      float64
	Marray    []float64
}

type PolyLine struct {
	Box       [4]float64
	NumParts  uint32
	NumPoints uint32
	Parts     []uint32
	Points    []Point
}

type PolyLineM struct {
	Box       [4]float64
	NumParts  uint32
	NumPoints uint32
	Parts     []uint32
	Points    []Point
	Mmin      float64
	Mmax      float64
	Marray    []float64
}

type PolyLineZ struct {
	Box       [4]float64
	NumParts  uint32
	NumPoints uint32
	Parts     []uint32
	Points    []Point
	Zmin      float64
	Zmax      float64
	Zarray    []float64
	Mmin      float64
	Mmax      float64
	Marray    []float64
}

// counter clockwise direction
type Polygon struct {
	Box       [4]float64
	NumParts  uint32   //number of rings in the polygon
	NumPoints uint32   //total number of points for all rings
	Parts     []uint32 //index to the first point in each ring
	Points    []Point
}

type PolygonM struct {
	Box       [4]float64
	NumParts  uint32
	NumPoints uint32
	Parts     []uint32
	Points    []Point
	Mmin      float64
	Mmax      float64
	Marray    []float64
}

type PolygonZ struct {
	Box       [4]float64
	NumParts  uint32
	NumPoints uint32
	Parts     []uint32
	Points    []Point
	Zmin      float64
	Zmax      float64
	Zarray    []float64
	Mmin      float64
	Mmax      float64
	Marray    []float64
}

type MultiPatch struct {
	Box       [4]float64
	NumParts  uint32
	NumPoints uint32
	Parts     []uint32
	PartTypes []uint32
	Points    []Point
	Zmin      float64
	Zmax      float64
	Zarray    []float64
	Mmin      float64
	Mmax      float64
	Marray    []float64
}

const (
	ShapeTypeNull        = 0
	ShapeTypePoint       = 1
	ShapeTypePolyLine    = 3
	ShapeTypePolygon     = 5
	ShapeTypeMultiPoint  = 8
	ShapeTypePointZ      = 11
	ShapeTypePolyLineZ   = 13
	ShapeTypePolygonZ    = 15
	ShapeTypeMultiPointZ = 18
	ShapeTypePointM      = 21
	ShapeTypePolyLineM   = 23
	ShapeTypePolygonM    = 25
	ShapeTypeMultiPointM = 28
	ShapeTypeMultiPatch  = 31
)
