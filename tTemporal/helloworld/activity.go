package helloworld

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/activity"
)

/**
@file:
@author: levi.Tang
@time: 2024/5/21 14:42
@description:
**/

func GetNameAsync(ctx context.Context, name string) (string, error) {
	greeting := fmt.Sprintf("Hello %s, async!", name)
	ac := activity.GetInfo(ctx)
	fmt.Printf("ac: %s", ac.ActivityID)
	return greeting, activity.ErrResultPending
}

func SayHello(ctx context.Context, name string) (string, error) {
	greeting2 := fmt.Sprintf("second hello: %s", name)
	return greeting2, nil
}
