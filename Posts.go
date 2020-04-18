package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
