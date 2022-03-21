package database

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/c2fo/testify/suite"
	file "github.com/toolboxfile"
)

type SafeJsonFileTestSuite struct {
	suite.Suite
}

type elem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (sjfts *SafeJsonFileTestSuite) TestCreateDB() {
	fileName, rmFn := file.NewTestFile(sjfts.T(), 0600)
	sjfts.T().Logf("Temp file: %s", fileName)
	defer rmFn()

	db, err := NewSafeJsonFile(fileName)
	sjfts.NoError(err)
	sjfts.NotNil(db)

	sjfts.NoError(db.Save(elem{
		ID:   1,
		Name: "TestElemA",
	}))

	var testElem elem

	sjfts.NoError(db.Load(&testElem))

	sjfts.Equal(1, testElem.ID)

}

func (sjfts *SafeJsonFileTestSuite) TestCannotCreateDB() {
	fileName := filepath.Join(os.TempDir(), fmt.Sprintf("test-%d.json", time.Now().Unix()))
	sjfts.T().Logf("Temp folder: %s", fileName)

	sjfts.NoError(os.Mkdir(fileName, os.ModePerm))
	defer func() {
		sjfts.NoError(os.Remove(fileName))
	}()

	db, err := NewSafeJsonFile(fileName)

	sjfts.Nil(db)

	sjfts.EqualError(err, fmt.Sprintf("Failed to ensure Json DB file: open %s: is a directory", fileName))
}

func (sjfts *SafeJsonFileTestSuite) TestMalformedDBJson() {
	fileName, rmFn := file.NewTestFile(sjfts.T(), 0600, "A")
	sjfts.T().Logf("Temp file: %s", fileName)
	defer rmFn()

	db, err := NewSafeJsonFile(fileName)
	sjfts.NoError(err)

	sjfts.EqualError(db.Load(""), "invalid character 'A' looking for beginning of value")
}

func (sjfts *SafeJsonFileTestSuite) TestUnreadableDBJson() {
	fileName, rmFn := file.NewTestFile(sjfts.T(), 0600)
	sjfts.T().Logf("Temp file: %s", fileName)

	db, err := NewSafeJsonFile(fileName)
	sjfts.NoError(err)

	rmFn()

	sjfts.EqualError(db.Load(""), fmt.Sprintf("Failed to read DB Json: open %s: no such file or directory", fileName))
}

func (sjfts *SafeJsonFileTestSuite) TestMalformedJsonInput() {
	fileName, rmFn := file.NewTestFile(sjfts.T(), 0600)
	sjfts.T().Logf("Temp file: %s", fileName)
	defer rmFn()

	db, err := NewSafeJsonFile(fileName)
	sjfts.NoError(err)

	sjfts.EqualError(db.Save(make(chan bool)), "Failed to encode data as Json: json: unsupported type: chan bool")
}

func (sjfts *SafeJsonFileTestSuite) TestPing() {
	fileName, rmFn := file.NewTestFile(sjfts.T(), 0600)
	sjfts.T().Logf("Temp file: %s", fileName)

	db, err := NewSafeJsonFile(fileName)
	sjfts.NoError(err)

	sjfts.NoError(db.Ping())

	rmFn()

	sjfts.EqualError(db.Ping(), fmt.Sprintf("Failed to check Json DB: stat %s: no such file or directory", fileName))
}

// TestJsonDB runs the whole test suite
func TestJsonDB(t *testing.T) {
	suite.Run(t, new(SafeJsonFileTestSuite))
}
