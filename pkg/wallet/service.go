package wallet

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/akhrorov/wallet/pkg/types"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

type testExampleAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultExampleTestAccount = testExampleAccount{
	phone:   "+992900000001",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

var defaultExampleTestAccount2 = testExampleAccount{
	phone:   "+992900000002",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func (s *Service) addAccount(data testExampleAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return nil, ErrAccountNotFound
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return ErrAccountNotFound
	}

	// зачисление средств пока не рассматриваем как платёж
	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return ErrPaymentNotFound
	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return ErrAccountNotFound
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	return s.Pay(payment.AccountID, payment.Amount, payment.Category)
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Name:      name,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}

	return nil, ErrFavoriteNotFound
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, ErrFavoriteNotFound
	}

	return s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
}

func (s *Service) ExportToFile(path string) error {

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	var str string
	for _, v := range s.accounts {
		str += fmt.Sprint(v.ID) + ";" + string(v.Phone) + ";" + fmt.Sprint(v.Balance) + "|"
	}
	_, err = file.WriteString(str)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	buf := make([]byte, 1)
	content := make([]byte, 0)
	for {
		read, err := file.Read(buf)

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		content = append(content, buf[:read]...)
	}
	data := strings.Split(string(content), "|")
	var row []string
	for _, v := range data {
		row = append(row, strings.ReplaceAll(v, ";", " "))
	}
	for _, acc := range row {
		if len(acc) == 0 {
			continue
		}

		accountSplit := strings.Split(acc, " ")

		id, err := strconv.ParseInt(accountSplit[0], 10, 64)
		if err != nil {
			return err
		}
		balance, err := strconv.ParseInt(accountSplit[2], 10, 64)
		if err != nil {
			return err
		}
		account := &types.Account{ID: id, Balance: types.Money(balance), Phone: types.Phone(accountSplit[1])}
		s.accounts = append(s.accounts, account)

	}
	return nil
}

