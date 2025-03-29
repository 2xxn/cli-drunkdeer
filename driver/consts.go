package driver

const (
	KEYBOARD_REPORT_ID   = 0x04
	KEYS_PER_ROW         = 59   // According to their dumb layout
	DEFAULT_ACTUATION    = 0x14 // 0x14 is 2.0mm, there should be always 59 keys!!! (unless 3rd row)
	CUSTOM_COLOR_PADDING = 0x80 // TODO: support
	COLORS_PER_PACKET    = 13   // hex, max 13 keys per packet TODO: support
)

const (
	PACKET_IDENTITY    = 0xA0
	PACKET_LEDMODESEL  = 0xAE
	PACKET_MODIFYKEY   = 0xB6
	PACKET_TURBORT     = 0xB5
	PACKET_KEYTRACKING = 0xB7
)

const (
	KEYBOARD_A75PRO = "A75" // We don't give a damn about pro, they work pretty much the same
	KEYBOARD_A75    = "A75"
	KEYBOARD_G65    = "G65"
	KEYBOARD_G60    = "G60"
	KEYBOARD_G75    = "G75"
)

const (
	SEQUENCE_OFF = 0x00

	SEQUENCE_ALWAYS          = 0x02
	SEQUENCE_SPECTRUM        = 0x03
	SEQUENCE_BREATH          = 0x04
	SEQUENCE_PRESS           = 0x05
	SEQUENCE_STARS           = 0x06
	SEQUENCE_WAVE            = 0x07
	SEQUENCE_SURF            = 0x08
	SEQUENCE_SURFDOWN        = 0x09
	SEQUENCE_RIPPLE          = 0x0A
	SEQUENCE_FISH            = 0x0B
	SEQUENCE_FOUNTAIN        = 0x0C
	SEQUENCE_TRAFFIC         = 0x0D
	SEQUENCE_SNAKE           = 0x0E
	SEQUENCE_SURF_REPEAT     = 0x0F
	SEQUENCE_SURF_CROSS      = 0x10
	SEQUENCE_LASER_KEY       = 0x11
	SEQUENCE_FOUNTAIN_RANDOM = 0x12
	SEQUENCE_CUSTOM          = 0x13
)

// Imagine this is a const
var KEYBOARD_LAYOUT = []string{
	"ESC", "", "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10", "F11", "F12", "KP7", "KP8", "KP9", "", "", "", "",
	"TILDE", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "MINUS", "PLUS", "BACK", "KP4", "KP5", "KP6", "", "", "", "",
	"TAB", "Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P", "BRKTS_L", "BRKTS_R", "SLASH_K29", "KP1", "KP2", "KP3", "", "", "", "",
	"CAPS", "A", "S", "D", "F", "G", "H", "J", "K", "L", "COLON", "QOTATN", "", "RETURN", "", "KP0", "KP_DEL", "", "", "", "",
	"SHF_L", "EUR_K45", "Z", "X", "C", "V", "B", "N", "M", "COMMA", "PERIOD", "SLASH", "", "SHF_R", "ARR_UP", "", "NUMS", "", "", "", "",
	"CTRL_L", "WIN_L", "ALT_L", "", "", "", "SPACE", "", "", "", "ALT_R", "FN1", "APP", "ARR_L", "ARR_DW", "ARR_R", "CTRL_R", "", "", "", "",
}

var WASD_KEYS = []int{44, 64, 65, 66}
var NUMERALS_KEYS = []int{22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
var CHARACTER_KEYS = []int{
	43, 44, 45, 46, 47, 48, 49, 50, 51, 52,
	64, 65, 66, 67, 68, 69, 70, 71, 72,
	86, 87, 88, 89, 90, 91, 92,
}
