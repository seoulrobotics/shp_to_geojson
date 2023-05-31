package main

import (
	"fmt"
	"io/ioutil"

	shp "github.com/julianhahn/shp_to_geojson"
)

func main() {
	// Read file to byte array
	file, readERr := ioutil.ReadFile("./B3_SURFACEMARK.shp")
	if readERr != nil {
		fmt.Println(readERr)
		return
	}
	// Parse
	shapefile, parseErr := shp.ParseFromBytes(file)
	if parseErr != nil {
		fmt.Println(parseErr)
		return
	}
	// Print
	fmt.Printf("%+v\n", shapefile)
}
