package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	shp "github.com/julianhahn/shp_to_geojson"
)

func main() {

	var collection shp.FeatureCollection = shp.FeatureCollection{
		Type:     "FeatureCollection",
		Features: []shp.Feature{},
	}
	files, dirErr := ioutil.ReadDir("./")
	if dirErr != nil {
		fmt.Println(dirErr)
		return
	}
	for _, file := range files {
		if file.Name()[len(file.Name())-4:] == ".shp" {
			// Read file to byte array
			file, readERr := ioutil.ReadFile("./" + file.Name())
			if readERr != nil {
				fmt.Println(readERr)
				return
			}
			// Parse
			var feature shp.Feature

			shapefile, parseErr := shp.ParseFromBytes(file)
			if parseErr != nil {
				fmt.Println(parseErr)
				return
			}

			// unmarshal to add to collection
			jsonErr := json.Unmarshal([]byte(shapefile), &feature)
			if jsonErr != nil {
				fmt.Println(jsonErr)
				return
			}

			collection.Features = append(collection.Features, feature)
		}
	}
	// marshal back to json
	json, _ := json.Marshal(collection)
	fmt.Printf("%+v\n", string(json))
}
