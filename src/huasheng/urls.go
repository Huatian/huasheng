package huasheng

import (
	"net/http"
)

type Handler struct {
	URL         string
	HandlerFunc http.HandlerFunc
}

var (
	handlers = []Handler{
		{"/", indexHandler},
		{"/index/", indexHandler},
		{"/login/", loginHandler},
		{"/register/", registerHandler},
		{"/goodslist/", goodslistHandler},
		{"/goods/", regexpGoodsHandler(goodsHandler)},
	}
)