package main

import "fmt"

//prints all prime divisors of a num and returns their count
func PrimeDivisors(n int) int {

	cnt := 0

	for i := 2; i <= n && n != 1; i++ {

		for ; n%i == 0; cnt++ {
			fmt.Printf("/%d", i)
			n /= i
		}

	}

	fmt.Printf("\n")
	return cnt
}

func main() {

	var num int

	fmt.Printf("enter a num :")
	fmt.Scanf("%d", &num)

	fmt.Printf("divisors count is %d\n", PrimeDivisors(num))

}
