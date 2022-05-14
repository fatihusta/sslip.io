// Copyright 2021 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nn

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"

	"github.com/nlpodyssey/spago/mat"
)

// init registers the param implementation with the gob subsystem - so that it knows how to encode and decode
// values of type nn.Param
func init() {
	gob.Register(&BaseParam[float32]{})
	gob.Register(&BaseParam[float64]{})
}

// MarshalBinary marshals a param into binary form.
func (p *BaseParam[T]) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := mat.MarshalBinaryMatrix[T](p.value, buf)
	if err != nil {
		return nil, err
	}

	if p.payload == nil {
		buf.WriteByte(0)
	} else {
		buf.WriteByte(1)
		pBin, err := p.payload.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(pBin)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary unmarshals a param from binary form.
func (p *BaseParam[T]) UnmarshalBinary(data []byte) error {
	var err error
	buf := bytes.NewReader(data)

	p.value, err = mat.UnmarshalBinaryMatrix[T](buf)
	if err != nil {
		return err
	}

	hasPayload, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if hasPayload == 0 {
		p.payload = nil
		return nil
	}

	p.payload = new(Payload[T])

	pBin, err := io.ReadAll(buf)
	if err != nil {
		return err
	}

	return p.payload.UnmarshalBinary(pBin)
}

// MarshalBinaryParam encodes a Param into binary form.
func MarshalBinaryParam[T mat.DType](p *BaseParam[T], w io.Writer) error {
	if p == nil {
		_, err := w.Write([]byte{0})
		return err
	}

	_, err := w.Write([]byte{1})
	if err != nil {
		return err
	}

	bin, err := p.MarshalBinary()
	if err != nil {
		return err
	}

	binLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(binLen, uint32(len(bin)))
	_, err = w.Write(binLen)
	if err != nil {
		return err
	}

	_, err = w.Write(bin)
	return err
}

// UnmarshalBinaryParam decodes a Param from binary form.
// TODO: add a "withBacking" optional argument to remove the need of UnmarshalBinaryParamWithReceiver().
func UnmarshalBinaryParam[T mat.DType](r io.Reader) (*BaseParam[T], error) {
	isPresent := make([]byte, 1)
	_, err := r.Read(isPresent)
	if err != nil {
		return nil, err
	}
	if isPresent[0] == 0 {
		return nil, nil
	}

	binLenBytes := make([]byte, 4)
	_, err = r.Read(binLenBytes)
	if err != nil {
		return nil, err
	}
	binLen := int(binary.LittleEndian.Uint32(binLenBytes))
	bin := make([]byte, binLen)
	_, err = r.Read(bin)
	if err != nil {
		return nil, err
	}

	p := new(BaseParam[T])
	err = p.UnmarshalBinary(bin)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// UnmarshalBinaryParamWithReceiver decodes a Param from binary form into the receiver.
func UnmarshalBinaryParamWithReceiver[T mat.DType](r io.Reader, dest *BaseParam[T]) error {
	isPresent := make([]byte, 1)
	_, err := r.Read(isPresent)
	if err != nil {
		return err
	}
	if isPresent[0] == 0 {
		return nil
	}

	binLenBytes := make([]byte, 4)
	_, err = r.Read(binLenBytes)
	if err != nil {
		return err
	}
	binLen := int(binary.LittleEndian.Uint32(binLenBytes))
	bin := make([]byte, binLen)
	_, err = r.Read(bin)
	if err != nil {
		return err
	}

	err = dest.UnmarshalBinary(bin)
	if err != nil {
		return err
	}
	return nil
}
