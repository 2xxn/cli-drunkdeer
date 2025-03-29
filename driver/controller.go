package driver

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
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

func (d *DrunkDeerController) GetIdentity() *DDKeyboardIdentity {
	if d.identity == nil {
		d.SendIdentity()
		for d.identity == nil && !d.shouldClose {
			time.Sleep(50 * time.Millisecond)
		}
	}

	return d.identity
}

func (d *DrunkDeerController) GetActuations() []byte {
	// No way to get actuations from device because it only echoes whatever you throw at it (if it's valid)

	return d.actuations
}

func (d *DrunkDeerController) sendReport(p []byte) {
	d.Log("Sending report: %x", p)
	report := make([]byte, 64)
	report[0] = KEYBOARD_REPORT_ID // Report ID

	if len(p) > 63 {
		p = p[:63]
	}

	copy(report[1:], p)

	_, err := d.device.Write(report)
	if err != nil {
		panic(err)
	}
}

func (d *DrunkDeerController) SetDebug(debug bool) {
	d.debug = debug
}

func (d *DrunkDeerController) Log(str string, v ...interface{}) {
	if str[len(str)-1] != '\n' {
		str += "\n"
	}

	if d.debug {
		fmt.Printf(color.HiGreenString("[DEBUG] ")+str, v...)
	}
}

