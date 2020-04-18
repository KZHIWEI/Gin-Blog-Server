package main

import (
	"encoding/json"
	"errors"
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

func ViewPost(id string) (map[string]string,error) {
	query := "SELECT idPosts,Title,Content,Images,PostDate from posts where idPosts = ?"
	rows,err:=SqlDB.Query(query, id)
	if err != nil {
		return nil, err
	}
	idPost := ""
	title := ""
	content := ""
	images := ""
	postDate := ""
	for rows.Next() {
		err = rows.Scan(&idPost,&title,&content,&images,&postDate)
		if err != nil {
			return nil, err
		}
		break
	}
	if idPost == "" {
		rows.Close()
		return nil, errors.New("post does not exist")
	}
	maps := make(map[string]string)
	maps["id"] = idPost
	maps["title"] = title
	maps["content"] = content
	maps["images"] = images
	maps["date"] = postDate
	return maps,rows.Close()

}
