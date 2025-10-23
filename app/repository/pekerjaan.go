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

// PekerjaanRepository implementasi IPekerjaanRepository menggunakan MongoDB.
type PekerjaanRepository struct {
	collection *mongo.Collection
}

// NewPekerjaanRepository membuat instance baru dari PekerjaanRepository dan mengaitkannya dengan koleksi MongoDB.
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
	// Konversi ID string menjadi primitive.ObjectID
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
	defer cursor.Close(ctx) // Pastikan kursor ditutup

	var pekerjaan []models.Pekerjaan
	if err = cursor.All(ctx, &pekerjaan); err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

// UpdatePekerjaan memperbarui dokumen pekerjaan di MongoDB.
func (r *PekerjaanRepository) UpdatePekerjaan(ctx context.Context, id string, pekerjaan *models.Pekerjaan) (*models.Pekerjaan, error) {
	// Konversi ID string menjadi primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Setel UpdatedAt di layer Service, bukan di sini.
	// Kita buat filter dan update operation.
	filter := bson.M{"_id": objID, "deleted_at": primitive.DateTime(time.Time{}.UnixMilli())}
	
	// Data yang akan di-set
	update := bson.M{"$set": pekerjaan}

	// Lakukan update
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Dapatkan kembali data yang diperbarui untuk dikembalikan ke Service/Controller
	var updatedPekerjaan models.Pekerjaan
	err = r.collection.FindOne(ctx, filter).Decode(&updatedPekerjaan)
	if err != nil {
		return nil, err // Mungkin terjadi jika dokumen hilang setelah update (sangat jarang)
	}

	return &updatedPekerjaan, nil
}

// SoftDeletePekerjaan menetapkan nilai DeletedAt pada dokumen yang sesuai dengan ID.
func (r *PekerjaanRepository) SoftDeletePekerjaan(ctx context.Context, id string) error {
	// Konversi ID string menjadi primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Setel DeletedAt menjadi waktu saat ini (soft delete)
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(), // time.Now() akan diubah MongoDB menjadi tipe BSON Date
		},
	}

	filter := bson.M{"_id": objID, "deleted_at": primitive.DateTime(time.Time{}.UnixMilli())} // Hanya hapus yang belum dihapus

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments // Tidak ada dokumen yang dimodifikasi, mungkin ID tidak ditemukan atau sudah dihapus
	}

	return nil
}