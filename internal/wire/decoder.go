package wire

type Decoder struct{ buf []byte }

func (d *Decoder) Decode(raw string) []byte {
	d.buf = append(d.buf[:0], raw...)
	return d.buf
}
