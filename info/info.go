package info

import (
	"fmt"
	"runtime"
	"strconv"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"../configs"
)

type Info struct {
	Warning			bool		`json:"warning"`
	Hostname 		string		`json:"hostname"`
	Platform 		string		`json:"platform"`
	TotalMemory		string		`json:"memory"`
	UsedMemory		string		`json:"used_memory"`
	Cpus			string		`json:"cpus"`
	TotalDisk		string		`json:"total_disk"`
	UsedDisk		string		`json:"used_disk"`
	UsedDiskPercent	string 		`json:"used_disk_percent"`
	Uptime 			string		`json:"uptime"`
}

func unixtimeToDays (number uint64) string {
	days := int32(number / (86400))
	hours := int32((number % 86400) / 3600)
	minutes := int32((number % 3600) / 60)
	seconds := number % 60
	return fmt.Sprintf("%d day(s) - %d hour(s) - %d minute(s) - %d second(s)", days, hours, minutes, seconds)
}

func diskThreshold (total uint64, used uint64, threshold_percent float64) (float64, bool) {
	diskUsedPercent := (float64(used) / float64(total)) * 100
	if  diskUsedPercent >= threshold_percent {
		return diskUsedPercent, true
	}
	return diskUsedPercent, false
}

func GetInfo() *Info {
	os, _ := host.Info()
	vm, _ := mem.VirtualMemory()
	disk, _ := disk.Usage("/")
	totalDisk := disk.Total
	usedDisk := disk.Used
	config := configs.GetConfigs()
	diskUsedPercent, isWarning := diskThreshold(totalDisk, usedDisk, config.DiskWarning)
	info := Info {
		Hostname: 			os.Hostname,
		Platform: 			os.Platform,
		TotalMemory:		fmt.Sprintf("%dMB", vm.Total / (1024*1024)),
		UsedMemory: 		fmt.Sprintf("%dMB", vm.Used / (1024*1024)),
		Cpus: 				strconv.Itoa(runtime.NumCPU()),
		TotalDisk:			fmt.Sprintf("%dGB", totalDisk / (1024*1024*1024)),
		UsedDisk:			fmt.Sprintf("%dGB", usedDisk / (1024*1024*1024)),
		UsedDiskPercent:	fmt.Sprintf("%.2f", diskUsedPercent) + "%",
		Warning:			isWarning,
		Uptime: 			unixtimeToDays(os.Uptime), 
	}
	return &info
}