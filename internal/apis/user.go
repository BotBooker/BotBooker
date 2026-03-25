package apis

import (
	"net/http"

	"github.com/botbooker/botbooker/internal/controllers"
	"github.com/botbooker/botbooker/internal/database"

	"github.com/gin-gonic/gin"
)

type ApiUsers struct {
	group          *gin.RouterGroup
	userController *controllers.UserController
}

type NewUser struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

func AddApiUsers(group *gin.RouterGroup, userController *controllers.UserController) *ApiUsers {
	api := &ApiUsers{
		group:          group,
		userController: userController,
	}
	api.addRoutes()
	return api
}

func (api *ApiUsers) addRoutes() {
	users := api.group.Group("/users")

	users.POST("/", api.createUser)
}

// ---------- Handlers ----------

// Create new user
func (api *ApiUsers) createUser(c *gin.Context) {
	var req NewUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := database.User{
		Email:        req.Email,
		FullName:     req.FullName,
		PasswordHash: req.Password,
	}

	ok, err := api.userController.CreateUser(&user)
	if !ok || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}
