package main

type Post struct {
	Title string `json:"title"`
	Content string `json:"content"`
	ImagesUrl []string `json:"images"`
}