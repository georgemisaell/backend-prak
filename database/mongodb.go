package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() *mongo.Database {

	mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        mongoURI = "mongodb://localhost:27017"
        log.Println("Peringatan: MONGO_URI tidak disetel. Menggunakan default:", mongoURI)
    }

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Koneksi database
    clientOptions := options.Client().ApplyURI(mongoURI)
    // Membuat konteks dengan batas waktu
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatalf("Koneksi ke MongoDB gagal: %v", err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatalf("Ping ke MongoDB gagal: %v", err)
    }

    fmt.Println("Berhasil terhubung ke database MongoDB!")
    return client.Database("alumnipedia")
}