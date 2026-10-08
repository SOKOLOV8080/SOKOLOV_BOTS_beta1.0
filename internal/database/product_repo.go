package database

import (
	"SOKOLOV_BOTS/internal/models"
)

func GetAllProducts() ([]models.Product, error) {
	rows, err := DB.Query("SELECT id, seller_id, title, description, price, github_url, contact_info, file_path, created_at FROM products ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.SellerID, &p.Title, &p.Description, &p.Price, &p.GithubURL, &p.ContactInfo, &p.FilePath, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func CreateProduct(p models.Product) error {
	_, err := DB.Exec(
		"INSERT INTO products (seller_id, title, description, price, github_url, contact_info, file_path) VALUES (?, ?, ?, ?, ?, ?, ?)",
		p.SellerID, p.Title, p.Description, p.Price, p.GithubURL, p.ContactInfo, p.FilePath,
	)
	return err
}
func GetProductByID(id int) (models.Product, error) {
	var p models.Product
	err := DB.QueryRow(
		"SELECT id, seller_id, title, description, price, github_url, contact_info, file_path, created_at FROM products WHERE id = ?",
		id,
	).Scan(&p.ID, &p.SellerID, &p.Title, &p.Description, &p.Price, &p.GithubURL, &p.ContactInfo, &p.FilePath, &p.CreatedAt)
	return p, err
}
func GetProductsBySellerID(sellerID int) ([]models.Product, error) {
	rows, err := DB.Query("SELECT id, seller_id, title, description, price, github_url, contact_info, file_path, created_at FROM products WHERE seller_id = ? ORDER BY created_at DESC", sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.SellerID, &p.Title, &p.Description, &p.Price, &p.GithubURL, &p.ContactInfo, &p.FilePath, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func DeleteProduct(id, sellerID int) error {
	_, err := DB.Exec("DELETE FROM products WHERE id = ? AND seller_id = ?", id, sellerID)
	return err
}
