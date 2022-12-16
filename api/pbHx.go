package api

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"goFix/pb"
	"io/ioutil"
	"net/http"
)

func ReceivePb(c *gin.Context) {
	loginSuccess := &pb.UserInfoREPLY{
		Code:    200,
		Message: "success login",
	}
	c.ProtoBuf(http.StatusOK, loginSuccess)
}

func SendPb(c *gin.Context) {
	resp, err := http.Get("http://127.0.0.1:8080/receivePb")
	if err != nil {
		basicLog.Errorf("failed to get response:%v", err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			basicLog.Errorf("failed to read response body:%v", err)
		} else {
			userInfo := &pb.UserInfoREPLY{}
			err := proto.UnmarshalMerge(body, userInfo)
			if err != nil {
				basicLog.Errorf("failed to unmarshal error response:%v", err)
			}
			basicLog.Debugf("userInfo: %+v", userInfo)
			c.JSON(http.StatusOK, userInfo)
		}
	}
}
