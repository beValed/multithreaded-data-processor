package mms

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"multithreaded-data-processor/internal/entities"
	"multithreaded-data-processor/internal/sms"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/jinzhu/copier"
)

func GetResultMMSData(wg *sync.WaitGroup) ([][]entities.MMSData, error) {
	out := make(chan [][]entities.MMSData)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)
		first, err := mmsRequest()
		if err != nil {
			log.Println(err)
			out <- nil
			return
		}
		second := []entities.MMSData{}
		err = copier.Copy(&second, &first)
		if err != nil {
			log.Print(err)
		}

		sort.Slice(first, func(i, j int) bool {
			return first[i].Provider < first[j].Provider
		})

		sort.Slice(second, func(i, j int) bool {
			return second[i].Country < second[j].Country
		})
		result := [][]entities.MMSData{
			first, second,
		}
		out <- result
	}()
	var result = <-out
	return result, nil
}

func mmsRequest() ([]entities.MMSData, error) {
	var result []entities.MMSData
	resp, err := http.Get("http://127.0.0.1:8383/mms")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		i := 0
		for _, value := range result {
			if checkMMS(value) {
				value.Country = sms.CountryFromAlpha(value.Country)
				result[i] = value
				i++
			}
		}
		result = result[:i]
		return result, nil
	} else {
		return nil, errors.New("Failed to retrieve MMS data")
	}
}

func checkMMS(value entities.MMSData) bool {
	if value.Country == sms.CountryAlpha2()[value.Country] {
		percentValue, err := strconv.Atoi(value.Bandwidth)
		if err != nil {
			log.Printf("The channel bandwidth value %v does not match the expected value.", value)
			return false
		}
		if -1 < percentValue && percentValue < 101 {
			_, err := strconv.Atoi(value.ResponseTime)
			if err == nil {
				providers := map[string]string{"Topolo": "Topolo", "Rond": "Rond", "Kildy": "Kildy"}
				if value.Provider == providers[value.Provider] {
					return true
				} else {
					log.Printf("The provider's %v value does not match what is expected.", value)
					return false
				}
			} else {
				log.Printf("The response value in ms %v does not match the expected value.", value)
				return false
			}
		} else {
			log.Printf("The channel bandwidth value %v does not match the expected value.", value)
			return false
		}
	} else {
		log.Printf("The country value alpha-2 %v does not match the expected value.", value)
		return false
	}
}
