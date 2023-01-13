package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const testUrl = "https://swapi.dev/api/films/%d/?format=json"

func NewApiServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		ListenAddress: listenAddress,
		store:         store,
	}

}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/films/{id}", MakeHttpHandleFunc(s.Getfilms))
	router.HandleFunc("/films", MakeHttpHandleFunc(s.GetAllfilms))
	router.HandleFunc("/characters", MakeHttpHandleFunc(s.GetCharacters))
	router.HandleFunc("/films/{id}/characters", MakeHttpHandleFunc(s.GetCharactersByFilm))
	router.HandleFunc("/comments", MakeHttpHandleFunc(s.HandleGetComments))
	router.HandleFunc("/films/{id}/comments", MakeHttpHandleFunc(s.HandleComments))

	log.Println("json api server is running on port:", s.ListenAddress)
	http.ListenAndServe(s.ListenAddress, router)

}
func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-type", "application/json")

	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)

}

func MakeHttpHandleFunc(f Apifunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}

}
func (s *APIServer) Getfilms(w http.ResponseWriter, r *http.Request) error {

	id, err := GetId(r)

	if err != nil {
		return err
	}
	url := fmt.Sprintf(testUrl, id)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("permission denied")
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}

	defer resp.Body.Close()

	var films Film
	err = json.NewDecoder(resp.Body).Decode(&films)
	if err != nil {
		return err

	}

	return WriteJson(w, http.StatusOK, films)

}

func GetId(r *http.Request) (int, error) {
	//mux.vars allows me to pass an id in my route
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("permission denied")
	}
	return id, nil

}

func (s *APIServer) GetAllfilms(w http.ResponseWriter, r *http.Request) error {
	//how to fetch url

	req, err := http.Get("https://swapi.dev/api/films/?format=json")
	if err != nil {
		return err

	}
	defer req.Body.Close()
	//instance
	var films FilmsEnc

	resp, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(resp, &films)

	WriteJson(w, http.StatusOK, films)

	return err

}

func (s *APIServer) GetCharacters(w http.ResponseWriter, r *http.Request) error {
	const characUrl = "https://swapi.dev/api/people/?format=json"

	resp, err := http.Get(characUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("permission denied")
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var characData Character
	json.Unmarshal(respBody, &characData)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, characData)

}

func (s *APIServer) GetCharactersByFilm(w http.ResponseWriter, r *http.Request) error {
	id, err := GetId(r)
	if err != nil {
		return err
	}
	newUrls := fmt.Sprintf(testUrl, id)
	resp, err := http.Get(newUrls)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("permission denied")
	}
	newData, _ := io.ReadAll(resp.Body)
	var newFilmData Film
	err = json.Unmarshal(newData, &newFilmData)
	if err != nil {
		return err
	}

	for _, People := range newFilmData.Characters {
		var p Person
		newData, err := convertUrl(People, p)
		if err != nil {
			return err
		}
		WriteJson(w, http.StatusOK, newData)
	}

	return err

}

func convertUrl(path string, out interface{}) (interface{}, error) {
	resp, err := http.Get(path)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return out, err
	}
	return out, err

}

func GetCharactersById(w http.ResponseWriter, r *http.Request) error {
	id, err := GetId(r)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://swapi.dev/api/people/%d", id)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("permission denied")
	}
	var person Person
	err = json.NewDecoder(resp.Body).Decode(&person)
	if err != nil {
		return err
	} else {
		WriteJson(w, http.StatusOK, person)
	}
	return err
}

func (s *APIServer) HandleComments(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return fmt.Errorf("permission denied")
	} else if r.Method == "POST" {
		return s.HandleFilmComments(w, r)

	}
	if r.Response.StatusCode != http.StatusOK {
		return fmt.Errorf("permission denied")
	}
	return fmt.Errorf("permission denied")

}

func (s *APIServer) HandleGetComments(w http.ResponseWriter, r *http.Request) error {
	comment, err := s.store.Getcomments()
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, comment)

}

func (s *APIServer) HandleFilmComments(w http.ResponseWriter, r *http.Request) error {

	id, err := GetId(r)

	if err != nil {
		return err
	}

	commenturl := fmt.Sprintf(testUrl, id)

	resp, err := http.Get(commenturl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("permission denied")
	}

	var filmdata CommentFilm
	//var filmlist Film

	json.NewDecoder(resp.Body).Decode(&filmdata.Result)

	if err := json.NewDecoder(r.Body).Decode(&filmdata.Comment); err != nil {
		return err
	}
	//var createCommentReq CreateCommentRequest

	comment := NewComment(filmdata.Comment.Text, filmdata.Comment.IpAddress)

	if err := s.store.CreateComment(comment); err != nil {
		return err
	}

	fmt.Println(filmdata)
	return WriteJson(w, http.StatusOK, filmdata)

}

func (s *APIServer) HandleGetfilmComments(w http.ResponseWriter, r *http.Request) error {

	return nil
}
