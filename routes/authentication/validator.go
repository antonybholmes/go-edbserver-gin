package authenticationroutes

import (
	"fmt"
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server-gin/routes"
	"github.com/gin-gonic/gin"
)

//
// Standardized data checkers for checking header and body contain
// the correct data for a route
//

type Validator struct {
	c            *gin.Context
	Address      *mail.Address
	LoginBodyReq *auth.LoginBodyReq

	AuthUser *auth.AuthUser
	Claims   *auth.TokenClaims
	Err      error
}

func NewValidator(c *gin.Context) *Validator {
	return &Validator{c: c, Address: nil, LoginBodyReq: nil, AuthUser: nil, Claims: nil, Err: nil}

}

func (validator *Validator) Ok() (*Validator, error) {
	if validator.Err != nil {
		return nil, validator.Err
	} else {
		return validator, nil
	}
}

// If the validator does not encounter errors, it will run the success function
// allowing you to extract data from the validator, otherwise it returns an error
// without running the function
func (validator *Validator) Success(success func(validator *Validator)) {

	if validator.Err != nil {
		return
	}

	success(validator)
}

func (validator *Validator) ParseLoginRequestBody() *Validator {
	if validator.Err != nil {
		return validator
	}

	if validator.LoginBodyReq == nil {
		var req auth.LoginBodyReq

		err := validator.c.Bind(&req)

		if err != nil {
			validator.Err = err
		} else {
			validator.LoginBodyReq = &req
		}
	}

	return validator
}

func (validator *Validator) CheckUsernameIsWellFormed() *Validator {
	validator.ParseLoginRequestBody()

	if validator.Err != nil {
		return validator
	}

	//address, err := auth.CheckEmailIsWellFormed(validator.Req.Email)

	err := auth.CheckUsername(validator.LoginBodyReq.Username)

	if err != nil {
		validator.Err = err
	}

	return validator
}

func (validator *Validator) CheckEmailIsWellFormed() *Validator {
	validator.ParseLoginRequestBody()

	if validator.Err != nil {
		return validator
	}

	//address, err := auth.CheckEmailIsWellFormed(validator.Req.Email)

	address, err := mail.ParseAddress(validator.LoginBodyReq.Email)

	if err != nil {
		validator.Err = err
	} else {
		validator.Address = address
	}

	return validator
}

func (validator *Validator) LoadAuthUserFromUuid() *Validator {

	if validator.Err != nil {
		return validator
	}

	authUser, err := userdbcache.FindUserByUuid(validator.LoginBodyReq.Uuid)

	if err != nil {
		validator.Err = fmt.Errorf(routes.ERROR_USER_DOES_NOT_EXIST)
	} else {
		validator.AuthUser = authUser
	}

	return validator

}

func (validator *Validator) LoadAuthUserFromEmail() *Validator {
	validator.CheckEmailIsWellFormed()

	if validator.Err != nil {
		return validator
	}

	authUser, err := userdbcache.FindUserByEmail(validator.Address)

	if err != nil {
		validator.Err = fmt.Errorf(routes.ERROR_USER_DOES_NOT_EXIST)
	} else {
		validator.AuthUser = authUser
	}

	return validator

}

func (validator *Validator) LoadAuthUserFromUsername() *Validator {
	validator.ParseLoginRequestBody()

	if validator.Err != nil {
		return validator
	}

	authUser, err := userdbcache.FindUserByUsername(validator.LoginBodyReq.Username)

	//log.Debug().Msgf("beep2 %s", authUser.Username)

	if err != nil {
		validator.Err = fmt.Errorf(routes.ERROR_USER_DOES_NOT_EXIST)
	} else {
		validator.AuthUser = authUser
	}

	return validator

}

func (validator *Validator) LoadAuthUserFromSession() *Validator {
	validator.ParseLoginRequestBody()

	if validator.Err != nil {
		return validator
	}

	sessionData, err := ReadSessionInfo(validator.c)

	if err != nil {
		validator.Err = fmt.Errorf("user not in session")
		validator.CheckIsValidRefreshToken().CheckUsernameIsWellFormed()
	}

	validator.AuthUser = sessionData.AuthUser

	return validator
}

func (validator *Validator) CheckAuthUserIsLoaded() *Validator {
	if validator.Err != nil {
		return validator
	}

	if validator.AuthUser == nil {
		validator.Err = fmt.Errorf("no auth user")
	}

	return validator
}

func (validator *Validator) CheckUserHasVerifiedEmailAddress() *Validator {
	validator.CheckAuthUserIsLoaded()

	if validator.Err != nil {
		return validator
	}

	if validator.AuthUser.EmailVerifiedAt == 0 {
		validator.Err = fmt.Errorf("email address not verified")
	}

	return validator
}

// If using jwt middleware, token is put into user variable
// and we can extract data from the jwt
func (validator *Validator) LoadTokenClaims() *Validator {
	if validator.Err != nil {
		return validator
	}

	if validator.Claims == nil {
		user, ok := validator.c.Get("user")

		if ok {
			validator.Claims = user.(*auth.TokenClaims)
		}
	}

	return validator
}

// Extracts public id from token, checks user exists and calls success function.
// If claims argument is nil, function will search for claims automatically.
// If claims are supplied, this step is skipped. This is so this function can
// be nested in other call backs that may have already extracted the claims
// without having to repeat this part.
func (validator *Validator) LoadAuthUserFromToken() *Validator {
	validator.LoadTokenClaims()

	if validator.Err != nil {
		return validator
	}

	authUser, err := userdbcache.FindUserByUuid(validator.Claims.Uuid)

	if err != nil {
		validator.Err = fmt.Errorf(routes.ERROR_USER_DOES_NOT_EXIST)
	} else {
		validator.AuthUser = authUser
	}

	return validator
}

func (validator *Validator) CheckIsValidRefreshToken() *Validator {
	validator.LoadTokenClaims()

	if validator.Err != nil {
		return validator
	}

	if validator.Claims.Type != auth.REFRESH_TOKEN {
		validator.Err = fmt.Errorf("no refresh token")
	}

	return validator

}

func (validator *Validator) CheckIsValidAccessToken() *Validator {
	validator.LoadTokenClaims()

	if validator.Err != nil {
		return validator
	}

	if validator.Claims.Type != auth.ACCESS_TOKEN {
		validator.Err = fmt.Errorf("no access token")
	}

	return validator
}
