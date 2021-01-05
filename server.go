package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type ServerRequest struct {
	Type     string
	Metadata map[string]string
}

func list(res http.ResponseWriter, _ *http.Request) {
	serverData, _ := json.Marshal(servers)
	res.Header().Add("Content-Type", "application/json")
	_, _ = res.Write(serverData)
}

func add(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData ServerRequest
	err := json.NewDecoder(req.Body).Decode(&reqData)
	if err != nil {
		http.Error(res, "invalid post body", http.StatusBadRequest)
		return
	}

	server, err := AddServer(reqData.Type, reqData.Metadata)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("added server %s of type %s", server.IP, server.Type)

	jsonServer, _ := json.Marshal(server)
	res.Header().Add("Content-Type", "application/json")
	_, _ = res.Write(jsonServer)
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

		decoded, err := hex.DecodeString(token)
		if err != nil || !config.Secret.Compare(Secret{decoded}) {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func StartServer() {
	http.Handle("/list", auth(http.HandlerFunc(list)))
	http.Handle("/add", auth(http.HandlerFunc(add)))

	log.Printf("starting listener on %s", config.Listen)

	err := http.ListenAndServe(config.Listen, nil)
	if err != nil {
		log.Fatalln("failed to start listener on port", err)
	}
}
