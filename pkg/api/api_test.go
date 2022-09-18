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

var repo []*Item

func init() {
	repo = make([]*Item, 100)
	for i := 0; i < 100; i++ {
		repo[i] = &Item{
			ID:       int64(10000 + (100 - i)),
			UID:      10001,
			WishType: StandardWish,
		}
	}
}

func handler() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

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

		list := lo.Filter(repo, func(item *Item, _ int) bool {
			return item.WishType.Str() == wishType
		})

		beginID, _ := strconv.ParseInt(c.Query("begin_id"), 10, 64)
		endID, _ := strconv.ParseInt(c.Query("end_id"), 10, 64)

		_, j, beginOK := lo.FindIndexOf(list, func(item *Item) bool {
			return item.ID == beginID
		})
		_, i, endOK := lo.FindIndexOf(list, func(item *Item) bool {
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
				Page:  0,
				Size:  size,
				Total: 0,
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
	assert.Equal(t, []*Item{}, items)
	assert.Nil(t, err)
}

func TestContext_FetchALL(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	ctx, err := New(Base{
		URL:   server.URL + "/event/gacha_info/api/getGachaLog",
		Query: BaseQuery{AuthKeyVer: "1", AuthKey: "x", Lang: "y"},
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx.SetInterval(0)

	items, err := ctx.WishType(SStandardWish).Size(10).FetchAll()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, repo, items)

	items, err = ctx.Begin(repo[10].ID).FetchAll()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, lo.Reverse(repo[:10]), items)

	items, err = ctx.End(repo[9].ID).FetchAll()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, repo[10:], items)
}

func TestContext_GetUID(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	ctx, err := New(Base{
		URL:   server.URL + "/event/gacha_info/api/getGachaLog",
		Query: BaseQuery{AuthKeyVer: "1", AuthKey: "x", Lang: "y"},
	})
	if err != nil {
		t.Fatal(err)
	}

	uid, err := ctx.WishType(SStandardWish).GetUID()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10001, uid)
}
