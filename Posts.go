package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func NewPostsHandler(c *gin.Context) {
	var post Post
	err := c.ShouldBindBodyWith(&post, binding.JSON)
	if err != nil {
		ResponseError(c, err)
		return
	}
	id,err:=post.Create()
	c.JSON(200, gin.H{
		"title":   post.Title,
		"id":id,
		"content": post.Content,
		"images":  post.ImagesUrl,
	})
}

func ViewPostHandler(c *gin.Context)  {
	//claims,err := GetClaims(c)
	//var id int
	//if err== nil {
	//	id = claims["id"].(int)
	//}
	postId := c.Param("id")
	maps,err:=ViewPost(postId)
	if err != nil {
		ResponseError(c,err)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"id":maps["id"],
		"title":maps["title"],
		"content":maps["context"],
		"images":maps["images"],
		"date":maps["date"],
	})

}
