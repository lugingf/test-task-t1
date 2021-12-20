package lib

import (
	"io/ioutil"

	"github.com/talon-one/assignment-props/rqctx"
)

func postToDatabaseHandler(ctx *rqctx.Context) error {
	requestBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	return ctx.RequestBody.Insert(requestBody)
}
