package imports

import (
	"strings"
	"testing"
)

func TestParseValid(t *testing.T) {
	r, e := Parse(strings.NewReader(`<nmaprun><host><address addr="1.2.3.4" addrtype="ipv4"/><ports><port protocol="tcp" portid="80"><state state="open"/><service name="http"/></port></ports></host></nmaprun>`))
	if e != nil || len(r.Hosts) != 1 || r.Hosts[0].Address() != "1.2.3.4" {
		t.Fatalf("unexpected parse: %#v %v", r, e)
	}
}
func TestParseMalformed(t *testing.T) {
	if _, e := Parse(strings.NewReader(`<nmaprun><host>`)); e == nil {
		t.Fatal("expected malformed XML error")
	}
}
func TestOptionalFields(t *testing.T) {
	r, e := Parse(strings.NewReader(`<nmaprun><host><ports/></host></nmaprun>`))
	if e != nil || r.Hosts[0].Hostname() != "" {
		t.Fatal("optional fields should be accepted")
	}
}
