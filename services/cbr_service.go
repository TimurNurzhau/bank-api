package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/beevik/etree"
)

type CBRService struct {
	client *http.Client
}

func NewCBRService() *CBRService {
	return &CBRService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *CBRService) GetKeyRate() (float64, error) {
	soapRequest := s.buildSOAPRequest()
	rawBody, err := s.sendRequest(soapRequest)
	if err != nil {
		return 0, err
	}

	return s.parseXMLResponse(rawBody)
}

func (s *CBRService) buildSOAPRequest() string {
	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
<soap12:Body>
<KeyRate xmlns="http://web.cbr.ru/">
<fromDate>%s</fromDate>
<ToDate>%s</ToDate>
</KeyRate>
</soap12:Body>
</soap12:Envelope>`, fromDate, toDate)
}

func (s *CBRService) sendRequest(soapRequest string) ([]byte, error) {
	req, err := http.NewRequest(
		"POST",
		"https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx",
		bytes.NewBuffer([]byte(soapRequest)),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://web.cbr.ru/KeyRate")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CBR request error: %v", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (s *CBRService) parseXMLResponse(rawBody []byte) (float64, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(rawBody); err != nil {
		return 0, fmt.Errorf("XML parse error: %v", err)
	}

	krElements := doc.FindElements("//diffgram/KeyRate/KR")
	if len(krElements) == 0 {
		// Fallback: возвращаем дефолтную ставку
		return 16.0, nil
	}

	latestKR := krElements[0]
	rateElement := latestKR.FindElement("./Rate")
	if rateElement == nil {
		return 16.0, nil
	}

	var rate float64
	if _, err := fmt.Sscanf(rateElement.Text(), "%f", &rate); err != nil {
		return 16.0, nil
	}

	return rate, nil
}