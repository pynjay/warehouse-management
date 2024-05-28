package errors

import "github.com/labstack/echo/v4"

func SendJson(err ResponseError, ctx echo.Context) {
	if validationErr, ok := err.(ValidationError); ok {
		validationErr.ValidationErrors()
		ctx.JSON(
			err.ErrorStatus(),
			&ValidationProblemJson{
				ProblemJson{
					Title:  err.ErrorTitle(),
					Status: err.ErrorStatus(),
					Detail: err.ErrorDetail(),
				},
				validationErr.ValidationErrors(),
			},
		)
	} else {
		ctx.JSON(
			err.ErrorStatus(),
			&ProblemJson{
				Title:  err.ErrorTitle(),
				Status: err.ErrorStatus(),
				Detail: err.ErrorDetail(),
			},
		)
	}
}
