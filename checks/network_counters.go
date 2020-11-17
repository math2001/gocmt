package checks

import (
	"github.com/math2001/gocmt/cmt"
	"github.com/shirou/gopsutil/net"
)

func NetworkCounters(
	c *cmt.CheckResult,
	_ map[string]interface{},
) {
	interfacesStat, err := net.IOCounters(false)
	if err != nil {
		panic(err)
	}

	statsSum := interfacesStat[0]

	if prevBytesSent, ok := c.DB["prev_bytes_sent"]; ok {
		c.AddItem(&cmt.CheckItem{
			Name:        "net_bytes_sent_diff",
			Value:       statsSum.BytesSent - uint64(prevBytesSent.(float64)),
			Description: "Number of bytes sent since the last run",
			Unit:        "bytes",
		})
	}
	c.DB["prev_bytes_sent"] = statsSum.BytesSent

	if prevBytesRecv, ok := c.DB["prev_bytes_recv"]; ok {
		c.AddItem(&cmt.CheckItem{
			Name:        "net_bytes_recv_diff",
			Value:       statsSum.BytesRecv - uint64(prevBytesRecv.(float64)),
			Description: "Number of bytes recieved since last run",
			Unit:        "bytes",
		})
	}
	c.DB["prev_bytes_recv"] = statsSum.BytesRecv
}
