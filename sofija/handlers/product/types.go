package product

type InsertRequestData struct {
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	Price            int64 `json:"price"`
}

type InsertResponseData struct {
	Product EgwProductModel
}

type UpdateRequestData struct {
	ShortDescription string
	Description      string
	Price            int64
}

type UpdateResponseData struct {
	Product EgwProductModel
}
