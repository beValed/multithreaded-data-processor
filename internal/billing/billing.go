package billing

import (
	"io/ioutil"
	"multithreaded-data-processor/internal/entities"
	"os"
	"strconv"
	"sync"
)

const (
	CreateCustomerMask int64 = 1 << iota
	PurchaseMask
	PayoutMask
	RecurringMask
	FraudControlMask
	CheckoutPageMask
)

func BillingDataReader(path string, wg *sync.WaitGroup) (entities.BillingData, error) {
	out := make(chan entities.BillingData)
	errChan := make(chan error)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)
		defer close(errChan)

		file, err := os.Open(path)
		if err != nil {
			errChan <- err
			return
		}
		defer file.Close()

		reader, err := ioutil.ReadAll(file)
		if err != nil {
			errChan <- err
			return
		}

		mask, err := strconv.ParseInt(string(reader), 2, 0)
		if err != nil {
			errChan <- err
			return
		}

		result := entities.BillingData{
			CreateCustomer: mask&CreateCustomerMask != 0,
			Purchase:       mask&PurchaseMask != 0,
			Payout:         mask&PayoutMask != 0,
			Recurring:      mask&RecurringMask != 0,
			FraudControl:   mask&FraudControlMask != 0,
			CheckoutPage:   mask&CheckoutPageMask != 0,
		}
		out <- result
	}()

	select {
	case result := <-out:
		return result, nil
	case err := <-errChan:
		return entities.BillingData{}, err
	}
}
