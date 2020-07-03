package main

import (
	"database/sql"
    "time"
)


type Order struct {
	ID	int `json:"id"`
	Created_at time.Time	`json:"order_date"`
    Order_name string `json:"order_name"`
	Total_amount float64 `json:"total_amount"`
	Delivered_amount float64 `json:"delivered_amount"`
    Name string `json:"customer_name"`
    Company_name string `json:"customer_company"`
    Rows int `json:"total_row"`
    
}

type Pagination struct {
	Start int `json:"start"`
    End int `json:"end"`
}


func getOrders(db *sql.DB, start, count int) ([]Order, error) {
	rows, err := db.Query(
        "SELECT a.id, a.created_at, a.order_name, COALESCE(SUM(b.price_per_unit * b.quantity), 0) AS total_amount" +
        ", COALESCE(SUM(c.delivered_quantity * b.price_per_unit), 0) AS delivered_amount, d.name, e.company_name, count(*) OVER() AS rows " +
        "FROM orders a " +
        "JOIN order_items b ON a.id = b.order_id " +
		"JOIN customers d ON a.customer_id = d.user_id " +
		"JOIN customer_companies e ON d.company_id = e.company_id " +
		"LEFT JOIN deliveries c ON b.id = c.order_item_id GROUP BY a.id, d.name, e.company_name ORDER BY a.id LIMIT $1 OFFSET $2",
    count, start)

    if err != nil {
        return nil, err
    }
	defer rows.Close()
    
    var orders []Order = []Order{}
    for rows.Next() {
        var o Order
        if err := rows.Scan(&o.ID, &o.Created_at, &o.Order_name, &o.Total_amount,&o.Delivered_amount, &o.Name, &o.Company_name, &o.Rows); err != nil {
            return nil, err
        }
        orders = append(orders, o)
    }

    return orders, nil
}