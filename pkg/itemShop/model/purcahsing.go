package itemShopModel

type BuyingReq struct {
	PlayerID string
	ItemID   uint64 `json:"itemID" validate:"required,gt=0"`
	Quantity uint   `json:"quantity" validate:"required,gt=0"`
}

type SellingReq struct {
	PlayerID string
	ItemID   uint64 `json:"itemID" validate:"required,gt=0"`
	Quantity uint   `json:"quantity" validate:"required,gt=0"`
}
