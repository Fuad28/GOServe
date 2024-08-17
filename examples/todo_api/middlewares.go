package main

import (
	"strconv"

	"github.com/Fuad28/GOServe.git/goserve"
	"github.com/Fuad28/GOServe.git/goserve/status"
)

func authenticationMiddlware(req *goserve.Request, res goserve.IResponse) goserve.IResponse {

	if token, exists := req.Headers().Get("Authorization"); exists {

		// Token authentication logic
		userId, _ := strconv.Atoi(token)
		req.Store.Set("userId", userId)

	} else {
		return res.SetStatus(status.HTTP_401_UNAUTHORIZED).Send(
			goserve.JSON{"message": "unauthorized"},
		)
	}

	// pass control to the next handler in the handlerChain
	return req.Next(res)
}
