package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
)

type ServerCreationRequest struct {
	Type     string
	Metadata map[string]string
}

type ServerDeletionRequest struct {
	IP string
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

	var reqData ServerCreationRequest
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

func remove(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData ServerDeletionRequest
	err := json.NewDecoder(req.Body).Decode(&reqData)
	if err != nil {
		http.Error(res, "invalid post body", http.StatusBadRequest)
		return
	}

	err = DeleteServer(IPAddr{
		&net.IPAddr{
			IP: net.ParseIP(reqData.IP),
			Zone: "",
		},
	})

	if err != nil {
		http.Error(res, "invalid ip address", http.StatusBadRequest)

		return
	}

	jsonResponse, _ := json.Marshal(map[string]bool {
		"deleted": true,
	})
	res.Header().Add("Content-Type", "application/json")
	_, _ = res.Write(jsonResponse)
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
	http.Handle("/remove", auth(http.HandlerFunc(remove)))

	log.Printf("starting listener on %s", config.Listen)

	err := http.ListenAndServe(config.Listen, nil)
	if err != nil {
		log.Fatalln("failed to start listener on port", err)
	}
}
