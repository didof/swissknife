package sniffer

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"

	"github.com/alexeyco/simpletable"
)

func Run(iface string) {
	device, _, err := getDevice(iface)
	if err != nil {
		log.Print(err)
		PrintDevicesInfo()
		os.Exit(1)
	}

	fmt.Println(device)
}

func getDevice(iface string) (device pcap.Interface, devices []pcap.Interface, err error) {
	devices, err = pcap.FindAllDevs()
	if err != nil {
		return device, devices, err
	}

	for _, d := range devices {
		if d.Name == iface {
			device = d
			break
		}
	}

	if device.Name == "" {
		return device, devices, fmt.Errorf("the device '%s' is not present on this machine", iface)
	}

	return device, devices, nil
}

func PrintDevicesInfo() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		errors.Wrap(err, "could not retrieve devices")
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "IP"},
			{Align: simpletable.AlignCenter, Text: "Netmask"},
		},
	}

	for _, device := range devices {
		var (
			ip      net.IP
			netmask net.IPMask
		)
		for _, address := range device.Addresses {
			ip = address.IP
			netmask = address.Netmask
		}
		row := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: device.Name},
			{Align: simpletable.AlignRight, Text: ip.String()},
			{Align: simpletable.AlignRight, Text: netmask.String()},
		}

		table.Body.Cells = append(table.Body.Cells, row)
	}

	fmt.Println(table.String())
}
