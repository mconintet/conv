package conv

import (
	"errors"
	"io"
)

func unicodeToUtf8(u uint32) (u8 []byte, err error) {
	switch {
	case u <= 0x7F:
		u8 = append(u8, byte(u))
	case u >= 0x80 && u <= 0x7FF:
		u8 = append(u8, byte((u>>6)+0xC0), byte(u&0x3F+0x80))
	case u >= 0x0080 && u <= 0xFFFF:
		u8 = append(u8, byte(u>>12+0xE0), byte((u&0xFFF)>>6+0x80), byte(u&0x3F+0x80))
	case u >= 0x10000 && u <= 0x10FFFF:
		u8 = append(u8,
			byte(u>>18+0xF0),
			byte((u&0x3FFFF)>>12+0x80),
			byte((u&0xFFF)>>6+0x80),
			byte(u&0x3F+0x80),
		)
	default:
		return nil, errors.New("invalid range.")
	}

	return u8, nil
}

func gbkToUtf8(gbk uint16) (u8 []byte, err error) {
	var (
		u  uint32
		ok bool
	)

	if u, ok = mapGbkUnicode[gbk]; !ok {
		return nil, errors.New("invalid code point.")
	}

	if u8, err = unicodeToUtf8(u); err != nil {
		return nil, err
	}

	return u8, nil
}

func GbkToUtf8(r io.Reader, w io.Writer, c bool) error {
	var (
		buf []byte
		gbk uint16
		i   int
		err error
		u8  []byte
	)

	buf = make([]byte, 1)

	for {
		i, err = r.Read(buf)

		if i == 1 {
			if gbk == 0 {
				gbk = uint16(buf[0])

				if gbk <= 0x7F {
					w.Write([]byte{byte(gbk)})
					gbk = 0
				}
			} else {
				gbk = gbk<<8 + uint16(buf[0])
				if u8, err = gbkToUtf8(gbk); err != nil && !c {
					return err
				} else {
					w.Write(u8)
				}

				gbk = 0
			}

			continue
		}

		if i == 0 || err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}
