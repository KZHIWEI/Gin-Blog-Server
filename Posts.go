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
	postId := c.Param("id")
	post,err:=ViewPost(postId)
	if err != nil {
		ResponseError(c,err)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"id":post.Id,
		"title":post.Title,
		"content":post.Content,
		"images":post.ImagesUrl,
		"date":post.PostDate,
	})
}
