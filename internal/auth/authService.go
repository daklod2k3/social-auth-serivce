package auth

import (
	"auth/internal/global"
	"errors"
	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	sharedEntity "shared/entity"
	authEntity "shared/entity/auth"
)

type service struct {
	goTrue gotrue.Client
	user   *sharedEntity.UserRepoRead
	*sharedEntity.Service
}

func NewService() *service {

	// init supabase supabase
	s := &service{}
	s.Service = sharedEntity.NewService(global.Config.Database.ConnectString)

	sp := global.Config.Supabase
	client, err := supabase.NewClient(sp.Url, sp.Key, &supabase.ClientOptions{})
	if err != nil {
		panic("supabase supabase init failed" + err.Error())
	}
	s.goTrue = client.Auth

	//init user repo
	s.user = sharedEntity.NewRepoRead(s.Service.Db.GetSchema())

	return s
}

func (s *service) Health() (*authEntity.HealthResponse, error) {
	_, err := s.goTrue.HealthCheck()
	if err != nil {
		return nil, s.Error(err)
	}

	return &authEntity.HealthResponse{
		Message: "ok",
	}, nil
}

func (s *service) GetSession(form *authEntity.SessionRequest) (*authEntity.AuthResponse, error) {

	cl := s.goTrue.WithToken(form.AccessToken)

	var (
		response = authEntity.AuthResponse{
			form.AccessToken,
			form.RefreshToken,
			nil,
			nil,
		}
		user *types.User = nil
	)

	if &form.AccessToken == nil {
		return nil, s.Error(errors.New("invalid access token"))
	}

	res, err := cl.GetUser()

	if err != nil {
		global.Logger.Error(err.Error())
		if form.RefreshToken == "" {
			return nil, s.Error(errors.New("invalid refresh token"))
		}
		tokenRes, err := cl.RefreshToken(form.RefreshToken)
		if err != nil {
			return nil, s.Error(err)
		}
		response.AccessToken = tokenRes.AccessToken
		user = &tokenRes.User
		response.RefreshToken = tokenRes.RefreshToken
		response.UserId = &tokenRes.User.ID
	} else {
		user = &res.User
		response.UserId = &res.User.ID
	}

	if user == nil {
		return nil, s.Error(errors.New("invalid token"))
	}

	response.User, err = s.user.FindUser(response.UserId)
	return &response, nil
}

func (s *service) Login(form *authEntity.LoginMail) (*authEntity.AuthResponse, error) {
	res, err := s.goTrue.SignInWithEmailPassword(form.Email, form.Password)
	if err != nil {
		return nil, s.Error(err)
	}
	user, err := s.user.FindUser(&res.User.ID)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, s.Error(err)
		}
	}
	return &authEntity.AuthResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		User:         user,
		UserId:       &res.User.ID,
	}, nil
}

func (s *service) Register(form *authEntity.RegisterMail) (*authEntity.AuthResponse, error) {
	auth, err := s.goTrue.Signup(types.SignupRequest{
		Email:    form.Email,
		Password: form.Password,
	})
	if err != nil {
		return nil, s.Error(err)
	}

	return &authEntity.AuthResponse{
		AccessToken:  auth.AccessToken,
		RefreshToken: auth.RefreshToken,
		UserId:       &auth.User.ID,
	}, nil

}

func (s *service) Error(err error) error {
	//global.Logger.Error(err.Error())
	return status.Error(codes.Aborted, err.Error())
}
