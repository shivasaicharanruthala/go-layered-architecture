package brand

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"layeredArchitecture/model"
	"testing"
)

func TestBrandStore_Get(t *testing.T) {
	fdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fdb.Close()
	BrandMockDb := brandStore{Db: fdb}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "brand1")

	mock.ExpectQuery("SELECT * FROM brand WHERE id=\\?").WithArgs(1).WillReturnRows(rows)

	output, err := BrandMockDb.GetByIdBrand(1)
	if err != nil {
		t.Errorf("wanted nil err but got %v", err)
	}

	expectedOutput := model.Brand{Id: 1, Name: "brand1"}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.Equal(t, expectedOutput, output)

	// test case 2
	mock.ExpectQuery("SELECT * FROM brand WHERE id=\\?").WithArgs(5).WillReturnError(errors.New("Brand Id not found"))
	output2, err2 := BrandMockDb.GetByIdBrand(5)
	if err2 != errors.New("Brand Id not found") {
		t.Errorf("Got %v but wanted %v", err2, errors.New("Brand Id not Found"))
	}
	expectedOutput2 := &model.Brand{Id: 0, Name: ""}
	if expectedOutput2 != output2 {
		t.Errorf("Got %v but wanted %v", expectedOutput2, output2)
	}
}
