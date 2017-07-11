package main

import "fmt"

var data = []int{5, 11, 23, 2, 15}

func print() {

	for i := 0; i < len(data); i++ {

		fmt.Printf("%d/", data[i])

	}

	fmt.Printf("\n")
}

func swap(a, b *int) {

	temp := *a
	*a = *b
	*b = temp

}

func InsertionSort() {

	for i := 1; i < len(data); i++ {
		for j := i; j > 0; j-- {
			if data[j] < data[j-1] {

				swap(&data[j], &data[j-1])

			}

		}
	}
}

func main() {

	InsertionSort()

	print()
}
