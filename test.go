package main


import (


	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"fmt"
)

func main() {

	cs,_ := cpu.Info()

	for _, c := range cs{

		fmt.Println(c.ModelName)
	}

	host, _ := host.Info()
	fmt.Println(host.Hostname, host.String())

	swapmem, _ := mem.SwapMemory()
	fmt.Println(swapmem.Total)
	virtual, _ := mem.VirtualMemory()
	fmt.Println(virtual.Total)

	disk, _ := disk.Usage("/Volumes/Macintosh HD")
	fmt.Println(disk.Total)















}