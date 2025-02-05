package authEntity

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shared/entity"
	"shared/internal/global"
	"shared/utils"
)

type SupabaseError struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"error_code"`
	Msg       string `json:"msg"`
}

type DefaultError struct {
	Code  int
	Error string
	Msg   string
}

func (e *DefaultError) WriteError(c *gin.Context) {
	entity.ResponseJson{Status: e.Code, Error: e.Error, Message: e.Msg}.WriteError(c)
}

func (err SupabaseError) Error() string {
	//byte, _ := json.Marshal(err)
	return fmt.Sprint(err)
}

// ParseError parse auth error, defaultCode will use when no status code included in error, use -1 for pre-set inside function
func ParseError(err error, defaultCode int) *DefaultError {
	var supabaseError SupabaseError
	var grpcError *status.Status

	if defaultCode == -1 {
		defaultCode = http.StatusInternalServerError
	}

	//global.Logger.Error(err.Error())
	jsonErr := utils.Deserialize(err.Error(), &supabaseError)
	if jsonErr != nil {
		global.Logger.Error(jsonErr.Error())
	}

	switch {

	case jsonErr == nil:
		global.Logger.Error("supabase error")
		return &DefaultError{supabaseError.Code, supabaseError.ErrorCode, supabaseError.Msg}

	case status.Code(err) == codes.Aborted:
		grpcError = status.Convert(err)
		global.Logger.Error(grpcError.Message())
		ParseError(errors.New(grpcError.Message()), defaultCode)
	}

	return &DefaultError{Code: defaultCode, Msg: "Unknown Error"}
}
