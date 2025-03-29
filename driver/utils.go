package driver

import "bytes"

func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func DetectKeyboardModel(modelBytes []byte) (string, int) { // Model and type
	A75 := [][]byte{
		{0x0b, 0x01, 0x01},
		{0x0b, 0x04, 0x01},
	}

	A75Pro := [][]byte{
		{0x0b, 0x04, 0x03},
	}

	G65 := [][]byte{
		{0x0f, 0x01, 0x01},
		{0x0b, 0x02, 0x01},
	}

	G60 := [][]byte{
		{0x0b, 0x03, 0x01},
	}

	G75 := [][]byte{
		{0x0b, 0x04, 0x05},
	}

	for _, modelA75 := range A75 {
		if bytes.Equal(modelA75, modelBytes) {
			return KEYBOARD_A75, 75
		}
	}

	for _, modelA75Pro := range A75Pro {
		if bytes.Equal(modelA75Pro, modelBytes) {
			return KEYBOARD_A75PRO, 750
		}
	}

	for _, modelG75 := range G75 {
		if bytes.Equal(modelG75, modelBytes) {
			return KEYBOARD_G75, 754
		}
	}

	for _, modelG65 := range G65 {
		if bytes.Equal(modelG65, modelBytes) {
			return KEYBOARD_G65, 65
		}
	}

	for _, modelG60 := range G60 {
		if bytes.Equal(modelG60, modelBytes) {
			return KEYBOARD_G60, 60
		}
	}

	// Default to A75 if no match found
	return KEYBOARD_A75, 75
}

func GetKeyByIndex(index int) string {
	if index >= 0 && index < len(KEYBOARD_LAYOUT) {
		return KEYBOARD_LAYOUT[index]
	}
	return ""
}

func GetIndexByKey(key string) int {
	for i, v := range KEYBOARD_LAYOUT {
		if v == key {
			return i
		}
	}
	return -1
}

func GetRowByIndex(index int) int {
	if index >= 0 && index < len(KEYBOARD_LAYOUT) {
		return index / KEYS_PER_ROW
	}
	return -1
}

func ActuationFloatToByte(actuation float32) byte {
	var normalized int = int(actuation * 10)
	if normalized < 0 {
		normalized = 1
	} // can't afford values this low

	if normalized > 255 {
		normalized = 255
	} // WHY WOULD ANYONE EVEN TRY TO SET THIS

	return byte(normalized)
}
