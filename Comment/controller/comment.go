package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/auth"
	"github.com/xiwen1/mini-tiktok-comment/Comment/service"
	"log"
	"net/http"
	"strconv"
	time2 "time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg, omitepty"`
}

type CommentListResponse struct {
	StatusCode  int32                 `json:"status_code"`
	StatusMsg   string                `json:"status_msg,omitempty"`
	CommentList []service.CommentInfo `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	StatusCode int32               `json:"status_code"`
	StatusMsg  string              `json:"status_msg,omitempty"`
	Comment    service.CommentInfo `json:"comment"`
}

func CommentAction(r *gin.Context) {
	log.Println("commentController-comment-action: running")
	token, _ := r.Get("token")
	statusCode, userId, err := service.Auth(token.(string))

	if statusCode != int(auth.AuthResponse_SUCCESS) || err != nil {
		r.JSON(200, CommentActionResponse{
			StatusMsg:  "comment auth fail",
			StatusCode: 0,
		})
	}

	videoId, err := strconv.ParseInt(r.Query("video_id"), 10, 64)
	if err != nil {
		r.JSON(http.StatusOK, CommentActionResponse{
			StatusMsg:  "comment video_id invalid",
			StatusCode: 0,
		})
		return
	}

	actionType, err := strconv.ParseInt(r.Query("action_type"), 10, 32)
	if err != nil || actionType < 1 || actionType > 2 {
		r.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "comment actionType invalid",
		})
		return
	}

	if actionType == 1 {
		content := r.Query("comment_text")
		time := time2.Now().Format("01-02")
		sendComment := service.CommentData{}
		sendComment.User = uint32(userId)
		sendComment.Video_id = videoId
		sendComment.Content = content
		sendComment.CreateDate = time
		commentInfo, err := service.Send(sendComment, token.(string))
		if err != nil {
			r.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: 0,
				StatusMsg:  "sending to database fail",
			})
			return
		}
		r.JSON(http.StatusOK, CommentActionResponse{
			StatusMsg:  "comment action success",
			StatusCode: 1,
			Comment:    commentInfo,
		})
	} else {
		commentId, err := strconv.ParseInt(r.Query("comment_id"), 10, 64)
		if err != nil {
			r.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: 0,
				StatusMsg:  "commentAction-delete: commentId invalid",
			})
			return
		}
		err = service.Delete(commentId)
		if err != nil {
			r.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: 0,
				StatusMsg:  "commentAcion-delete: delete fail",
			})
			return
		}
		r.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 1,
			StatusMsg:  "commentAction-delete: success",
		})
		fmt.Println("commentAction: success")
		return
	}

}

func CommentList(r *gin.Context) {
	fmt.Println("commentList: running")
	id, _ := r.Get("userId")
	token, _ := r.Get("token")
	userId, err := strconv.ParseInt(id.(string), 10, 64)
	videoId, err := strconv.ParseInt(r.Query("video_id"), 10, 64)

	// 错误处理:
	if err != nil {
		r.JSON(http.StatusOK, Response{
			StatusMsg:  "comment videoId or userId invalid",
			StatusCode: 0,
		})
		return
	}
	fmt.Printf("userId: %v, videoId: %v", userId, videoId)
	commentList, err := service.SearchComment(uint32(videoId), token.(string))
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	r.JSON(http.StatusOK, CommentListResponse{
		StatusCode: 1,
		StatusMsg: "get commentList successfully",
		CommentList: *commentList,
	})
	return
}
