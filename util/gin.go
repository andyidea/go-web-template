package util

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/mcuadros/go-defaults"
)

func GinBind(c *gin.Context, obj interface{}) error {
	contentType := c.ContentType()
	if contentType == "application/json" {
		return GinBindJSON(c, obj)
	}
	defaults.SetDefaults(obj)

	err := c.ShouldBind(obj)
	if err != nil {
		return err
	}
	valid := validation.Validation{}
	b, err := valid.Valid(obj)
	if err != nil {
		return nil
	}
	if !b {
		if valid.HasErrors() {
			var errstrs []string
			for _, err := range valid.Errors {
				errstrs = append(errstrs, err.Field+" "+err.String())
			}

			return errors.New(strings.Join(errstrs, ";"))
		}
	}

	log.Debug().Interface("param", obj)

	return nil

}

func GinBindJSON(c *gin.Context, obj interface{}) error {
	defaults.SetDefaults(obj)

	err := c.ShouldBindJSON(obj)
	if err != nil {
		return err
	}
	valid := validation.Validation{}
	b, err := valid.Valid(obj)
	if err != nil {
		return nil
	}
	if !b {
		if valid.HasErrors() {
			var errstrs []string
			for _, err := range valid.Errors {
				errstrs = append(errstrs, err.Field+" "+err.String())
			}

			return errors.New(strings.Join(errstrs, ";"))
		}
	}

	return nil

}
