package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cucumber/godog"

	"github.com/Flynn81/tkdo/model"
)

var resp *http.Response
var userID string

func makeGetRequest(endpoint string, includeUserID bool) error {
	host := os.Getenv(envHost)

	url := fmt.Sprintf("http://%s:7056/%s", host, endpoint)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if includeUserID {
		req.Header.Set("uid", userID)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err = client.Do(req)
	return err
}

func makePostRequest(endpoint string, includeUserID bool, v interface{}) error {
	var err error
	body, err := json.Marshal(v)

	if err != nil {
		return err
	}

	host := os.Getenv(envHost)
	url := fmt.Sprintf("http://%s:7056/%s", host, endpoint)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if includeUserID {
		req.Header.Set("uid", userID)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err = client.Do(req)
	return err
}

func aResponseCodeIsReturnedAndThereIsNoBody(arg1 int) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response code was %d", resp.StatusCode)
	}
	if resp.ContentLength > 0 {
		return fmt.Errorf("response body not empt")
	}
	return nil
}

func iGetAResponse() error {
	if resp == nil {
		return fmt.Errorf("response was nil")
	}
	return nil
}

func iMakeARequestToTheHealthCheck() error {
	return makeGetRequest("hc", false)
}

func getBackTheUserISentWithAnId() error {
	var user model.User
	err := json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return fmt.Errorf("error decoding response, %e, response content length %d, response code %d ", err, resp.ContentLength, resp.StatusCode)
	}
	if user.ID == "" {
		return fmt.Errorf("user id returned was empty")
	}
	return nil
}

func iCreateAUser() error {
	email := strconv.Itoa(time.Now().Nanosecond())
	user := model.User{Name: "new user", Email: email}
	return makePostRequest("users", false, user)
}

func iGetBackAnEmptyArray() error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response code was %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var tasks []model.Task
	err := json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return fmt.Errorf("error decoding response, %e", err)
	}
	if tasks == nil || len(tasks) > 0 {
		return fmt.Errorf("tasks is nil or len > 0")
	}
	return nil
}

func theUserHasCreatedATask() error {
	return aTaskIsCreated()
}

func tasksAreRequestedByAUser() error {
	var err error
	err = iCreateAUser()
	if err != nil {
		return err
	}

	var user model.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return fmt.Errorf("error decoding response, %e, response content length %d, response code %d ", err, resp.ContentLength, resp.StatusCode)
	}

	userID = user.ID
	err = makeGetRequest("tasks", true)

	if err != nil {
		return err
	}

	return nil
}

func theApiIsUp() error {
	return iMakeARequestToTheHealthCheck()
}

func aTaskIsCreated() error {
	task := model.Task{Name: "new task", TaskType: "basic"}
	return makePostRequest("tasks", true, task)
}

func iGetBackCreatedStatus() error {
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("response code was %d", resp.StatusCode)
	}
	return nil
}

func theOneTaskIsReturned() error {
	var tasks []model.Task
	err := json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return fmt.Errorf("error decoding response, %e", err)
	}
	if tasks == nil {
		return fmt.Errorf("tasks is nil")
	}
	if len(tasks) != 1 {
		return fmt.Errorf("tasks len = %d", len(tasks))
	}
	return nil
}

func theUserHasCreatedOneTask() error {
	return godog.ErrPending
}

func tasksAreRequested() error {
	return makeGetRequest("tasks", false)
}

func iGetBackAnError() error {
	if resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("response code was %d", resp.StatusCode)
	}
	return nil
}

func aUserIsCreated() error {
	err := iCreateAUser()
	if err != nil {
		return err
	}
	return nil
}

func tasksAreRequestedByTheUser() error {
	return makeGetRequest("tasks", true)
}

func theUserCreatesATask() error {
	var user model.User
	err := json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return fmt.Errorf("error decoding response, %e, response content length %d, response code %d ", err, resp.ContentLength, resp.StatusCode)
	}

	userID = user.ID
	return aTaskIsCreated()
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a (\d+) response code is returned And there is no body$`, aResponseCodeIsReturnedAndThereIsNoBody)
	ctx.Step(`^a task is created$`, aTaskIsCreated)
	ctx.Step(`^a user is created$`, aUserIsCreated)
	ctx.Step(`^get back the user I sent with an id$`, getBackTheUserISentWithAnId)
	ctx.Step(`^I create a user$`, iCreateAUser)
	ctx.Step(`^I get a response$`, iGetAResponse)
	ctx.Step(`^I get back an empty array$`, iGetBackAnEmptyArray)
	ctx.Step(`^I get back an error$`, iGetBackAnError)
	ctx.Step(`^I get back created status$`, iGetBackCreatedStatus)
	ctx.Step(`^I make a request to the health check$`, iMakeARequestToTheHealthCheck)
	ctx.Step(`^tasks are requested$`, tasksAreRequested)
	ctx.Step(`^tasks are requested by a user$`, tasksAreRequestedByAUser)
	ctx.Step(`^tasks are requested by the user$`, tasksAreRequestedByTheUser)
	ctx.Step(`^the api is up$`, theApiIsUp)
	ctx.Step(`^the one task is returned$`, theOneTaskIsReturned)
	ctx.Step(`^the user creates a task$`, theUserCreatesATask)
}
