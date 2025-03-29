package main

import (
	"github.com/sstallion/go-hid"
)

type HIDDeviceData struct {
	VendorId  uint16
	ProductId uint16
	UsagePage uint16
	Usage     uint16
}

func GetDDProducts() []HIDDeviceData {
	return []HIDDeviceData{
		{0x352D, 0x2383, 0xFF00, 0x00},
		{0x352D, 0x2382, 0xFF00, 0x00},
		{0x352D, 0x2384, 0xFF00, 0x00},
		{0x352D, 0x2386, 0xFF00, 0x00},
		{0x05AC, 0x024F, 0xFF00, 0x00},
	}
}

func FindDrunkDeerDevices() []hid.DeviceInfo {
	devices := GetDDProducts()
	foundDevices := make([]hid.DeviceInfo, 0)
	hid.Enumerate(hid.VendorIDAny, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		for _, device := range devices {
			if info.VendorID == device.VendorId && info.ProductID == device.ProductId && info.UsagePage == device.UsagePage && info.Usage == device.Usage {
				foundDevices = append(foundDevices, *info)
				return nil
			}
		}
		return nil
	})
	return foundDevices
}
