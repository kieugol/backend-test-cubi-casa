package unit

import (
	"context"
	"encoding/json"
	"flag"
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
	"github.com/backend-test-cubi-casa/helpers/resp"
	"github.com/backend-test-cubi-casa/helpers/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

// Init system
func init() {
	baseDir := flag.String("e", os.Getenv("HOME_DIR"), "")
	log.Println("env:", baseDir)
	gin.SetMode(gin.TestMode)
	config.Init("development", "/go/src/github.com/backend-test-cubi-casa/config")
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// Mock function
func (m *UserServiceMock) HandleSearch(ctx context.Context, req models.UserSearchReq) ([]*models.UserResp, error) {
	agrs := m.Called(ctx, req)

	var data []*models.UserResp
	byteD, _ := json.Marshal(agrs.Get(0))
	util.ParseJSON(byteD, &data)

	return data, agrs.Error(1)
}

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

func Test_Case_UserController_1_Success(t *testing.T) {
	// Mock data
	userSrvMock := new(UserServiceMock)
	req := userSrvMock.initParams("krol", "nhu")
	pathF1, _ := filepath.Abs("../mock_data/user/data_user.json")
	userDataMock := string(util.ReadFile(pathF1))

	var respData []*models.UserResp
	var user *models.UserResp
	util.ParseJSON([]byte(userDataMock), &user)
	if user != nil {
		respData = append(respData, user)
	}

	var c *gin.Context
	userSrvMock.On("HandleSearch", c, req).Return(respData, nil)

	dataExpected, _ := json.Marshal(resp.Success(respData, http.StatusOK))

	// Execute test Controller
	ctrlTest := controllers.NewUserController(userSrvMock)
	r := SetUpRouter()
	r.GET("/v1/users", ctrlTest.Search)
	reqTest, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users?name=%s&mail=%s", req.Name, req.Email), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqTest)

	// Assert test result
	responseData, _ := ioutil.ReadAll(w.Body)
	log.Println(string(dataExpected), string(responseData))
	//assert.Equal(t, string(dataExpected), string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

//func Test_Case_UserController_2_Failed_NotFoundData(t *testing.T) {
//	// Mock data
//	userSrvMock := new(UserServiceMock)
//	req := userSrvMock.initParams(2)
//	userSrvMock.On("HandleDetail", req).Return(nil, http.StatusNotFound)
//	dataExpected, _ := json.Marshal(resp.NotFound())
//
//	// Execute test Controller
//	userCtrlTest := controllers.NewUserController(userSrvMock)
//	r := SetUpRouter()
//	r.GET("/v1/users/:id", userCtrlTest.Detail)
//	reqTest, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users/%d", req.ID), nil)
//	w := httptest.NewRecorder()
//	r.ServeHTTP(w, reqTest)
//
//	// Assert test result
//	responseData, _ := ioutil.ReadAll(w.Body)
//	assert.Equal(t, string(dataExpected), string(responseData))
//	assert.Equal(t, http.StatusNotFound, w.Code)
//}
