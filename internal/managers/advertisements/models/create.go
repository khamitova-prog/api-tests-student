package models

// PhotoFile содержит имя и данные фотографии для загрузки.
type PhotoFile struct {
	Name string
	Data []byte
}

// CreateRequest содержит поля формы и фотографии объявления.
// Map позволяет пропускать поля в негативных тестах.
type CreateRequest struct {
	Fields map[string]string
	Photos []PhotoFile
}

// CreatedPhoto описывает загруженную фотографию.
type CreatedPhoto struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
}

// CreateResponse описывает созданное объявление.
type CreateResponse struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Quantity    int            `json:"quantity"`
	UserID      string         `json:"user_id"`
	Photos      []CreatedPhoto `json:"photos"`
}

// APIError описывает ответ сервера при ошибке создания объявления.
type APIError struct {
	Details string `json:"details"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
