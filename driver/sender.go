package driver

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
	var packet []byte
	if len(p) != 63 {
		packet = make([]byte, 63)
		copy(packet[:], p[:])
	} else {
		packet = p
	}

	d.packetQueue <- packet
}
