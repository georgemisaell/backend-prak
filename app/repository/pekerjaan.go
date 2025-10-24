package repository

import (
	"context"
	"latihan_uts_2/app/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IPekerjaanRepository interface {
	CreatePekerjaan(ctx context.Context, user *models.Pekerjaan) (*models.Pekerjaan, error)
	FindPekerjaanByID(ctx context.Context, id string) (*models.Pekerjaan, error)
	FindAllPekerjaan(ctx context.Context) ([]models.Pekerjaan, error)
	UpdatePekerjaan(ctx context.Context, id string, pekerjaan *models.Pekerjaan) (*models.Pekerjaan, error)
	SoftDeletePekerjaan(ctx context.Context, id string) error
}

// PekerjaanRepository implementasi IPekerjaanRepository 
type PekerjaanRepository struct {
	collection *mongo.Collection
}

// NewPekerjaanRepository 
func NewPekerjaanRepository(db *mongo.Database) IPekerjaanRepository {
	return &PekerjaanRepository{
		collection: db.Collection("pekerjaan"),
	}
}

// CreatePekerjaan menyimpan pengguna baru ke MongoDB.
func (r *PekerjaanRepository) CreatePekerjaan(ctx context.Context, pekerjaan *models.Pekerjaan) (*models.Pekerjaan, error) {
	// Pastikan ID tidak disetel saat InsertOne
	pekerjaan.ID = primitive.NilObjectID

	result, err := r.collection.InsertOne(ctx, pekerjaan)
	if err != nil {
		return nil, err
	}
	// Dapatkan ID yang baru dibuat dari hasil insert
	pekerjaan.ID = result.InsertedID.(primitive.ObjectID)
	return pekerjaan, nil
}

// FindPekerjaanByID mengambil pengguna berdasarkan ID string.
func (r *PekerjaanRepository) FindPekerjaanByID(ctx context.Context, id string) (*models.Pekerjaan, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err // ID tidak valid
	}

	var pekerjaan models.Pekerjaan
	filter := bson.M{"_id": objID}

	err = r.collection.FindOne(ctx, filter).Decode(&pekerjaan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Dokumen tidak ditemukan, bukan error fatal
		}
		return nil, err
	}
	return &pekerjaan, nil
}

// FindAllPekerjaan mengambil semua pengguna dari koleksi.
func (r *PekerjaanRepository) FindAllPekerjaan(ctx context.Context) ([]models.Pekerjaan, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaan []models.Pekerjaan
	if err = cursor.All(ctx, &pekerjaan); err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

// UpdatePekerjaan memperbarui dokumen pekerjaan di MongoDB.
func (r *PekerjaanRepository) UpdatePekerjaan(ctx context.Context, id string, pekerjaan *models.Pekerjaan) (*models.Pekerjaan, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "deleted_at": primitive.DateTime(time.Time{}.UnixMilli())}
	
	update := bson.M{"$set": pekerjaan}

	// Lakukan update
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	var updatedPekerjaan models.Pekerjaan
	err = r.collection.FindOne(ctx, filter).Decode(&updatedPekerjaan)
	if err != nil {
		return nil, err
	}

	return &updatedPekerjaan, nil
}

// SoftDeletePekerjaan menetapkan nilai DeletedAt pada dokumen yang sesuai dengan ID.
func (r *PekerjaanRepository) SoftDeletePekerjaan(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}

	filter := bson.M{"_id": objID, "deleted_at": primitive.DateTime(time.Time{}.UnixMilli())}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}