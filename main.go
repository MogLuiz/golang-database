package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Product struct {
	ID    string
	Name  string
	Price float64
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	product := NewProduct("Macbook", 6999.90)
	err = insertProduct(db, product)
	if err != nil {
		panic(err)
	}

	product.Price = 12990.90
	err = updateProduct(db, product)
	if err != nil {
		panic(err)
	}

	searchedProduct, err := showProduct(db, "96114a79-addb-4f45-ac55-881108077bcd")
	if err != nil {
		panic(err)
	}
	fmt.Printf("product: %v, possui o preço de %.2f", searchedProduct.Name, searchedProduct.Price)

	productList, err := listProducts(db)
	if err != nil {
		panic(err)
	}

	for _, p := range productList {
		fmt.Printf("product: %v, possui o preço de %.2f\n", p.Name, p.Price)
	}
}

func insertProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("insert into products(id, name, price) values(?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.ID, product.Name, product.Price)
	if err != nil {
		panic(err)
	}
	return nil
}

func updateProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("update products set name = ?, price = ? where id = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Price, product.ID)
	if err != nil {
		panic(err)
	}
	return nil
}

func showProduct(db *sql.DB, id string) (*Product, error) {
	stmt, err := db.Prepare("select id, name, price from products where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var productData Product

	// QueryRow busca apenas uma linha do meu DB
	// O scan pega o valor de cada coluna e atribui a meu ponteiro
	err = stmt.QueryRow(id).Scan(&productData.ID, &productData.Name, &productData.Price)
	if err != nil {
		return nil, err
	}

	return &productData, nil
}

func listProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("select id, name, price from products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
