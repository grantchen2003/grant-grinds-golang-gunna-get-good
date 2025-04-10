package utils

func Chunkify[T any](arr []T, chunkSize int) [][]T {
	var chunks [][]T

	for i := 0; i < len(arr); i += chunkSize {
		chunks = append(chunks, arr[i:min(i+chunkSize, len(arr))])
	}

	return chunks
}

func min(a int, b int) int {
	if a >= b {
		return b
	}

	return a
}
