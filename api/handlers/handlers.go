package handlers

import (
	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
)

type Handler struct {
	KeyCache  *keys.Cache
	UrlStore  *urls.UrlStore
	ApiDomain string
}
