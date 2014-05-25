// Parse Exif4film xml.
// And output the result
//
// See LICENSE

package main

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/lucsky/go-exml.v3"
)

type E4fObject interface {
	Id() int
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

	GpsMap     map[int]*GpsLocation
	LensMap    map[int]*Lens
	ObjectsMap map[int]E4fObject
}

// Build the id -> data maps for the various elements
func (db *E4fDb) buildMaps() {
	db.ObjectsMap = make(map[int]E4fObject)

	for _, cam := range db.Cameras {
		o, present := db.ObjectsMap[cam.id]
		if present {
			fmt.Printf("Object %d present already of type %s\n",
				cam.id, o.Type())
		}
		db.ObjectsMap[cam.id] = cam
	}

	for _, mk := range db.Makes {
		o, present := db.ObjectsMap[mk.id]
		if present {
			fmt.Printf("Object %d present already of type %s\n",
				mk.id, o.Type())
		}
		db.ObjectsMap[mk.id] = mk
	}

	db.GpsMap = make(map[int]*GpsLocation)
	for _, gps := range db.GpsLocations {
		db.GpsMap[gps.id] = gps
	}

	for _, roll := range db.ExposedRolls {
		o, present := db.ObjectsMap[roll.id]
		if present {
			fmt.Printf("Object %d present already of type %s\n",
				roll.id, o.Type())
		}
		db.ObjectsMap[roll.id] = roll
	}

	for _, exp := range db.Exposures {
		o, present := db.ObjectsMap[exp.id]
		if present {
			fmt.Printf("Object %d present already of type %s\n",
				exp.id, o.Type())
		}
		db.ObjectsMap[exp.id] = exp
	}

	for _, film := range db.Films {
		o, present := db.ObjectsMap[film.id]
		if present {
			fmt.Printf("Object %d present already of type %s\n",
				film.id, o.Type())
		}
		db.ObjectsMap[film.id] = film
	}

	db.LensMap = make(map[int]*Lens)
	for _, lens := range db.Lenses {
		db.LensMap[lens.id] = lens
	}
}

func (db *E4fDb) exposuresForRoll(id int) (exposures []*Exposure) {
	obj := db.ObjectsMap[id]
	if obj.Type() != "ExposedRoll" {
		fmt.Printf("Found type %s\n", obj.Type())
		return nil
	}

	for _, exp := range db.Exposures {
		if exp.RollId == id {
			exposures = append(exposures, exp)
		}
	}
	return
}

type Camera struct {
	id                int
	DefaultFrameCount int
	MakeId            int
	SerialNumber      string
	DefaultFilmType   string
	Title             string
}

func (o *Camera) Id() int {
	return o.id
}
func (o *Camera) Type() string {
	return "Camera"
}

type Make struct {
	id   int
	Name string
}

func (o *Make) Id() int {
	return o.id
}
func (o *Make) Type() string {
	return "Make"
}

type GpsLocation struct {
	id             int
	Long, Lat, Alt float64
}

func (o *GpsLocation) Id() int {
	return o.id
}
func (o *GpsLocation) Type() string {
	return "GpsLocation"
}

type ExposedRoll struct {
	id           int
	FilmType     string
	CameraId     int
	Iso          int
	FrameCount   int
	TimeUnloaded string
	TimeLoaded   string
	FilmId       int
}

func (o *ExposedRoll) Id() int {
	return o.id
}
func (o *ExposedRoll) Type() string {
	return "ExposedRoll"
}

type Exposure struct {
	id           int
	FlashOn      bool
	Desc         string
	Number       int
	GpsLocId     int
	ExpComp      int
	RollId       int
	FocalLength  int
	LightSource  string
	TimeTaken    string
	ShutterSpeed string
	LensId       int
	Aperture     string
	MeteringMode string
}

func (o *Exposure) Id() int {
	return o.id
}
func (o *Exposure) Type() string {
	return "Exposure"
}

type Film struct {
	id        int
	Process   string
	Title     string
	ColorType string
	Iso       int
	MakeId    int
}

