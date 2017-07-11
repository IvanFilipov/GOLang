package main

import "fmt"

func SquareSumDifference(n uint64) uint64 {

	var sum, SquareSum, i uint64 = 0, 0, 1

	for ; i <= n; i++ {

		sum += i
		SquareSum += i * i

	}

	sum *= sum

	return sum - SquareSum

}

func main() {

	var n uint64

	fmt.Scanf("%d", &n)

	fmt.Printf("%d\n", SquareSumDifference(n))

}
