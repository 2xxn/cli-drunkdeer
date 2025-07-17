package driver

import (
	"sync"

	"github.com/sstallion/go-hid"
)

type DrunkDeerController struct {
	device   *hid.Device
	identity *DDKeyboardIdentity

	actuations  []byte
	downstrokes []byte
	upstrokes   []byte

	turbo        bool
	rapidTrigger bool
	debug        bool

	wg          sync.WaitGroup
	packetChan  chan DDPacket
	packetQueue chan []byte

	Light *DDLight

	shouldClose bool
	closeOnce   sync.Once
}

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