func (s *Service) Export(dir string) error {

	if len(s.accounts) > 0 {
		var accountItem string
		for _, account := range s.accounts {
			accountItem += fmt.Sprint(account.ID) + ";" + string(account.Phone) + ";" + fmt.Sprint(account.Balance) + "\n"
		}

		file, err := os.Create(dir + "/accounts.dump")
		if err != nil {
			log.Print(err)
			return err
		}
		defer func() {
			if cerr := file.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		_, err = file.WriteString(accountItem)
	}
	if len(s.payments) > 0 {
		var paymentItem string
		for _, payment := range s.payments {
			paymentItem += fmt.Sprint(payment.ID) + ";" + fmt.Sprint(payment.Amount) + ";" + fmt.Sprint(payment.Category) + ";" + fmt.Sprint(payment.AccountID) + ";" + fmt.Sprint(payment.Status) + "\n"
		}

		file, err := os.Create(dir + "/payments.dump")
		if err != nil {
			log.Print(err)
			return err
		}
		defer func() {
			if cerr := file.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		_, err = file.WriteString(paymentItem)
	}
	if len(s.favorites) > 0 {
		var favoriteItem string
		for _, favorite := range s.favorites {
			favoriteItem += fmt.Sprint(favorite.ID) + ";" + fmt.Sprint(favorite.Amount) + ";" + fmt.Sprint(favorite.Category) + ";" + fmt.Sprint(favorite.AccountID) + ";" + fmt.Sprint(favorite.Name) + "\n"
		}

		file, err := os.Create(dir + "/favorites.dump")
		if err != nil {
			log.Print(err)
			return err
		}
		defer func() {
			if cerr := file.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		_, err = file.WriteString(favoriteItem)
	}
	return nil
}

func (s *Service) Import(dir string) error {

	_, err := os.Stat(dir + "/accounts.dump")
	if err == nil {
		file, err := os.Open(dir + "/accounts.dump")
		if err != nil {
			return err
		}
		defer func() {
			if cerr := file.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')

			if err == io.EOF {
				log.Print(line)
				break
			}

			if err != nil {
				log.Print(err)
				return err
			}
			line = strings.ReplaceAll(line, "\n", "")
			splitedItem := strings.Split(line, ";")
			id, err := strconv.ParseInt(splitedItem[0], 10, 64)
			if err != nil {
				return err
			}
			balance, err := strconv.ParseInt(splitedItem[2], 10, 64)
			if err != nil {
				return err
			}
			account := &types.Account{ID: id, Phone: types.Phone(splitedItem[1]), Balance: types.Money(balance)}
			_, err = s.FindAccountByID(account.ID)
			if err == ErrAccountNotFound {
				s.accounts = append(s.accounts, account)
			}
			s.nextAccountID = int64(len(s.accounts))
		}
	}
	_, err = os.Stat(dir + "/payments.dump")
	if err == nil {
		file, err := os.Open(dir + "/payments.dump")
		if err != nil {
			return err
		}
		defer func() {
			if cerr := file.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')

			if err == io.EOF {
				log.Print(line)
				break
			}

			if err != nil {
				log.Print(err)
				return err
			}
			line = strings.ReplaceAll(line, "\n", "")
			splitedItem := strings.Split(line, ";")
			amount, err := strconv.ParseInt(splitedItem[1], 10, 64)
			if err != nil {
				return err
			}
			accountID, err := strconv.ParseInt(splitedItem[3], 10, 64)
			if err != nil {
				return err
			}
			payment := &types.Payment{ID: splitedItem[0], Amount: types.Money(amount), Category: types.PaymentCategory(splitedItem[2]), AccountID: accountID, Status: types.PaymentStatus(splitedItem[4])}
			_, err = s.FindPaymentByID(payment.ID)
			if err == ErrPaymentNotFound {
				s.payments = append(s.payments, payment)
			}
		}
	}
	_, err = os.Stat(dir + "/favorites.dump")
	if err == nil {
		file, err := os.Open(dir + "/favorites.dump")
		if err != nil {
			return err
		}
		defer func() {
			if cerr := file.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')

			if err == io.EOF {
				log.Print(line)
				break
			}

			if err != nil {
				log.Print(err)
				return err
			}
			line = strings.ReplaceAll(line, "\n", "")
			splitedItem := strings.Split(line, ";")

			amount, err := strconv.ParseInt(splitedItem[1], 10, 64)
			if err != nil {
				return err
			}
			accountID, err := strconv.ParseInt(splitedItem[3], 10, 64)
			if err != nil {
				return err
			}
			favorite := &types.Favorite{ID: splitedItem[0], Amount: types.Money(amount), Category: types.PaymentCategory(splitedItem[2]), AccountID: accountID, Name: splitedItem[4]}
			_, err = s.FindFavoriteByID(favorite.ID)
			if err == ErrFavoriteNotFound {
				s.favorites = append(s.favorites, favorite)
			}
		}
	}
	return nil
}

func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	findedPayments := []types.Payment{}

	for _, payment := range s.payments {
		if payment.AccountID == account.ID {
			findedPayments = append(findedPayments, types.Payment{
				ID:        payment.ID,
				AccountID: payment.AccountID,
				Status:    payment.Status,
				Category:  payment.Category,
				Amount:    payment.Amount,
			})
		}
	}

	return findedPayments, nil
}

func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	if len(payments) > 0 {
		if records >= len(payments) {
			var paymentItem string

			for _, payment := range payments {
				paymentItem += fmt.Sprint(payment.ID) + ";" + fmt.Sprint(payment.Amount) + ";" + fmt.Sprint(payment.Category) + ";" + fmt.Sprint(payment.AccountID) + ";" + fmt.Sprint(payment.Status) + "\n"
			}

			file, err := os.Create(dir + "/payments.dump")
			if err != nil {
				log.Print(err)
				return err
			}

			defer func() {
				if cerr := file.Close(); cerr != nil {
					log.Print(cerr)
				}
			}()

			_, err = file.WriteString(paymentItem)
		} else {
			var paymentItem string
			counter := 1
			counterForPayments := 0
			for _, payment := range payments {
				paymentItem += fmt.Sprint(payment.ID) + ";" + fmt.Sprint(payment.Amount) + ";" + fmt.Sprint(payment.Category) + ";" + fmt.Sprint(payment.AccountID) + ";" + fmt.Sprint(payment.Status) + "\n"

				file, err := os.Create(dir + "/payments" + fmt.Sprint(counter) + ".dump")
				if err != nil {
					log.Print(err)
					return err
				}

				defer func() {
					if cerr := file.Close(); cerr != nil {
						log.Print(cerr)
					}
				}()

				_, err = file.WriteString(paymentItem)
				counterForPayments++
				if counterForPayments >= records {
					counter++
					counterForPayments = 0
					paymentItem = ""
				}
			}
		}
	}
	return nil
}
func (s *Service) SumPayments(goroutines int) types.Money {

	mu := sync.Mutex{}
	sum := types.Money(0)

	if goroutines == 0 || goroutines == 1 {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			val := types.Money(0)
			for _, payment := range s.payments {
				val += payment.Amount
			}
			mu.Lock()
			defer mu.Unlock()
			sum += val
		}()
		wg.Wait()
		return sum
	}
	wg := sync.WaitGroup{}

	for _, vp := range s.payments {
		wg.Add(1)
		go func(vp *types.Payment) {
			defer wg.Done()
			val := types.Money(0)

			val += vp.Amount

			mu.Lock()
			defer mu.Unlock()
			sum += val
		}(vp)
	}

	wg.Wait()
	return sum
}

