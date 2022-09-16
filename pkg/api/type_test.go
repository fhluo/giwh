package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

const resultString = `{
  "data": {
    "list": [
      {
        "count": "1",
        "gacha_type": "200",
        "id": "1000000000000000006",
        "item_id": "",
        "item_type": "武器",
        "lang": "zh-cn",
        "name": "沐浴龙血的剑",
        "rank_type": "3",
        "time": "2022-09-16 00:00:06",
        "uid": "10001"
      },
      {
        "count": "1",
        "gacha_type": "200",
        "id": "1000000000000000005",
        "item_id": "",
        "item_type": "武器",
        "lang": "zh-cn",
        "name": "翡玉法球",
        "rank_type": "3",
        "time": "2022-09-16 00:00:05",
        "uid": "10001"
      },
      {
        "count": "1",
        "gacha_type": "200",
        "id": "1000000000000000004",
        "item_id": "",
        "item_type": "武器",
        "lang": "zh-cn",
        "name": "匣里灭辰",
        "rank_type": "4",
        "time": "2022-09-16 00:00:04",
        "uid": "10001"
      },
      {
        "count": "1",
        "gacha_type": "200",
        "id": "1000000000000000003",
        "item_id": "",
        "item_type": "武器",
        "lang": "zh-cn",
        "name": "飞天御剑",
        "rank_type": "3",
        "time": "2022-09-16 00:00:02",
        "uid": "10001"
      },
      {
        "count": "1",
        "gacha_type": "200",
        "id": "1000000000000000002",
        "item_id": "",
        "item_type": "武器",
        "lang": "zh-cn",
        "name": "沐浴龙血的剑",
        "rank_type": "3",
        "time": "2022-09-16 00:00:02",
        "uid": "10001"
      },
      {
        "count": "1",
        "gacha_type": "200",
        "id": "1000000000000000001",
        "item_id": "",
        "item_type": "武器",
        "lang": "zh-cn",
        "name": "讨龙英杰谭",
        "rank_type": "3",
        "time": "2022-09-16 00:00:01",
        "uid": "10001"
      }
    ],
    "page": "0",
    "region": "cn_gf01",
    "size": "6",
    "total": "0"
  },
  "message": "OK",
  "retcode": 0
}`

func TestResult(t *testing.T) {
	var result Result
	err := json.Unmarshal([]byte(resultString), &result)
	if err != nil {
		t.Fatal(err)
	}
	
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, resultString, string(b))
}
