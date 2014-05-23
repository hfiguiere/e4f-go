package main

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/lucsky/go-exml.v3"
)

type E4fDb struct {
	Version      string
	Cameras      []*Camera
	Makes        []*Make
	GpsLocations []*GpsLocation
	ExposedRolls []*ExposedRoll
	Exposures    []*Exposure
	Films        []*Film
	Lenses       []*Lens
	Artists      []*Artist
}

type Camera struct {
	Id                int
	DefaultFrameCount int
	MakeId            int
	SerialNumber      string
	DefaultFilmType   string
	Title             string
}

type Make struct {
	Id   int
	Name string
}

type GpsLocation struct {
	Id int
	Long, Lat, Alt float64
}

type ExposedRoll struct {
	Id           int
	FilmType     string
	CameraId     int
	Iso          int
        FrameCount   int
	TimeUnloaded string
	TimeLoaded   string
	FilmId       int
}

type Exposure struct {
}

type Film struct {
}

type Lens struct {
}

type Artist struct {
	Name string
}

func toInt(dst *int) exml.TextCallback {
	return func(c exml.CharData) {
		n, err := strconv.ParseInt(string(c), 0, 32)
		if err == nil {
			*dst = int(n)
		}
	}
}

func toFloat(dst *float64) exml.TextCallback {
	return func(c exml.CharData) {
		f, err := strconv.ParseFloat(string(c), 64)
		if err == nil {
			*dst = float64(f)
		}
	}
}

func parse(file string) (*E4fDb) {

	reader, _ := os.Open(file)
	defer reader.Close()

	e4fDb := &E4fDb{}
	decoder := exml.NewDecoder(reader)

	decoder.On("Exif4Film", func(attrs exml.Attrs) {

		e4fDb.Version, _= attrs.Get("version")

		decoder.On("Camera/dk.codeunited.exif4film.model.Camera",
			func(attrs exml.Attrs) {
				camera := &Camera{}
				e4fDb.Cameras = append(e4fDb.Cameras, camera)
				decoder.OnTextOf("camera_default_frame_count",
					toInt(&camera.DefaultFrameCount))
				decoder.OnTextOf("id", toInt(&camera.Id))
				decoder.OnTextOf("camera_make_id",
					toInt(&camera.MakeId))
				decoder.OnTextOf("camera_serial_number",
					exml.Assign(&camera.SerialNumber))
				decoder.OnTextOf("camera_default_film_type",
					exml.Assign(&camera.DefaultFilmType))
				decoder.OnTextOf("camera_title",
					exml.Assign(&camera.Title))
			})
		decoder.On("Make/dk.codeunited.exif4film.model.Make",
			func(attrs exml.Attrs) {
				make := &Make{}
				e4fDb.Makes = append(e4fDb.Makes, make)
				decoder.OnTextOf("id", toInt(&make.Id))
				decoder.OnTextOf("make_name",
					exml.Assign(&make.Name))
			})
		decoder.On(
			"GpsLocation/dk.codeunited.exif4film.model.GpsLocation",
			func(attrs exml.Attrs) {
				gps := &GpsLocation{}
				e4fDb.GpsLocations = append(e4fDb.GpsLocations,
					gps)
				decoder.OnTextOf("id", toInt(&gps.Id))
				decoder.OnTextOf("gps_latitude",
					toFloat(&gps.Lat))
				decoder.OnTextOf("gps_longitude",
					toFloat(&gps.Long))
				decoder.OnTextOf("gps_altitude",
					toFloat(&gps.Alt))
			})
		decoder.On(
			"ExposedRoll/dk.codeunited.exif4film.model.ExposedRoll",
			func(attrs exml.Attrs) {
				roll := &ExposedRoll{}
				e4fDb.ExposedRolls = append(e4fDb.ExposedRolls,
					roll)
				decoder.OnTextOf("id", toInt(&roll.Id))
				decoder.OnTextOf("exposedroll_film_type",
					exml.Assign(&roll.FilmType))
				decoder.OnTextOf("exposedroll_camera_id",
					toInt(&roll.CameraId))
				decoder.OnTextOf("exposedroll_film_id",
					toInt(&roll.FilmId))
				decoder.OnTextOf("exposedroll_iso",
					toInt(&roll.Iso))
				decoder.OnTextOf("exposedroll_frame_count",
					toInt(&roll.FrameCount))
				decoder.OnTextOf("exposedroll_time_unloaded",
					exml.Assign(&roll.TimeUnloaded))
				decoder.OnTextOf("exposedroll_time_loaded",
					exml.Assign(&roll.TimeLoaded))
			})
		decoder.On("Artist/dk.codeunited.exif4film.model.Artist",
			func(attrs exml.Attrs) {
				artist := &Artist{}
				e4fDb.Artists = append(e4fDb.Artists, artist)
				decoder.OnTextOf("artist_name",
					exml.Assign(&artist.Name))
			})
	})
	decoder.Run()

	return e4fDb
}


func main() {
	e4fDb := parse("samples/export-Roll-20130630_203650.xml");

	fmt.Printf("E4f: %s\n", e4fDb.Version)
}
