package email

import (
	"fmt"
	"io/ioutil"
	"log"
	"multithreaded-data-processor/internal/entities"
	"multithreaded-data-processor/internal/sms"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func GetResultEmailData(path string, wg *sync.WaitGroup) (map[string][][]entities.EmailData, error) {
	out := make(chan map[string][][]entities.EmailData)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)

		result := make(map[string][][]entities.EmailData)
		data, err := emailDataReader(path)
		if err != nil {
			out <- nil
			return
		}
		countrySlice := []string{}
		for _, value := range data {
			countrySlice = append(countrySlice, value.Country)
		}
		countrySlice = removeDuplicates(countrySlice)
		for _, valueCountry := range countrySlice {
			var fastProviders []entities.EmailData
			var slowProviders []entities.EmailData
			for _, value := range data {
				if value.Country == valueCountry {
					fastProviders = append(fastProviders, value)
					slowProviders = append(slowProviders, value)
				}
			}
			country := sms.CountryFromAlpha(valueCountry)
			sort.Slice(fastProviders, func(i, j int) bool {
				return fastProviders[i].DeliveryTime < fastProviders[j].DeliveryTime
			})
			sort.Slice(slowProviders, func(i, j int) bool {
				return slowProviders[i].DeliveryTime > slowProviders[j].DeliveryTime
			})
			result[country] = [][]entities.EmailData{
				fastProviders[:3],
				slowProviders[:3],
			}
		}
		out <- result
	}()

	errorMsg := "Failed to retrieve result from channel"
	result, ok := <-out
	if !ok {
		return nil, fmt.Errorf(errorMsg)
	}
	return result, nil
}

func emailDataReader(path string) ([]entities.EmailData, error) {
	var result []entities.EmailData
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(reader), "\n")
	for _, value := range lines {
		splitVal := strings.Split(value, ";")
		if len(splitVal) == 3 {
			if checkEmailData(splitVal) {
				deliveryValue, _ := strconv.Atoi(splitVal[2])
				res := entities.EmailData{
					Country:      splitVal[0],
					Provider:     splitVal[1],
					DeliveryTime: deliveryValue,
				}
				result = append(result, res)
			}
		}
	}
	return result, nil
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, value := range slice {
		if _, ok := keys[value]; !ok {
			keys[value] = true
			list = append(list, value)
		}
	}
	return list
}

func checkEmailData(value []string) bool {
	if value[0] == sms.CountryAlpha2()[value[0]] {
		providers := map[string]string{"Gmail": "Gmail", "Yahoo": "Yahoo", "Hotmail": "Hotmail", "MSN": "MSN", "Orange": "Orange", "Comcast": "Comcast", "AOL": "AOL", "Live": "Live", "RediffMail": "RediffMail", "GMX": "GMX", "Protonmail": "Protonmail", "Yandex": "Yandex", "Mail.ru": "Mail.ru"}
		if value[1] == providers[value[1]] {
			_, err := strconv.Atoi(value[2])
			if err != nil {
				log.Printf("The delivery time value %v does not match the expected format.", value)
				return false
			} else {
				return true
			}
		} else {
			log.Printf("The provider value %v does not match the expected format.", value)
			return false
		}
	} else {
		log.Printf("The country alpha-2 value %v does not match the expected format.", value)
		return false
	}
}
