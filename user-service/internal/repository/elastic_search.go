package repository

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"user-service/internal/models"

	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type UserRepoElasticSearchInterface interface {
	Insert(ctx context.Context, user *models.User) error
	Search(ctx context.Context, value string) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, ID int64) error
}

type UserRepoElasticSearchObj struct {
	client *elasticsearch.Client
}

func NewUserRepoElasticSearchObj() *UserRepoElasticSearchObj {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("ES_HOST")},
		Username:  os.Getenv("ES_USERNAME"),
		Password:  os.Getenv("ES_PASSWORD"),
		Transport: transport,
	})
	if err != nil {
		return nil
	}

	return &UserRepoElasticSearchObj{client: client}
}

const (
	IndexName = "users"
	TimeOut   = time.Second * 15
)

func (es *UserRepoElasticSearchObj) Insert(ctx context.Context, user *models.User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("[Insert]: marshall data: %w", err)
	}

	req := esapi.CreateRequest{
		Index:      IndexName,
		DocumentID: fmt.Sprint(user.Id),
		Body:       bytes.NewReader(body),
	}

	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("[Insert] request: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {

		return fmt.Errorf("[Insert]: response: %s", res.String())
	}

	return nil
}

func (es *UserRepoElasticSearchObj) Search(ctx context.Context, value string) ([]models.User, error) {
	users := make([]models.User, 0)
	mapResp := make(map[string]interface{})

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]interface{}{
				"username": "*" + strings.ToLower(value) + "*",
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("[Search] encode  %s", err)
	}

	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()

	req := esapi.SearchRequest{
		Index:          []string{IndexName},
		Body:           &buf,
		TrackTotalHits: true,
		Pretty:         true,
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return nil, fmt.Errorf("[Search] request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("[Search] response: %s", res.String())
	}

	if err := json.NewDecoder(res.Body).Decode(&mapResp); err != nil {
		return nil, fmt.Errorf("[Search] decode: %s", err)
	}

	for _, hit := range mapResp["hits"].(map[string]interface{})["hits"].([]interface{}) {
		user := models.User{}
		doc := hit.(map[string]interface{})
		source := doc["_source"]

		byteData, _ := json.Marshal(source)
		json.Unmarshal(byteData, &user)
		users = append(users, user)

	}

	return users, nil
}

func (es *UserRepoElasticSearchObj) Update(ctx context.Context, user *models.User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("[Update]: marshall data: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      IndexName,
		DocumentID: fmt.Sprint(user.Id),
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, body))),
	}

	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("[Update] request: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("[Update]: response: %s", res.String())
	}

	return nil
}

func (es *UserRepoElasticSearchObj) Delete(ctx context.Context, ID int64) error {
	req := esapi.DeleteRequest{
		Index:      IndexName,
		DocumentID: fmt.Sprint(ID),
	}

	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("[Delete] request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("[Delete] response: %s", res.String())
	}

	return nil
}
