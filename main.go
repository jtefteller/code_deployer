package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"cloud.google.com/go/pubsub"
	"gopkg.in/yaml.v3"
)

const (
	ProjectID       = "project_id"
	SubscriptionID  = "subscription_id"
	DockerRepo      = "docker_repo"
	SuccessfulBuild = "Build successful"
)

var (
	home           = os.Getenv("HOME")
	configPath     = home + "/.code_deployer/"
	configYaml     = configPath + "config.yaml"
	serviceAccount = configPath + "service-account.json"
)

type Message struct {
	Status string   `json:"status"`
	Paths  []string `json:"paths"`
}

func main() {
	_ = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", serviceAccount)

	config := map[string]string{}
	yamlFile, err := os.Open(configYaml)
	if err != nil {
		log.Fatalf("Error reading yaml file: %v", err)
	}
	yaml.NewDecoder(yamlFile).Decode(&config)
	projectID, subscriptionID, dockerRepo := mustCompile(config)

	dockerLogin(dockerRepo)

	ctx := context.Background()
	go heathCheck()
	pubsubClient, _ := pubsub.NewClient(ctx, projectID)
	pubsubClient.Subscription(subscriptionID).Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", msg.Data)
		msg.Ack()
		msgStruct := &Message{}
		json.Unmarshal(msg.Data, msgStruct)
		if msgStruct.Status != SuccessfulBuild {
			return
		}

		runMakeDeploy(msgStruct.Paths)
	})
}

func runMakeDeploy(paths []string) {
	for _, path := range paths {
		cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && make deploy", path))
		b, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error: %v\nOutput: %s", err, b)
		} else {
			log.Printf("Output: %s", b)
		}
	}
}

func dockerLogin(dockerRepo string) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cat %s | base64 | docker login -u _json_key_base64 --password-stdin %s", serviceAccount, dockerRepo))

	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, b)
	}

	log.Printf("Output: %s", b)
}

func mustCompile(config map[string]string) (string, string, string) {
	projectID, ok := config[ProjectID]
	if !ok {
		log.Fatalf("Missing %s in config.yaml", ProjectID)
	}

	subscriptionID, ok := config[SubscriptionID]
	if !ok {
		log.Fatalf("Missing %s in config.yaml", SubscriptionID)
	}

	dockerRepo, ok := config[DockerRepo]
	if !ok {
		log.Fatalf("Missing %s in config.yaml", DockerRepo)
	}

	return projectID, subscriptionID, dockerRepo
}

func heathCheck() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":1337", nil)
}
