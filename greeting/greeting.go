package greeting

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func Greet(ctx workflow.Context, name string) (string, error) {
	return "hello " + name, nil
}

// Activity
// context is first arg, it's optional but recommended to pass imp info from temporal workflow
func GreetInSpanish(ctx context.Context, name string) (string, error) {
	base := "http://localhost:9999/get-spanish-greeting?name=%s"
	url := fmt.Sprintf(base, url.QueryEscape(name))

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	translation := string(body)
	status := resp.StatusCode
	if status >= 400 {
		message := fmt.Sprintf("HTTP Error %d: %s", status, translation)
		return "", errors.New(message)
	}
	return translation, nil
}

// workflow which calls activity
func GreetSomeone(ctx workflow.Context, name string) (string, error) {
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    15 * time.Second, //first retry will occur after 15 sec
		BackoffCoefficient: 2.0,              // double the delay after each retry -> 2, 4, 8, 64
		MaximumInterval:    time.Second * 60, // delay upto 60 sec in above logic
		MaximumAttempts:    100,              //fail activity after 100 attempts
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy:         retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	var spanishGreeting string

	//ctx -> workflow.WithActivityOptions, nameOfFunc, Input param
	//get -> on future object
	err := workflow.ExecuteActivity(ctx, GreetInSpanish, name).Get(ctx, &spanishGreeting)
	if err != nil {
		return "", err
	}
	return spanishGreeting, nil
}
