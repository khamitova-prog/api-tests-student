
package models

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}

type Advertisement struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
