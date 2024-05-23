package helloworld

import (
	"go.temporal.io/sdk/workflow"
	"time"
)

/**
@file:
@author: levi.Tang
@time: 2024/5/21 12:10
@description:
**/

func GreetingWorkflow(ctx workflow.Context, name string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}

	ctx = workflow.WithActivityOptions(ctx, options)
	logger := workflow.GetLogger(ctx)
	var result string
	err := workflow.ExecuteActivity(ctx, GetNameAsync, name).Get(ctx, &result)
	if err != nil {
		logger.Error("get Nameasync error ", err)
		return "", err
	}
	err = workflow.ExecuteActivity(ctx, SayHello, result).Get(ctx, &result)
	if err != nil {
		logger.Error("say hello err", err)
		return "", err
	}
	return result, err
}
