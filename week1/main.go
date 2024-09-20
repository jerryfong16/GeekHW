package main

import "errors"

func main() {
}

func removeIndexInSlice[T any](s1 []T, i int) ([]T, T, error) {
	var deletedEle T
	if i < 0 || i > len(s1) || len(s1) <= 1 {
		return s1, deletedEle, errors.New("invalid index")
	}
	deletedEle = s1[i]
	for j := i; j < len(s1)-1; j++ {
		s1[j] = s1[j+1]
	}
	s1 = s1[:len(s1)-1]
	if cap(s1) > 2048 && cap(s1)/len(s1) > 2 {
		s2 := make([]T, 0, cap(s1)/2)
		for j := 0; j < len(s1); j++ {
			s2[j] = s1[j]
		}
		return s2, deletedEle, nil
	}
	return s1, deletedEle, nil
}
