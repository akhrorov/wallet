package wallet

import (
	"github.com/akhrorov/wallet/pkg/types"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func BenchmarkService_SumPayments(b *testing.B) {
	service := &Service{}
	account, _, err := service.addAccount(defaultExampleTestAccount)
	_, _, err = service.addAccount(defaultExampleTestAccount2)
	if err != nil {
		b.Fatalf("SumPayments(): can't addAccount, %v", err)
		return
	}
	_,err = service.Pay(account.ID, 10_000, "food")
	_,err = service.Pay(account.ID, 30_000, "food")
	_,err = service.Pay(account.ID, 30_000, "food")
	_,err = service.Pay(account.ID, 30_000, "food")
	_,err = service.Pay(account.ID, 30_000, "food")
	_,err = service.Pay(account.ID, 30_000, "food")
	_,err = service.Pay(account.ID, 30_000, "food")
	_,err = service.Pay(account.ID, 30_000, "food")
	_,err = service.Pay(account.ID, 30_000, "food")
	_,err = service.Pay(account.ID, 30_000, "food")

	want := types.Money(480000)
	for i := 0; i < b.N; i++ {
		result := service.SumPayments(1)
		if result != want {
			b.Fatalf("Invalid result, got %v, want %v", result, want)
		}
	}
}

func TestService_Import(t *testing.T) {
	service := &Service{}
	account, payments, err := service.addAccount(defaultExampleTestAccount)
	_, _, err = service.addAccount(defaultExampleTestAccount2)
	if err != nil {
		t.Errorf("Import(): can't addAccount, %v", err)
		return
	}

	payment := payments[0]
	_, err = service.FavoritePayment(payment.ID, "something")
	if err != nil {
		t.Errorf("PayFromFavorite(): can't find favorite payment = %v", err)
		return
	}

	//err = service.Export("data")
	//if err != nil {
	//	t.Errorf("Import(): can't Import, %v", err)
	//}
	
	paymentss, err := service.ExportAccountHistory(account.ID)
	if err != nil {
		t.Errorf("Import(): can't Export Account History, %v", err)
	}

	err = service.HistoryToFiles(paymentss, "data", 4)
	if err != nil {
		t.Errorf("Import(): can't Export Account History, %v", err)
	}

	//err = service.Import("accounts.dump")
	//if err != nil {
	//	t.Errorf("Import(): can't find accounts file, %v", err)
	//}
}

func TestService_Export(t *testing.T) {
	service := &Service{}
	_, _, err := service.addAccount(defaultExampleTestAccount)
	_, _, err = service.addAccount(defaultExampleTestAccount2)
	if err != nil {
		t.Errorf("Export(): can't addAccount, %v", err)
		return
	}

	err = service.Export("data")
	if err != nil {
		t.Errorf("Export(): can't Export, %v", err)
	}
}

func TestService_ImportFromFile(t *testing.T) {
	service := &Service{}
	_, err := service.RegisterAccount("+992001010522")
	_, err = service.RegisterAccount("+992001010523")
	_, err = service.RegisterAccount("+992001010524")

	if err != nil {
		t.Errorf("ImportFromFile(): can't register account, %v", err)
	}

	err = service.ExportToFile("accounts.txt")
	if err != nil {
		t.Errorf("ImportFromFile(): can't export accounts, %v", err)
	}

	err = service.ImportFromFile("accounts.txt")
	if err != nil {
		t.Errorf("ImportFromFile(): can't import file")
	}
}

func TestService_ExportToFile(t *testing.T) {
	service := &Service{}

	_, err := service.RegisterAccount("+992001010522")
	_, err = service.RegisterAccount("+992001010523")
	_, err = service.RegisterAccount("+992001010524")
	if err != nil {
		t.Errorf("ExportToFile(): can't register account, %v", err)
	}

	err = service.ExportToFile("accounts.txt")
	if err != nil {
		t.Errorf("ExportToFile(): can't export accounts, %v", err)
	}
}

func TestService_RegisterAccount_success(t *testing.T) {
	service := &Service{}

	_,err := service.RegisterAccount("+992001010522")
	if err != nil {
		t.Errorf("RegisterAccount(): can't register account, %v", err)
		return
	}

}

func TestService_RegisterAccount_fail(t *testing.T) {
	service := &Service{}

	_,err := service.RegisterAccount("+992001010522")
	_,err = service.RegisterAccount("+992001010522")
	if err == nil {
		t.Errorf("RegisterAccount(): must return error, returned nil = %v", err)
		return
	}
	if err != ErrPhoneRegistered {
		t.Errorf("RegisterAccount(): must return ErrPhoneRegistered, returned nil")
	}
}

func TestService_FindAccountByID_success(t *testing.T) {
	service := &Service{}

	account, err := service.RegisterAccount("+992001010522")
	if err != nil {
		t.Errorf("FindAccountByID(): can't register account, %v", err)
		return
	}

	_, err = service.FindAccountByID(account.ID)
	if err != nil {
		t.Errorf("FindAccountByID(): can't find account, %v", err)
		return
	}
}


func TestService_FindAccountByID_fail(t *testing.T) {
	service := &Service{}

	_, err := service.RegisterAccount("+992001010522")
	if err != nil {
		t.Errorf("FindAccountByID(): can't register account, %v", err)
		return
	}

	_, err = service.FindAccountByID(2)
	if err == nil {
		t.Errorf("FindAccountByID(): can't find account, %v", err)
		return
	}
	if err != ErrAccountNotFound {
		t.Errorf("FindAccountByID(): must return ErrAccountNotFound, but returned:, %v", err)
		return
	}
}

func TestService_Deposit_success(t *testing.T) {
	service := &Service{}

	account, err := service.RegisterAccount("+992001010522")
	if err != nil {
		t.Errorf("Deposit(): can't register account, error = %v", err)
		return
	}

	err = service.Deposit(account.ID, 10_000)
	if err != nil {
		t.Errorf("Deposit(): can't deposit account, error = %v", err)
		return
	}
}

func TestService_Deposit_fail(t *testing.T) {
	service := &Service{}

	_, err := service.RegisterAccount("+992001010522")
	if err != nil {
		t.Errorf("Deposit(): can't register account, error = %v", err)
		return
	}

	err = service.Deposit(2, 10_000)
	if err == nil {
		t.Errorf("Deposit(): must return error, returned nil = %v", err)
		return
	}

	if err != ErrAccountNotFound {
		t.Errorf("Deposit(): must return ErrAccountNotFound, returned nil = %v", err)
		return
	}
}


func TestService_Pay_success(t *testing.T) {
	service := &Service{}

	account, err := service.RegisterAccount("+992001010522")
	service.Deposit(account.ID, 10_000)
	if err != nil {
		t.Errorf("RegisterAccount(): can't register account = %v", err)
		return
	}

	_, err = service.Pay(account.ID, 2000, "food")
	if err != nil {
		t.Errorf("Pay(): can't make payment = %v", err)
		return
	}
}

func TestService_Pay_fail_not_enough_balance(t *testing.T) {
	service := &Service{}

	account, err := service.RegisterAccount("+992001010522")
	if err != nil {
		t.Errorf("RegisterAccount(): can't register account = %v", err)
		return
	}

	_, err = service.Pay(account.ID, 100_000_000, "food")
	if err == nil {
		t.Errorf("Pay(): must return error, returned nil = %v", err)
		return
	}

	if err != ErrNotEnoughBalance {
		t.Errorf("Pay(): must return ErrNotEnoughBalance, returned nil = %v", err)
		return
	}
}

func TestService_Pay_fail(t *testing.T) {
	service := &Service{}

	_, err := service.RegisterAccount("+992001010522")
	if err != nil {
		t.Errorf("RegisterAccount(): can't register account = %v", err)
		return
	}

	_, err = service.Pay(2, 100_000_000, "food")
	if err == nil {
		t.Errorf("Pay(): must return error, returned nil = %v", err)
		return
	}

	if err != ErrAccountNotFound {
		t.Errorf("Pay(): must return ErrAccountNotFound, returned nil = %v", err)
		return
	}
}


func TestService_Reject_success(t *testing.T) {
	service := &Service{}

	account, payments, err := service.addAccount(defaultExampleTestAccount)
	payment := payments[0]
	if payment == nil {
		t.Errorf("Reject(): returned ErrPaymentNotFound ")
		return
	}
	err = service.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't reject payment, error = %v", err)
		return
	}

	savedPayment, err := service.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment, error = %v", err)
		return
	}

	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v", savedPayment)
		return
	}

	savedAccount, err := service.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account, error = %v", err)
		return
	}

	if savedAccount.Balance != account.Balance {
		t.Errorf("Reject(): balance didn't changed, account = %v", savedAccount)
		return
	}
}

