package seeder

import (
	"bytes"
	"encoding/json"
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
	"net/http"
	"strconv"
	"sync"
)

func createTag(t mdls.Tag) {
	url := baseUrl + "/tag/create"
	body, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	bodyReader := bytes.NewBuffer(body)
	resp, err := http.Post(url, "application/json", bodyReader)
	if err != nil {
		panic(err)
	}
	handleBadResponse(resp)
	s := readResponse(resp)
	_, err = strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}

func seedTags() {
	cmn.Log("seeding tags", cmn.LogLevels.Info)
	var wg sync.WaitGroup
	tagCreater := func(t mdls.Tag) {
		defer wg.Done()
		createTag(t)
	}
	for _, t := range tags {
		wg.Add(1)
		tagCreater(t)
	}
	wg.Wait()
}
