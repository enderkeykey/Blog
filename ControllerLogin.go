package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func Register(c *gin.Context) {
	var user Claims
	db, _ := gorm.Open(sqlite.Open("userInfo.db"), &gorm.Config{})
	res := db.Where("Username=?", c.PostForm("username")).First(&user)
	if err := res.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			str, _ := Hash(c.PostForm("password"))
			db.Create(&Claims{Username: c.PostForm("username"), Password: str})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "注冊成功"})
			return
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "用户已注冊"})
	return
}

func Login(c *gin.Context) {
	var user Claims
	db, _ := gorm.Open(sqlite.Open("userInfo.db"), &gorm.Config{})

	res := db.Where("Username=?", c.PostForm("username")).First(&user)
	if err := res.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed!"})
			return
		}
	}
	if CompareHash(user.Password, c.PostForm("password")) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed"})
		return
	}

	token, err := GenerateToken(user.Username, user.Uid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "err"})
		return
	}

	c.SetCookie("Gincookie", token, 3000, "/",
		"localhost", false, true)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "login success"})
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端cookie并校验
		cookie, err := c.Cookie("Gincookie")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "cookie error"})
			c.Abort()
			return
		}
		_, err = ParseToken(cookie)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "claims error"})
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login success"})
		c.Next()
		return
	}
}

func Photo(c *gin.Context) {
	//获取文件头
	file, err := c.FormFile("upload")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "請求失敗"})
		return
	}
	//获取文件名
	fileName := file.Filename
	fmt.Println("文件名：", fileName)
	//保存文件到服务器本地
	//SaveUploadedFile(文件头，保存路径)

	if err := c.SaveUploadedFile(file, fileName); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "保存失败"})
		return
	}
	var user Claims
	cookie, _ := c.Cookie("Gincookie")
	claims, _ := ParseToken(cookie)
	db, _ := gorm.Open(sqlite.Open("userInfo.db"), &gorm.Config{})
	db.Where("Username=?", claims.Username).First(&user)
	db.Model(&user).Update("Avatar", fileName)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "上传文件成功"})
}

func Comment(c *gin.Context) { //創建帖子
	var t, p CommentInfo
	cookie, _ := c.Cookie("Gincookie")
	claims, _ := ParseToken(cookie)
	db, err := gorm.Open(sqlite.Open("UserInfo.db"), &gorm.Config{}) //打開db
	if err != nil {
		panic("failed to connect database")
	}

	t.ParentId = -1
	t.Uid = claims.Uid
	t.Comment = c.PostForm("comment")
	db.Create(&CommentInfo{Comment: t.Comment, ParentId: t.ParentId, Uid: t.Uid})
	//更新Root_id
	db.Last(&p)
	db.Model(&p).Update("root_id", p.Id)
	c.JSON(http.StatusUnauthorized, gin.H{"comment": p})
}

func Response(c *gin.Context) {
	var reT, t CommentInfo
	db, err := gorm.Open(sqlite.Open("UserInfo.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	cookie, _ := c.Cookie("Gincookie")
	claims, _ := ParseToken(cookie)

	id, _ := strconv.Atoi(c.PostForm("id"))
	fmt.Println(id)
	db.Where("id=?", id).First(&t)
	fmt.Println(id)
	reT.RootId = t.RootId
	reT.ParentId = id
	fmt.Println(reT.ParentId)
	reT.Uid = claims.Uid
	reT.Comment = c.PostForm("comment")

	db.Create(&CommentInfo{RootId: reT.RootId, ParentId: reT.ParentId, Uid: reT.Uid, Comment: reT.Comment})
	c.JSON(http.StatusUnauthorized, gin.H{"content": reT})
}

func ShowComment(c *gin.Context) {
	var t, p CommentInfo
	db, err := gorm.Open(sqlite.Open("UserInfo.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	postId, _ := strconv.Atoi(c.PostForm("postId"))
	res := db.Where("id=? AND parent_id =?", postId, -1).First(&t)
	if err = res.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Post not found"})
			return
		}
	}
	t = p

	db.Where("root_id =?", postId).Last(&t)
	lastId := t.Id
	parentId := -1
	rootId := t.RootId
	t = p

	for true {
		res = db.Where("root_Id=? AND parent_id =?", rootId, parentId).First(&t)
		if err = res.Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Post not found"})
				break
			}
		}

		if t.Id == lastId {
			c.JSON(http.StatusUnauthorized, gin.H{"content": t})
			break
		}
		parentId = t.Id
		c.JSON(http.StatusUnauthorized, gin.H{"content": t})
		t = p
	}
}
