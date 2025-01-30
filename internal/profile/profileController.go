package profile

import (
	"auth/internal/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shared/database"
	"shared/entity"
	authEntity "shared/entity/auth"
	"shared/helper"
)

type controller struct {
	service  *profileService
	profiles *entity.UserRepo
}

func NewController() *controller {

	db := database.New(global.Config.Database.ConnectString)

	return &controller{
		profiles: entity.NewRepo(db.GetSchema()),
	}
}

func (s *controller) CreateHdl(c *gin.Context) {
	var form Post
	if err := c.Bind(&form); err != nil {
		authEntity.ParseError(err, 400).WriteError(c)
		return
	}

	global.Logger.Info(fmt.Sprintf("%+v", form))

	session, err := helper.GetSessionFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
	}

	createdUser, err := s.profiles.CreateUser(session.UserId, &form.DisplayName, &form.AvatarPath)
	if err != nil {
		authEntity.ParseError(err, 400).WriteError(c)
		return
	}
	c.JSON(http.StatusCreated, createdUser)
}

func (s *controller) GetHdl(c *gin.Context) {
	user, err := helper.GetUserFromContext(c)
	if user == nil || err != nil {
		c.AbortWithStatus(403)
		c.JSON(-1, gin.H{
			"message": "User not found",
		})
		return
	}
	c.JSON(http.StatusOK, user)
}
