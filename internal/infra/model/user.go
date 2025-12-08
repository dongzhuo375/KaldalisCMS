package model

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 2. 用户名
	// unique: 唯一
	// not null: 非 NULL
	// check: 禁止空串和纯空格
	// json: 前端看到的是 "username"
	Username string `gorm:"unique;not null;check:char_length(TRIM(username)) > 0" json:"username"`

	Email string `gorm:"unique;not null;check:char_length(TRIM(email)) > 0" json:"email"`
	Password string `gorm:"not null;check:char_length(TRIM(password)) > 0" json:"-"`
	Role string `gorm:"not null;default:'user'" json:"role"`

	// 新增：一对多关系
	// foreignKey:AuthorID 指明 Post 表里是用哪个字段关联回来的
	// constraint:OnUpdate:CASCADE,OnDelete:SET NULL; 指明外键约束行为
	// json:"-" 强烈建议加上！防止查询用户信息时带出几千篇文章，导致 JSON 爆炸
	Posts []Post `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

}