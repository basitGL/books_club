package models

type CartItem struct {
	Book Book 
	Quantity int
}

type Cart struct {
	CartItems CartItem
	TotalPrice float32
}