func TestService_Reject_fail(t *testing.T) {
	service := &Service{}
	_, _, err := service.addAccount(defaultExampleTestAccount)
	if err != nil {
		t.Errorf("Reject(): can't register account, error = %v", err)
		return
	}

	_, err = service.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentByID(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	service := &Service{}

	_, payments, err := service.addAccount(defaultExampleTestAccount)
	if err != nil {
		t.Errorf("addAccount(): can't add account, error = %v", err)
		return
	}

	payment := payments[0]
	got, err := service.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): can't find payment, error = %v", err)
		return
	}

	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned, error = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	service := &Service{}

	_, _, err := service.addAccount(defaultExampleTestAccount)
	if err != nil {
		t.Errorf("addAccount(): can't add account, error = %v", err)
		return
	}

	_, err = service.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentByID(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}

func TestService_Repeat_success(t *testing.T) {
	service := &Service{}

	_, payments, err := service.addAccount(defaultExampleTestAccount)
	if err != nil {
		t.Errorf("Repeat(): can't register account = %v", err)
		return
	}

	payment := payments[0]
	_, err = service.Repeat(payment.ID)
	if err != nil {

		t.Errorf("Repeat(): can't repeat payment, error = %v", err)
		return
	}

}

func TestService_Repeat_fail(t *testing.T) {
	service := &Service{}

	_, _, err := service.addAccount(defaultExampleTestAccount)
	if err != nil {
		t.Errorf("Repeat(): can't register account = %v", err)
		return
	}

	_, err = service.Repeat(uuid.New().String())
	if err == nil {
		t.Errorf("Repeat(): must return error, returned nil = %v", err)
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("Repeat(): must return ErrPaymentNotFound, returned nil = %v", err)
		return
	}
}

