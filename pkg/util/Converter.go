package util

func BinaryConverter(number int, bits int) []int {
	factor := number
	result := make([]int, bits)

	for factor >= 0 && number > 0 {
		factor = number % 2
		number = number / 2
		result[bits-1] = factor
		bits--
	}

	return result
}