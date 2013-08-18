package identicon

import (
        "testing"
)

func BenchmarkIdenticon(b *testing.B) {

	for i := 0; i < b.N; i++ {
		item := "dgryski"
		key := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
		icon := New7x7(key)
		data := []byte(item)
		icon.Render(data)
	}
}
