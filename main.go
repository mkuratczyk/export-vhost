package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func vhostDefinitions(vhost string) string {
	appEnv, _ := cfenv.Current()
	sharedService, err := appEnv.Services.WithName("shared-instance-admin")
	if err != nil {
		return fmt.Sprintf("{\"error\": \"Cannot find shared service credentials: %s\"}", err)
	}
	url, _ := sharedService.CredentialString("URL")
	username, _ := sharedService.CredentialString("username")
	password, _ := sharedService.CredentialString("password")

	tlsConfig := new(tls.Config)
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	rmqc, _ := rabbithole.NewTLSClient(url, username, password, transport)
	definitions, err := rmqc.ListVhostDefinitions(vhost)
	if err != nil {
		return fmt.Sprintf("{\"error\": \"Cannot retrieve the definitions: %s\"}", err)
	}
	definitionsJSON, _ := json.Marshal(definitions)
	return string(definitionsJSON)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	vhosts, ok := r.URL.Query()["vhost"]
	if !ok {
		appEnv, _ := cfenv.Current()
		fmt.Fprintf(w, "Please specify `vhost` in the URL. For example: https://%s/?vhost=example\n", appEnv.ApplicationURIs[0])
		return
	}

	vhost, err := url.QueryUnescape(vhosts[0])
	if err != nil || len(vhost) < 1 {
		fmt.Fprintln(w, "Incorrect vhost")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, vhostDefinitions(vhost))
}

func main() {
	http.HandleFunc("/", IndexHandler)

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
