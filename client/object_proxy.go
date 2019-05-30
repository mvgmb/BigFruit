package client

// Upload writes a chunk of bytes into a file
func Upload(filePath string, start int64, bytes []byte) error {
	return nil
}

// Download returns a chunk of bytes from a file
func Download(filePath string, start int64, offset int) ([]byte, error) {
	return nil, nil
}
