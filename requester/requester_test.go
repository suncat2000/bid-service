package requester

import (
	"testing"
	"sync"
	"log"
	"io/ioutil"
)

func fillSlices() (sourcePrices1 []SourcePrice, sourcePrices2 []SourcePrice, sourcePrices3 []SourcePrice) {
	sourcePrices1 = make([]SourcePrice, 3)
	sourcePrices1[0] = SourcePrice{12, "http://example.com/12"}
	sourcePrices1[1] = SourcePrice{28, "http://example.com/28"}
	sourcePrices1[2] = SourcePrice{31, "http://example.com/31"}

	sourcePrices2 = make([]SourcePrice, 6)
	sourcePrices2[0] = SourcePrice{17, "http://example.com/17"}
	sourcePrices2[1] = SourcePrice{26, "http://example.com/26"}
	sourcePrices2[2] = SourcePrice{28, "http://example.com/28"}
	sourcePrices2[3] = SourcePrice{62, "http://example.com/62"}
	sourcePrices2[4] = SourcePrice{44, "http://example.com/44"}
	sourcePrices2[5] = SourcePrice{13, "http://example.com/13"}

	sourcePrices3 = make([]SourcePrice, 5)
	sourcePrices3[0] = SourcePrice{11, "http://example.com/11"}
	sourcePrices3[1] = SourcePrice{33, "http://example.com/33"}
	sourcePrices3[2] = SourcePrice{27, "http://example.com/27"}
	sourcePrices3[3] = SourcePrice{37, "http://example.com/37"}
	sourcePrices3[4] = SourcePrice{23, "http://example.com/23"}

	return sourcePrices1, sourcePrices2, sourcePrices3
}

// Test Bid calculation method
func TestBidCalculation(t *testing.T) {
	sourcePrices1, sourcePrices2, sourcePrices3 := fillSlices()
	requestHandler := NewRequester()

	response1, _ := requestHandler.bidCalculation(sourcePrices1)
	availableResponse1 := "{\"price\":28,\"source\":\"http://example.com/28\"}"
	if string(response1) != availableResponse1 {
		t.Errorf("Bid not equeals %s -> %s", availableResponse1, response1)
	}
	response2, _ := requestHandler.bidCalculation(sourcePrices2)
	availableResponse2 := "{\"price\":44,\"source\":\"http://example.com/44\"}"
	if string(response2) != availableResponse2 {
		t.Errorf("Bid not equeals %s -> %s", availableResponse2, response2)
	}
	response3, _ := requestHandler.bidCalculation(sourcePrices3)
	availableResponse3 := "{\"price\":33,\"source\":\"http://example.com/33\"}"
	if string(response3) != availableResponse3 {
		t.Errorf("Bid not equeals %s -> %s", availableResponse3, response3)
	}
}

// Test Bid calculation method
func TestMaxPriceInt(t *testing.T) {
	maxPrices := make([]Price, 7)
	maxPrices[0] = Price{11}
	maxPrices[1] = Price{33}
	maxPrices[2] = Price{27}
	maxPrices[3] = Price{37}
	maxPrices[4] = Price{23}
	maxPrices[5] = Price{63}
	maxPrices[6] = Price{46}

	requestHandler := NewRequester()
	maxPrice := requestHandler.maxPrice(maxPrices)
	availableMaxPrice := 63
	if maxPrice != availableMaxPrice {
		t.Errorf("Max price not equeals %d -> %d", availableMaxPrice, maxPrice)
	}
}

func TestHandle(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	var aggregatedPrices []SourcePrice
	wg := new(sync.WaitGroup)
	urls := [3]string{
		"http://localhost:8081/primes",
		"http://localhost:8081/rand",
		"http://localhost:8081/fact",
	}
	ch := make(chan SourcePrice, len(urls))
	wg.Add(len(urls))
	//
	requestHandler := NewRequester()
	// make concurrent requests
	for _, url := range urls {
		go requestHandler.makeRequest(url, ch, wg)
	}

	wg.Wait()
	close(ch)

	// aggregate prices
	for sourcePrice := range ch {
		aggregatedPrices = append(aggregatedPrices, sourcePrice)
	}

	_, err := requestHandler.bidCalculation(aggregatedPrices)
	if err != nil {
		t.Logf("Response error: %s", err)
	}

}