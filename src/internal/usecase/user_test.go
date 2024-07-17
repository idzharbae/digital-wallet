package usecase_test

import (
	"context"
	"testing"

	"github.com/idzharbae/digital-wallet/src/internal/repository"
	"github.com/idzharbae/digital-wallet/src/internal/repository/repomock"
	"github.com/idzharbae/digital-wallet/src/internal/usecase"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestUser(t *testing.T) {
	Convey("Test Register User", t, func() {
		ctrl := gomock.NewController(t)
		userTokenRepoMock := repomock.NewMockUserTokenRepository(ctrl)
		userBalanceRepoMock := repomock.NewMockUserBalanceRepository(ctrl)
		transactionHandlerMock := repomock.NewMockTransactionHandler(ctrl)
		ctx := context.Background()
		uc := usecase.NewUser(userTokenRepoMock, userBalanceRepoMock, transactionHandlerMock)
		Convey("Must generate random token before inserting", func() {
			username := "test"
			txMock := repomock.NewMockTx(ctrl)
			transactionHandlerMock.EXPECT().ExecuteTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, f repository.TransactionFunction) error {
					return f(txMock)
				})

			userTokenRepoMock.EXPECT().WithTransaction(txMock).Return(userTokenRepoMock)
			userTokenRepoMock.EXPECT().InsertUserToken(ctx, username, gomock.Any()).Return(nil).Times(1)
			userBalanceRepoMock.EXPECT().WithTransaction(txMock).Return(userBalanceRepoMock)
			userBalanceRepoMock.EXPECT().CreateUserBalance(ctx, username).Return(nil).Times(1)
			token, err := uc.RegisterUser(ctx, username)
			So(err, ShouldBeNil)
			So(len(token), ShouldEqual, 36)
		})
	})
}
