package main

import (
	"os"
	"fmt"
	"log"
	"math"
	"strconv"

	"e4f"
	"xmp"
)

// Convert a float to a GpsCoord from XMP. dir is either 'N' or 'E'.
// Sign will change it.
func floatToGpsCoord(f float64, dir byte) string {

	negative := math.Signbit(f)
	if negative {
		switch dir {
		case 'N':
			dir = 'S'
		case 'E':
			dir = 'W'
		}
	}

	f = math.Abs(f)
	degs := math.Floor(f)

	minutes := (f - degs) * 60
	return fmt.Sprintf("%d,%f%c", int(degs), minutes, dir)
}

func exposureToXmp(db *e4f.E4fDb, roll *e4f.ExposedRoll, exp *e4f.Exposure, index int) {

//			fmt.Printf("Exposure %d: ", exp.Number)
	fmt.Println(exp)

	x := xmp.NewEmpty()
	defer xmp.Free(x)

	xmp.SetProperty(x, xmp.NS_EXIF_AUX, "ImageNumber",
		fmt.Sprintf("%d", index + 1), 0)
	if exp.Desc != "" {
		xmp.SetProperty(x, xmp.NS_DC, "description", exp.Desc, 0)
	}
	// ISO
	xmp.AppendArrayItem(x, xmp.NS_EXIF, "ISOSpeedRatings",
		xmp.PROP_VALUE_IS_ARRAY, fmt.Sprintf("%d", roll.Iso), 0)
	// Shutter speed
	xmp.SetProperty(x, xmp.NS_EXIF, "ShutterSpeedValue",
		exp.ShutterSpeed, 0)
	// Aperture
	f, err := strconv.ParseFloat(exp.Aperture, 64)
	if err == nil {
		xmp.SetProperty(x, xmp.NS_EXIF, "FNumber",
			fmt.Sprintf("%d/10", int(f * 10)), 0)
	}

	// FocalLength
	if exp.FocalLength != 0 {
		xmp.SetProperty(x, xmp.NS_EXIF, "FocalLength",
			fmt.Sprintf("%d", exp.FocalLength), 0)
	}
	// Lens
	lens, found := db.LensMap[exp.LensId]
	if lens != nil && found {
		// in Exif MaxApertureValue is the widest aperture,
		// ie the lowest number. Unlike in e4f
		f, err := strconv.ParseFloat(lens.ApertureMin, 64)
		if err == nil {
			xmp.SetProperty(x, xmp.NS_EXIF, "MaxApertureValue",
				fmt.Sprintf("%d/10", int(f * 10)), 0)
		}
		maker := db.MakeMap[lens.MakeId]
		xmp.SetProperty(x, xmp.NS_EXIF_AUX, "Lens",
			fmt.Sprintf("%s %s", maker.Name, lens.Title), 0)
	}

	// Flash
	flash := "false"
	if exp.FlashOn {
		flash = "true"
	}
	xmp.SetProperty(x, xmp.NS_EXIF, "Flash/exif:Fired", flash, 0)
	// Metering
	var meteringMode = 0
	switch exp.MeteringMode {
	case "Unknown":
		meteringMode = 0
	case "Average":
		meteringMode = 1
		// TODO finish
	}
	xmp.SetProperty(x, xmp.NS_EXIF, "MeteringMode",
		fmt.Sprintf("%d", meteringMode), 0)

	// Light source
	var lightSource = 0
	switch exp.LightSource {
	case "Daylight":
		lightSource = 1
	}
	xmp.SetProperty(x, xmp.NS_EXIF, "LightSource",
		fmt.Sprintf("%d", lightSource), 0)
	// Gps
	gps, found := db.GpsMap[exp.GpsLocId]
	if gps != nil && found {
		// create a fraction. Assume 1/10th of meter precision
		alt := gps.Alt * 10
		xmp.SetProperty(x, xmp.NS_EXIF, "GPSAltitude",
			fmt.Sprintf("%d/10", int(alt)), 0)

		coord := floatToGpsCoord(gps.Lat, 'N');
		xmp.SetProperty(x, xmp.NS_EXIF, "GPSLatitude", coord, 0)

		coord = floatToGpsCoord(gps.Long, 'E');
		xmp.SetProperty(x, xmp.NS_EXIF, "GPSLongitude", coord, 0)
	}

	buffer := xmp.StringNew()
	defer xmp.StringFree(buffer)

	xmp.Serialize(x, buffer, xmp.SERIAL_OMITPACKETWRAPPER, 0)

	fmt.Println(xmp.StringGo(buffer))
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Not enough arg")
	}

	path := os.Args[1]
	e4fDb := e4f.Parse(path)

	for _, roll := range e4fDb.ExposedRolls {
		id := roll.Id
		fmt.Println("Roll:")
		e4fDb.Print(roll)
		exps := e4fDb.ExposuresForRoll(id)
		fmt.Println(exps)
		for i, exp := range exps {
			exposureToXmp(e4fDb, roll, exp, i)
		}
	}
}
