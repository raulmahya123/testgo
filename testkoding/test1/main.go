package main

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type Wallet struct {
	PersonName string
	Credit     float64
	mu         sync.Mutex // Mutex untuk mencegah data race
}

// Fungsi Withdrawal untuk menarik kredit dari wallet
func (w *Wallet) Withdrawal(wdAmount float64) error {
	w.mu.Lock() // Mengunci akses ke wallet untuk mencegah data race
	defer w.mu.Unlock()

	if wdAmount < minWd || wdAmount > maxWd {
		return errors.New("withdrawal amount out of range")
	}
	if w.Credit-wdAmount < 0 {
		return errors.New("insufficient funds")
	}

	w.Credit -= wdAmount // Mengurangi kredit sesuai jumlah penarikan
	return nil
}

// Fungsi GetWallet untuk mendapatkan informasi wallet
func (w *Wallet) GetWallet() Wallet {
	w.mu.Lock() // Mengunci akses ke wallet untuk mencegah data race
	defer w.mu.Unlock()

	return *w
}

// ===== DO NOT EDIT. =====
const minWd = 1
const maxWd = 20

var wallet = &Wallet{
	PersonName: "John Doe",
	Credit:     maxWd,
}

func TestCaseWallet(t *testing.T) {
	iteration := 20
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())

	for i := 0; i <= iteration; i++ {
		wdAmount := float64(rand.Intn((maxWd-minWd)/0.05))*0.05 + minWd

		wg.Add(1)
		go func(x int) {
			remaining, err := atm(&wg, x, wdAmount, wallet)
			if err == nil {
				log.Printf("Withdraw Amount: %.2f, Remaining Credit: %.2f", wdAmount, remaining)
			}
		}(i)
	}

	wg.Wait()
	log.Println("+------------+")
	log.Printf("%s's final credit: %.2f", wallet.PersonName, wallet.Credit)

	switch {
	case wallet.Credit < 0:
		t.Fail()

	case wallet.Credit >= maxWd:
		t.Fail()
	}
}

func atm(
	wg *sync.WaitGroup,
	c int,
	wdAmount float64,
	wal *Wallet,
) (float64, error) {
	defer wg.Done()

	err := wal.Withdrawal(wdAmount)
	w := wal.GetWallet()
	if err != nil {
		return w.Credit, err
	}

	return w.Credit, nil
}

func main() {
	testSuite := []testing.InternalTest{
		{
			Name: "TestCaseWallet",
			F:    TestCaseWallet,
		},
	}

	testing.Main(nil, testSuite, nil, nil)
}

// ===== DO NOT EDIT. =====
