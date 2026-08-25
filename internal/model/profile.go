package model

type CodecProfile struct {
	Name   string
	Stages map[string]bool
	Ready  bool
}
