package xmp

// #cgo pkg-config: exempi-2.0
// #include <exempi/xmp.h>
// #include <exempi/xmpconsts.h>
import "C"
import "unsafe"

var (
	NS_EXIF = (*C.char)(unsafe.Pointer(&C.NS_EXIF))
	NS_TIFF = (*C.char)(unsafe.Pointer(&C.NS_TIFF))
	NS_XAP = (*C.char)(unsafe.Pointer(&C.NS_XAP))
	NS_XAP_RIGHTS = (*C.char)(unsafe.Pointer(&C.NS_XAP_RIGHTS))
	NS_DC = (*C.char)(unsafe.Pointer(&C.NS_DC))
	NS_EXIF_AUX = (*C.char)(unsafe.Pointer(&C.NS_EXIF_AUX))
	NS_CRS = (*C.char)(unsafe.Pointer(&C.NS_CRS))
	NS_LIGHTROOM = (*C.char)(unsafe.Pointer(&C.NS_LIGHTROOM))
	NS_PHOTOSHOP = (*C.char)(unsafe.Pointer(&C.NS_PHOTOSHOP))
	NS_CAMERA_RAW_SETTINGS = (*C.char)(unsafe.Pointer(&C.NS_CAMERA_RAW_SETTINGS))
	NS_CAMERA_RAW_SAVED_SETTINGS = (*C.char)(unsafe.Pointer(&C.NS_CAMERA_RAW_SAVED_SETTINGS))
	NS_IPTC4XMP = (*C.char)(unsafe.Pointer(&C.NS_IPTC4XMP))
	NS_TPG = (*C.char)(unsafe.Pointer(&C.NS_TPG))
	NS_DIMENSIONS_TYPE = (*C.char)(unsafe.Pointer(&C.NS_DIMENSIONS_TYPE))
	NS_CC = (*C.char)(unsafe.Pointer(&C.NS_CC))
	NS_PDF = (*C.char)(unsafe.Pointer(&C.NS_PDF))
)


const (
	PROP_VALUE_IS_STRUCT = 0x00000100
	PROP_VALUE_IS_ARRAY  = 0x00000200
	PROP_ARRAY_IS_UNORDERED = PROP_VALUE_IS_ARRAY
	PROP_ARRAY_IS_ORDERED = 0x00000400
	PROP_ARRAY_IS_ALT    = 0x00000800
)

const (
	SERIAL_OMITPACKETWRAPPER   = 0x0010
	SERIAL_READONLYPACKET      = 0x0020
	SERIAL_USECOMPACTFORMAT    = 0x0040
	SERIAL_INCLUDETHUMBNAILPAD = 0x0100
	SERIAL_EXACTPACKETLENGTH   = 0x0200
	SERIAL_WRITEALIASCOMMENTS  = 0x0400
	SERIAL_OMITALLFORMATTING   = 0x0800
)

func NewEmpty() C.XmpPtr {
	return C.xmp_new_empty()
}

func Free(x C.XmpPtr) {
	C.xmp_free(x)
}

func GetError() int {
	return int(C.xmp_get_error())
}

func Serialize(x C.XmpPtr, buffer C.XmpStringPtr, options C.uint32_t,
	padding C.uint32_t) bool {

	ret := C.xmp_serialize(x, buffer, options, padding)

	return bool(ret)
}

func SetProperty(x C.XmpPtr, schema *C.char, name string, value string,
	optionBits C.uint32_t) bool {

	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	valueC := C.CString(value)
	defer C.free(unsafe.Pointer(valueC))

	ret := C.xmp_set_property(x, schema, nameC, valueC, optionBits)
	return bool(ret)
}

func SetArrayItem(x C.XmpPtr, schema *C.char, name string, index int32,
	value string, optionBits C.uint32_t) bool {

	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	valueC := C.CString(value)
	defer C.free(unsafe.Pointer(valueC))

	ret := C.xmp_set_array_item(x, schema, nameC, C.int32_t(index), valueC,
		optionBits)
	return bool(ret)
}

func AppendArrayItem(x C.XmpPtr, schema *C.char, name string,
	arrayOptions C.uint32_t, value string, optionBits C.uint32_t) bool {

	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	valueC := C.CString(value)
	defer C.free(unsafe.Pointer(valueC))

	ret := C.xmp_append_array_item(x, schema, nameC, arrayOptions, valueC,
		optionBits)
	return bool(ret)
}

func StringNew() C.XmpStringPtr {
	return C.xmp_string_new()
}

func StringFree(s C.XmpStringPtr) {
	C.xmp_string_free(s)
}

// Convert an XmpString to Go string
func StringGo(str C.XmpStringPtr) string {
	return C.GoString(C.xmp_string_cstr(str))
}



func init () {
	C.xmp_init()
}
