package payment

import (
	"context"
	"fmt"
	mq "github.com/apache/rocketmq-client-go/v2"
	wePayment "github.com/medivhzhan/weapp/payment"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"net/http"
)

// Service is cutomer service
type Service interface {
	// Pay successfully, notify other services via rocketMQ.
	PaySucceed(ctx context.Context, w http.ResponseWriter, req *http.Request) error
}

type service struct {
	repo    Repository
	rocket  mq.Producer
	message Message
	logger  logger.Provider
}

// NewService is to create new service
func NewService(repo Repository, rocketMQ mq.Producer, tracing tracing.Provider, logger logger.Provider) Service {
	message := NewMessage(rocketMQ, tracing)
	return &service{repo, rocketMQ, message, logger}
}

func (s *service) PaySucceed(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	err := wePayment.HandlePaidNotify(w, req, s.payCallback())
	return err
}

func (s *service) payCallback() func(ntf wePayment.PaidNotify) (bool, string) {
	return func(ntf wePayment.PaidNotify) (bool, string) {
		// Check the postgres db to find whether this order has been processed.
		//@todo, the flag should include the two flags opts
		flag, err := s.repo.GetProcessedFlag(ctx, ntf.OutTradeNo)
		if err != nil {
			msg := fmt.Sprintf("order id: %s with error: %s", ntf.OutTradeNo, err.Error())
			s.logger.Errorw("payment get process flag", msg)

			return false, err.Error()
		}

		// Return if this order has been processed already
		if flag {
			return true, ""
		}

		// Send notification to rocketmq.

		// @todo use transaction for below two messages. If one message fails to send,
		// rollback the other one.
		paymentMethod, err := s.repo.GetPaymentMethod(ctx, ntf.OutTradeNo)
		if err != nil {
			msg := fmt.Sprintf("order id: %s with error: %s", ntf.OutTradeNo, err.Error())
			s.logger.Errorw("payment get method", msg)

			return false, err.Error()
		}

		err = s.message.ProduceOrderPaid(ctx, ntf.OutTradeNo)
		if err != nil {
			msg := fmt.Sprintf("order id: %s with error: %s", ntf.OutTradeNo, err.Error())
			s.logger.Errorw("payment send order complete message", msg)

			return false, err.Error()
		}

		if paymentMethod.Method == "" {
			err = s.message.ProduceBalanceDecrease(ctx, ntf.OutTradeNo, paymentMethod.CustomerID, paymentMethod.Amount)
			if err != nil {
				msg := fmt.Sprintf("order id: %s with error: %s", ntf.OutTradeNo, err.Error())
				s.logger.Errorw("payment send balance complete message", msg)

				return false, err.Error()
			}
		}

		// Save the order id processed flag, so this callback will not be invoked at the same order again.
		err = s.repo.SetProcessedFlag(ctx, ntf.OutTradeNo)
		if err != nil {
			msg := fmt.Sprintf("order id: %s with error: %s", ntf.OutTradeNo, err.Error())
			s.logger.Errorw("payment set process flag", msg)

			return false, err.Error()
		}

		return true, ""
	}
}
