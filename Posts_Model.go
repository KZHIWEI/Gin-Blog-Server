package main

import (
	"encoding/json"
)

type Post struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	ImagesUrl []string `json:"images"`
}

func (post *Post) Create() (int64, error){
	query := "INSERT IGNORE  posts (Title, Content, Images) VALUES (?,?,?)"
	imagesJson , err := json.Marshal(post.ImagesUrl)
	rs, err := SqlDB.Exec(query, post.Title, post.Content,string(imagesJson))
	return HandleSQLResponse(rs, err)
}
