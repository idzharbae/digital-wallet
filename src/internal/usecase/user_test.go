package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/idzharbae/digital-wallet/src/internal/repository"
	"github.com/idzharbae/digital-wallet/src/internal/repository/repomock"
	"github.com/idzharbae/digital-wallet/src/internal/usecase"
	"github.com/palantir/stacktrace"
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

		Convey("Must return error when an error occured", func() {
			username := "test"
			txMock := repomock.NewMockTx(ctrl)
			transactionHandlerMock.EXPECT().ExecuteTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, f repository.TransactionFunction) error {
					return f(txMock)
				})

			userTokenRepoMock.EXPECT().WithTransaction(txMock).Return(userTokenRepoMock)
			userTokenRepoMock.EXPECT().InsertUserToken(ctx, username, gomock.Any()).Return(nil).Times(1)
			userBalanceRepoMock.EXPECT().WithTransaction(txMock).Return(userBalanceRepoMock)
			userBalanceRepoMock.EXPECT().CreateUserBalance(ctx, username).Return(errors.New("error")).Times(1)

			_, err := uc.RegisterUser(ctx, username)
			So(stacktrace.RootCause(err), ShouldResemble, errors.New("error"))
		})
	})

	Convey("Test Top Up Balance", t, func() {
		ctrl := gomock.NewController(t)
		userBalanceRepoMock := repomock.NewMockUserBalanceRepository(ctrl)
		transactionHandlerMock := repomock.NewMockTransactionHandler(ctrl)
		ctx := context.Background()
		uc := usecase.NewUser(nil, userBalanceRepoMock, transactionHandlerMock)
		Convey("Must start transaction and update balance based on user's current balance", func() {
			username := "test"
			topUpAmount := 5432
			currentBalance := 1234

			txMock := repomock.NewMockTx(ctrl)
			transactionHandlerMock.EXPECT().ExecuteTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, f repository.TransactionFunction) error {
					return f(txMock)
				})

			userBalanceRepoMock.EXPECT().WithTransaction(txMock).Return(userBalanceRepoMock)
			userBalanceRepoMock.EXPECT().GetUserBalanceForUpdate(ctx, username).Return(currentBalance, nil).Times(1)
			userBalanceRepoMock.EXPECT().UpdateBalance(ctx, username, topUpAmount+currentBalance).Return(nil).Times(1)

			newBalance, err := uc.TopUpUserBalance(ctx, username, topUpAmount)
			So(err, ShouldBeNil)
			So(newBalance, ShouldEqual, topUpAmount+currentBalance)
		})

		Convey("Must return error when an error occured", func() {
			username := "test"
			topUpAmount := 5432
			currentBalance := 1234

			txMock := repomock.NewMockTx(ctrl)
			transactionHandlerMock.EXPECT().ExecuteTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, f repository.TransactionFunction) error {
					return f(txMock)
				})

			userBalanceRepoMock.EXPECT().WithTransaction(txMock).Return(userBalanceRepoMock)
			userBalanceRepoMock.EXPECT().GetUserBalanceForUpdate(ctx, username).Return(currentBalance, nil).Times(1)
			userBalanceRepoMock.EXPECT().UpdateBalance(ctx, username, topUpAmount+currentBalance).Return(errors.New("error")).Times(1)

			newBalance, err := uc.TopUpUserBalance(ctx, username, topUpAmount)
			So(stacktrace.RootCause(err), ShouldResemble, errors.New("error"))
			So(newBalance, ShouldEqual, 0)
		})

		Convey("Must return error when amount >= 10000000", func() {
			username := "test"
			topUpAmount := 10000000

			_, err := uc.TopUpUserBalance(ctx, username, topUpAmount)
			So(stacktrace.RootCause(err), ShouldResemble, usecase.ErrTopUpTooLarge)
		})
	})
}
