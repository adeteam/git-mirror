package main

import (
	"bytes"
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

	webhook_id := request.Header.Get("X-Hook-UUID")

	var payload map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.WithError(err).Errorf("[%s] unable to read post body", webhook_id)
		http.Error(response, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Debugf("[%s] JSON: %s", webhook_id, body)
	json.Unmarshal(body, &payload)
	if repository, ok := payload["repository"].(map[string]interface{}); ok {
		repository_name := repository["full_name"].(string)
		log.Debugf("[%s] received webhook by for %s", webhook_id, repository_name)
		// run this in a separate routine and allow send a response back
		go func() {
			forward_headers := []string{
				"Content-Type",
				"User-Agent",
				"X-Attempt-Number",
				"X-B3-Sampled",
				"X-B3-SpanId",
				"X-B3-TracedId",
				"X-Event-Key",
				"X-Event-Type",
				"X-Hook-UUID",
				"X-Request-UUID",
			}

			// ensure this repo is mirror
			service.Git().HandleRepoMirror(repository_name)

			// forward this webhook request
			if service.Config().Current.WebhookTrigger != "" {
				var sent_err error
				client := http.Client{}

				// we will try up to 3 times
				for i := 0; i < 3; i++ {
					forward_body := bytes.NewBuffer(body)
					forward_request, err := http.NewRequest("POST", service.Config().Current.WebhookTrigger, forward_body)

					for _, fheader := range forward_headers {
						forward_request.Header.Add(
							fheader,
							request.Header.Get(fheader),
						)
					}

					forward_response, err := client.Do(forward_request)
					if err == nil {
						forward_response_body, _ := ioutil.ReadAll(forward_response.Body)
						// we have sent this request
						log.Infof(
							"[%s] webhook request forwarded, %s - %s",
							webhook_id,
							forward_response.Status,
							forward_response_body,
						)
						sent_err = nil
						break
					} else {
						sent_err = err
					}
				}

				if sent_err != nil {
					log.WithError(sent_err).Errorf("[%s] unable to forward webhook request", webhook_id)
				}
			}
		}()

		message := map[string]string{
			"message": "thanks",
		}
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(message)
	} else {
		log.WithError(err).Errorf("[%s] unable to extract repository info from post body", webhook_id)
		http.Error(response, "400 Bad Request", http.StatusBadRequest)
	}
}
