package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	HeaderSize   = 4
	MaxFrameSize = 1 << 20
	//  x << y === x * (2^y) ???
)

func WriteFrame(w io.Writer, payload []byte) error {
	if len(payload) > MaxFrameSize {
		return fmt.Errorf("frame too large: %d bytes", len(payload))
	}

	var header [HeaderSize]byte

	binary.BigEndian.PutUint32(header[:], uint32(len(payload)))

	if err := writeFull(w, header[:]); err != nil {
		return fmt.Errorf("write frame header: %w", err)
	}

	if err := writeFull(w, payload); err != nil {
		return fmt.Errorf("write frame payload: %w", err)
	}

	return nil
}

func ReadFrame(r io.Reader) ([]byte, error) {
	var header [HeaderSize]byte

	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, fmt.Errorf("read frame header: %w", err)
	}

	length := binary.BigEndian.Uint32(header[:])

	if length > MaxFrameSize {
		return nil, fmt.Errorf("frame too large: %d bytes", length)
	}

	payload := make([]byte, int(length))

	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, fmt.Errorf("read frame payload: %w", err)
	}

	return payload, nil
}

// using io.writer is known as programming to an interface rather than a concrete implementation
func writeFull(w io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := w.Write(data)

		if n > 0 {
			data = data[n:]
		}

		if err != nil {
			return err
		}

		if n == 0 {
			return io.ErrShortWrite
			// return when writing didn't make enough progress
		}
	}

	return nil
}
