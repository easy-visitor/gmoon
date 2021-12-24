package gmoon

import (
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
	Client, err := elastic.NewClient(elastic.SetURL("http://192.168.50.128:12000/"),elastic.SetSniff(false))
	if err != nil {
		log.Fatal(err)
	}
	return &ElasticAdapter{
		Client: Client,
	}
}

//	name := "people2"
//	Client.CreateIndex(name).Do(context.Background())
