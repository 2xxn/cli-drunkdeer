package driver

type DDPacket struct {
	Packet uint8
	Data   []byte
}

type DDLight struct {
	Direction  byte
	Speed      byte
	Sequence   byte
	Brightness byte
}

type DDKeyboardIdentity struct {
	KeyboardModel string
	KeyboardType  uint8

	FirmwareVersion string
	Turbo           bool
	RapidTrigger    bool
}