func (o *Film) Id() int {
	return o.id
}
func (o *Film) Type() string {
	return "Film"
}

type Lens struct {
	id             int
	Title          string
	SerialNumber   string
	MakeId         int
	ApertureMin    string
	ApertureMax    string
	FocalLengthMin int
	FocalLengthMax int
}

func (o *Lens) Id() int {
	return o.id
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

func toBool(dst *bool) exml.TextCallback {
	return func(c exml.CharData) {
		if string(c) == "true" {
			*dst = true
		} else {
			*dst = false
		}
	}
}

func parse(file string) *E4fDb {

	reader, _ := os.Open(file)
	defer reader.Close()

	e4fDb := &E4fDb{}
	decoder := exml.NewDecoder(reader)

	decoder.On("Exif4Film", func(attrs exml.Attrs) {

		e4fDb.Version, _ = attrs.Get("version")

		decoder.On("Camera/dk.codeunited.exif4film.model.Camera",
			func(attrs exml.Attrs) {
				camera := &Camera{}
				e4fDb.Cameras = append(e4fDb.Cameras, camera)
				decoder.OnTextOf("camera_default_frame_count",
					toInt(&camera.DefaultFrameCount))
				decoder.OnTextOf("id", toInt(&camera.id))
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
				decoder.OnTextOf("id", toInt(&m.id))
				decoder.OnTextOf("make_name",
					exml.Assign(&m.Name))
			})
		decoder.On(
			"GpsLocation/dk.codeunited.exif4film.model.GpsLocation",
			func(attrs exml.Attrs) {
				gps := &GpsLocation{}
				e4fDb.GpsLocations = append(e4fDb.GpsLocations,
					gps)
				decoder.OnTextOf("id", toInt(&gps.id))
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
				decoder.OnTextOf("id", toInt(&roll.id))
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
				decoder.OnTextOf("id", toInt(&exp.id))
				decoder.OnTextOf("exposure_flash_on",
					toBool(&exp.FlashOn))
				decoder.OnTextOf("exposure_description",
					exml.Assign(&exp.Desc))
				decoder.OnTextOf("exposure_number",
					toInt(&exp.Number))
				decoder.OnTextOf("exposure_gps_location",
					toInt(&exp.GpsLocId))
				decoder.OnTextOf("exposure_compensation",
					toInt(&exp.ExpComp))
				decoder.OnTextOf("exposure_roll_id",
					toInt(&exp.RollId))
				decoder.OnTextOf("exposure_focal_length",
					toInt(&exp.FocalLength))
				decoder.OnTextOf("exposure_light_source",
					exml.Assign(&exp.LightSource))
				decoder.OnTextOf("exposure_time_taken",
					exml.Assign(&exp.TimeTaken))
				decoder.OnTextOf("exposure_shutter_speed",
					exml.Assign(&exp.ShutterSpeed))
				decoder.OnTextOf("exposure_lens_id",
					toInt(&exp.LensId))
				decoder.OnTextOf("exposure_aperture",
					exml.Assign(&exp.Aperture))
				decoder.OnTextOf("exposure_metering_mode",
					exml.Assign(&exp.MeteringMode))
			})
		decoder.On("Film/dk.codeunited.exif4film.model.Film",
			func(attrs exml.Attrs) {
				film := &Film{}
				e4fDb.Films = append(e4fDb.Films, film)
				decoder.OnTextOf("id", toInt(&film.id))
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
				decoder.OnTextOf("id", toInt(&lens.id))
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
					toInt(&lens.FocalLengthMin))
				decoder.OnTextOf("lens_focal_length_max",
					toInt(&lens.FocalLengthMax))
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
	e4fDb := parse("samples/export-Roll-20130630_203650.xml")

	e4fDb.buildMaps()

	for _, roll := range e4fDb.ExposedRolls {
		id := roll.Id()
		fmt.Println("Roll:")
		fmt.Println(roll)
		exps := e4fDb.exposuresForRoll(id)
		for _, exp := range exps {
			fmt.Printf("Exposure %d: ", exp.Number)
			fmt.Println(exp)
		}
	}
}
