package api

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var repo []Item

func init() {
	repo = make([]Item, 100)
	for i := 99; i >= 0; i-- {
		repo[i].WishType = "200"
		repo[i].ID = strconv.Itoa(10000 + i)
	}
}

func handler() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/event/gacha_info/api/getGachaLog", func(c *gin.Context) {
		if c.Query("authkey_ver") == "" || c.Query("authkey") == "" {
			c.JSON(http.StatusOK, Result{
				Data:    nil,
				Message: "authkey error",
				RetCode: -100,
			})
			return
		}

		if c.Query("lang") == "" {
			c.JSON(http.StatusOK, Result{
				Data:    nil,
				Message: "language error",
				RetCode: -108,
			})
			return
		}

		wishType := c.Query("gacha_type")
		size := 6
		if n, err := strconv.Atoi(c.Query("size")); err == nil && n >= 1 {
			size = n
		}

		list := lo.Filter(repo, func(item Item, _ int) bool {
			return item.WishType == wishType
		})

		beginID := c.Query("begin_id")
		endID := c.Query("end_id")

		_, j, beginOK := lo.FindIndexOf(list, func(item Item) bool {
			return item.ID == beginID
		})
		_, i, endOK := lo.FindIndexOf(list, func(item Item) bool {
			return item.ID == endID
		})

		switch {
		case beginOK:
			list = list[:j]
			if len(list) > size {
				list = list[len(list)-size:]
			}
		case endOK:
			list = list[i+1:]
			fallthrough
		default:
			if len(list) > size {
				list = list[:size]
			}
		}

		c.JSON(http.StatusOK, Result{
			Data: &Data{
				Page:  "0",
				Size:  strconv.Itoa(size),
				Total: "0",
				List:  list,
			},
			Message: "OK",
			RetCode: 0,
		})
		return
	})

	return r
}

func TestGetWishHistory(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	items, err := GetWishHistory(server.URL + "/event/gacha_info/api/getGachaLog")
	assert.Nil(t, items)
	assert.Equal(t, "authkey error", err.Error())

	items, err = GetWishHistory(server.URL + "/event/gacha_info/api/getGachaLog?authkey_ver=1&authkey=x")
	assert.Nil(t, items)
	assert.Equal(t, "language error", err.Error())

	items, err = GetWishHistory(server.URL + "/event/gacha_info/api/getGachaLog?authkey_ver=1&authkey=x&lang=y")
	assert.Equal(t, []Item{}, items)
	assert.Nil(t, err)
}
