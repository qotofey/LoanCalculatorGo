package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	amount     float64
	annualRate float64
	loanDate   time.Time
)

const TERM = 3

func main() {
	annualRate = 0.55
	fmt.Println("Добро пожаловать в МФО! Мы выдадим вам займ на следующих условиях:")
	fmt.Println(" - сроком только на 3 месяца;")
	fmt.Println(" - под 55% годовых.")
	fmt.Print("Введите сумму займа: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	input = strings.TrimSpace(input)
	amount, err = strconv.ParseFloat(input, 64)
	if err != nil {
		log.Fatal(err)
	}

	loanDate = time.Date(2022, time.Month(12), 17, 0, 0, 0, 0, time.UTC)
	fmt.Println("Дата оформления займа:", loanDate.Format("2006-01-02"))

	payment := math.Floor(monthlyPayment()*100) / 100

	fmt.Println("Ежемесячный платёж составит", payment, "₽")
	fmt.Println("График платежей")
	fmt.Println("----------------------------------------------------------------------------------------------------------")
	fmt.Println(" # | дата платежа | дней в периоде | суммарый платёж | уйдёт в процент | в основную часть | остаток долга ")
	fmt.Println("----------------------------------------------------------------------------------------------------------")

	var totalPayment float64
	var totalInterest float64
	var totalPrincipal float64
	balance := amount
	var interest float64
	var principal float64
	paymentDate := loanDate
	for i := 0; i < TERM; i++ {
		paymentDate = paymentDate.AddDate(0, 1, 0)
		interest = math.Ceil(calculateInterest(paymentDate, balance*100)) / 100
		if i == TERM-1 {
			principal = balance
			balance = 0.0

			payment = principal + interest
		} else {
			principal = payment - interest
			balance -= principal
		}
		// balanceOwed := 11.03 if amount == 35_000.0

		totalInterest += interest
		totalPrincipal += principal
		totalPayment += payment

		fmt.Printf("%2d | %12s | %14d | %13.2f ₽ | %13.2f ₽ | %14.2f ₽ | %11.2f ₽ \n", i+1, paymentDate.Format("2006-01-02"), daysPerMonth(paymentDate), payment, interest, principal, balance)
	}
	fmt.Println("----------------------------------------------------------------------------------------------------------")
	fmt.Printf("%36s %13.2f ₽ | %13.2f ₽ | %14.2f ₽ \n", "", totalPayment, totalInterest, totalPrincipal)
}

func monthlyRate() float64 {
	return annualRate / 12.0
}

func annuityRatio() float64 {
	x := math.Pow(1+monthlyRate(), float64(TERM))
	return monthlyRate() * x / (x - 1)
}

func daysPerMonth(paymentDate time.Time) int {
	lastDay := time.Date(paymentDate.Year(), paymentDate.Month(), 0, 0, 0, 0, 0, paymentDate.Location()).Day()
	return lastDay
}

func daysPerYear() float64 {
	return 365.0
}

func monthlyRatio(paymentDate time.Time) float64 {
	return float64(daysPerMonth(paymentDate)) / daysPerYear()
}

func monthlyPayment() float64 {
	return amount * annuityRatio()
}

func calculateInterest(paymentDate time.Time, balance float64) float64 {
	return balance * annualRate * monthlyRatio(paymentDate)
}
