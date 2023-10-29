package unit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// Mock function
func (m *UserServiceMock) HandleSearch(req models.UserSearchReq) ([]*models.UserResp, error) {
	agrs := m.Called(req)
	var data []*models.UserResp
	byteD, _ := json.Marshal(agrs.Get(0))
	log.Println("byteD:", string(byteD))
	util.ParseJSON(byteD, &data)

	return data, agrs.Error(1)
}

func (m *UserServiceMock) HandleCreate(req models.UserCreateReq) (*models.User, error) {
	agrs := m.Called(req)
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
	req := userSrvMock.initParams("krol", "nhu")
	pathF1, _ := filepath.Abs("../mock_data/user/data_user.json")
	userDataMock := util.ReadFile(pathF1)

	var users []*models.UserResp
	var user *models.UserResp
	util.ParseJSON(userDataMock, &user)
	if user != nil {
		users = append(users, user)
	}

	userSrvMock.On("HandleSearch", req).Return(users, nil)
	respExpected, _ := json.Marshal(resp.Success(users, http.StatusOK))

	// Execute test Controller
	r := SetUpRouter()
	r.GET("/v1/users", controllers.NewUserController(userSrvMock).Search)
	reqTest, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users?name=%s&email=%s", req.Name, req.Email), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqTest)

	// Assert test result
	respActual, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, string(respExpected), string(respActual))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchUser_Empty_Params(t *testing.T) {
	// Mock data
	userSrvMock := new(UserServiceMock)
	req := userSrvMock.initParams("", "")

	userSrvMock.On("HandleSearch", req).Return(nil, nil)
	respExpected, _ := json.Marshal(resp.Success(nil, http.StatusOK))

	// Execute test Controller
	r := SetUpRouter()
	r.GET("/v1/users", controllers.NewUserController(userSrvMock).Search)
	reqTest, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users?name=%s&email=%s", req.Name, req.Email), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqTest)

	// Assert test result
	respActual, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, string(respExpected), string(respActual))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchUser_Failed_InternalServerError(t *testing.T) {
	// Mock data
	userSrvMock := new(UserServiceMock)
	req := userSrvMock.initParams("nhu", "kieu")

	userSrvMock.On("HandleSearch", req).Return(nil, gorm.ErrInvalidDB)
	respExpected, _ := json.Marshal(resp.InternalServerError())

	// Execute test Controller
	r := SetUpRouter()
	r.GET("/v1/users", controllers.NewUserController(userSrvMock).Search)
	reqTest, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users?name=%s&email=%s", req.Name, req.Email), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqTest)

	// Assert test result
	respActual, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, string(respExpected), string(respActual))
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

	userSrvMock.On("HandleCreate", *userCreateReq).Return(user, nil)
	respExpected, _ := json.Marshal(resp.Success(user, http.StatusOK))

	// Execute test Controller
	r := SetUpRouter()
	r.POST("/v1/users", controllers.NewUserController(userSrvMock).Create)
	payload, _ := json.Marshal(userCreateReq)
	reqTest, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqTest)

	// Assert test result
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, string(respExpected), string(responseData))
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

	userSrvMock.On("HandleCreate", userCreateReq).Return(nil, nil)
	var errV []*errc.ValidationError
	strErr := `[{"field":"name","reason":"required"},{"field":"email","reason":"email"},{"field":"phone","reason":"e164"}]`
	util.ParseJSON([]byte(strErr), &errV)

	respExpected, _ := json.Marshal(resp.BadRequest(errV))

	// Execute test Controller
	r := SetUpRouter()
	r.POST("/v1/users", controllers.NewUserController(userSrvMock).Create)
	body, _ := json.Marshal(userCreateReq)
	reqTest, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqTest)

	// Assert test result
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, string(respExpected), string(responseData))
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

	userSrvMock.On("HandleCreate", userCreateReq).Return(nil, gorm.ErrDuplicatedKey)
	respExpected, _ := json.Marshal(resp.BadRequest(gorm.ErrDuplicatedKey.Error()))

	// Execute test Controller
	r := SetUpRouter()
	r.POST("/v1/users", controllers.NewUserController(userSrvMock).Create)
	body, _ := json.Marshal(userCreateReq)
	reqTest, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqTest)

	// Assert test result
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, string(respExpected), string(responseData))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_Failed_InternalServerError(t *testing.T) {
	// Mock data
	userSrvMock := new(UserServiceMock)
	pathF1, _ := filepath.Abs("../mock_data/user/data.user.req.json")
	userReqMock := util.ReadFile(pathF1)
	var userCreateReq *models.UserCreateReq
	util.ParseJSON(userReqMock, &userCreateReq)

	userSrvMock.On("HandleCreate", *userCreateReq).Return(nil, gorm.ErrInvalidDB)
	respExpected, _ := json.Marshal(resp.InternalServerError())

	// Execute test Controller
	r := SetUpRouter()
	r.POST("/v1/users", controllers.NewUserController(userSrvMock).Create)
	reqTest, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(userReqMock))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqTest)

	// Assert test result
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, string(respExpected), string(responseData))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
