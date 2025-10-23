package repository

import (
	"context"
	"latihan_uts_2/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IAlumniRepository interface {
	CreateAlumni(ctx context.Context, user *models.Alumni) (*models.Alumni, error)
	FindAlumniByID(ctx context.Context, id string) (*models.Alumni, error)
	FindAllAlumni(ctx context.Context) ([]models.Alumni, error)
}

// UserRepository implementasi IUserRepository menggunakan MongoDB.
type AlumniRepository struct {
    collection *mongo.Collection
}

// NewUserRepository membuat instance baru dari UserRepository dan mengaitkannya dengan koleksi MongoDB.
func NewAlumniRepository(db *mongo.Database) IAlumniRepository {
    return &AlumniRepository{
        collection: db.Collection("alumni"), // Ganti 'users' dengan nama koleksi Anda
    }
}

// CreateUser menyimpan pengguna baru ke MongoDB.
func (r *AlumniRepository) CreateAlumni(ctx context.Context, alumni *models.Alumni) (*models.Alumni, error) {
    // Pastikan ID tidak disetel saat InsertOne
    alumni.ID = primitive.NilObjectID

    result, err := r.collection.InsertOne(ctx, alumni)
    if err != nil {
        return nil, err
    }
    // Dapatkan ID yang baru dibuat dari hasil insert
    alumni.ID = result.InsertedID.(primitive.ObjectID)
    return alumni, nil
}

// FindAlumniByID mengambil pengguna berdasarkan ID string.
func (r *AlumniRepository) FindAlumniByID(ctx context.Context, id string) (*models.Alumni, error) {
    // Konversi ID string menjadi primitive.ObjectID
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err // ID tidak valid
    }

    var alumni models.Alumni
    filter := bson.M{"_id": objID}

    err = r.collection.FindOne(ctx, filter).Decode(&alumni)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil // Dokumen tidak ditemukan, bukan error fatal
        }
        return nil, err
    }
    return &alumni, nil
}

// FindAllUsers mengambil semua pengguna dari koleksi.
func (r *AlumniRepository) FindAllAlumni(ctx context.Context) ([]models.Alumni, error) {
    cursor, err := r.collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx) // Pastikan kursor ditutup

    var alumni []models.Alumni
    if err = cursor.All(ctx, &alumni); err != nil {
        return nil, err
    }
    return alumni, nil
}
