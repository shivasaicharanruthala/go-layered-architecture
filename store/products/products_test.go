package products

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"layeredArchitecture/model"
	"reflect"
	"regexp"
	"testing"
)

func TestProductStore_Get(t *testing.T) {
	fdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fdb.Close()

	ProductMockDb := productStore{Db: fdb}

	rows := sqlmock.NewRows([]string{"id", "name", "brand"}). // define a struct in place of brand
									AddRow(1, "biscuts", 1)

	mock.ExpectQuery("SELECT id,name,brand FROM products*").WithArgs(1).WillReturnRows(rows)

	output, err := ProductMockDb.GetByIdProducts(1)
	if err != nil {
		t.Errorf("wanted nil err but got %v", err)
	}
	expectedOutput := &model.Products{Id: 1, Name: "biscuts", Brand: model.Brand{Id: 1}}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("test failed, got %v but want %v", output, expectedOutput)
	}


	mock.ExpectQuery("SELECT id,name,brand FROM products WHERE id=\\?").WithArgs(5).WillReturnError(errors.New("Product Id not Found"))
	output2, err2 := ProductMockDb.GetByIdProducts(5)

	if err2.Error() != errors.New("Product Id not Found").Error() {
		t.Errorf("wanted %v but got %v", errors.New("Product Id not Found"), err2)
	}
	expectedOutput2 := &model.Products{}
	if !reflect.DeepEqual(expectedOutput2, output2) {
		t.Errorf("got %v but wanted %v", expectedOutput2, output2)
	}
}

func TestProductStore_PostProduct(t *testing.T) {
	fdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fdb.Close()

	productMockDB := productStore{Db: fdb}

	// test case 1 - correct body
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO products (name, brand) VALUES (?, ?)`)).WithArgs("oreo", 1).WillReturnResult(sqlmock.NewResult(7, 1))

	testProduct := model.Products{Name: "oreo", Brand: model.Brand{Id: 1}}
	output, err := productMockDB.PostProduct(testProduct)

	if err != nil {
		t.Errorf("wanted nil err but got %v", err)
	}

	var expectedOutput int = 7

	if !reflect.DeepEqual(expectedOutput, output.Id) {
		t.Errorf("Expected ")
	}

	//	test case2
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO products (name, brand) VALUES (?, ?)`)).WithArgs("oreo", 15).WillReturnError(errors.New("brand id missing"))
	testProduct2 := model.Products{Name: "oreo", Brand: model.Brand{Id: 15}}
	output2, err2 := productMockDB.PostProduct(testProduct2)
	if !reflect.DeepEqual(err2, errors.New("brand id missing")) {
		t.Errorf("wanted nil err but got %v", err)
	}

	var expectedOutput2  = &model.Products{}

	if !reflect.DeepEqual(expectedOutput2, output2) {
		t.Errorf("Expected ")
	}
}

func TestProductStore_DeleteById(t *testing.T) {
	fdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fdb.Close()

	productMockDB := productStore{Db: fdb}

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM products WHERE id=?`)).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	output, err := productMockDB.DeleteById(1)
	if !reflect.DeepEqual(err, nil){
		t.Errorf("wanted nil err but got %v", err)
	}
	expectedOutput := 1
	if !reflect.DeepEqual(expectedOutput, output){
		t.Errorf("test failed got %v but wanted %v", output,expectedOutput)
	}
}

func TestProductStore_UpdateProduct(t *testing.T) {
	fdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fdb.Close()

	productMockDB := productStore{Db: fdb}
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET name=?, brand=? WHERE id=?`)).WithArgs("good day", 1, 3).WillReturnResult(sqlmock.NewResult(0, 1))

	inp := model.Products{Id: 3, Name: "good day", Brand: model.Brand{Id: 1}}
	output, err := productMockDB.UpdateProduct(inp)
	if !reflect.DeepEqual(err, nil){
		t.Errorf("wanted nil err but got %v", err)
	}
	expectedOutput := 1
	if !reflect.DeepEqual(expectedOutput, output){
		t.Errorf("test failed got %v but wanted %v", output,expectedOutput)
	}

}
