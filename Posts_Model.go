package main

import (
	"encoding/json"
	"errors"
)

type Post struct {
	Id        string   `json:"idPosts"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	ImagesUrl []string `json:"images"`
	PostDate  string   `json:"date"`
}

func (post *Post) Create() (int64, error){
	query := "INSERT IGNORE  posts (Title, Content, Images) VALUES (?,?,?)"
	imagesJson , err := json.Marshal(post.ImagesUrl)
	rs, err := SqlDB.Exec(query, post.Title, post.Content,string(imagesJson))
	return HandleSQLResponse(rs, err)
}

func ViewPost(id string) (Post,error) {
	query := "SELECT idPosts,Title,Content,Images,PostDate from posts where idPosts = ?"
	rows,err:=SqlDB.Query(query, id)
	var post Post
	if err != nil {
		return post, err
	}
	idPost := ""
	title := ""
	content := ""
	images := ""
	postDate := ""
	for rows.Next() {
		err = rows.Scan(&idPost,&title,&content,&images,&postDate)
		if err != nil {
			return post, err
		}
		break
	}
	if idPost == "" {
		rows.Close()
		return post, errors.New("post does not exist")
	}
	post.Id = idPost
	post.Title = title
	post.Content = content
	_ = json.Unmarshal([]byte(images), &post.ImagesUrl)
	post.PostDate = postDate
	return post,rows.Close()

}
