package voice

import (
	"io/ioutil"
	"log"
	"multithreaded-data-processor/internal/entities"
	"multithreaded-data-processor/internal/sms"
	"os"
	"strconv"
	"strings"
	"sync"
)

func VoiceCallReader(path string, wg *sync.WaitGroup) ([]entities.VoiceCallData, error) {
	out := make(chan []entities.VoiceCallData)
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

		var resultIn []entities.VoiceCallData
		lines := strings.Split(string(reader), "\n")
		for _, value := range lines {
			splitVal := strings.Split(value, ";")
			if len(splitVal) == 8 {
				if checkVoiceCall(splitVal) {
					percentValue, _ := strconv.Atoi(splitVal[1])
					responseTime, _ := strconv.Atoi(splitVal[2])
					connVal, _ := strconv.ParseFloat(splitVal[4], 32)
					connVal32 := float32(connVal)
					ttfbVal, _ := strconv.Atoi(splitVal[5])
					voicePurVal, _ := strconv.Atoi(splitVal[6])
					medianVoicVal, _ := strconv.Atoi(splitVal[7])
					res := entities.VoiceCallData{
						Country:             splitVal[0],
						Bandwidth:           percentValue,
						ResponseTime:        responseTime,
						Provider:            splitVal[3],
						ConnectionStability: connVal32,
						TTFB:                ttfbVal,
						VoicePurity:         voicePurVal,
						MedianOfCallsTime:   medianVoicVal,
					}
					resultIn = append(resultIn, res)
				}
			}
		}
		out <- resultIn
	}()
	select {
	case result := <-out:
		return result, nil
	case err := <-errChan:
		return nil, err
	}
}

func checkVoiceCall(value []string) bool {
	if value[0] == sms.CountryAlpha2()[value[0]] {
		percentValue, err := strconv.Atoi(value[1])
		if err != nil {
			log.Printf("The channel bandwidth value %v does not match the expected value.", value)
			return false
		}
		if -1 < percentValue && percentValue < 101 {
			_, err := strconv.Atoi(value[2])
			if err == nil {
				providers := map[string]string{"TransparentCalls": "TransparentCalls", "E-Voice": "E-Voice", "JustPhone": "JustPhone"}
				if value[3] == providers[value[3]] {
					_, err := strconv.ParseFloat(value[4], 32)
					if err != nil {
						log.Printf("The stability value of the %v connection does not match the expected value.", value)
						return false
					} else {
						_, err := strconv.Atoi(value[5])
						if err == nil {
							_, err := strconv.Atoi(value[6])
							if err == nil {
								_, err := strconv.Atoi(value[7])
								if err == nil {
									return true
								} else {
									log.Printf("The median value of the %v call does not match the expected value.", value)
									return false
								}
							} else {
								log.Printf("The purity value of the %v bond does not match the expected value.", value)
								return false
							}
						} else {
							log.Printf("The TTFB %v value does not match the expected value.", value)
							return false
						}
					}
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
