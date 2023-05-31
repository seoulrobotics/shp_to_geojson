# shp to geojson

## Install

```bash
go get github.com/seoulrobotics/shp_to_geojson
```

```go
import (
   shp "github.com/seoulrobotics/shp_to_geojson"
  )
```

## Usage

There is only one method that expect's a byte array from a shpfile and the result will be a string in the format of geo json.

To save some time for users, I included two examples:

1. [parse one file](https://github.com/seoulrobotics/shp_to_geojson/tree/main/examples/read_parse_print)
2. [parse multiple files](https://github.com/seoulrobotics/shp_to_geojson/tree/main/examples/read_multiple)

### Background

In my quick research I found only two packages which handles shp packages. Neither of both could directly convert to geojson, or simply extract the coordinates. Therefore I created this package with the simple goal to convert a shp file to the more modern standard [geosjon](https://datatracker.ietf.org/doc/html/rfc7946#section-3.1.4).

For this I used the [technial description of shapefiles](https://www.esri.com/content/dam/esrisites/sitecore-archive/Files/Pdfs/library/whitepapers/pdfs/shapefile.pdf) and created for each shapetype a struct which then get's read directly from the binary of the file. The Flow is as follows:

1. read the first 100 bytes as fileheader and extract the shapetype and file length
2. loop until the byteCounter is bigger then the filezise
3. in every loop, read the record header and extract the record size
4. after we know the record size, we can extract the needed information from the record body and then increment the byteIndex by the size of the record size to jump to the next record
5. once the loop is finished we can be sure to extracted all data from the shp file.
   ( optional ) - handle a conversion of the point type
6. loop again over the data and fill it in a geojson struct
