package main

import "fmt"

func main() {
	var months int = 60
	var rate = 550
	var principal = 2500000

	var displayPrincipal = principal / 100
	var displayRate = rate / 100

	fmt.Println("Smelter Loancalc Pro XT 3000")
	fmt.Println("-------------------------------")
	fmt.Printf("Principal: $%d\n", displayPrincipal)
	fmt.Printf("Rate: %d%%\n", displayRate)
	fmt.Printf("Months: %d\n", months)
	fmt.Println("-------------------------------")
	for i := 1; i <= months; i++ {
		var interest = principal * rate / 10000 / 12
		var payment = principal / (months - i + 1)
		principal -= payment
		var displayPayment = payment / 100
		var displayInterest = interest / 100
		fmt.Printf("Month %2d: Payment: $%6d Interest: $%5d Remaining Principal: $%7d\n", i, displayPayment, displayInterest, principal/100)
	}

}
