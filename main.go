package shp_to_geojson

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

const (
	shpHeaderSize       = 100 // size of the shp file header in bytes
	shpRecordHeaderSize = 8   // size of the shp record header in bytes
)

func printObject(input interface{}) {
	content, _ := json.Marshal(input)
	fmt.Println(string(content))
}

// ReadFileHeader uses the reader to read bytes into de AppleIcon structure
func ReadFileHeader(r *bytes.Reader) *FileHeader {
	var header FileHeader
	binary.Read(r, binary.BigEndian, &header.FileCode)
	binary.Read(r, binary.BigEndian, &header.Unused)
	binary.Read(r, binary.BigEndian, &header.FileLength)
	binary.Read(r, binary.LittleEndian, &header.Version)
	binary.Read(r, binary.LittleEndian, &header.ShapeType)
	binary.Read(r, binary.LittleEndian, &header.BoundingBoxInfo)
	return &header
}

func wordSizeToBytes(wordCount int) int {
	return wordCount * 16 / 8
}

func convertObjToGeo(recordsArray []interface{}) (string, error) {
	Xmax, Ymax, Zmax := math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64
	Xmin, Ymin, Zmin := math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
	var first = recordsArray[0]

	switch first.(type) {
	case Point:
		return "", errors.New("Point not supported")
	case PointM:
		return "", errors.New("PointM not supported")
	case PointZ:
		return "", errors.New("PointZ not supported")
	case MultiPoint:
		return "", errors.New("MultiPoint not supported")
	case MultiPointM:
		return "", errors.New("MultiPointM not supported")
	case MultiPointZ:
		var Geom = GeoJSON_MultiPoint{
			Type:        "MultiPoint",
			Coordinates: make([]GeoJSON_base_point, 0),
		}

		feature := Feature{
			Type:       "Feature",
			Properties: map[string]interface{}{},
			Geometry:   &Geom,
		}
		for _, point := range recordsArray {
			var p = point.(MultiPointZ)
			if p.Box[0] < Xmin {
				Xmin = p.Box[0]
			}
			if p.Box[1] < Ymin {
				Ymin = p.Box[1]
			}
			if p.Zmin < Zmin {
				Zmin = p.Zmin
			}
			if p.Box[2] > Xmax {
				Xmax = p.Box[2]
			}
			if p.Box[3] > Ymax {
				Ymax = p.Box[3]
			}
			if p.Zmax > Zmax {
				Zmax = p.Zmax
			}
			for index, c := range p.Points {
				var geoPoint = GeoJSON_base_point{c.X, c.Y, p.Zarray[index]}
				Geom.Coordinates = append(Geom.Coordinates, geoPoint)
			}
		}
		feature.Properties["Xmax"] = Xmax
		feature.Properties["Ymax"] = Ymax
		feature.Properties["Zmax"] = Zmax
		feature.Properties["Xmin"] = Xmin
		feature.Properties["Ymin"] = Ymin
		feature.Properties["Zmin"] = Zmin
		content, _ := json.Marshal(feature)
		return string(content), nil
	case PolyLine:
		return "", errors.New("PolyLine not supported")
	case PolyLineM:
		return "", errors.New("PolyLineM not supported")
	case PolyLineZ:
		Geom := GeoJSON_MultiLineString{
			Type:        "MultiLineString",
			Coordinates: make([][]GeoJSON_base_point, 0),
		}
		feature := Feature{
			Type:       "Feature",
			Properties: map[string]interface{}{},
			Geometry:   &Geom,
		}
		for _, line := range recordsArray {
			var l = line.(PolyLineZ)
			if l.Box[0] < Xmin {
				Xmin = l.Box[0]
			}
			if l.Box[1] < Ymin {
				Ymin = l.Box[1]
			}
			if l.Zmin < Zmin {
				Zmin = l.Zmin
			}
			if l.Box[2] > Xmax {
				Xmax = l.Box[2]
			}
			if l.Box[3] > Ymax {
				Ymax = l.Box[3]
			}
			if l.Zmax > Zmax {
				Zmax = l.Zmax
			}
			geo_line := GeoJSON_LineStrings{
				Type:        "LineString",
				Coordinates: make([]GeoJSON_base_point, 0),
			}
			for index, c := range l.Points {
				var geo_point = GeoJSON_base_point{c.X, c.Y, l.Zarray[index]}
				geo_line.Coordinates = append(geo_line.Coordinates, geo_point)
			}
			Geom.Coordinates = append(Geom.Coordinates, geo_line.Coordinates)
		}
		feature.Properties["Xmax"] = Xmax
		feature.Properties["Ymax"] = Ymax
		feature.Properties["Zmax"] = Zmax
		feature.Properties["Xmin"] = Xmin
		feature.Properties["Ymin"] = Ymin
		feature.Properties["Zmin"] = Zmin
		content, jsonErr := json.Marshal(feature)
		if jsonErr != nil {
			return "", jsonErr
		}
		return string(content), nil
	case Polygon:
		return "", errors.New("Polygon not supported")
	case PolygonM:
		return "", errors.New("PolygonM not supported")
	case PolygonZ:
		Geom := GeoJSON_MultiPolygon{
			Type:        "MultiPolygon",
			Coordinates: make([][][]GeoJSON_base_point, 0),
		}

		feature := Feature{
			Type:       "Feature",
			Properties: map[string]interface{}{},
			Geometry:   &Geom,
		}
		for _, polygon := range recordsArray {
			var p = polygon.(PolygonZ)
			Geo_poly := GeoJSON_Polygon{
				Type:        "Polygon",
				Coordinates: make([][]GeoJSON_base_point, 0),
			}
			if p.Box[0] < Xmin {
				Xmin = p.Box[0]
			}
			if p.Box[1] < Ymin {
				Ymin = p.Box[1]
			}
			if p.Zmin < Zmin {
				Zmin = p.Zmin
			}
			if p.Box[2] > Xmax {
				Xmax = p.Box[2]
			}
			if p.Box[3] > Ymax {
				Ymax = p.Box[3]
			}
			if p.Zmax > Zmax {
				Zmax = p.Zmax
			}
			// each polygon has one long points array which we need to split up into it's parts depending on the parts array which points to the starting index of each part
			for partCounter := 0; partCounter < int(p.NumParts); partCounter++ {
				partsArray := make([]GeoJSON_base_point, 0)

				// if it's the last part the just run from the index to the end
				if int(p.NumParts)-1 == partCounter {
					partStartIndex := int(p.Parts[partCounter])
					for pointIndex := partStartIndex; pointIndex < int(p.NumPoints); pointIndex++ {
						point := p.Points[pointIndex]
						geo_point := GeoJSON_base_point{point.X, point.Y, p.Zarray[pointIndex]}
						partsArray = append(partsArray, geo_point)
					}
				}
				Geo_poly.Coordinates = append(Geo_poly.Coordinates, partsArray)
			}
			Geom.Coordinates = append(Geom.Coordinates, Geo_poly.Coordinates)
		}
		feature.Properties["Xmax"] = Xmax
		feature.Properties["Ymax"] = Ymax
		feature.Properties["Zmax"] = Zmax
		feature.Properties["Xmin"] = Xmin
		feature.Properties["Ymin"] = Ymin
		feature.Properties["Zmin"] = Zmin

		content, _ := json.Marshal(feature)
		return string(content), nil
	default:
		return "", errors.New("unknown type")
	}
}

