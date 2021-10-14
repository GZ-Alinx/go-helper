package resp

// custom errors

const (
	Ok                  = 201
	NotOk               = 405
	Unauthorized        = 401
	Forbidden           = 403
	InternalServerError = 500
)

const (
	OkMsg                      = "success"
	NotOkMsg                   = "failed"
	UnauthorizedMsg            = "login expired, please login again"
	LoginCheckErrorMsg         = "wrong username or password"
	ForbiddenMsg               = "no permission to access this resource"
	InternalServerErrorMsg     = "server internal error"
	IdempotenceTokenEmptyMsg   = "idempotent token is empty"
	IdempotenceTokenInvalidMsg = "idempotent token expired"
	UserDisabledMsg            = "the account has been disabled"
)

var CustomError = map[int]string{
	Ok:                  OkMsg,
	NotOk:               NotOkMsg,
	Unauthorized:        UnauthorizedMsg,
	Forbidden:           ForbiddenMsg,
	InternalServerError: InternalServerErrorMsg,
}
