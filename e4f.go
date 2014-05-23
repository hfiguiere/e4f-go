package main

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/lucsky/go-exml.v3"
)

type E4fObject interface {
	GetId() int
	Type() string
}


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

	ObjectsMap   map[int] E4fObject
}
// Build the id -> data maps for the various elements
func (db *E4fDb) buildMap() {
	db.ObjectsMap = make(map[int] E4fObject)

	for _, cam := range db.Cameras {
		_, present := db.ObjectsMap[cam.Id]
		if present {
			fmt.Printf("Object %d present already\n", cam.Id)
		}
		db.ObjectsMap[cam.Id] = cam
	}

	for _, mk := range db.Makes {
		_, present := db.ObjectsMap[mk.Id]
		if present {
			fmt.Printf("Object %d present already\n", mk.Id)
		}
		db.ObjectsMap[mk.Id] = mk
	}

	for _, gps := range db.GpsLocations {
		_, present := db.ObjectsMap[gps.Id]
		if present {
			fmt.Printf("Object %d present already\n", gps.Id)
		}
		db.ObjectsMap[gps.Id] = gps
	}

	for _, roll := range db.ExposedRolls {
		_, present := db.ObjectsMap[roll.Id]
		if present {
			fmt.Printf("Object %d present already\n", roll.Id)
		}
		db.ObjectsMap[roll.Id] = roll
	}

	for _, exp := range db.Exposures {
		_, present := db.ObjectsMap[exp.Id]
		if present {
			fmt.Printf("Object %d present already\n", exp.Id)
		}
		db.ObjectsMap[exp.Id] = exp
	}

	for _, film := range db.Films {
		_, present := db.ObjectsMap[film.Id]
		if present {
			fmt.Printf("Object %d present already\n", film.Id)
		}
		db.ObjectsMap[film.Id] = film
	}

	for _, lens := range db.Lenses {
		_, present := db.ObjectsMap[lens.Id]
		if present {
			fmt.Printf("Object %d present already\n", lens.Id)
		}
		db.ObjectsMap[lens.Id] = lens
	}
}

type Camera struct {
	Id                int
	DefaultFrameCount int
	MakeId            int
	SerialNumber      string
	DefaultFilmType   string
	Title             string
}

func (o *Camera) GetId() int {
	return o.Id
}
func (o *Camera) Type() string {
	return "Camera"
}

type Make struct {
	Id   int
	Name string
}

func (o *Make) GetId() int {
	return o.Id
}
func (o *Make) Type() string {
	return "Make"
}

type GpsLocation struct {
	Id int
	Long, Lat, Alt float64
}
func (o *GpsLocation) GetId() int {
	return o.Id
}
func (o *GpsLocation) Type() string {
	return "GpsLocation"
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
func (o *ExposedRoll) GetId() int {
	return o.Id
}
func (o *ExposedRoll) Type() string {
	return "ExposedRoll"
}

type Exposure struct {
	Id           int
}
func (o *Exposure) GetId() int {
	return o.Id
}
func (o *Exposure) Type() string {
	return "Exposure"
}

type Film struct {
	Id             int
	Process        string
	Title          string
	ColorType      string
	Iso            int
	MakeId         int
}
func (o *Film) GetId() int {
	return o.Id
}
func (o *Film) Type() string {
	return "Film"
}

type Lens struct {
	Id             int
	Title          string
	SerialNumber   string
	MakeId         int
	ApertureMin    string
	ApertureMax    string
	FocalLengthMin int
	FocalLengthMax int
}
func (o *Lens) GetId() int {
	return o.Id
}
func (o *Lens) Type() string {
	return "Lens"
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
				m := &Make{}
				e4fDb.Makes = append(e4fDb.Makes, m)
				decoder.OnTextOf("id", toInt(&m.Id))
				decoder.OnTextOf("make_name",
					exml.Assign(&m.Name))
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
		decoder.On("Exposure/dk.codeunited.exif4film.model.Exposure",
			func(attrs exml.Attrs) {
				exp := &Exposure{}
				e4fDb.Exposures = append(e4fDb.Exposures, exp)
				decoder.OnTextOf("id", toInt(&exp.Id))
			})
		decoder.On("Film/dk.codeunited.exif4film.model.Film",
			func(attrs exml.Attrs) {
				film := &Film{}
				e4fDb.Films = append(e4fDb.Films, film)
				decoder.OnTextOf("id", toInt(&film.Id))
				decoder.OnTextOf("film_title",
					exml.Assign(&film.Title))
				decoder.OnTextOf("film_make_process",
					exml.Assign(&film.Process))
				decoder.OnTextOf("film_color_type",
					exml.Assign(&film.ColorType))
				decoder.OnTextOf("film_iso", toInt(&film.Iso))
				decoder.OnTextOf("film_make_id",
					toInt(&film.MakeId))
			})
		decoder.On("Lens/dk.codeunited.exif4film.model.Lens",
			func(attrs exml.Attrs) {
				lens := &Lens{}
				e4fDb.Lenses = append(e4fDb.Lenses, lens)
				decoder.OnTextOf("id", toInt(&lens.Id))
				decoder.OnTextOf("lens_title",
					exml.Assign(&lens.Title))
				decoder.OnTextOf("lens_serial_number",
					exml.Assign(&lens.SerialNumber))
				decoder.OnTextOf("lens_make_id",
					toInt(&lens.MakeId))
				decoder.OnTextOf("lens_aperture_min",
					exml.Assign(&lens.ApertureMin))
				decoder.OnTextOf("lens_aperture_max",
					exml.Assign(&lens.ApertureMax))
				decoder.OnTextOf("lens_focal_length_min",
					toInt(&lens.FocalLengthMin));
				decoder.OnTextOf("lens_focal_length_max",
					toInt(&lens.FocalLengthMax));
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

	e4fDb.buildMap()

	fmt.Printf("E4f: %s\n", e4fDb.Version)
	for i := 0; i < len(e4fDb.Cameras); i++ {
		fmt.Println(e4fDb.Cameras[i])
	}
	for i := 0; i < len(e4fDb.Makes); i++ {
		fmt.Println(e4fDb.Makes[i])
	}
	fmt.Println(e4fDb.ObjectsMap)
}
