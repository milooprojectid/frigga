package driver

import (
	"context"
	"fmt"
	"os"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"

	"google.golang.org/api/option"
)

// CT ...
var CT *cloudtasks.Client

// FSContext ...
// var FSContext *context.Context

// InitializeCloudTask ...
func InitializeCloudTask() {
	ctx := context.Background()
	sa := option.WithCredentialsFile(os.Getenv("FIREBASE_ACCOUNT_URL"))
	client, err := cloudtasks.NewClient(ctx, sa)
	if err != nil {
		fmt.Errorf("NewClient: %v", err)
		return
	}

	CT = client
	// FSContext = &ctx
}
