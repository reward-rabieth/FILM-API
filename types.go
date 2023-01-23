package main

import (
	"net/http"
	"time"
)

type APIServer struct {
	Store        Storage
	listenadress string
}

type Film struct {
	Title        string    `json:"title"`
	EpisodeID    int       `json:"episode_id"`
	OpeningCrawl string    `json:"opening_crawl"`
	Director     string    `json:"director"`
	Producer     string    `json:"producer"`
	ReleaseDate  string    `json:"release_date"`
	Characters   []string  `json:"characters"`
	Planets      []string  `json:"planets"`
	Starships    []string  `json:"starships"`
	Vehicles     []string  `json:"vehicles"`
	Species      []string  `json:"species"`
	Created      time.Time `json:"created"`
	Edited       time.Time `json:"edited"`
	URL          string    `json:"url"`
}
type Person struct {
	Name      string        `json:"name"`
	Height    string        `json:"height"`
	Mass      string        `json:"mass"`
	HairColor string        `json:"hair_color"`
	SkinColor string        `json:"skin_color"`
	EyeColor  string        `json:"eye_color"`
	BirthYear string        `json:"birth_year"`
	Gender    string        `json:"gender"`
	Homeworld string        `json:"homeworld"`
	Films     []string      `json:"films"`
	Species   []interface{} `json:"species"`
	Vehicles  []string      `json:"vehicles"`
	Starships []string      `json:"starships"`
	Created   time.Time     `json:"created"`
	Edited    time.Time     `json:"edited"`
	URL       string        `json:"url"`
}
type FilmsEnc struct {
	Count int    `json:"count"`
	Films []Film `json:"results"`
}

type Apifunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}
type Character struct {
	Count   int
	Results []Person
}

//	type SavedComment struct {
//		Result    Comment   `json:"results"`
//		FilmData  Film      `json:"film_data"`
//		CreatedAt time.Time `json:"created_at"`
//	}
type CommentFilm struct {
	Comment Comment `json:"comment"`
	Result  Film    `json:"results"`
}

type CreateCommentRequest struct {
	Text      string    `json:"text"`
	IpAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	Text      string    `json:"text"`
	IpAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
	//Result    Film      `json:"results"`S
}

func NewComment(text, ipaddress string) *Comment {

	return &Comment{
		Text:      text,
		IpAddress: ipaddress,
	}

}
