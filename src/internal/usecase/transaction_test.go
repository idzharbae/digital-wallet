package usecase_test

import (
	"context"
	"testing"

	"github.com/idzharbae/digital-wallet/src/internal/constants"
	"github.com/idzharbae/digital-wallet/src/internal/entity"
	"github.com/idzharbae/digital-wallet/src/internal/gateway/gatewaymock"
	"github.com/idzharbae/digital-wallet/src/internal/repository"
	"github.com/idzharbae/digital-wallet/src/internal/repository/repomock"
	"github.com/idzharbae/digital-wallet/src/internal/usecase"
	"github.com/palantir/stacktrace"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestTransaction(t *testing.T) {
	Convey("Test Transfer Balance", t, func() {
		var mockTransactionHandler *repomock.MockTransactionHandler
		var mockUserTransaction *repomock.MockUserTransactionRepository
		var mockUserBalance *repomock.MockUserBalanceRepository
		var mockRabbitMq *gatewaymock.MockRabbitMqGateway

		setup := func(t *testing.T) (*gomock.Controller, usecase.TransactionUC) {
			ctrl := gomock.NewController(t)
			mockTransactionHandler = repomock.NewMockTransactionHandler(ctrl)
			mockUserTransaction = repomock.NewMockUserTransactionRepository(ctrl)
			mockUserBalance = repomock.NewMockUserBalanceRepository(ctrl)
			mockRabbitMq = gatewaymock.NewMockRabbitMqGateway(ctrl)
			return ctrl, usecase.NewTransaction(mockTransactionHandler, mockUserTransaction, mockUserBalance, mockRabbitMq)
		}
		Convey("Must return error if transaction amount is not positive integer", func() {
			ctrl, uc := setup(t)
			defer ctrl.Finish()
			ctx := context.Background()

			err := uc.TransferBalance(ctx, "testsender", "testrecipient", 0)
			So(stacktrace.RootCause(err), ShouldResemble, usecase.ErrInvalidTransferAmount)
			err = uc.TransferBalance(ctx, "testsender", "testrecipient", -1000)
			So(stacktrace.RootCause(err), ShouldResemble, usecase.ErrInvalidTransferAmount)
		})

		Convey("Must execute transaction in order of username", func() {
			Convey("Sender first, then recipient", func() {
				ctrl, uc := setup(t)
				defer ctrl.Finish()
				ctx := context.Background()

				txMock := repomock.NewMockTx(ctrl)
				mockTransactionHandler.EXPECT().ExecuteTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f repository.TransactionFunction) error {
						return f(txMock)
					})
				mockUserBalance.EXPECT().WithTransaction(txMock).Return(mockUserBalance)
				mockUserTransaction.EXPECT().WithTransaction(txMock).Return(mockUserTransaction)

				gomock.InOrder(
					mockUserBalance.EXPECT().GetUserBalanceForUpdate(gomock.Any(), "ayu").Return(12345, nil).Times(1),
					mockUserBalance.EXPECT().UpdateBalance(gomock.Any(), "ayu", 0).Return(nil).Times(1),
					mockUserBalance.EXPECT().GetUserBalanceForUpdate(gomock.Any(), "dewi").Return(0, nil).Times(1),
					mockUserBalance.EXPECT().UpdateBalance(gomock.Any(), "dewi", 12345).Return(nil).Times(1),
					mockUserTransaction.EXPECT().InsertTransaction(gomock.Any(), "ayu", "dewi", entity.TransactionTypeDebit, 12345).Return(nil).Times(1),
					mockUserTransaction.EXPECT().InsertTransaction(gomock.Any(), "dewi", "ayu", entity.TransactionTypeCredit, 12345).Return(nil).Times(1),
					mockRabbitMq.EXPECT().PublishMessage(constants.TOTAL_DEBIT_QUEUE, "application/json", gomock.Any()).Return(nil).Times(1),
				)

				err := uc.TransferBalance(ctx, "ayu", "dewi", 12345)
				So(err, ShouldBeNil)
			})
			Convey("Recipient first, then sender", func() {
				ctrl, uc := setup(t)
				defer ctrl.Finish()
				ctx := context.Background()

				txMock := repomock.NewMockTx(ctrl)
				mockTransactionHandler.EXPECT().ExecuteTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f repository.TransactionFunction) error {
						return f(txMock)
					})
				mockUserBalance.EXPECT().WithTransaction(txMock).Return(mockUserBalance)
				mockUserTransaction.EXPECT().WithTransaction(txMock).Return(mockUserTransaction)

				gomock.InOrder(
					mockUserBalance.EXPECT().GetUserBalanceForUpdate(gomock.Any(), "ayu").Return(0, nil).Times(1),
					mockUserBalance.EXPECT().UpdateBalance(gomock.Any(), "ayu", 12345).Return(nil).Times(1),
					mockUserBalance.EXPECT().GetUserBalanceForUpdate(gomock.Any(), "dewi").Return(12345, nil).Times(1),
					mockUserBalance.EXPECT().UpdateBalance(gomock.Any(), "dewi", 0).Return(nil).Times(1),
					mockUserTransaction.EXPECT().InsertTransaction(gomock.Any(), "dewi", "ayu", entity.TransactionTypeDebit, 12345).Return(nil).Times(1),
					mockUserTransaction.EXPECT().InsertTransaction(gomock.Any(), "ayu", "dewi", entity.TransactionTypeCredit, 12345).Return(nil).Times(1),
					mockRabbitMq.EXPECT().PublishMessage(constants.TOTAL_DEBIT_QUEUE, "application/json", gomock.Any()).Return(nil).Times(1),
				)

				err := uc.TransferBalance(ctx, "dewi", "ayu", 12345)
				So(err, ShouldBeNil)
			})
		})
	})
}
