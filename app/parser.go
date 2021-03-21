package main

import(
	"strings"
	"fmt"
	"strconv"
)

//ParseVersion gets version information from string
func parseVersion(input string) (*Version, error) {
	v := Version{}
	a := strings.Split(trim(input), "\n")
	if len(a) != 3 {
		return nil, fmt.Errorf("Wrong number of lines, expected %d, got %d", 3, len(a))
	}
	v.OpenVPN = stripPrefix(a[0], "OpenVPN Version: ")
	v.Management = stripPrefix(a[1], "Management Version: ")

	return &v, nil
}
//ParseStats gets stats from string
func parseStats(input string) (*LoadStats, error) {
	ls := LoadStats{}
	a := strings.Split(trim(input), "\n")

	if len(a) != 1 {
		return nil, fmt.Errorf("Wrong number of lines, expected %d, got %d", 1, len(a))
	}
	line := a[0]
	if !isSuccess(line) {
		return nil, fmt.Errorf("Bad response: %s", line)
	}

	dString := stripPrefix(line, "SUCCESS: ")
	dElements := strings.Split(dString, ",")
	var err error
	ls.NClients, err = getLStatsValue(dElements[0])
	if err != nil {
		return nil, err
	}
	ls.BytesIn, err = getLStatsValue(dElements[1])
	if err != nil {
		return nil, err
	}
	ls.BytesOut, err = getLStatsValue(dElements[2])
	if err != nil {
		return nil, err
	}
	return &ls, nil
}

func getLStatsValue(s string) (int64, error) {
	a := strings.Split(s, "=")
	if len(a) != 2 {
		return int64(-1), fmt.Errorf("Parsing error")
	}
	return strconv.ParseInt(a[1], 10, 64)
}
//ParseStatus gets status information from string
func parseStatus(input string) (*Status, error) {
	s := &Status{}
	s.ClientList = make([]*OVClient, 0, 0)
	a := strings.Split(trim(input), "\n")
	for _, line := range a {
		fields := strings.Split(trim(line), ",")
		c := fields[0]
		switch {
		case c == "TITLE":
			s.Title = fields[1]
		case c == "TIME":
			s.Time = fields[1]
			s.TimeT = fields[2]
		case c == "TCP/UDP read bytes":
			bytes, _ := strconv.ParseUint(fields[1], 10, 64)
			s.ReadBytes = bytes
		case c == "TCP/UDP write bytes":
			bytes, _ := strconv.ParseUint(fields[1], 10, 64)
			s.WriteBytes = bytes
		case c == "CLIENT_LIST":
			bytesR, _ := strconv.ParseUint(fields[4], 10, 64)
			bytesS, _ := strconv.ParseUint(fields[5], 10, 64)
			item := &OVClient{
				CommonName:      fields[1],
				RealAddress:     fields[2],
				VirtualAddress:  fields[3],
				BytesReceived:   bytesR,
				BytesSent:       bytesS,
				ConnectedSince:  fields[6],
				ConnectedSinceT: fields[7],
				Username:        fields[8],
			}
			s.ClientList = append(s.ClientList, item)
		}
	}
	return s, nil
}
func trim(s string) string {
	return strings.Trim(strings.Trim(s, "\r\n"), "\n")
}

func stripPrefix(s, prefix string) string {
	return trim(strings.Replace(s, prefix, "", 1))
}

func isSuccess(s string) bool {
	if strings.HasPrefix(s, "SUCCESS: ") {
		return true
	}
	return false
}

