package requester

import (
	"net/http"
	"encoding/json"
	"log"
	"sync"
	"time"
	"io/ioutil"
	"sort"
	"errors"
	"fmt"
)
// Error response
type ErrorResponse struct {
	Code	int	`json:"code"`
	Error 	string	`json:"error"`
}
// Encode error response
func processErrorResponse(code int, error string) []byte {
	jsonString, err := json.Marshal(ErrorResponse{code, error})
	if err != nil {
		return []byte(fmt.Sprintf("{\"code\": %d, \"error\": \"%s\"}", code, error))
	}
	return jsonString
}
// Incoming Price
type Price struct {
	Price    int    `json:"price"`
}
// Bid price
type SourcePrice struct {
	Price	int	`json:"price"`
	Source 	string	`json:"source"`
}
// Requester interface
type IRequester interface {
	Handle(writer http.ResponseWriter, request *http.Request)
	bidCalculation(sourcePrices []SourcePrice) ([]byte, error)
	maxPrice(prices []Price) int
	makeRequest(url string, ch chan SourcePrice, wg *sync.WaitGroup)
}

// Requester
type Requester struct {
	client				http.Client		// http client
	timeout				time.Duration 	// timeout for outcoming request
}

// New Requester.
func NewRequester() IRequester {
	//	log.Printf("New\n")
	requester := new(Requester)
	requester.timeout = time.Duration(100 * time.Millisecond)
	requester.client = http.Client{
		Timeout: requester.timeout,
	}

	return requester
}

// Process function
func (self *Requester) bidCalculation(sourcePrices []SourcePrice) ([]byte, error) {
	if (len(sourcePrices) == 0) {
		return []byte{}, errors.New("Prices not exists")
	}
	// Sorting max prices
	sort.Slice(sourcePrices, func(i, j int) bool {
		return sourcePrices[i].Price < sourcePrices[j].Price
	})

	//log.Printf("Source prices %s\n", sourcePrices)
	var secondMax SourcePrice
	if (len(sourcePrices) > 1) {
		secondMax = sourcePrices[len(sourcePrices)-2]
	} else {
		secondMax = sourcePrices[0]
	}

	jsonString, err := json.Marshal(secondMax)
	if err != nil {
		return []byte{}, err
	}
	return jsonString, nil
}

// Get max price
func (self *Requester) maxPrice(prices []Price) int {
	var max Price = prices[0]
	for _, price := range prices {
		if max.Price < price.Price {
			max = price
		}
	}
	return max.Price
}

// Make outcoming request
func (self *Requester) makeRequest(url string, ch chan SourcePrice, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := self.client.Get(url)

	if err != nil {
		log.Printf("Error: %s\n\n", err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Response body error %s -> %s\n", url, err)
		} else if (resp.StatusCode >= 200 && resp.StatusCode < 300) {
			var prices []Price
			errJson := json.Unmarshal(body, &prices)
			if errJson != nil {
				log.Printf("Json deconding error %s -> %s\n", url, err)
			}
			if (len(prices) > 0) {
				ch <- SourcePrice{self.maxPrice(prices), url}
			}
		}
	}
}

// Handle request
func (self *Requester) Handle(writer http.ResponseWriter, request *http.Request) {
	var aggregatedPrices []SourcePrice
	wg := new(sync.WaitGroup)
	urls := request.URL.Query()["s"]
	ch := make(chan SourcePrice, len(urls))
	wg.Add(len(urls))
	//
	log.Printf("Urls %s\n", urls)
	writer.Header().Set("Content-Type", "application/json")

	if len(urls) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(processErrorResponse(http.StatusBadRequest, "Bad Request"))
	}

	// make concurrent requests
	for _, url := range urls {
		go self.makeRequest(url, ch, wg)
	}

	wg.Wait()
	close(ch)

	// aggregate prices
	for sourcePrice := range ch {
		log.Printf("Url: %s", sourcePrice.Source)
		aggregatedPrices = append(aggregatedPrices, sourcePrice)
	}

	// If bid not found -> 404
	if (len(aggregatedPrices) == 0) {
		log.Printf("Error: %s", "Bid not found")
		writer.WriteHeader(http.StatusNotFound)
		writer.Write(processErrorResponse(http.StatusNotFound, "Bid not found"))
		return
	}

	//log.Printf("Aggregated prices %s\n", aggregatedPrices)
	response, err := self.bidCalculation(aggregatedPrices)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		writer.WriteHeader(http.StatusConflict)
		writer.Write(processErrorResponse(http.StatusConflict, err.Error()))
		return
	}

	writer.Write(response)
}

