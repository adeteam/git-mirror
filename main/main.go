package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/adeteam/git-mirror/definition"
	"github.com/adeteam/git-mirror/service"
)

func main() {
	parser := argparse.NewParser("main", "Start Git Mirror Service")
	storage := parser.String("s", "storage-path", &argparse.Options{
		Required: true,
		Help:     "Local path to where to store repository mirrors",
	})
	webhook_trigger := parser.String("w", "webhook-trigger", &argparse.Options{
		Required: false,
		Help:     "The next Webhook to forward the request",
	})
	username := parser.String("u", "username", &argparse.Options{
		Required: false,
		Help:     "The username to use for git mirror of the remote repository",
	})
	password := parser.String("p", "password", &argparse.Options{
		Required: false,
		Help:     "The password to use for git mirror of the remote repository",
	})
	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Help: "Enable verbose logging",
	})

	port := parser.Int("P", "port", &argparse.Options{
		Required: false,
		Default:  4000,
		Help:     "The port to listen on for Webhooks",
	})

	// parser the inputs
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	logrus.SetLevel(logrus.InfoLevel)
	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	config := definition.Config{
		StoragePath:    *storage,
		WebhookTrigger: *webhook_trigger,
		GitUsername:    *username,
		GitPassword:    *password,
		Port:           *port,
	}
	service.Config().Current = config

	log.Infof("Starting Git-Mirror Service on %d", config.Port)
	http.HandleFunc("/", DefaultHandler)
	http.ListenAndServe(
		fmt.Sprintf(":%d", config.Port),
		nil,
	)
}

func DefaultHandler(response http.ResponseWriter, request *http.Request) {
	log.Debugf("received request method %s on path %s", request.Method, request.RequestURI)
	if request.Method != "POST" {
		// we do not support any other methods
		http.Error(response, "400 Bad Request", http.StatusBadRequest)
		return
	}

	var payload map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.WithError(err).Errorf("unable to read post body")
		http.Error(response, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Debugf("JSON: %s", body)
	json.Unmarshal(body, &payload)
	if repository, ok := payload["repository"].(map[string]interface{}); ok {
		repository_name := repository["full_name"].(string)
		log.Debugf("received webhook by for %s", repository_name)
	} else {
		log.WithError(err).Errorf("unable to extract repository info from post body")
		http.Error(response, "400 Bad Request", http.StatusBadRequest)
	}
}
