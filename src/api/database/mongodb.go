package database

import (
  "context"
  "go.mongodb.org/mongo-driver/v2/bson"
  "go.mongodb.org/mongo-driver/v2/mongo"
  "go.mongodb.org/mongo-driver/v2/mongo/options"
  "os"
)

type MongoDatabase struct {
  Context           context.Context
  Database          *mongo.Database
  ArchiveCollection *mongo.Collection
}

func NewMongoDatabase() Database {
  opts := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
  client, err := mongo.Connect(opts)
  if err != nil {
    panic(err) // TODO better error handling -> disable archive
  }

  database := client.Database(os.Getenv("MONGO_DB"))
  archiveCollection := database.Collection("archive")

  return &MongoDatabase{
    Context:           context.Background(),
    Database:          database,
    ArchiveCollection: archiveCollection,
  }
}

func (db MongoDatabase) UpsertArchiveEntry(entry ArchiveEntry[any]) error {
  filter := bson.M{
    "_id": entry.Id,
  }
  update := bson.M{
    "$set": entry,
  }
  print("Upserting archive entry with id: ", entry.Id, " and data: ", entry.Data)
  _, err := db.ArchiveCollection.UpdateOne(db.Context, filter, update, options.UpdateOne().SetUpsert(true))
  return err
}

func (db MongoDatabase) FindArchiveEntry(id string, dest any) error {
  filter := bson.M{
    "_id": id,
  }
  return db.ArchiveCollection.FindOne(db.Context, filter).Decode(&dest)
}
