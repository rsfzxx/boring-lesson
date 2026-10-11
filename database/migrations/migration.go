package migrations

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Migration struct {
	Version int
	Name    string
	Up      func(ctx context.Context, db *mongo.Database) error
	Down    func(ctx context.Context, db *mongo.Database) error
}

var all = []Migration{
	createRoles,
	createGroups,
	createUsers,
}

const historyCollection = "schema_migrations"

type record struct {
	Version   int       `bson:"_id"`
	Name      string    `bson:"name"`
	AppliedAt time.Time `bson:"applied_at"`
}

func Up(ctx context.Context, db *mongo.Database) error {
	history := db.Collection(historyCollection)
	for _, m := range all {
		err := history.FindOne(ctx, bson.M{"_id": m.Version}).Err()
		if err == nil {
			continue
		}
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
		if err := m.Up(ctx, db); err != nil {
			return fmt.Errorf("migration %06d_%s: %w", m.Version, m.Name, err)
		}
		rec := record{Version: m.Version, Name: m.Name, AppliedAt: time.Now().UTC()}
		if _, err := history.InsertOne(ctx, rec); err != nil {
			return err
		}
		log.Printf("up   %06d_%s", m.Version, m.Name)
	}
	return nil
}

func Down(ctx context.Context, db *mongo.Database) error {
	history := db.Collection(historyCollection)

	var last record
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})
	err := history.FindOne(ctx, bson.M{}, opts).Decode(&last)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Println("tidak ada migration untuk dibatalkan")
		return nil
	}
	if err != nil {
		return err
	}

	for _, m := range all {
		if m.Version != last.Version {
			continue
		}
		if err := m.Down(ctx, db); err != nil {
			return fmt.Errorf("migration %06d_%s: %w", m.Version, m.Name, err)
		}
		if _, err := history.DeleteOne(ctx, bson.M{"_id": m.Version}); err != nil {
			return err
		}
		log.Printf("down %06d_%s", m.Version, m.Name)
		return nil
	}
	return fmt.Errorf("migration versi %d tidak ditemukan di kode", last.Version)
}

func createCollection(ctx context.Context, db *mongo.Database, name string, schema bson.M) error {
	names, err := db.ListCollectionNames(ctx, bson.M{"name": name})
	if err != nil {
		return err
	}
	if len(names) > 0 {
		return nil
	}
	opts := options.CreateCollection().SetValidator(bson.M{"$jsonSchema": schema})
	return db.CreateCollection(ctx, name, opts)
}

const emailPattern = `^[^@\s]+@[^@\s]+\.[^@\s]+$`