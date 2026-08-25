package payload

type PayloadParser struct {
	buffer []byte
}

func (p *PayloadParser) Parse(body []byte) []byte {
	p.buffer = append(p.buffer[:0], body...)
	return p.buffer
}
