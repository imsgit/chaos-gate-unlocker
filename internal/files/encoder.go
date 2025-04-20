package files

import (
	"bytes"
	"sync"
)

func encode(data []byte) []byte {
	var (
		result  = make([][]byte, len(data))
		wg      sync.WaitGroup
		unicode int
	)

	for i, b := range data {
		switch b {
		case 34, 51, 68, 85, 102, 119:
			result[i] = []byte{b}
			continue
		default:
			if b > 183 {
				unicode = 3
			}
			if unicode > 0 {
				unicode--
				result[i] = []byte{b}
				continue
			}
		}

		wg.Add(1)
		go func(i int, b byte) {
			defer wg.Done()
			b = (b<<4)&0xF0 | (b>>4)&0x0F
			result[i] = []byte(string(b))
		}(i, b)
	}

	wg.Wait()

	return bytes.Join(result, []byte{})
}

func decode(data []byte) []byte {
	var (
		result  = make([]byte, len(data))
		wg      sync.WaitGroup
		shift   bool
		skip    int
		unicode int
	)

	for i, b := range data {
		switch b {
		case 194:
			skip++
			continue
		case 195:
			shift = true
			skip++
			continue
		case 34, 51, 68, 85, 102, 119:
			result[i-skip] = b
			continue
		default:
			if b > 183 {
				unicode = 3
			}
			if unicode > 0 {
				unicode--
				result[i-skip] = b
				continue
			}
		}

		wg.Add(1)
		go func(i int, b byte, shift bool) {
			defer wg.Done()
			b = (b<<4)&0xF0 | (b>>4)&0x0F
			if shift {
				b += 4
			}
			result[i] = b
		}(i-skip, b, shift)

		shift = false
	}

	wg.Wait()

	return result[:len(result)-skip]
}
