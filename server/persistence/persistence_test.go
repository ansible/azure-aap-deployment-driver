package persistence_test

import (
	"os"
	"server/persistence"
	"server/telemetry"
	"server/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func testDb(t *testing.T, db *persistence.Database) {
	type Entity struct {
		gorm.Model
		Value string
	}
	db.Instance.AutoMigrate(&Entity{})
	testData := Entity{
		Value: "test",
	}
	db.Instance.Save(&testData)
	retrieved := &Entity{}
	db.Instance.First(retrieved)
	assert.Equal(t, "test", retrieved.Value)
	sqlDb, _ := db.Instance.DB()
	sqlDb.Close()
}
func TestRealDatabase(t *testing.T) {
	const DB_FILENAME = "tmp.db"
	db := persistence.NewPersistentDB(DB_FILENAME)
	testDb(t, db)
	assert.FileExists(t, DB_FILENAME)
	os.Remove(DB_FILENAME)
	assert.NoFileExists(t, DB_FILENAME)
}

func TestInMemoryDatabase(t *testing.T) {
	db := persistence.NewInMemoryDB()
	testDb(t, db)
}
func TestTelemetryTable(t *testing.T) {

	db := persistence.NewInMemoryDB()
	telemetry.SetMetric(db.Instance, telemetry.DeployStatus, "SUCCESS")
	telemetry.SetMetric(db.Instance, telemetry.AccessType, "PRIVATE")
	retrieved := telemetry.Telemetry{}
	retrieved = telemetry.Metric(db.Instance, telemetry.DeployStatus)
	assert.Equal(t, "SUCCESS", retrieved.MetricValue)
	retrieved = telemetry.Metric(db.Instance, telemetry.CustomerSubscriptionID)
	assert.Equal(t, "", retrieved.MetricValue)
	sqlDb, _ := db.Instance.DB()
	sqlDb.Close()
}

// TestMain wraps the tests.  Setup is done before the call to m.Run() and any
// needed teardown after that.
func TestMain(m *testing.M) {
	test.SetEnvironment()
	m.Run()
}