func ParseFromBytes(data []byte) (string, error) {
	if len(data) < shpHeaderSize {
		return "", errors.New("invalid shp file: file header too small what")
	}

	reader := bytes.NewReader(data)
	header := ReadFileHeader(reader)

	// Verify the file code and version
	if header.FileCode != 9994 || header.Version != 1000 {
		return "", errors.New("invalid shp file: unsupported version or file code")
	}

	var fileZise = wordSizeToBytes(int(header.FileLength))
	var bytesRead int = 100
	var shapeType = int(header.ShapeType)

	// check if shape type is in the const list of shapetypes otherwise return string message with shape type and error
	if shapeType != ShapeTypePoint && shapeType != ShapeTypePolyLine && shapeType != ShapeTypePolygon && shapeType != ShapeTypeMultiPoint && shapeType != ShapeTypePointZ && shapeType != ShapeTypePolyLineZ && shapeType != ShapeTypePolygonZ && shapeType != ShapeTypeMultiPointZ {
		return "", errors.New(fmt.Sprint("Shape type %d not supported", shapeType))
	}

	content := make([]interface{}, 0)

	for bytesRead < fileZise {
		var recordHeader RecordHeader
		binary.Read(bytes.NewReader(data[bytesRead:bytesRead+shpRecordHeaderSize]), binary.BigEndian, &recordHeader)
		bytesRead += 8
		var recordContentSize int = wordSizeToBytes(int(recordHeader.ContentLength))
		recordReader := bytes.NewReader(data[bytesRead : bytesRead+recordContentSize])
		// allways skip the first 4 bytes because they are the shape type and we already know that from the header
		recordReader.Seek(4, 0)

		if shapeType == ShapeTypeMultiPointZ {
			var point MultiPointZ
			binary.Read(recordReader, binary.LittleEndian, &point.Box)
			binary.Read(recordReader, binary.LittleEndian, &point.NumPoints)
			point.Points = make([]Point, point.NumPoints)
			binary.Read(recordReader, binary.LittleEndian, point.Points)
			binary.Read(recordReader, binary.LittleEndian, &point.Zmin)
			binary.Read(recordReader, binary.LittleEndian, &point.Zmax)
			point.Zarray = make([]float64, point.NumPoints)
			binary.Read(recordReader, binary.LittleEndian, point.Zarray)
			binary.Read(recordReader, binary.LittleEndian, &point.Mmin)
			binary.Read(recordReader, binary.LittleEndian, &point.Mmax)
			point.Marray = make([]float64, point.NumPoints)
			binary.Read(recordReader, binary.LittleEndian, point.Marray)
			content = append(content, point)
		} else if shapeType == ShapeTypePolygonZ {
			var polygon PolygonZ
			binary.Read(recordReader, binary.LittleEndian, &polygon.Box)
			binary.Read(recordReader, binary.LittleEndian, &polygon.NumParts)
			binary.Read(recordReader, binary.LittleEndian, &polygon.NumPoints)
			polygon.Parts = make([]uint32, polygon.NumParts)
			binary.Read(recordReader, binary.LittleEndian, &polygon.Parts)
			polygon.Points = make([]Point, polygon.NumPoints)
			binary.Read(recordReader, binary.LittleEndian, polygon.Points)
			binary.Read(recordReader, binary.LittleEndian, &polygon.Zmin)
			binary.Read(recordReader, binary.LittleEndian, &polygon.Zmax)
			polygon.Zarray = make([]float64, polygon.NumPoints)
			binary.Read(recordReader, binary.LittleEndian, &polygon.Zarray)
			binary.Read(recordReader, binary.LittleEndian, &polygon.Mmin)
			binary.Read(recordReader, binary.LittleEndian, &polygon.Mmax)
			polygon.Marray = make([]float64, polygon.NumPoints)
			binary.Read(recordReader, binary.LittleEndian, &polygon.Marray)
			content = append(content, polygon)
		} else if shapeType == ShapeTypePolyLineZ {
			var polyline PolyLineZ
			binary.Read(recordReader, binary.LittleEndian, &polyline.Box)
			binary.Read(recordReader, binary.LittleEndian, &polyline.NumParts)
			binary.Read(recordReader, binary.LittleEndian, &polyline.NumPoints)
			polyline.Parts = make([]uint32, polyline.NumParts)
			binary.Read(recordReader, binary.LittleEndian, &polyline.Parts)
			polyline.Points = make([]Point, polyline.NumPoints)
			binary.Read(recordReader, binary.LittleEndian, polyline.Points)
			binary.Read(recordReader, binary.LittleEndian, &polyline.Zmin)
			binary.Read(recordReader, binary.LittleEndian, &polyline.Zmax)
			polyline.Zarray = make([]float64, polyline.NumPoints)
			binary.Read(recordReader, binary.LittleEndian, &polyline.Zarray)
			binary.Read(recordReader, binary.LittleEndian, &polyline.Mmin)
			binary.Read(recordReader, binary.LittleEndian, &polyline.Mmax)
			polyline.Marray = make([]float64, polyline.NumPoints)
			binary.Read(recordReader, binary.LittleEndian, &polyline.Marray)
			content = append(content, polyline)
		}
		bytesRead += recordContentSize
	}

	if len(content) > 0 {
		json_string, convert_err := convertObjToGeo(content)
		if convert_err != nil {
			return "", convert_err
		} else {
			return json_string, nil
		}
	}
	return "", errors.New("no content found in shp file")
}
