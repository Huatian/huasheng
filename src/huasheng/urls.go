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
		{"/index", indexHandler},
		{"/signin", signinHandler},
		{"/signup", signupHandler},
		{"/goodslist", goodslistHandler},
		{"/goods", regexpGoodsHandler(goodsHandler)},
	}
)