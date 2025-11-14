package term

import (
	"colortime-service/pkg/constants"
	"colortime-service/pkg/consul"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
)

type TermService interface {
	GetTermByID(ctx context.Context, id string) (*TermInfor, error)
}

type termService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	mainService = "term-service"
)

func NewTermService(client *api.Client) TermService {
	mainServiceAPI := NewServiceAPI(client, mainService)
	return &termService{
		client: mainServiceAPI,
	}
}

func NewServiceAPI(client *api.Client, serviceName string) *callAPI {
	sd, err := consul.NewServiceDiscovery(client, serviceName)
	if err != nil {
		fmt.Printf("Error creating service discovery: %v\n", err)
		return nil
	}

	var service *api.CatalogService

	for i := 0; i < 10; i++ {
		service, err = sd.DiscoverService()
		if err == nil && service != nil {
			break
		}
		fmt.Printf("Waiting for service %s... retry %d/10\n", serviceName, i+1)
		time.Sleep(3 * time.Second)
	}

	if service == nil {
		fmt.Printf("Service %s not found after retries, continuing anyway...\n", serviceName)
	}

	if os.Getenv("LOCAL_TEST") == "true" {
		fmt.Println("Running in LOCAL_TEST mode — overriding service address to localhost")
		service.ServiceAddress = "localhost"
	}

	return &callAPI{
		client:       sd,
		clientServer: service,
	}
}

func (s *termService) GetTermByID(ctx context.Context, id string) (*TermInfor, error) {

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := s.client.getTermByID(token, id)
	if err != nil {
		log.Printf("[ERROR] termService.GetTermByID failed (id=%s): %v", id, err)
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("term not found")
	}

	termID, _ := data["id"].(string)
	startDate, _ := data["start_date"].(string)
	endDate, _ := data["end_date"].(string)

	if termID == "" || startDate == "" || endDate == "" {
		log.Printf("[ERROR] termService.GetTermByID invalid data: %+v", data)
		return nil, fmt.Errorf("invalid term data")
	}

	return &TermInfor{
		ID:        termID,
		StartDate: startDate,
		EndDate:   endDate,
	}, nil

}

func (c *callAPI) getTermByID(token, id string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/api/v1/gateway/terms/%s", id)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	response, err := c.client.CallAPI(c.clientServer, endpoint, "GET", nil, headers)
	if err != nil {
		log.Printf("[ERROR] CallAPI failed: %v", err)
		return nil, fmt.Errorf("call api term service failed: %w", err)
	}

	var parse map[string]interface{}
	if err := json.Unmarshal([]byte(response), &parse); err != nil {
		log.Printf("[ERROR] JSON unmarshal failed: %v | raw=%s", err, response)
		return nil, fmt.Errorf("invalid JSON response from term service: %w", err)
	}

	dataRaw, ok := parse["data"].(map[string]interface{})
	if !ok {
		statusCode, _ := parse["status_code"].(float64)
		errorMsg, _ := parse["error"].(string)
		log.Printf("[ERROR] Unexpected response format from term service (id=%s). status_code=%v, error=%s, raw=%+v", id, statusCode, errorMsg, parse)
		return nil, fmt.Errorf("term service returned error (status_code=%v, error=%s)", statusCode, errorMsg)
	}
	fmt.Printf("dataRaw: %v\n", dataRaw)
	return dataRaw, nil
}
