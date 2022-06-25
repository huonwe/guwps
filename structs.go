package main

type DeviceMsg struct {
	Addr string
	Msg  []byte
}

type UDPData struct {
	Length     uint64
	Version    string
	DeviceType string
	DataType   string
	ID         string
	MM         string
	Tag        string

	STA string

	Identifier string
}
