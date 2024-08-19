package main

import (
	"api_gateway/pkg/logger"
	"fmt"
	"github.com/casbin/casbin/v2"
	_ "github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
)

var (
	user   = "postgres"
	port   = 5433
	dbname = "cabin"
)

func CasbinEnforcer(log logger.ILogger) (*casbin.Enforcer, error) {

	adapter, err := xormadapter.NewAdapter("postgres", fmt.Sprintf("postgres3=postgres port=%d user=postgres dbname=casbin password=1111 sslmode=disable", 5433))
	if err != nil {
		log.Error("this error is adapter to connect to", logger.Error(err))
		return nil, err
	}
	enforcer, err := casbin.NewEnforcer("casbin/model.conf", adapter)
	if err != nil {
		log.Error("this error is acasbin read from casbin/model.conf", logger.Error(err))
		return nil, err
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Error("this error is casbin Load policy", logger.Error(err))
		return nil, err
	}

	policies := [][]string{
		{"admin", "/budget_service/v1/account", "POST"},
		{"admin", "/budget_service/v1/account", "GET"},
		{"admin", "/budget_service/v1/accounts", "GET"},
		{"admin", "/budget_service/v1/account", "PUT"},
		{"admin", "/budget_service/v1/account", "DELETE"},

		{"admin", "/budget_service/v2/budget", "POST"},
		{"admin", "/budget_service/v2/budget", "GET"},
		{"admin", "/budget_service/v2/budgets", "GET"},
		{"admin", "/budget_service/v2/budget", "PUT"},
		{"admin", "/budget_service/v2/budget", "DELETE"},

		{"admin", "/budget_service/v3/category", "POST"},
		{"admin", "/budget_service/v3/category", "GET"},
		{"admin", "/budget_service/v3/categories", "GET"},
		{"admin", "/budget_service/v3/category", "PUT"},
		{"admin", "/budget_service/v3/category", "DELETE"},

		{"admin", "/budget_service/v4/goal", "POST"},
		{"admin", "/budget_service/v4/goal", "GET"},
		{"admin", "/budget_service/v4/goals", "GET"},
		{"admin", "/budget_service/v4/goal", "PUT"},
		{"admin", "/budget_service/v4/goal", "DELETE"},
		{"admin", "/budget_service/v4/goals/report-progress", "GET"},

		{"admin", "/budget_service/v5/transaction", "POST"},
		{"admin", "/budget_service/v5/transaction", "GET"},
		{"admin", "/budget_service/v5/transactions", "GET"},
		{"admin", "/budget_service/v5/transaction", "PUT"},
		{"admin", "/budget_service/v5/transaction", "DELETE"},
		{"admin", "/budget_service/v5/transactions/spend", "GET"},
		{"admin", "/budget_service/v5/transactions/income", "GET"},

		{"admin", "/user_service/v6/user", "GET"},
		{"admin", "/user_service/v6/user/update", "PUT"},
	}
	_, err = enforcer.AddPolicies(policies)
	if err != nil {
		log.Error("this error is add polices", logger.Error(err))
		return nil, err
	}
	err = enforcer.SavePolicy()
	if err != nil {
		log.Error("Error saving Casbin policy", logger.Error(err))
		return nil, err
	}
	return enforcer, nil

}
