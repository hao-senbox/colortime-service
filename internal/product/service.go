package product

import (
	"colortime-service/pkg/consul"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService interface {
	GetProductInfor(productID string) (*Product, error)
}

type productService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	mainService = "product-service"
)

func NewUserService(client *api.Client) ProductService {
	mainServiceAPI := NewServiceAPI(client, mainService)
	return &productService{
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

func (s *productService) GetProductInfor(productID string) (*Product, error) {

	productRes, err := s.client.getProductInfor(productID)
	if err != nil {
		return nil, err
	}

	rawData, ok := productRes["data"]
	if !ok || rawData == nil {
		return nil, fmt.Errorf("product not found for id=%s", productID)
	}

	product, ok := rawData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid product format for id=%s", productID)
	}

	name, _ := product["product_name"].(string)
	id, _ := product["id"].(string)

	objIDProduct, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID format: %v", err)
	}

	priceStore, _ := product["original_price_store"].(float64)
	priceService, _ := product["original_price_service"].(float64)

	var imageURL string
	if v, ok := product["cover_image"].(string); ok {
		imageURL = v
	}

	var topic string
	if rawTopic, ok := product["topic"].(map[string]interface{}); ok {
		topic, _ = rawTopic["topic_name"].(string)
	}

	var category string
	if rawCategory, ok := product["category"].(map[string]interface{}); ok {
		category, _ = rawCategory["category_name"].(string)
	}

	return &Product{
		ID:                   objIDProduct,
		ProductName:          name,
		OriginalPriceStore:   priceStore,
		OriginalPriceService: priceService,
		ProductImage:         imageURL,
		TopicName:            topic,
		CategoryName:         category,
	}, nil
}

func (c *callAPI) getProductInfor(productID string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/api/v1/products/%s", productID)

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, nil)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	var productData interface{}

	err = json.Unmarshal([]byte(res), &productData)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	myMap := productData.(map[string]interface{})

	return myMap, nil

}
