package gcf

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	firestore "google.golang.org/api/firestore/v1beta2"
)

var projectID string

type PubSubMessage struct {
	Data string `json:"data"`
}

func init() {
	projectID = os.Getenv("GCP_PROJECT")
}

func BackupFirestore(ctx context.Context, m PubSubMessage) error {
	client, err := google.DefaultClient(ctx,
		"https://www.googleapis.com/auth/datastore",
		"https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return errors.Wrap(err, "Failed to create a Google client")
	}

	svc, err := firestore.New(client)
	if err != nil {
		return errors.Wrap(err, "Failed to create Firestore service")
	}

	req := &firestore.GoogleFirestoreAdminV1beta2ExportDocumentsRequest{
		OutputUriPrefix: fmt.Sprintf("gs://%s-backup-firestore", projectID),
	}
	_, err = firestore.NewProjectsDatabasesService(svc).ExportDocuments(
		fmt.Sprintf("projects/%s/databases/(default)", projectID), req,
	).Context(ctx).Do()
	if err != nil {
		return errors.Wrap(err, "Failed to export Firestore")
	}

	fmt.Println("Backup Successfully")

	return nil
}