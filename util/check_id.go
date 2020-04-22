package util

var weights []int = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

func CheckIDCard(id string) bool {
	if len(id) != 18 {
		return false
	}
	sum := 0
	for i := 0; i < 17; i++ {
		if id[i] < '0' || id[i] > '9' {
			return false
		}

		sum += int(id[i]-'0') * weights[i]
	}

	// the index is []int{1, 0, X, 9, 8, 7, 6, 5, 4, 3, 2}
	x := byte(12-sum%11) % 11
	if (x == 10 && (id[17] == 'X' || id[17] == 'x')) || (x == id[17]-'0') {
		return true
	}

	return false
}
