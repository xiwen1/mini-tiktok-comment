package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type Comment struct {
	ID         int64
	user       int64
	content    string
	createDate string // 格式为 mm-dd
	video_id   int64
}

type User struct {
	ID             int64
	name           string
	follow_count   int64
	follower_count int64
	is_follow      bool
}

type CommentActionRequest struct {
	token        string
	video_id     int64
	type_code    int
	comment_text string
	comment_id   int64
}

type CommentActionResponse struct {
	status_code int
	status_msg  string
	comment     Comment
}

var (
	//connStr = "postgres:RpB27iLmDV4z7ZU5tpkn0UPLQWTQx1zFGaUJixDZQhPght7WWLzfZ8PLhZjavGUZ@tcp(srv.paracraft.club:31294)/nicognaw"
	connStr = "postgres://root:zkw030813@127.0.0.1:5432/root?sslmode=disable"
	pool    *sql.DB
)

func CommentAction(request CommentActionRequest) (response CommentActionResponse, err error) {
	//token := request.token
	if request.type_code == 1 {
		c := Comment{}
		c.user = 1
		c.video_id = request.video_id
		c.content = request.comment_text
		c.createDate = time.Now().Format("01 02")
		fmt.Println(c.createDate)
		c.insertComment()
		checkerr(err)
		response.comment = c
		response.status_msg = "success"
		response.status_code = 1
	} else {
		delete(request.comment_id)
		checkerr(err)
		response.status_msg = "sucess"
		response.status_code = 1
	}
	return
}

func delete(commentId int64) error {
	statement, err := pool.Prepare("DELETE FROM public.comment WHERE id=$1")
	checkerr(err)
	_, err = statement.Exec(commentId)
	return err
}

func checkerr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (c *Comment) insertComment() (err error) {
	var statement *sql.Stmt
	statement, err = pool.Prepare("INSERT INTO public.comment(user_id, content, video_id) VALUES ($1,$2, $3);")
	_, err = statement.Exec(c.user, c.content, c.video_id)
	if err != nil {
		println("error")
		log.Fatal(err.Error())
	}
	return
}

func (u *User) insertUser() error {
	statement, err := pool.Prepare("INSERT INTO public.users(name, follow_count, follower_count, is_follow) VALUES ($1, $2, $3, $4);")
	checkerr(err)
	_, err = statement.Exec(u.name, u.follow_count, u.follower_count, u.is_follow)
	checkerr(err)
	return err
}

func search(id int64) (c Comment, err error) {
	c = Comment{}
	err = pool.QueryRow("SELECT id, content, video_id FROM public.comment WHERE id=$1", 2).Scan(&c.ID, &c.content, &c.video_id)
	return
}

func main() {
	var err error
	pool, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err.Error())
	}
	c := Comment{3, 1, "hello", "1", 1}
	err1 := c.insertComment()
	u1 := User{1, "xiwen", 2, 2, true}
	err = u1.insertUser()
	if err1 != nil {
		return
	}
	req := CommentActionRequest{"hello", 1, 1, "xiwen", 0}
	req2 := CommentActionRequest{"hello", 1, 0, "xiwen", 5}
	resp, err := CommentAction(req)
	checkerr(err)
	resp2, err := CommentAction(req2)
	checkerr(err)
	fmt.Println(resp, resp2)
}
