package structure

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

func TransmitAndReceive(structure interface{}, rw io.ReadWriter) error {
	if err := Transmit(structure, rw); err != nil {
        return err
    }
    return Receive(structure, rw)
}

func Transmit(structure interface{}, w io.Writer) error {
	data, err := json.Marshal(structure)
	if err != nil {
		return err
	}
	l := len(data)
	tx := make([]byte, l+8)
	binary.BigEndian.PutUint64(tx, uint64(l))
	copy(tx[8:], data)
	n := 0
	for n < l+8 {
        x, err := w.Write(tx[n:])
		if err != nil {
            return err
        }
		n += x
    }
	return nil
}

func Receive(structure interface{}, r io.Reader) error {
	sb := make([]byte, 8)
	_, err := io.ReadFull(r, sb)
	if err != nil {
        return err
    }
    l := binary.BigEndian.Uint64(sb)
    if err != nil {
        return err
    }
    rx := make([]byte, l)
    n := 0
    for n < int(l) {
        x, err := r.Read(rx[n:])
        if err == io.EOF {
            break
        } else if err != nil {
            return err
        }
        n += x
    }
    err = json.Unmarshal(rx, structure)
    if err != nil {
        return err
    }
    return nil
}

