package modules

type ChanWriter struct {
	Chan chan []byte
}

func (cw ChanWriter) Write(b []byte) (n int, err error) {
	cw.Chan <- b
	return 0, nil
}

/*
func NewChanWriter() ChanWriter {
        c := make(map[string](chan []byte), 0)
        cw := ChanWriter{c}
        return cw
}
*/
