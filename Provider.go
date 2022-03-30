package main

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CommentInfo struct {
	Id       int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ParentId int    `json:"parentId"`
	RootId   int    `json:"rootId"`
	Uid      int    `json:"uid"`
	Comment  string `json:"comment"`
}

type JwtClalims struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type Claims struct {
	Uid      int    `json:"uid" gorm:"primary_key;AUTO_INCREMENT"`
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

func DataInit() {
	db, err := gorm.Open(sqlite.Open("UserInfo.db"), &gorm.Config{}) //打開db
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&Claims{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&CommentInfo{})
	if err != nil {
		return
	}

	//db.Create(&Claims{Username: "Lili", Password: "111"})
}
