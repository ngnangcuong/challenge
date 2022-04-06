package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"fmt"
	"encoding/json"
	"strconv"
	"context"
	"bytes"
	"challenge3/models"
)

type postSearchRepoImpl struct {
	esClient *elasticsearch.Client
}

var esIndex = "post"

func NewPostSearchRepo(esClient *elasticsearch.Client) models.PostSearchRepo {
	return &postSearchRepoImpl{
		esClient: esClient,
	}
}

func (p *postSearchRepoImpl) Search(keyword string) ([]models.Post, error) {
	var (
		buf bytes.Buffer
		r map[string]interface{}
		postList []models.Post
	)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase_prefix": map[string]interface{}{
				"content": keyword,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := p.esClient.Search(
		p.esClient.Search.WithContext(context.Background()),
		p.esClient.Search.WithIndex(esIndex),
		p.esClient.Search.WithBody(&buf),
		p.esClient.Search.WithTrackTotalHits(true),
		p.esClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Response from ESClient got errors")
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var post = models.Post{
			ID: uint(hit.(map[string]interface{})["_source"].(map[string]interface{})["id"].(float64)),
			Email: hit.(map[string]interface{})["_source"].(map[string]interface{})["email"].(string),
			Content: hit.(map[string]interface{})["_source"].(map[string]interface{})["content"].(string),
		}
		postList = append(postList, post)
	}

	return postList, nil
}

func (p *postSearchRepoImpl) Index(post models.Post) error {
	body, err := json.Marshal(post)
	if err != nil {
		return err
	}

	var id = strconv.FormatUint(uint64(post.ID), 10)

	req := esapi.IndexRequest{
		Index: esIndex,
		DocumentID: id,
		Body: bytes.NewReader(body),
	}

	res, err := req.Do(context.Background(), p.esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 400 {
		return fmt.Errorf("Bad Request")
	}

	if res.IsError() {
		return fmt.Errorf("response: %s", res.String())
	}

	return nil
}

func (p *postSearchRepoImpl) Update(post models.Post) error {
	body, err := json.Marshal(post)
	if err != nil {
		return err
	}

	var id = strconv.FormatUint(uint64(post.ID), 10)

	req := esapi.UpdateRequest{
		Index: esIndex,
		DocumentID: id,
		Body: bytes.NewReader([]byte(fmt.Sprintf(`{"doc": %s }`, body))),
	}

	res, err := req.Do(context.Background(), p.esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return fmt.Errorf("Not exist post")
	}

	if res.IsError() {
		return fmt.Errorf("response: %s", res.String())
	}

	return nil
}

func (p *postSearchRepoImpl) Delete(id string) error {
	req := esapi.DeleteRequest{
		Index: esIndex,
		DocumentID: id,
	}

	res, err := req.Do(context.Background(), p.esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return fmt.Errorf("Not exist post")
	}

	if res.IsError() {
		return fmt.Errorf("reponse: %s", res.String())
	}

	return nil
}