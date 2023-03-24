package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Book struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Published int      `json:"published"`
	Pages     int      `json:"pages"`
	Genres    []string `json:"genres"`
	Rating    float32  `json:"rating"`
}

type BookResponse struct {
	Book *Book `json:"book"`
}

type BooksResponse struct {
	Books *[]Book `json:"books"`
}

type ReadinglistModel struct {
	Endpoint string
}

func (m *ReadinglistModel) GetAll() (*[]Book, error) {
	resp, err := http.Get(m.Endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var booksResp BooksResponse
	err = json.Unmarshal(data, &booksResp)
	if err != nil {
		return nil, err
	}

	return booksResp.Books, nil
}

func (m *ReadinglistModel) Get(id int64) (*Book, error) {
	url := fmt.Sprintf("%s/%d", m.Endpoint, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var bookResp BookResponse
	err = json.Unmarshal(data, &bookResp)
	if err != nil {
		return nil, err
	}

	return bookResp.Book, nil
}