// #region Packet senders
func (d *DrunkDeerController) SendIdentity() {
	report := BuildIdentity()
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendLEDModeSelect(direction, sequence, speed, brightness, rgb byte) {
	report := BuildLEDModeSelect(direction, sequence, speed, brightness, rgb)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendLEDModeSelectTurbo(direction, sequence, speed, brightness, rgb byte) {
	report := BuildLEDModeSelectTurbo(direction, sequence, speed, brightness, rgb)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendModifyRow(row uint8, keys []byte) {
	report := BuildModifyRowActuation(row, keys)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendRapidTriggerTurbo(rt, turbo bool) {
	report := BuildRapidTriggerTurbo(rt, turbo)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendDownstrokes(row uint8, keys []byte) {
	report := BuildModifyRowDownstroke(row, keys)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendUpstrokes(row uint8, keys []byte) {
	report := BuildModifyRowUpstroke(row, keys)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) QueuePacket(p []byte) {
	d.packetQueue <- p
}

// #endregion

func (d *DrunkDeerController) LoadActuations(actuations []byte) {
	if len(actuations) != len(KEYBOARD_LAYOUT) {
		panic("Actuations length does not match keyboard layout length")
	}

	d.actuations = actuations

	for i := 0; i < len(d.actuations); i += KEYS_PER_ROW {
		end := i + KEYS_PER_ROW
		if end > len(d.actuations) {
			end = len(d.actuations)
		}
		row := d.actuations[i:end]

		d.SendModifyRow(uint8(i/KEYS_PER_ROW), row)
	}
}

func (d *DrunkDeerController) LoadDownstrokes(downstrokes []byte) {
	if len(downstrokes) != len(KEYBOARD_LAYOUT) {
		panic("Downstrokes length does not match keyboard layout length")
	}

	d.downstrokes = downstrokes

	for i := 0; i < len(d.downstrokes); i += KEYS_PER_ROW {
		end := i + KEYS_PER_ROW
		if end > len(d.downstrokes) {
			end = len(d.downstrokes)
		}

		row := d.downstrokes[i:end]
		rowIndex := uint8(i / KEYS_PER_ROW)
		d.SendDownstrokes(rowIndex, row)
	}
}

func (d *DrunkDeerController) LoadUpstrokes(upstrokes []byte) {
	if len(upstrokes) != len(KEYBOARD_LAYOUT) {
		panic("Upstrokes length does not match keyboard layout length")
	}

	d.upstrokes = upstrokes

	for i := 0; i < len(d.upstrokes); i += KEYS_PER_ROW {
		end := i + KEYS_PER_ROW
		if end > len(d.upstrokes) {
			end = len(d.upstrokes)
		}

		row := d.upstrokes[i:end]
		rowIndex := uint8(i / KEYS_PER_ROW)
		d.SendUpstrokes(rowIndex, row)
	}
}

// #region Modifiers
func (d *DrunkDeerController) ModifyActuationsByNames(names []string, actuations byte) {
	for _, name := range names {
		index := GetIndexByKey(name)
		if index != -1 {
			d.actuations[index] = actuations
		}
	}
}

func (d *DrunkDeerController) ModifyActuationsByIndexes(indexes []int, actuation byte) {
	for _, index := range indexes {
		if index >= 0 && index < len(d.actuations) {
			d.actuations[index] = actuation
		}
	}
}

func (d *DrunkDeerController) ModifyAllActuations(actuation byte) {
	for i := range d.actuations {
		d.actuations[i] = actuation
	}
}

// #endregion

func (d *DrunkDeerController) WriteDefaults() {
	actuations := make([]byte, KEYS_PER_ROW)
	for i := 0; i < KEYS_PER_ROW; i++ {
		actuations[i] = DEFAULT_ACTUATION
	}

	d.Log("Writing defaults")
	d.SendLEDModeSelect(0, SEQUENCE_OFF, 5, 9, 0xff)
	d.SendRapidTriggerTurbo(false, false)

	for i := 0; i < len(d.actuations); i += KEYS_PER_ROW {
		end := i + KEYS_PER_ROW
		if end > len(d.actuations) {
			end = len(d.actuations)
		}
		rowIndex := uint8(i / KEYS_PER_ROW)

		row := d.actuations[i:end]
		d.SendModifyRow(rowIndex, row)

		downstrokes := d.downstrokes[i:end]
		d.SendDownstrokes(rowIndex, downstrokes)

		upstrokes := d.upstrokes[i:end]
		d.SendUpstrokes(rowIndex, upstrokes)
	}

	d.Log("Defaults written")
}

func (d *DrunkDeerController) Close() error {
	var closeErr error
	d.closeOnce.Do(func() {
		d.shouldClose = true

		// Wait for all goroutines to finish.
		done := make(chan struct{})
		go func() {
			d.wg.Wait()
			close(done)
		}()
		select {
		case <-done:
			// All goroutines have finished.
		case <-time.After(5 * time.Second):
			closeErr = fmt.Errorf("timeout waiting for goroutines to exit")
		}

		close(d.packetQueue)
		close(d.packetChan)
	})
	return closeErr
}

func NewDrunkDeerController(device *hid.Device) *DrunkDeerController {
	controller := &DrunkDeerController{
		device:      device,
		packetChan:  make(chan DDPacket),
		packetQueue: make(chan []byte, 10),
		Light:       &DDLight{},
	}

	controller.actuations = make([]byte, len(KEYBOARD_LAYOUT))
	controller.downstrokes = make([]byte, len(KEYBOARD_LAYOUT))
	controller.upstrokes = make([]byte, len(KEYBOARD_LAYOUT))

	for i := range KEYBOARD_LAYOUT {
		controller.actuations[i] = DEFAULT_ACTUATION
		controller.downstrokes[i] = 0x00
		controller.upstrokes[i] = 0x00
	}

	go controller.drunkDeerMessageReceiver()
	controller.wg.Add(1)

	// Start a goroutine to read from the device and send packets to the channel
	go func() {
		defer controller.wg.Done()
		for {
			if controller.shouldClose {
				return // Exit when shouldClose is set
			}

			buf := make([]byte, 64)
			n, err := device.Read(buf)
			if err != nil {
				return // Exit if reading fails
			}

			if buf[0] != KEYBOARD_REPORT_ID {
				continue
			}

			if n > 0 {
				packet := DDPacket{
					Packet: buf[1],
					Data:   buf[2:n],
				}
				select {
				case controller.packetChan <- packet:
				case <-time.After(100 * time.Millisecond): // Prevent blocking
				}
			}
		}
	}()
	controller.wg.Add(1)

	go controller.drunkDeerReporter()
	controller.wg.Add(1)

	return controller
}

func (d *DrunkDeerController) drunkDeerReporter() {
	defer d.wg.Done()
	for {
		select {
		case p, ok := <-d.packetQueue:
			if !ok {
				return // Exit if channel is closed
			}
			d.sendReport(p)
			time.Sleep(100 * time.Millisecond) // Add a small delay to avoid overwhelming the device
		case <-time.After(100 * time.Millisecond): // Check shouldClose periodically
			if d.shouldClose {
				return
			}
		}
	}
}

func (d *DrunkDeerController) drunkDeerMessageReceiver() {
	i := 0
	defer d.wg.Done()
	for p := range d.packetChan {
		data := bytes.NewBuffer(p.Data)
		expectedValue := data.Next(1)
		i += 1

		d.Log("%d Packet received: %x", i, p.Data)
		switch p.Packet {
		case PACKET_IDENTITY:
			// fmt.Printf("Identity packet received: %x\n", p.Data)
			theNullByte, _ := data.ReadByte()
			data.ReadByte()
			modelBytes := data.Next(3)

			if theNullByte != 0x00 {
				d.Log("Unknown byte 3: %x", theNullByte)
				break
			}

			if expectedValue[0] != 0x02 {
				d.Log("Unknown expected value: %x", expectedValue[0])
			}

			versionBytes := data.Next(2)
			model, keyboardType := DetectKeyboardModel(modelBytes)
			version := int(versionBytes[0] | versionBytes[1]<<8)
			ident := DDKeyboardIdentity{
				KeyboardModel:   model,
				KeyboardType:    uint8(keyboardType),
				FirmwareVersion: fmt.Sprintf("0.0%v", version),
				RapidTrigger:    p.Data[15] != 0,
				Turbo:           p.Data[14] != 0,
			}
			d.identity = &ident
			break
		case PACKET_LEDMODESEL:
			d.Light.Direction = p.Data[2]
			d.Light.Sequence = p.Data[3]
			d.Light.Speed = p.Data[4]
			d.Light.Brightness = p.Data[5]

			break
		case PACKET_TURBORT:
			d.Log("Turbo packet received: %x", p.Data)
			d.turbo = p.Data[6] != 0
			d.rapidTrigger = p.Data[7] != 0
			break
		case PACKET_MODIFYKEY:
			d.Log("Modify key packet received: %x", p.Data)
			break
		case PACKET_KEYTRACKING:

			break
		default:
			d.Log("Unknown packet type: %x", p.Packet)
			d.Log("Data: %x", p.Data)
			break
		}
	}
}
