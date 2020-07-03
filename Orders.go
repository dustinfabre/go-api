package main

import (
	"database/sql"
    "time"
)


type Order struct {
	ID	int `json:"id"`
	Created_at time.Time	`json:"created_ate"`
	Order_name string `json:"order_name"`
	Customer_id string `json:"customer_id"`
}

var orders []Order = []Order{}

func getOrders(db *sql.DB, start, count int) ([]Order, error) {
	rows, err := db.Query(
        "SELECT * FROM public.orders LIMIT $1 OFFSET $2",
        count, start)

    if err != nil {
        return nil, err
    }
	defer rows.Close()
	
    for rows.Next() {
        var o Order
        if err := rows.Scan(&o.ID, &o.Created_at, &o.Order_name, &o.Customer_id); err != nil {
            return nil, err
        }
        orders = append(orders, o)
    }

    return orders, nil
}