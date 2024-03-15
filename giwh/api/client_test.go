package api

import (
	"fmt"
	"github.com/fhluo/giwh/giwh/api/gacha"
	"github.com/fhluo/giwh/hoyo-auth/clients"
	"testing"
)

func TestClient(t *testing.T) {
	genshin := clients.Latest()
	if genshin == nil {
		t.Fatal("no client available")
	}

	auths := genshin.Auths()
	if len(auths) == 0 {
		t.Fatal("failed to find auth info")
	}

	c := NewClient(auths[len(auths)-1])
	fmt.Print(c.NewFetcher(gacha.CharacterEventWish).NextURL())
}
