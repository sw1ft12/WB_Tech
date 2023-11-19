package repository

type Order struct {
	Id                string  `json:"order_uid" validate:"required"`
	TrackNumber       string  `json:"track_number" validate:"required"`
	Entry             string  `json:"entry"`
	User              User    `json:"delivery" validate:"dive"`
	Payment           Payment `json:"payment" validate:"dive"`
	Items             []Item  `json:"items" validate:"dive"`
	Locale            string  `json:"locale"`
	InternalSignature string  `json:"internal_signature"`
	CustomerId        string  `json:"customer_id"`
	DeliveryService   string  `json:"delivery_service"`
	ShardKey          string  `json:"shard_key"`
	SmId              int     `json:"sm_id"`
	DateCreated       string  `json:"date_created"`
	OofShar           string  `json:"oof_shar"`
}

type User struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int    `json:"amount" validate:"gt=0"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost" validate:"ltfield=Amount"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"gt=0"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price" validate:"gt=0"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
