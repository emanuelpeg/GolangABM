package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/zcalusic/sysinfo"
	"net/http"
	"runtime"
	"strconv"
)

func info(c echo.Context) error {
	var si sysinfo.SysInfo
	si.GetSysInfo()

	return c.JSON(http.StatusOK, si)
}

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

type StatusInfo struct {
	RuntimeOS            string `json:"runtimeOS"`
	TotalMemory          string `json:"total_memory"`
	FreeMemory           string `json:"free_memory"`
	PercentageUsedMemory string `json:"percentage_used_memory"`

	TotalDiskSpace           string `json:"total_disk_space"`
	UsedDiskSpace            string `json:"used_disk_space"`
	FreeDiskSpace            string `json:"free_disk_space"`
	PercentageDiskSpaceUsage string `json:"percentage_disk_space_usage"`

	CpuIndexNumber string `json:"cpu_index_number"`
	VendorId       string `json:"vendorId"`
	Family         string `json:"family"`
	NumberOfCores  string `json:"number_of_cores"`
	ModelName      string `json:"model_name"`
	Speed          string `json:"speed"`

	Hostname                 string `json:"hostname"`
	Uptime                   string `json:"uptime"`
	NumberOfProcessesRunning string `json:"number_of_processes_running"`

	Os       string `json:"os"`
	Platform string `json:"platform"`

	HostId string `json:"host_id"`
}

func status(c echo.Context) error {
	runtimeOS := runtime.GOOS
	// memory
	vmStat, err := mem.VirtualMemory()
	dealwithErr(err)

	diskStat, err := disk.Usage("/")
	dealwithErr(err)

	// cpu - get CPU number of cores and speed
	cpuStat, err := cpu.Info()
	dealwithErr(err)

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	dealwithErr(err)

	var sInfo StatusInfo

	sInfo.RuntimeOS = runtimeOS
	sInfo.TotalMemory = strconv.FormatUint(vmStat.Total, 10) + " bytes "
	sInfo.FreeMemory = strconv.FormatUint(vmStat.Free, 10) + " bytes"
	sInfo.PercentageUsedMemory = strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%"

	sInfo.TotalDiskSpace = strconv.FormatInt(int64(diskStat.Total), 10) + " bytes "
	sInfo.UsedDiskSpace = strconv.FormatInt(int64(diskStat.Used), 10) + " bytes"
	sInfo.FreeDiskSpace = strconv.FormatInt(int64(diskStat.Free), 10) + " bytes"
	sInfo.PercentageDiskSpaceUsage = strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64) + "%"

	sInfo.CpuIndexNumber = strconv.FormatInt(int64(cpuStat[0].CPU), 10)
	sInfo.VendorId = cpuStat[0].VendorID
	sInfo.Family = cpuStat[0].Family
	sInfo.NumberOfCores = strconv.FormatInt(int64(cpuStat[0].Cores), 10)
	sInfo.ModelName = cpuStat[0].ModelName
	sInfo.Speed = strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64)

	sInfo.Hostname = hostStat.Hostname
	sInfo.Uptime = strconv.FormatUint(hostStat.Uptime, 10)
	sInfo.NumberOfProcessesRunning = strconv.FormatUint(hostStat.Procs, 10)

	sInfo.Os = hostStat.OS
	sInfo.Platform = hostStat.Platform

	sInfo.HostId = hostStat.HostID

	return c.JSON(http.StatusOK, sInfo)
}

func main() {
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})

	e.GET("/info", func(c echo.Context) error {
		return info(c)
	})

	e.GET("/status", func(c echo.Context) error {
		return status(c)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
