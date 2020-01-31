package pointers_errors

import (
	"testing"
)

func TestWallet(t *testing.T) {
	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(10)
		got := wallet.Balance()
		assertBalance(t, got, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		err := wallet.Withdraw(8)
		got := wallet.Balance()
		assertNotError(err, t)
		assertBalance(t, got, Bitcoin(12))
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(8)
		wallet := Wallet{balance: startingBalance}
		err := wallet.Withdraw(20)
		got := wallet.Balance()
		assertError(err, ErrInsufficientFunds, t)
		assertBalance(t, got, startingBalance)
	})

}

func assertBalance(t *testing.T, got, want Bitcoin) {
	t.Helper()
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func assertError(got, want error, t *testing.T) {
	t.Helper()
	if got == nil {
		t.Fatal("wanted an error but didn't get one")
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertNotError(got error, t *testing.T) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}
