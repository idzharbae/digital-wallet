package http_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/idzharbae/digital-wallet/src/internal/delivery/http"
	"github.com/idzharbae/digital-wallet/src/internal/usecase/ucmock"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestHttpHandler(t *testing.T) {
	Convey("Test Register User", t, func() {
		ctrl := gomock.NewController(t)
		ucMock := ucmock.NewMockUserUC(ctrl)
		server := http.NewServer(nil, nil, nil, ucMock)

		type Response struct {
			Message string `json:"message"`
			Token   string `json:"token"`
		}

		Convey("Must return error if username field is empty", func() {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/create_user", strings.NewReader(`{}`))

			server.RegisterUser(c)

			var response Response
			json.NewDecoder(w.Body).Decode(&response)
			So(response.Message, ShouldEqual, "username is empty")
			So(w.Code, ShouldEqual, 400)
		})

		Convey("Must return error if usecase returned error", func() {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/create_user", strings.NewReader(`{"username": "test"}`))

			ucMock.EXPECT().RegisterUser(c.Request.Context(), "test").Return("", errors.New("Error!"))

			server.RegisterUser(c)

			var response Response
			json.NewDecoder(w.Body).Decode(&response)
			So(response.Message, ShouldEqual, "failed to register user")
			So(w.Code, ShouldEqual, 500)
		})

		Convey("Must return token if usecase raised no error", func() {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/create_user", strings.NewReader(`{"username": "test"}`))

			ucMock.EXPECT().RegisterUser(c.Request.Context(), "test").Return("asdfg", nil)

			server.RegisterUser(c)

			var response Response
			json.NewDecoder(w.Body).Decode(&response)
			So(response.Token, ShouldEqual, "asdfg")
			So(w.Code, ShouldEqual, 200)
		})
	})
}
