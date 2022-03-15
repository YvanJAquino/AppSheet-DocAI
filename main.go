package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	appsheet "github.com/yaq-cc/go-appsheet"
	events "github.com/yaq-cc/go-appsheet/logevent"
)

var (
	parent               = context.Background()
	PORT                 = "8080"
	ApplicationId        = "fefe8383-e022-44d8-832c-4aab45054f44"
	ApplicationAccessKey = "V2-HJqrB-0xdlI-nJETq-7nAPt-2lR0A-5VkAj-ZLPTy-CXPNc"
)

type AppSheetHandler struct {
	*appsheet.AppSheetClient
}

func NewAppSheetHandler(id, accessKey string) *AppSheetHandler {
	return &AppSheetHandler{
		AppSheetClient: appsheet.NewAppSheetClient(id, accessKey),
	}
}

func (h *AppSheetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var msg events.PubSubMessage
	json.NewDecoder(r.Body).Decode(&msg)
	json.NewEncoder(w).Encode(msg.Message.Data.ProtoPayload)
}

func main() {
	// AppSheet Client Configuration
	handler := NewAppSheetHandler(ApplicationId, ApplicationAccessKey)

	notify, stop := signal.NotifyContext(parent, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	server := &http.Server{
		Addr:        ":" + PORT,
		Handler:     mux,
		BaseContext: func(net.Listener) context.Context { return parent },
	}
	fmt.Printf("Starting HTTP/S server on localhost:%s", PORT)
	go server.ListenAndServe()

	<-notify.Done()
	shutdown, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()
	server.Shutdown(shutdown)
}
