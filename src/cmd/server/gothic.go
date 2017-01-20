package main

import (
	"errors"
	"net/http"

	"github.com/gomiddleware/mux"
	"github.com/markbates/goth/gothic"
)

func init() {
	gothic.GetProviderName = getProviderName
}

func getProviderName(r *http.Request) (string, error) {
	vals := mux.Vals(r)
	provider := vals["provider"]

	if provider == "" {
		return provider, errors.New("you must select a provider")
	}
	return provider, nil
}
