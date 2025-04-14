package hdwallet

func eraseBytes(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
