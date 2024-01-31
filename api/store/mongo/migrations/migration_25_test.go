package migrations

import (
	"context"
	"testing"
	"time"

	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/stretchr/testify/assert"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestMigration25(t *testing.T) {

	migrations := GenerateMigrations()[:24]

	migrates := migrate.NewMigrate(mongoClient.Database("test"), migrations...)
	err := migrates.Up(migrate.AllAvailable)
	assert.NoError(t, err)

	version, _, err := migrates.Version()
	assert.NoError(t, err)
	assert.Equal(t, uint64(24), version)

	namespace := models.Namespace{
		Name:     "name",
		Owner:    "owner",
		TenantID: "tenant",
	}
	_, err = mongoClient.Database("test").Collection("namespaces").InsertOne(context.TODO(), namespace)
	assert.NoError(t, err)

	device := models.Device{
		Name:     "device",
		UID:      "uid",
		Identity: &models.DeviceIdentity{MAC: "mac"},
		TenantID: "tenant",
		LastSeen: time.Now(),
	}
	_, err = mongoClient.Database("test").Collection("devices").InsertOne(context.TODO(), device)
	assert.NoError(t, err)

	device = models.Device{
		Name:     "device2",
		UID:      "uid2",
		Identity: &models.DeviceIdentity{MAC: "mac"},
		TenantID: "tenant2",
		LastSeen: time.Now(),
	}
	_, err = mongoClient.Database("test").Collection("devices").InsertOne(context.TODO(), device)
	assert.NoError(t, err)

	device = models.Device{
		Name:     "device3",
		UID:      "uid3",
		Identity: &models.DeviceIdentity{MAC: "mac"},
		TenantID: "tenant3",
		LastSeen: time.Now(),
	}
	_, err = mongoClient.Database("test").Collection("devices").InsertOne(context.TODO(), device)
	assert.NoError(t, err)

	migration := GenerateMigrations()[24]

	migrates = migrate.NewMigrate(mongoClient.Database("test"), migration)
	err = migrates.Up(migrate.AllAvailable)
	assert.NoError(t, err)

	version, _, err = migrates.Version()
	assert.NoError(t, err)
	assert.Equal(t, uint64(25), version)

	var migratedDevice *models.Device
	err = mongoClient.Database("test").Collection("devices").FindOne(context.TODO(), bson.M{"tenant_id": "tenant"}).Decode(&migratedDevice)
	assert.NoError(t, err)
	assert.Equal(t, "device", migratedDevice.Name)

	err = mongoClient.Database("test").Collection("devices").FindOne(context.TODO(), bson.M{"tenant_id": "tenant2"}).Decode(&models.Namespace{})
	assert.EqualError(t, mongo.ErrNoDocuments, err.Error())

	err = mongoClient.Database("test").Collection("devices").FindOne(context.TODO(), bson.M{"tenant_id": "tenant3"}).Decode(&models.Namespace{})
	assert.EqualError(t, mongo.ErrNoDocuments, err.Error())
}
