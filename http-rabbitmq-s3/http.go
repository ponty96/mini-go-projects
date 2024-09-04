package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerConfig struct {
	host  string
	port  int
	Event Events
}

type ProfileJson struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Address   Address `json:"address"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

// when this server receives request to publish events, it will add it to this channel, which will listened to and used to publish events to RabbitMQ
type HttpServer struct {
	serverConfig *ServerConfig
}

func (s *HttpServer) serveHTTP() {
	r := mux.NewRouter()
	fmt.Println("Starting HTTP Server on port 8080")
	r.HandleFunc("/show-profile/{profileId}", s.showProfile).Methods("GET")
	r.HandleFunc("/publish-event", s.publishEventHandler).Methods("POST")
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.serverConfig.host, s.serverConfig.port), r)
}

func (s *HttpServer) showProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("i got a request showProfile")
	vars := mux.Vars(r)
	profileId := vars["profileId"]

	w.WriteHeader(http.StatusOK)
	// get profile template from S3
	w.Write([]byte(fmt.Sprintf("Showing profile %s", profileId)))
}

// We need to consider different event types
func (s *HttpServer) publishEventHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("i got a request publishEventHandler")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Panicf("Failure to read body! %s", err)
		return
	}

	var profile ProfileJson
	err = json.Unmarshal(body, profile)

	log.Printf(" [x] Body request %s", profile.FirstName)

	s.serverConfig.Event.handlePublishEvent(profile)

	response, _ := json.Marshal(map[string]string{"status": "success"})

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
