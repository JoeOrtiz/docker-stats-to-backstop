// vim: ts=2 nowrap
package translate

import (
	"github.com/JoeOrtiz/docker-stats-to-backstop/translate/backstop"
	"github.com/JoeOrtiz/docker-stats-to-backstop/translate/docker"
	"time"
)

func Translate(prefix string, stats *docker.ContainerStats) []backstop.Metric {
	c := collector{prefix: prefix, timestamp: stats.Timestamp, metrics: []backstop.Metric{}}
	c.add("cpu.system", stats.CpuStats.SystemUsage)
	c.add("cpu.cores", cpuCores(stats))
	c.add("cpu.percent", cpuPercent(stats))
	c.add("cpu.total", stats.CpuStats.CpuUsage.TotalUsage)
	c.add("cpu.kernel", stats.CpuStats.CpuUsage.UsageInKernelmode)
	c.add("cpu.user", stats.CpuStats.CpuUsage.UsageInUsermode)
	c.add("memory.usage", stats.MemoryStats.Usage)
	c.add("memory.cache", stats.MemoryStats.Stats.TotalCache)
	c.add("memory.active", activeMemory(stats))
	c.add("memory.max_usage", stats.MemoryStats.MaxUsage)
        c.add("memory.fail", stats.MemoryStats.Failcnt)
	c.add("memory.limit", stats.MemoryStats.Limit)
	c.add("network.rx_bytes", stats.Network.RxBytes)
	c.add("network.rx_packets", stats.Network.RxPackets)
	c.add("network.rx_errors", stats.Network.RxErrors)
	c.add("network.rx_dropped", stats.Network.RxDropped)
	c.add("network.tx_bytes", stats.Network.TxBytes)
	c.add("network.tx_packets", stats.Network.TxPackets)
	c.add("network.tx_errors", stats.Network.TxErrors)
	c.add("network.tx_dropped", stats.Network.TxDropped)
//	  c.add("blkio.io_serviced_bytes", stats.BlkioStats.IoServiceBytesRecursive.Value)
//        c.add("blkio.io_serviced", stats.BlkioStats.IoServicedRecursive.Value)
//        c.add("blkio.io_queued", stats.BlkioStats.IoQueuedRecursive.Value)
//        c.add("blkio.io_serviced_time", stats.BlkioStats.IoServiceTimeRecursive.Value)
//        c.add("blkio.io_wait_time", stats.BlkioStats.IoWaitTimeRecursive.Value)
//        c.add("blkio.io_merged", stats.BlkioStats.IoMergedRecursive.Value)
//        c.add("blkio.io_time", stats.BlkioStats.IoTimeRecursive.Value)
//        c.add("blkio.sectors", stats.BlkioStats.SectorsRecursive.Value)
	return c.metrics
}

type collector struct {
	prefix    string
	timestamp time.Time
	metrics   []backstop.Metric
}

func (c *collector) add(name string, value *uint64) {
	if value != nil {
		c.metrics = append(c.metrics, backstop.Metric{
			Name:      c.prefix + "." + name,
			Value:     *value,
			Timestamp: c.timestamp.Unix(),
		})
	}
}

func activeMemory(stats *docker.ContainerStats) *uint64 {
	usage := stats.MemoryStats.Usage
	cache := stats.MemoryStats.Stats.TotalCache
	if usage != nil && cache != nil {
		active := *usage - *cache
		return &active
	}
	return nil
}

func cpuCores(stats *docker.ContainerStats) *uint64 {
        containerCpucores := uint64(len(stats.CpuStats.CpuUsage.PercpuUsage))
        return &containerCpucores
}

func cpuPercent(stats *docker.ContainerStats) *uint64 {
	cpuDelta := stats.CpuStats.CpuUsage.TotalUsage
        systemDelta := stats.CpuStats.SystemUsage
	if systemDelta != nil && cpuDelta != nil {
                containercpuPercent := (*cpuDelta / *systemDelta) * uint64(len(stats.CpuStats.CpuUsage.PercpuUsage)) * 100.0
                return &containercpuPercent
	}
        return nil
}
