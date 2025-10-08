package language

import (
	"colortime-service/pkg/constants"
	"colortime-service/pkg/consul"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/consul/api"
)

type MessageLanguageGateway interface {
	UploadMessage(ctx context.Context, req UploadMessageRequest) error
	UploadMessages(ctx context.Context, req UploadMessageLanguagesRequest) error
	GetMessageLanguages(ctx context.Context, typeID string) ([]MessageLanguageResponse, error)
	// GetMessageLanguage(ctx context.Context, typeID string) (MessageLanguageResponse, error)
}

type messageLanguageGateway struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	mainService = "go-main-service"
)

func NewLanguageService(client *api.Client) MessageLanguageGateway {
	mainServiceAPI := NewServiceAPI(client, mainService)
	return &messageLanguageGateway{
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

	return &callAPI{
		client:       sd,
		clientServer: service,
	}
}

func (g *messageLanguageGateway) UploadMessage(ctx context.Context, req UploadMessageRequest) error {

	token, ok := ctx.Value(constants.Token).(string)

	if !ok {
		return fmt.Errorf("token not found in context")
	}

	err := g.client.uploadMessage(token, req)
	if err != nil {
		return err
	}

	return nil

}

func (g *messageLanguageGateway) UploadMessages(ctx context.Context, req UploadMessageLanguagesRequest) error {

	token, ok := ctx.Value(constants.Token).(string)

	if !ok {
		return fmt.Errorf("token not found in context")
	}

	err := g.client.uploadMessages(token, req)
	if err != nil {
		return err
	}

	return nil

}

func (g *messageLanguageGateway) GetMessageLanguages(ctx context.Context, typeID string) ([]MessageLanguageResponse, error) {

	token, ok := ctx.Value(constants.Token).(string)

	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	resp, err := g.client.getMessageLanguages(token, typeID)
	if err != nil {
		return nil, err
	}

	return resp, nil

}
func (c *callAPI) uploadMessage(token string, req UploadMessageRequest) error {

	endpoint := "/v1/gateway/messages"

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error marshalling request: %v\n", err)
		return err
	}

	_, err = c.client.CallAPI(c.clientServer, endpoint, http.MethodPost, jsonReq, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return err
	}

	return nil

}

func (c *callAPI) uploadMessages(token string, req UploadMessageLanguagesRequest) error {

	endpoint := "/v1/gateway/messages"

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error marshalling request: %v\n", err)
		return err
	}

	_, err = c.client.CallAPI(c.clientServer, endpoint, http.MethodPost, jsonReq, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return err
	}

	return nil

}

func (c *callAPI) getMessageLanguages(token string, typeID string) ([]MessageLanguageResponse, error) {

	endpoint := fmt.Sprintf("/v1/gateway/messages?type=%s&type_id=%s", "colortime", typeID)

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}
	fmt.Printf("res: %v\n", res)
	var data APIGateWayResponse[[]MessageLanguageResponse]
	err = json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	return data.Data, nil

}
