package imports

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type Report struct {
	XMLName xml.Name `xml:"nmaprun"`
	Hosts   []Host   `xml:"host"`
}
type Host struct {
	Status struct {
		State string `xml:"state,attr"`
	} `xml:"status"`
	Addresses []struct {
		Addr string `xml:"addr,attr"`
		Type string `xml:"addrtype,attr"`
	} `xml:"address"`
	Hostnames []struct {
		Name string `xml:"name,attr"`
	} `xml:"hostnames>hostname"`
	Ports []Port `xml:"ports>port"`
}
type Port struct {
	Protocol string `xml:"protocol,attr"`
	ID       int    `xml:"portid,attr"`
	State    struct {
		State string `xml:"state,attr"`
	} `xml:"state"`
	Service struct {
		Name string `xml:"name,attr"`
	} `xml:"service"`
	Scripts []struct {
		ID     string `xml:"id,attr"`
		Output string `xml:"output,attr"`
	} `xml:"script"`
}

func Parse(r io.Reader) (Report, error) {
	d := xml.NewDecoder(io.LimitReader(r, 10<<20))
	d.Strict = true
	var report Report
	if err := d.Decode(&report); err != nil {
		return report, fmt.Errorf("invalid Nmap XML: %w", err)
	}
	if report.XMLName.Local != "nmaprun" {
		return report, fmt.Errorf("unsupported XML: expected nmaprun root element")
	}
	return report, nil
}
func (h Host) Address() string {
	for _, a := range h.Addresses {
		if a.Type == "ipv4" || a.Type == "ipv6" {
			return a.Addr
		}
	}
	return ""
}
func (h Host) Hostname() string {
	if len(h.Hostnames) > 0 {
		return h.Hostnames[0].Name
	}
	return ""
}
func ScriptEvidence(p Port) string {
	var x []string
	for _, s := range p.Scripts {
		x = append(x, s.ID+": "+s.Output)
	}
	return strings.Join(x, "\n")
}
