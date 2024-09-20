package main

func main() {
}

func removeIndexInSlice[T any](s1 []T, i int) []T {
	if i < 0 || i > len(s1) || len(s1) <= 1 {
		return s1
	}
	s2 := make([]T, len(s1)-1)
	for j := 0; j < len(s1); j++ {
		if i == j {
			continue
		}
		s2[j] = s1[j]
	}
	return s2
}
