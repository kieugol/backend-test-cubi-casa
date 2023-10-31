package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/backend-test-cubi-casa/api/v1/controllers"
	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/config"
	"github.com/backend-test-cubi-casa/helpers/errc"
	"github.com/backend-test-cubi-casa/helpers/resp"
	"github.com/backend-test-cubi-casa/helpers/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type UserServiceMock struct {
	mock.Mock
}

// Init system
func init() {
	gin.SetMode(gin.TestMode)
	config.Init(os.Getenv("APP_ENV"), os.Getenv("HOME_DIR")+"/config")
}

// Mock function
func (m *UserServiceMock) HandleSearch(ctx context.Context, req models.UserSearchReq) ([]*models.UserResp, error) {
	agrs := m.Called(ctx, req)
	var data []*models.UserResp
	byteD, _ := json.Marshal(agrs.Get(0))
	util.ParseJSON(byteD, &data)

	return data, agrs.Error(1)
}

// Mock function
func (m *UserServiceMock) HandleCreate(ctx context.Context, req models.UserCreateReq) (*models.User, error) {
	agrs := m.Called(ctx, req)
	var data *models.User
	byteD, _ := json.Marshal(agrs.Get(0))
	util.ParseJSON(byteD, &data)

	return data, agrs.Error(1)
}

// Prepare params
func (m *UserServiceMock) initParams(name string, email string) models.UserSearchReq {
	req := models.UserSearchReq{
		Name:  name,
		Email: email,
	}
	return req
}

func TestSearchUser_Found_WithNameAndEmail(t *testing.T) {
	// Mock data
	userSrvMock := new(UserServiceMock)
	reqMock := userSrvMock.initParams("krol", "nhu")
	pathF1, _ := filepath.Abs("../mock_data/user/data.user.json")
	userDataMock := util.ReadFile(pathF1)

	var users []*models.UserResp
	var user *models.UserResp
	util.ParseJSON(userDataMock, &user)
	users = append(users, user)

	respExpected, _ := json.Marshal(resp.Success(users, http.StatusOK))

	// Setup Route
	r := gin.New()
	r.GET("/v1/users", controllers.NewUserController(userSrvMock).Search)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users?name=%s&email=%s", reqMock.Name, reqMock.Email), nil)
	c.Request = req
	// mock data from services
	userSrvMock.On("HandleSearch", c, reqMock).Return(users, nil)
	r.HandleContext(c)

	// Assert test result
	assert.Equal(t, string(respExpected), w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchUser_Null_WithEmptyParams(t *testing.T) {
	// Mock data
	userSrvMock := new(UserServiceMock)
	reqMock := userSrvMock.initParams("", "")
	respExpected, _ := json.Marshal(resp.Success(nil, http.StatusOK))

	// Setup Route and execute test
	r := gin.New()
	r.GET("/v1/users", controllers.NewUserController(userSrvMock).Search)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users?name=%s&email=%s", reqMock.Name, reqMock.Email), nil)
	c.Request = req
	userSrvMock.On("HandleSearch", c, reqMock).Return(nil, nil)
	r.HandleContext(c)

	// Assert test result
	respActual, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, string(respExpected), string(respActual))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchUser_Failed_InternalServerError(t *testing.T) {
	// Mock data
	userSrvMock := new(UserServiceMock)
	reqMock := userSrvMock.initParams("krol", "krol@gmail.com")

	respExpected, _ := json.Marshal(resp.InternalServerError())

	// Setup Route and execute test
	r := gin.New()
	r.GET("/v1/users", controllers.NewUserController(userSrvMock).Search)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users?name=%s&email=%s", reqMock.Name, reqMock.Email), nil)
	c.Request = req
	// mock data from services
	userSrvMock.On("HandleSearch", c, reqMock).Return(nil, gorm.ErrInvalidDB)
	r.HandleContext(c)

	// Assert test result
	assert.Equal(t, string(respExpected), w.Body.String())
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateUser_Success(t *testing.T) {
	// Mock data
	userSrvMock := new(UserServiceMock)
	pathF1, _ := filepath.Abs("../mock_data/user/data.user.req.json")
	userDataMock := util.ReadFile(pathF1)

	var userCreateReq *models.UserCreateReq
	var user *models.User
	util.ParseJSON(userDataMock, &user)
	util.ParseJSON(userDataMock, &userCreateReq)
	body, _ := json.Marshal(userCreateReq)

	respExpected, _ := json.Marshal(resp.Success(user, http.StatusOK))

	// Setup Route and execute test
	r := gin.New()
	r.POST("/v1/users", controllers.NewUserController(userSrvMock).Create)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	c.Request = req
	// mock data from services
	userSrvMock.On("HandleCreate", c, *userCreateReq).Return(user, nil)
	r.HandleContext(c)

	// Assert test result
	assert.Equal(t, string(respExpected), w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateUser_Failed_Validation_Params(t *testing.T) {
	// Mock data
	userSrvMock := new(UserServiceMock)
	userCreateReq := models.UserCreateReq{
		Email:   "krol",
		Phone:   "13333",
		Address: "HCM city",
		TeamID:  1,
	}
	body, _ := json.Marshal(userCreateReq)
	var errV []*errc.ValidationError
	strErr := `[{"field":"name","reason":"required"},{"field":"email","reason":"email"},{"field":"phone","reason":"e164"}]`
	util.ParseJSON([]byte(strErr), &errV)
	respExpected, _ := json.Marshal(resp.BadRequest(errV))

	// Setup Route and execute test
	r := gin.New()
	r.POST("/v1/users", controllers.NewUserController(userSrvMock).Create)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	c.Request = req
	// mock data from services
	userSrvMock.On("HandleCreate", c, userCreateReq).Return(nil, nil)
	r.HandleContext(c)

	// Assert test result
	assert.Equal(t, string(respExpected), w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_Failed_Validation_Duplicate(t *testing.T) {
	// Mock data
	userSrvMock := new(UserServiceMock)
	userCreateReq := models.UserCreateReq{
		Name:    "krol",
		Address: "HCM city",
		Email:   "krol@gmail.com",
		Phone:   "+8491747477",
		TeamID:  1,
	}
	respExpected, _ := json.Marshal(resp.BadRequest(gorm.ErrDuplicatedKey.Error()))
	body, _ := json.Marshal(userCreateReq)

	// Setup Route and execute test
	r := gin.New()
	r.POST("/v1/users", controllers.NewUserController(userSrvMock).Create)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	c.Request = req
	// mock data from services
	userSrvMock.On("HandleCreate", c, userCreateReq).Return(nil, gorm.ErrDuplicatedKey)
	r.HandleContext(c)

	// Assert test result
	assert.Equal(t, string(respExpected), w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_Failed_InternalServerError(t *testing.T) {
	// Mock data
	userSrvMock := new(UserServiceMock)
	userCreateReq := models.UserCreateReq{
		Name:    "krol",
		Email:   "krol@gmail.com",
		Phone:   "+8491747477",
		Address: "HCM city",
		TeamID:  1,
	}
	respExpected, _ := json.Marshal(resp.InternalServerError())
	body, _ := json.Marshal(userCreateReq)

	// Execute test Controller
	r := gin.New()
	r.POST("/v1/users", controllers.NewUserController(userSrvMock).Create)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	c.Request = req
	// mock data from services
	userSrvMock.On("HandleCreate", c, userCreateReq).Return(nil, gorm.ErrInvalidDB)
	r.HandleContext(c)

	// Assert test result
	assert.Equal(t, string(respExpected), w.Body.String())
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