func (s *Service) FilterPaymentsForGoroutines(goroutinesCount int, accountID int64) ([][]types.Payment, error) {
	_, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	pm := []types.Payment{}

	for _, p := range s.payments {

		if p.AccountID == accountID {

			pm = append(pm, *p)

		}
	}

	grouped := [][]types.Payment{}

	for i := 0; i < len(pm); i++ {

		if i+goroutinesCount > len(pm)-1 {

			grouped = append(grouped, pm[i:])

			break
		}

		grouped = append(grouped, pm[i:i+goroutinesCount])

		i += goroutinesCount - 1
	}

	return grouped, nil
}
func (s *Service) FilterPaymentsForG(goroutinesCount int) ([][]types.Payment, error) {

	pm := []types.Payment{}

	for _, p := range s.payments {

		pm = append(pm, *p)

	}

	grouped := [][]types.Payment{}

	for i := 0; i < len(pm); i++ {

		if i+goroutinesCount > len(pm)-1 {

			grouped = append(grouped, pm[i:])

			break
		}

		grouped = append(grouped, pm[i:i+goroutinesCount])

		i += goroutinesCount - 1
	}

	return grouped, nil
}

func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {
	if goroutines == 0 {
		mu := sync.Mutex{}
		payments := []types.Payment{}

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			val := []types.Payment{}
			for _, payment := range s.payments {
				if accountID == payment.AccountID {
					val = append(payments, *payment)
				}
			}
			mu.Lock()
			defer mu.Unlock()
			payments = append(payments, val...)

		}()

		wg.Wait()

		if len(payments) == 0 {

			return nil, ErrAccountNotFound
		}

		return payments, nil
	}

	wg := sync.WaitGroup{}

	mu := sync.Mutex{}
	payments := []types.Payment{}

	filteredPayments, err := s.FilterPaymentsForGoroutines(goroutines, accountID)
	if err != nil {
		return nil, err
	}
	if len(filteredPayments) == 0 {
		return nil, nil
	}
	for _, fp := range filteredPayments {
		wg.Add(1)
		go func(fp []types.Payment) {
			defer wg.Done()
			mu.Lock()
			payments = append(payments, fp...)
			defer mu.Unlock()
		}(fp)
	}

	wg.Wait()
	if len(payments) == 0 {
		return nil, nil
	}

	return payments, nil

}

func (s *Service) FilterPaymentsByFn(filter func(payment types.Payment) bool, goroutines int) ([]types.Payment, error) {
	if goroutines == 0 {
		mu := sync.Mutex{}
		payments := []types.Payment{}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			val := []types.Payment{}
			for _, payment := range s.payments {
				if filter(*payment) {
					val = append(payments, *payment)
				}
			}
			mu.Lock()
			defer mu.Unlock()
			payments = append(payments, val...)

		}()
		wg.Wait()
		if len(payments) == 0 {
			return nil, ErrAccountNotFound
		}
		return payments, nil
	}
	goroutinesCount := goroutines
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	payments := []types.Payment{}
	pm := []types.Payment{}
	for _, p := range s.payments {
		if filter(*p) {
			pm = append(pm, *p)
		}
	}
	grouped := [][]types.Payment{}
	for i := 0; i < len(pm); i++ {
		if i+goroutinesCount > len(pm)-1 {
			grouped = append(grouped, pm[i:])
			break
		}
		grouped = append(grouped, pm[i:i+goroutinesCount])
		i += goroutinesCount - 1
	}
	if len(grouped) == 0 {
		return nil, nil
	}
	for _, fp := range grouped {
		wg.Add(1)
		go func(fp []types.Payment) {
			defer wg.Done()
			mu.Lock()
			payments = append(payments, fp...)
			defer mu.Unlock()
		}(fp)
	}
	wg.Wait()
	if len(payments) == 0 {
		return nil, nil
	}
	return payments, nil
}

//SumPaymentsWithProgress ...
func (s *Service) SumPaymentsWithProgress() <-chan types.Progress {

	ch := make(chan types.Progress)
	size := 100_000
	parts := len(s.payments) / size
	wg := sync.WaitGroup{}
	// Если части меньше чем ноль то части = 1
	if parts < 1 {
		parts = 1
	}

	for i := 0; i < parts; i++ {
		wg.Add(1)
		var payments []*types.Payment

		if len(s.payments) < size {
			// если если платежей меньше чем
			// 100_000 то создаётся одна горутина и ей отдаём всё
			payments = s.payments
		} else {
			//  если если платежей больше чем чем 100_000 то
			//  отдаём платежи по равным частям, (последнему может достатся меньше)
			payments = s.payments[i*size : (i+1)*size]
		}

		go func(ch chan types.Progress, data []*types.Payment) {
			defer wg.Done()
			val := types.Money(0)
			//Сумируем то что передано по параметам в одну переменную
			for _, v := range data {
				val += v.Amount
			}
			// Записываем в канал
			ch <- types.Progress{
				Part:   len(data),
				Result: val,
			}

		}(ch, payments)
	}
	//Если платежей больше чем size
	if len(s.payments) > size {
		wg.Add(1)
		payments := s.payments
		go func(ch chan types.Progress, data []*types.Payment) {
			defer wg.Done()
			val := types.Money(0)
			for _, v := range data {
				val += v.Amount
			}
			ch <- types.Progress{
				Part:   len(data),
				Result: val,
			}

		}(ch, payments)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}
