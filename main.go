package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

type LoanResult struct {
	PrincipalCents    int64
	MonthlyPayment    int64  // in cents
	TotalPayment      int64  // in cents
	TotalInterest     int64  // in cents
	MonthlyPaymentStr string // formatted "$1,234.56"
	TotalPaymentStr   string
	TotalInterestStr  string
}

func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/calculate", calculateHandler)

	fmt.Println("Starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	amountStr := r.FormValue("amount")
	rateStr := r.FormValue("rate")
	yearsStr := r.FormValue("years")

	// Input validation
	amount, errA := strconv.ParseFloat(amountStr, 64)
	rate, errR := strconv.ParseFloat(rateStr, 64)
	years, errY := strconv.Atoi(yearsStr)

	if errA != nil || errR != nil || errY != nil || amount <= 0 || rate <= 0 || years <= 0 {
		http.Error(w, "Invalid input. Please enter positive numbers for amount, rate, and years.", http.StatusBadRequest)
		return
	}

	// Integer calcualations - everything in cents
	principalCents := int64(math.Round(amount * 100))

	// Monthly interest rate in basis points (1bp = 0.01%)
	monthlyRateBP := int64(math.Round(rate * 10000 / 12)) // e.g. 5.25% → 4375 bp monthly
	months := int64(years * 12)

	// Standard loan payment formula using integer math (approximation with scaling)
	// M = P × (r × (1+r)^n) / ((1+r)^n - 1)
	// We use fixed-point scaling with 1_000_000 (6 decimals) for precision

	scale := int64(1_000_000)
	rScaled := monthlyRateBP * scale / 10000 // monthly rate with 6 decimals
	// (1 + r)^n
	power := int64(1)
	base := scale + rScaled
	for i := int64(0); i < months; i++ {
		power = power * base / scale
		if power == 0 { // overflow/underflow protection
			power = 1
			break
		}
	}
	numerator := rScaled * power
	denominator := power - scale

	monthlyPaymentScaled := int64(0)
	if denominator != 0 {
		monthlyPaymentScaled = (principalCents * numerator / scale) / denominator
	}

	// Round to nearest cent (banker's rounding via +0.5 then floor)
	monthlyPaymentCents := (monthlyPaymentScaled + 50) / 100

	totalPaymentCents := monthlyPaymentCents * months
	totalInterestCents := totalPaymentCents - principalCents

	// Format for display
	formatCents := func(c int64) string {
		dollars := c / 100
		cents := c % 100
		if cents < 0 {
			cents = -cents
		}
		return fmt.Sprintf("$%d.%02d", dollars, cents)
	}
	result := LoanResult{
		PrincipalCents:    principalCents,
		MonthlyPayment:    monthlyPaymentCents,
		TotalPayment:      totalPaymentCents,
		TotalInterest:     totalInterestCents,
		MonthlyPaymentStr: formatCents(monthlyPaymentCents),
		TotalPaymentStr:   formatCents(totalPaymentCents),
		TotalInterestStr:  formatCents(totalInterestCents),
	}

	tmpl := template.Must(template.ParseFiles("templates/result.html"))
	tmpl.Execute(w, result)
}
