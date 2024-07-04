package utils

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func DeleteCollection(firestore *firestore.Client, collectionName string,
) error {

	ctx := context.Background()
	col := firestore.Collection(collectionName)
	bulkwriter := firestore.BulkWriter(ctx)

	for {
		// Get a batch of documents
		iter := col.Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to the BulkWriter.
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			bulkwriter.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			bulkwriter.End()
			break
		}

		bulkwriter.Flush()
	}
	return nil
}
