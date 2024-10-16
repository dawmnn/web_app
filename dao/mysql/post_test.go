package mysql

import (
	"testing"
	"web_app/models"
	"web_app/settings"
)

func init() {

	cfg := &settings.MysqlConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "123456",
		DbName:       "itcast",
		Port:         3306,
		MaxOpenConns: 10,
		MaxIdleConns: 10}

	err := Init(cfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:          10,
		AuthorID:    123,
		CommunityID: 1,
		Title:       "test",
		Content:     "just a test",
	}
	err := CreatePost(&post)

	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed,err:%v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}
