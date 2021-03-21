package main

type LoadStats struct {
	NClients int64
	BytesIn  int64
	BytesOut int64
}
//Status .
type Status struct {
	Title        string
	Time         string
	TimeT        string
	ReadBytes    uint64
	WriteBytes   uint64
	ClientList   []*OVClient
}

//OVClient .
type OVClient struct {
	CommonName      string
	RealAddress     string
	VirtualAddress  string
	BytesReceived   uint64
	BytesSent       uint64
	ConnectedSince  string
	ConnectedSinceT string
	Username        string
}
type Version struct {
	OpenVPN    string
	Management string
}

