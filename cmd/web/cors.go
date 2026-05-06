package main

import "net/http"

type CORSBridge struct {
}

func (b CORSBridge) Data(w http.ResponseWriter, r *http.Request, m map[string]any) (any, error) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	return nil, nil
}

func (b CORSBridge) Name() string {
	return "CORS"
}
