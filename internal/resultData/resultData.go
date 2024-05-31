package resultData

import (
	"multithreaded-data-processor/internal/billing"
	"multithreaded-data-processor/internal/email"
	"multithreaded-data-processor/internal/entities"
	"multithreaded-data-processor/internal/incident"
	"multithreaded-data-processor/internal/mms"
	"multithreaded-data-processor/internal/sms"
	"multithreaded-data-processor/internal/support"
	"multithreaded-data-processor/internal/voice"
	"reflect"
	"sync"
	"time"
)

type ResultDataStorage struct {
	Time    time.Time
	Storage entities.ResultSetT
	sync.Mutex
}

func NewStorage() *ResultDataStorage {
	return &ResultDataStorage{
		Time: time.Now().Add(-31 * time.Second),
	}
}

func (r *ResultDataStorage) GetResultData() (entities.ResultSetT, error) {
	r.Lock()
	defer r.Unlock()
	t := time.Now()
	difference := t.Sub(r.Time)
	if difference > time.Second*30 {
		var wg sync.WaitGroup

		sms, err := sms.GetResultSMSData("../simulator/sms.data", &wg)
		if err != nil {
			return entities.ResultSetT{}, err
		}

		mms, err := mms.GetResultMMSData(&wg)
		if err != nil {
			return entities.ResultSetT{}, err
		}

		voice, err := voice.VoiceCallReader("../simulator/voice.data", &wg)
		if err != nil {
			return entities.ResultSetT{}, err
		}

		billing, err := billing.BillingDataReader("../simulator/billing.data", &wg)
		if err != nil {
			return entities.ResultSetT{}, err
		}

		support := support.GetResultSupportData(&wg)

		incident := incident.GetResultIncidentData(&wg)

		email, err := email.GetResultEmailData("../simulator/email.data", &wg)
		if err != nil {
			return entities.ResultSetT{}, err
		}

		wg.Wait()

		result := entities.ResultSetT{
			SMS:       sms,
			MMS:       mms,
			VoiceCall: voice,
			Email:     email,
			Billing:   billing,
			Support:   support,
			Incidents: incident,
		}
		r.Storage = result
		r.Time = time.Now()
		return r.Storage, nil
	} else {
		return r.Storage, nil
	}
}

func (r *ResultDataStorage) IsFull() bool {
	r.Lock()
	defer r.Unlock()
	v := reflect.ValueOf(r.Storage)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsZero() == true {
			return false
		}
	}
	return true
}
