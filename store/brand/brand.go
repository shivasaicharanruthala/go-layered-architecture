package brand

import (
	"database/sql"
	"errors"
	"fmt"
	"layeredArchitecture/model"
	"layeredArchitecture/store"
)

type brandStore struct {
	Db *sql.DB
}

func New(db *sql.DB) store.BrandInterface {
	return &brandStore{Db: db}
}

func (b *brandStore) GetByIdBrand(id int) (*model.Brand, error) {
	var (
		rows *sql.Rows
		err  error
	)
	rows, err = b.Db.Query("SELECT * FROM brand WHERE id=?", id)
	if err != nil {
		return &model.Brand{}, errors.New("Brand Id Not Found") //nil, err
	}
	var brand model.Brand
	for rows.Next() {
		err := rows.Scan(&brand.Id, &brand.Name)
		if err != nil {
			return &model.Brand{}, err
		}
	}
	fmt.Println("datastore layer: ", brand)
	return &brand, nil
}
func (b *brandStore) PostBrand(brandDetails model.Brand) (*model.Brand, error) {
	resp, err := b.Db.Exec("INSERT INTO brand (name) VALUES (?)", brandDetails.Name)
	if err != nil {
		return &model.Brand{}, err
	}

	lastInsertedId, err := resp.LastInsertId()
	if err != nil {
		return &model.Brand{}, err
	}

	brandDetails.Id = int(lastInsertedId)
	return &brandDetails, nil
}

func (b *brandStore) CheckBrandExsists(bName string) (int, error) {
	resp, err := b.Db.Query("SELECT id FROM brand WHERE name=?", bName)
	if err != nil {
		return 0, err
	}
	var brandId int
	for resp.Next() {
		err := resp.Scan(&brandId)
		if err != nil {
			return 0, err
		}
	}
	return brandId, nil
}
