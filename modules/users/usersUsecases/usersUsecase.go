package usersusecases

import (
	"fmt"

	"github.com/MrXMMM/E-commerce-Project/config"
	"github.com/MrXMMM/E-commerce-Project/modules/users"
	usersrepositories "github.com/MrXMMM/E-commerce-Project/modules/users/usersRepositories"
	"github.com/MrXMMM/E-commerce-Project/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
	GetUserProfile(userId string) (*users.User, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersrepositories.IUserRepository
}

func UsersUseCase(cfg config.IConfig, usersRepository usersrepositories.IUserRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {

	// hashing password
	if err := req.BcrypeHashing(); err != nil {
		return nil, err
	}

	// Insert user
	result, err := u.usersRepository.InsertUSer(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error) {

	// hashing password
	if err := req.BcrypeHashing(); err != nil {
		return nil, err
	}

	// Insert user
	result, err := u.usersRepository.InsertUSer(req, true)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	// Find user
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	// SignToken
	accessToken, err := auth.NewAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	RefreshToken, err := auth.NewAuth(auth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	// Set Passport
	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: RefreshToken.SignToken(),
		},
	}

	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	// Parse token
	claims, err := auth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// check oauth
	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	//Find profile
	profile, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := auth.NewAuth(
		auth.Access,
		u.cfg.Jwt(),
		newClaims,
	)

	if err != nil {
		return nil, err
	}

	refreshToken := auth.RepeatToken(
		u.cfg.Jwt(),
		newClaims,
		claims.ExpiresAt.Unix(),
	)

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}
	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *usersUsecase) DeleteOauth(oauthId string) error {
	if err := u.usersRepository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil
}

func (u *usersUsecase) GetUserProfile(userId string) (*users.User, error) {

	profile, err := u.usersRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
