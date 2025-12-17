package util

func WipeBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}