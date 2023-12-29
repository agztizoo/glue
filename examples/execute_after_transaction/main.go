package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/agztizoo/glue/db"
	"github.com/agztizoo/glue/env"
	"github.com/agztizoo/glue/transaction"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 演示用例，使用 sqlite 代替.
	dial := func(opts *db.Options) (gorm.Dialector, error) {
		dsn := filepath.Join(env.WorkDir(), opts.DBName)
		dl := sqlite.Open(dsn)
		return dl, nil
	}
	opts := &db.Options{
		UserName: "xxx",
		Password: "xxx",
		DBName:   "execute_after_transaction.db",
	}
	source, err := opts.ToSource(dial, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Provider 实现了事务管理接口.
	provider := db.NewProvider(source)
	emailNotifier := &emailNotifier{
		tm: provider,
	}
	domainService := &CreateOrderService{
		notifier: emailNotifier,
	}
	appService := &CreateOrderAppService{
		tm:  provider,
		srv: domainService,
	}

	appService.CreateOrder(context.Background(), &CreateOrderCommand{
		User: "user_xxx",
	})
}

type CreateOrderAppService struct {
	tm  transaction.Manager
	srv *CreateOrderService
}

type CreateOrderCommand struct {
	User string
}

func (s *CreateOrderAppService) CreateOrder(ctx context.Context, cmd *CreateOrderCommand) {
	err := s.tm.Transaction(ctx, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		panic(err)
	}
}

type CreateOrderService struct {
	notifier Notifier
}

type CreateOrderParam struct {
	User string
}

func (s *CreateOrderService) Execute(ctx context.Context, param *CreateOrderParam) {
	// 1. Do business
	// TODO: TBD

	// 2. xxxx
	// TODO: TBD

	// 3. notify
	s.notifier.Notify(ctx, param.User)
}

type Notifier interface {
	Notify(ctx context.Context, userID string)
}

type emailNotifier struct {
	tm transaction.Manager
}

func (e *emailNotifier) Notify(ctx context.Context, userID string) {
	// Email Notifier 实现在事务公共后发送邮件.
	//
	// 该方法不保障邮件一定成功，但是一定是在业务成功后发送邮件.
	registered := e.tm.OnCommitted(ctx, func(ctx context.Context) {
		// TODO: TBD
		fmt.Printf("email has been sent to user %s", userID)
	})
	if !registered {
		panic("need to open a transaction")
	}
}
