package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"frigga/modules/driver"
	"os"

	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

// createHTTPTask creates a new task with a HTTP target then adds it to a Queue.
func createHTTPTask(projectID string, locationID string, queueID string, url string, message interface{}, token string) (*taskspb.Task, error) {

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	// Build the Task payload.
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
					Headers: map[string]string{
						"Authorization": token,
					},
				},
			},
		},
	}

	// Add a payload message if one is present.
	payload, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	req.Task.GetHttpRequest().Body = payload

	ctx := context.Background()
	createdTask, err := driver.CT.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}

	return createdTask, nil
}

// DispatchBotEventWorker ...
func DispatchBotEventWorker(event BotEvent) {
	projectID := os.Getenv("PROJECT_ID")
	locationID := os.Getenv("PROJECT_LOCATION")
	authToken := os.Getenv("APP_BASE_TOKEN")
	url := os.Getenv("PROJECT_URL") + "/worker"
	queueID := "bot-event"

	_, err := createHTTPTask(projectID, locationID, queueID, url, event, authToken)
	if err != nil {
		fmt.Println(err.Error())
	}

	return
}
