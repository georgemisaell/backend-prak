package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongoDB() *mongo.Client {

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable tidak disetel.")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Koneksi database
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Gagal koneksi ke MongoDB: %v", err)
	}

	// Tes Koneksi (Ping)
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer pingCancel()

	if err = client.Ping(pingCtx, readpref.Primary()); err != nil {
		// Log error dan keluar jika ping gagal
		log.Fatalf("Gagal ping MongoDB: %v", err)
	}

	log.Println("Sukses koneksi ke MongoDB!")
	return client
}