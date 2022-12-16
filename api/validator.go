package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translational "github.com/go-playground/validator/v10/translations/zh"
	"io"
)

type User struct {
	UserName string `json:"userName" validate:"required,min=0,max=6"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

func Validator(c *gin.Context) {
	validate := validator.New()
	// 中文翻译器
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator("zh")

	err := translational.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println(err)
		return
	}

	r := c.Request
	req := User{}
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		basicLog.Errorf("FlowInteraction getReq err %v bs [%s]\n", err, bs)
		return
	}
	err = json.Unmarshal(bs, &req)
	if err != nil {
		basicLog.Errorf("FlowInteraction Unmarshal err %v bs [%s]\n", err, bs)
		return
	}
	errs := validate.Struct(req)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			c.JSON(500, err.Translate(trans))
		}
	}
}
