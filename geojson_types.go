package shp_to_geojson

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Geometry   interface{}            `json:"geometry"`
}

type GeoJSON_base_point [3]float64

type GeoJSON_Point struct {
	Type        string               `json:"type"`
	Coordinates []GeoJSON_base_point `json:"coordinates"`
}

type GeoJSON_LineStrings struct {
	Type        string               `json:"type"`
	Coordinates []GeoJSON_base_point `json:"coordinates"`
}

type GeoJSON_Polygon struct {
	Type        string                 `json:"type"`
	Properties  map[string]interface{} `json:"properties"`
	Coordinates [][]GeoJSON_base_point `json:"coordinates"`
}

type GeoJSON_MultiPoint struct {
	Type        string               `json:"type"`
	Coordinates []GeoJSON_base_point `json:"coordinates"`
}

type GeoJSON_MultiLineString struct {
	Type        string                 `json:"type"`
	Coordinates [][]GeoJSON_base_point `json:"coordinates"`
}

type GeoJSON_MultiPolygon struct {
	Type        string                   `json:"type"`
	Coordinates [][][]GeoJSON_base_point `json:"coordinates"`
}
