package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
	monitor "github.com/kidinamoto01/tools/tm-monitor/monitor"
)

const (
	// Default refresh rate - 200ms
	defaultRefreshRate = time.Millisecond * 200
	totalSteaks = 1800

)

// Ton - table of nodes.
//
// It produces the unordered list of nodes and updates it periodically.
//
// Default output is stdout, but it could be changed. Note if you want for
// refresh to work properly, output must support [ANSI escape
// codes](http://en.wikipedia.org/wiki/ANSI_escape_code).
//
// Ton was inspired by [Linux top
// program](https://en.wikipedia.org/wiki/Top_(software)) as the name suggests.
type Ton struct {
	monitor *monitor.Monitor

	RefreshRate time.Duration
	Output      io.Writer
	quit        chan struct{}
}

func NewTon(m *monitor.Monitor) *Ton {
	return &Ton{
		RefreshRate: defaultRefreshRate,
		Output:      os.Stdout,
		quit:        make(chan struct{}),
		monitor:     m,
	}
}

func (o *Ton) Start() {
	clearScreen(o.Output)
	o.Print()
	go o.refresher()
}

func (o *Ton) Print() {
	moveCursor(o.Output, 1, 1)
	o.printHeader()
	fmt.Println()
	o.printTable()
}

func (o *Ton) Stop() {
	close(o.quit)
}

func (o *Ton) printHeader() {
	n := o.monitor.Network
	fmt.Fprintf(o.Output, "%v up %.2f%%\n", n.StartTime(), n.Uptime())
	fmt.Println("Total steaks in genesis:",totalSteaks)
	fmt.Println("Total steaks bonded:",n.PowerSum)
	//ratio :=int(n.PowerSum)*100/totalSteaks
	fmt.Println("Bonded Percentage: ",(n.PowerSum)*100/totalSteaks,"% ")
	fmt.Fprintf(o.Output, "Height: %d\n", n.Height)
	fmt.Fprintf(o.Output, "Power Online: %d %%\n", n.PowerOnline)

	fmt.Fprintf(o.Output, "Avg block time: %.3f ms\n", n.AvgBlockTime)
	fmt.Fprintf(o.Output, "Avg tx throughput: %.0f per sec\n", n.AvgTxThroughput)
	fmt.Fprintf(o.Output, "Avg block latency: %.3f ms\n", n.AvgBlockLatency)
	fmt.Fprintf(o.Output, "Active nodes: %d/%d (health: %s) Validators: %d\n", n.NumNodesMonitoredOnline, n.NumNodesMonitored, n.GetHealthString(), n.NumValidators)
}

func (o *Ton) printTable() {
	w := tabwriter.NewWriter(o.Output, 0, 0, 5, ' ', 0)
	nw := o.monitor.Network
	if nw.PowerSum > 0{
		fmt.Fprintln(w, "NAME\tHEIGHT\tBLOCK LATENCY\tONLINE\tVALIDATOR\tPOWER\tRatio %\tPrecommit_Sum\tPrecommit_Miss\t")
		for _, n := range o.monitor.Nodes {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%d\t%.3f ms\t%v\t%v\t%d\t%d\t%d\t%d\t", n.Name, n.Height, n.BlockLatency, n.Online, n.IsValidator,n.Power,n.Power*100/nw.PowerSum,n.PrecommitSum,n.PrecommitMiss))
		}
	}
	// else{
	//	fmt.Fprintln(w, "NAME\tHEIGHT\tBLOCK LATENCY\tONLINE\tVALIDATOR\tPOWER\tRatio\tPrecommit_Sum\tPrecommit_Miss\t")
	//	for _, n := range o.monitor.Nodes {
	//		fmt.Fprintln(w, fmt.Sprintf("%s\t%d\t%.3f ms\t%v\t%v\t%d\t%d\t%d\t%d\t", n.Name, n.Height, n.BlockLatency, n.Online, n.IsValidator,n.Power,0,n.PrecommitSum,n.PrecommitMiss))
	//	}
	//}

	w.Flush()
}

// Internal loop for refreshing
func (o *Ton) refresher() {
	for {
		select {
		case <-o.quit:
			return
		case <-time.After(o.RefreshRate):
			o.Print()
		}
	}
}

func clearScreen(w io.Writer) {
	fmt.Fprint(w, "\033[2J")
}

func moveCursor(w io.Writer, x int, y int) {
	fmt.Fprintf(w, "\033[%d;%dH", x, y)
}
