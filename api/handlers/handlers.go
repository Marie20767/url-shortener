package handlers

import (
	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
)

type Handler struct {
	KeyStore  *keys.KeyStore
	UrlStore  *urls.UrlStore
	ApiDomain string
}
