package products

import (
	"database/sql"
	"layeredArchitecture/model"
	"layeredArchitecture/store"
)

type productStore struct {
	Db *sql.DB
}

func New(db *sql.DB) store.ProductsInterface {
	return &productStore{Db: db}
}

func (p *productStore) GetByIdProducts(id int) (*model.Products, error) {
	var (
		rows *sql.Rows
		err  error
	)

	rows, err = p.Db.Query("SELECT id,name,brand FROM products WHERE id=?", id)

	if err != nil {
		return &model.Products{}, err //nil, err
	}

	var product model.Products

	for rows.Next() {
		err = rows.Scan(&product.Id, &product.Name, &product.Brand.Id)
		if err != nil {
			return &model.Products{}, err //nil, err
		}
	}
	return &product, nil
}

func (p *productStore) PostProduct(product model.Products) (*model.Products, error) {
	res, err := p.Db.Exec("INSERT INTO products (name, brand) VALUES (?, ?)", product.Name, product.Brand.Id)
	if err != nil {
		return &model.Products{}, err
	}
	lastInsertedId, _ := res.LastInsertId()

	product.Id = int(lastInsertedId)
	return &product, nil
}

func (p *productStore) DeleteById(id int) (int, error) {
	res, err := p.Db.Exec("DELETE FROM products WHERE id=?", id)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := res.RowsAffected()
	return int(rowsAffected), nil
}

func (p *productStore) UpdateProduct(prod model.Products) (int, error) {
	res, err := p.Db.Exec("UPDATE products SET name=?, brand=? WHERE id=?", prod.Name, prod.Brand.Id, prod.Id)
	if err != nil {
		return 0, err
	}
	rowsChanged, _ := res.RowsAffected()
	return int(rowsChanged), nil
}
