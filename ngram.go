package indexr

type Ngram uint32

// getNgrams returns a list of n-grams from a given string
func getNgrams(str string, n int) []Ngram {
	if len(str) == 0 {
		return nil
	}

	var result []Ngram
	for i := 0; i < len(str)-(n-1); i++ {
		var ngram Ngram
		for j := 0; j < n; j++ {
			shift := uint(j * 8)
			ngram = ngram + Ngram(uint32(str[i+(n-1-j)])<<shift)
		}
		result = append(result, ngram)
	}
	return result
}
