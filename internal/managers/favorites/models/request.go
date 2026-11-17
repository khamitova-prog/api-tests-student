package models

// AddFavoriteRequest содержит ID объявления для добавления в избранное.
type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}
