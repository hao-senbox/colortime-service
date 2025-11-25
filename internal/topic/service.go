package topic

import (
	"colortime-service/pkg/constants"
	"colortime-service/pkg/consul"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
)

type TopicService interface {
	GetTopicInfor(ctx context.Context, topicID string) (*Topic, error)
	GetVocabularyInforByTopicID(ctx context.Context, topicID string) ([]*Vocabulary, error)
}

type topicService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	mainService = "media-service"
)

func NewTopicService(client *api.Client) TopicService {
	mainServiceAPI := NewServiceAPI(client, mainService)
	return &topicService{
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

func (s *topicService) GetTopicInfor(ctx context.Context, topicID string) (*Topic, error) {

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	topictRes, err := s.client.getTopicInfor(topicID, token)
	if err != nil {
		return nil, err
	}

	rawData, ok := topictRes["data"]
	if !ok || rawData == nil {
		return nil, nil
	}

	topic, ok := rawData.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	name, _ := topic["title"].(string)
	id, _ := topic["id"].(string)
	mainImageUrl, _ := topic["main_image_url"].(string)
	videoUrl, _ := topic["video_url"].(string)

	return &Topic{
		ID:           id,
		Name:         name,
		MainImageUrl: mainImageUrl,
		VideoUrl:     videoUrl,
	}, nil
}

func (s *topicService) GetVocabularyInforByTopicID(ctx context.Context, topicID string) ([]*Vocabulary, error) {
	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	vocabularyRes, err := s.client.getVocabularyInforByTopicID(topicID, token)
	if err != nil {
		return nil, err
	}

	var vocabularies []*Vocabulary
	for _, item := range vocabularyRes {
		vocabMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		vocab := &Vocabulary{}

		if id, ok := vocabMap["id"].(string); ok {
			vocab.ID = id
		}

		if title, ok := vocabMap["title"].(string); ok {
			vocab.Title = title
		}

		if mainImageUrl, ok := vocabMap["main_image_url"].(string); ok {
			vocab.MainImageUrl = mainImageUrl
		}

		vocabularies = append(vocabularies, vocab)
	}

	return vocabularies, nil
}

func (c *callAPI) getVocabularyInforByTopicID(topicID, token string) ([]interface{}, error) {
	endpoint := fmt.Sprintf("/api/v2/gateway/%s/vocabularies", topicID)

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal([]byte(res), &apiResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	rawData, ok := apiResponse["data"]
	if !ok || rawData == nil {
		return nil, nil
	}

	vocabularyData, ok := rawData.([]interface{})
	if !ok {
		return nil, nil
	}

	return vocabularyData, nil
}

func (c *callAPI) getTopicInfor(topicID, token string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/api/v2/gateway/topics/%s", topicID)

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	var topicData interface{}

	err = json.Unmarshal([]byte(res), &topicData)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	myMap := topicData.(map[string]interface{})

	return myMap, nil

}
