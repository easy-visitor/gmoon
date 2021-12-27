package gmoon

import (
	"context"
	"errors"
	"github.com/olivere/elastic/v7"
	"log"
)

type ElasticAdapter struct {
	Client *elastic.Client
}

func (this *ElasticAdapter) Name() string {
	return "ElasticAdapter"
}

func NewElasticAdapter() *ElasticAdapter {
	Client, err := elastic.NewClient(elastic.SetURL("http://192.168.50.128:12000/"), elastic.SetSniff(false))
	if err != nil {
		log.Fatal(err)
	}
	return &ElasticAdapter{
		Client: Client,
	}
}

func (this *ElasticAdapter) CreateIndex(indices string, body string) error {

	ctx1 := context.Background()
	exists, err := this.Client.IndexExists(indices).Do(context.Background())
	if err != nil {
		return err
	}
	if !exists {
		// 如果不存在，就创建
		createIndex, err := this.Client.CreateIndex("user").BodyString(body).Do(ctx1)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return errors.New("创建失败")
			// Not acknowledged ,创建失败
		}
	}
	return nil
}
