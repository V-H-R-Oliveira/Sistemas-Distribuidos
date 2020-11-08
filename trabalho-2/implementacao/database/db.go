package database

import (
	"context"
	"fmt"
	"log"
	"sort"
	"trabalho-2/m/model"

	"cloud.google.com/go/firestore"
)

type fireDocs []*firestore.DocumentSnapshot

func (p fireDocs) Len() int {
	return len(p)
}

func (p fireDocs) Less(i, j int) bool {
	return p[i].CreateTime.Before(p[j].CreateTime)
}

func (p fireDocs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func getCollection(client *firestore.Client,
	collection string) *firestore.CollectionRef {
	return client.Collection(collection)
}

func addBlock(ctx context.Context,
	remoteCollection *firestore.CollectionRef, block *model.ConcurrentBlock) error {
	_, _, err := remoteCollection.Add(ctx, block.Block)
	return err
}

func sortByTimestamp(docs fireDocs) {
	sort.Sort(docs)
}

// GetBlocks -> Get all blocks from firestore.
func GetBlocks(ctx context.Context, client *firestore.Client, collection string) (*model.TemplateData, error) {
	remoteCollection := getCollection(client, collection)

	if remoteCollection == nil {
		return nil, fmt.Errorf("%s is not a collection", collection)
	}

	docs, err := remoteCollection.Documents(ctx).GetAll()

	if err != nil {
		return nil, err
	}

	sortByTimestamp(docs)
	fireDocs := []*model.Block{}

	for _, doc := range docs {
		block := model.NewBlock()

		if err := doc.DataTo(block); err != nil {
			log.Println(err)
			continue
		}

		fireDocs = append(fireDocs, block)
	}

	return model.NewTemplateData(fireDocs), nil
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

	sortByTimestamp(docs)

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