func TestService_FavoritePayment_success(t *testing.T) {
	service := &Service{}

	_, payments, err := service.addAccount(defaultExampleTestAccount)
	if err != nil {
		t.Errorf("Repeat(): can't register account = %v", err)
		return
	}

	payment := payments[0]
	_,err = service.FavoritePayment(payment.ID, "something")
	if err != nil {
		t.Errorf("FavoritePayment(): can't find favorite payment = %v", err)
		return
	}
}

func TestService_FavoritePayment_fail(t *testing.T) {
	service := &Service{}

	_, _, err := service.addAccount(defaultExampleTestAccount)
	if err != nil {
		t.Errorf("Repeat(): can't register account = %v", err)
		return
	}

	_,err = service.FavoritePayment(uuid.New().String(), "something")
	if err == ErrFavoriteNotFound {
		t.Errorf("FavoritePayment(): must return ErrFavoriteNotFound = %v", err)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	service := &Service{}

	_, payments, err := service.addAccount(defaultExampleTestAccount)
	if err != nil {
		t.Errorf("PayFromFavorite(): can't register account = %v", err)
		return
	}

	payment := payments[0]
	getFavPayment, err := service.FavoritePayment(payment.ID, "something")
	if err != nil {
		t.Errorf("PayFromFavorite(): can't find favorite payment = %v", err)
		return
	}

	_,err = service.PayFromFavorite(getFavPayment.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): can't pay with favorite payment = %v", err)
		return
	}

}

func TestService_PayFromFavorite_fail(t *testing.T) {
	service := &Service{}

	_, payments, err := service.addAccount(defaultExampleTestAccount)
	if err != nil {
		t.Errorf("PayFromFavorite(): can't register account = %v", err)
		return
	}

	payment := payments[0]
	_, err = service.FavoritePayment(payment.ID, "something")
	if err != nil {
		t.Errorf("PayFromFavorite(): can't find favorite payment = %v", err)
		return
	}

	_,err = service.PayFromFavorite(uuid.New().String())
	if err != ErrFavoriteNotFound {
		t.Errorf("PayFromFavorite(): must return ErrFavoriteNotFound = %v", err)
		return
	}

}