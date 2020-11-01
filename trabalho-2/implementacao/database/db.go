package database

import (
	"context"
	"fmt"
	"trabalho-2/m/model"

	"cloud.google.com/go/firestore"
)

func getCollection(client *firestore.Client,
	collection string) *firestore.CollectionRef {
	return client.Collection(collection)
}

func addBlock(ctx context.Context,
	remoteCollection *firestore.CollectionRef, block *model.ConcurrentBlock) error {
	_, _, err := remoteCollection.Add(ctx, block.Block)
	return err
}

// AddBlock -> Add a block to Firestore.
func AddBlock(ctx context.Context,
	client *firestore.Client, collection string, block *model.Block) error {

	remoteCollection := getCollection(client, collection)

	if remoteCollection == nil {
		return fmt.Errorf("%s is not a collection", collection)
	}

	docs, err := remoteCollection.Documents(ctx).GetAll()

	if err != nil {
		return err
	}

	fmt.Printf("Has %d documents.\n", len(docs))
	concurrentBlock := model.NewConcurrentBlock(block)
	concurrentBlock.Mu.Lock()
	defer concurrentBlock.Mu.Unlock()

	if len(docs) > 0 {
		parentBlock := model.NewBlock()

		if err := docs[len(docs)-1].DataTo(parentBlock); err != nil {
			return err
		}

		concurrentBlock.Block.Parent = &parentBlock.ID
		return addBlock(ctx, remoteCollection, concurrentBlock)
	}

	concurrentBlock.Block.Parent = nil
	return addBlock(ctx, remoteCollection, concurrentBlock)
}
