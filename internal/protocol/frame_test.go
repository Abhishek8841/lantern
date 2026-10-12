package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

func TestWriteAndReadFrame(t *testing.T) {
	var buffer bytes.Buffer

	want := []byte("Hello from LANtern!")

	if err := WriteFrame(&buffer, want); err != nil {
		t.Fatalf("WriteFrame() error = %v", err)
	}

	got, err := ReadFrame(&buffer)
	if err != nil {
		t.Fatalf("ReadFrame() error = %v", err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("ReadFrame() = %q; want %q", got, want)
	}
}

func TestMultipleFrames(t *testing.T) {
	messages := []string{"Bahh", "Chadd", "shiromni akali dal"}
	var buffer bytes.Buffer

	for _, message := range messages {
		if err := WriteFrame(&buffer, []byte(message)); err != nil {
			t.Fatalf("WriteFrame() error = %v", err)
		}
	}
	for _, want := range messages {
		got, err := ReadFrame(&buffer)
		if err != nil {
			t.Fatalf("ReadFrame() error %v", err)
		}
		if string(got) != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	}
	if buffer.Len() != 0 {
		t.Fatalf("Unexpected %d bytes remaining", buffer.Len())
	}
}

func TestEmptyFrame(t *testing.T) {
	var buffer bytes.Buffer

	if err := WriteFrame(&buffer, nil); err != nil {
		t.Fatalf("WriteFrame() error %v", err)
	}

	got, err := ReadFrame(&buffer)
	if err != nil {
		t.Fatalf("ReadFrame() error = %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("got %d bytes, want 0", len(got))
	}
}

func TestIncompletePayload(t *testing.T) {
	var buffer bytes.Buffer
	var header [HeaderSize]byte

	var payload = "HE"

	binary.BigEndian.PutUint32(header[:], uint32(4))

	buffer.Write(header[:])
	buffer.Write([]byte(payload))

	_, err := ReadFrame(&buffer)

	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected ErrUnexpectedEOF, got %v", err)
	}

}
