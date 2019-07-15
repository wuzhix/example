package database

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"io"
	"sync"
)

/**
 * es使用scroll分页获取数据
 */
func EsScroll() {
	url := []string{"http://127.0.0.1:9200", "http://127.0.0.1:9201"}
	index := "index_name"
	indexType := "index_type"
	size := 1000
	ctx := context.Background()
	es, _ := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(url...))
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewTermQuery("must", "must_value"))
	boolQuery.MustNot(elastic.NewTermQuery("must_not", "must_not_value"))
	boolQuery.Should(elastic.NewTermQuery("should", "should"))
	obj := es.Scroll(index).Type(indexType).Query(boolQuery).Size(size)
	var (
		wg sync.WaitGroup
		cs = make(chan int, 10)
	)
	for {
		res, err := obj.KeepAlive("10m").Do(ctx)
		if err == io.EOF {
			//一定要加这个，不然会导致退不出去
			break
		}
		if res == nil {
			break
		}
		cs <- 1
		wg.Add(1)
		go func(result *elastic.SearchResult) {
			for _, hit := range result.Hits.Hits {
				item := make(map[string]interface{})
				err := json.Unmarshal(*hit.Source, &item)
				if err != nil {
					continue
				}
				str := item["str"].(string)
				float := item["float"].(float32)
				integer := item["integer"].(int)
				println(str, float, integer)
			}
		}(res)
	}
}
