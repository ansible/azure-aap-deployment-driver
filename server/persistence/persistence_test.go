package persistence_test

import (
	"os"
	"server/model"
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
	testData := &telemetry.Telemetry{
		BaseModel:   model.BaseModel{},
		MetricName:  telemetry.DeployStatus,
		MetricValue: "SUCCESS",
		Step:        "1",
	}
	db.Instance.Save(&testData)
	retrieved := &telemetry.Telemetry{}
	db.Instance.First(retrieved)
	assert.Equal(t, "SUCCESS", retrieved.MetricValue)
	sqlDb, _ := db.Instance.DB()
	sqlDb.Close()
}

// TestMain wraps the tests.  Setup is done before the call to m.Run() and any
// needed teardown after that.
func TestMain(m *testing.M) {
	test.SetEnvironment()
	m.Run()
}
