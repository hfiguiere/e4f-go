// Parse Exif4film xml.
// And output the result
//
// See LICENSE

package e4f

import (
	"os"
	"strconv"
	"fmt"

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

	RollMap    map[int]*ExposedRoll
	MakeMap    map[int]*Make
	CameraMap  map[int]*Camera
	GpsMap     map[int]*GpsLocation
	LensMap    map[int]*Lens
	FilmMap    map[int]*Film
}

// Build the id -> data maps for the various elements
func (db *E4fDb) buildMaps() {
	db.CameraMap = make(map[int]*Camera)
	for _, cam := range db.Cameras {
		db.CameraMap[cam.Id] = cam
	}

	db.MakeMap = make(map[int]*Make)
	for _, mk := range db.Makes {
		db.MakeMap[mk.Id] = mk
	}

	db.GpsMap = make(map[int]*GpsLocation)
	for _, gps := range db.GpsLocations {
		db.GpsMap[gps.Id] = gps
	}

	db.RollMap = make(map[int]*ExposedRoll)
	for _, roll := range db.ExposedRolls {
		db.RollMap[roll.Id] = roll
	}

	db.FilmMap = make(map[int]*Film)
	for _, film := range db.Films {
		db.FilmMap[film.Id] = film
	}

	db.LensMap = make(map[int]*Lens)
	for _, lens := range db.Lenses {
		db.LensMap[lens.Id] = lens
	}
}

func (db *E4fDb) ExposuresForRoll(id int) (exposures []*Exposure) {
	for _, exp := range db.Exposures {
		if exp.RollId == id {
			exposures = append(exposures, exp)
		}
	}
	return
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
	Id             int
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
	Desc         string
}

type Exposure struct {
	Id           int
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

type Film struct {
	Id        int
	Process   string
	Title     string
	ColorType string
	Iso       int
	MakeId    int
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

func Parse(file string) *E4fDb {

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
				decoder.OnTextOf("exposedroll_description",
					exml.Assign(&roll.Desc))
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

	e4fDb.buildMaps()

	return e4fDb
}


func (db *E4fDb) Print(roll *ExposedRoll) {
	fmt.Printf("%s\n", roll.Desc)
	var filmLabel string
	film, found := db.FilmMap[roll.FilmId]
	if film != nil && found {
		mk := db.MakeMap[film.MakeId]
		var filmMake string
		if mk != nil && mk.Name != "" {
			filmMake = mk.Name
		}
		if film.Title != "" {
			if filmMake != "" {
				filmLabel = fmt.Sprintf("%s %s", mk.Name,
					film.Title)
			} else {
				filmLabel = film.Title
			}
		}
	}

	fmt.Printf("Type %s, %s, %d ISO\n\n", roll.FilmType, filmLabel, roll.Iso)
}

