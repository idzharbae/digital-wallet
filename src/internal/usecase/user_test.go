package usecase_test

import (
	"testing"

	"github.com/idzharbae/digital-wallet/src/internal/repository/repomock"
	"github.com/idzharbae/digital-wallet/src/internal/usecase"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestUser(t *testing.T) {
	Convey("Test Register User", t, func() {
		ctrl := gomock.NewController(t)
		repoMock := repomock.NewMockUserTokenRepository(ctrl)
		uc := usecase.NewUser(repoMock)
		Convey("Must generate random token before inserting", func() {
			username := "test"
			repoMock.EXPECT().InsertUserToken(username, gomock.Any()).Return(nil).Times(1)
			token, err := uc.RegisterUser(username)
			So(err, ShouldBeNil)
			So(len(token), ShouldEqual, 36)
		})
	})
}
