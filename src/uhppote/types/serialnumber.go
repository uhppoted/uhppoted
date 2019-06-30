package types

import (
	"encoding/binary"
	"fmt"
)

type SerialNumber uint32

func (s SerialNumber) String() string {
	return fmt.Sprintf("%-10d", s)
}

func (s SerialNumber) MarshalUT0311L0x() ([]byte, error) {
	bytes := make([]byte, 4)

	binary.LittleEndian.PutUint32(bytes, uint32(s))

	return bytes, nil
}

func (s *SerialNumber) UnmarshalUT0311L0x(bytes []byte) (interface{}, error) {
	v := SerialNumber(binary.LittleEndian.Uint32(bytes[:4]))

	return &v, nil
}
