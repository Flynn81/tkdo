package model

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type fakeTaskDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (d *fakeTaskDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, nil
}

func (d *fakeTaskDynamoDB) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {

	name := "alice"
	taskType := "basic"
	status := "new"
	userID := "test-user-id"
	id := "an-id"

	return &dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"userId": {
				S: aws.String(userID),
			},
			"taskType": {
				S: aws.String(taskType),
			},
			"status": {
				S: aws.String(status),
			},
			"name": {
				S: aws.String(name),
			},
		}}, nil
}

func (d *fakeTaskDynamoDB) UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return &dynamodb.UpdateItemOutput{}, nil
}

func (d *fakeTaskDynamoDB) DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return &dynamodb.DeleteItemOutput{}, nil
}

func TestCreateTask(t *testing.T) {

	db = &fakeTaskDynamoDB{}

	name := "alice"
	taskType := "basic"
	status := "new"
	userID := "test-user-id"

	ta := CockroachTaskAccess{}
	task := Task{Name: name, TaskType: taskType, Status: status, UserID: userID}
	returnedTask := ta.Create(&task)

	if returnedTask.Name != name {
		t.Errorf("name not set correctly")
	} else if returnedTask.TaskType != taskType {
		t.Errorf("task type not set correctly")
	} else if returnedTask.Status != status {
		t.Errorf("status not set correctly")
	} else if returnedTask.UserID != userID {
		t.Errorf("user id not set correctly")
	} else if returnedTask.ID == "" {
		t.Errorf("ID not set correctly")
	}
}

func TestGetTask(t *testing.T) {
	db = &fakeTaskDynamoDB{}

	name := "alice"
	taskType := "basic"
	status := "new"
	userID := "test-user-id"
	id := "an-id"

	ta := CockroachTaskAccess{}
	returnedTask, err := ta.Get(id, userID)

	if err != nil {
		t.Errorf("error getting task, %e", err)
	} else if returnedTask.Name != name {
		t.Errorf("name not set correctly")
	} else if returnedTask.TaskType != taskType {
		t.Errorf("task type not set correctly")
	} else if returnedTask.Status != status {
		t.Errorf("status not set correctly")
	} else if returnedTask.UserID != userID {
		t.Errorf("user id not set correctly")
	} else if returnedTask.ID == "" {
		t.Errorf("ID not set correctly")
	}
}

func TestUpdateTask(t *testing.T) {
	db = &fakeTaskDynamoDB{}

	name := "alice"
	taskType := "basic"
	status := "new"
	userID := "test-user-id"
	id := "an-id"

	task := Task{ID: id, Name: name, TaskType: taskType, Status: status, UserID: userID}
	ta := CockroachTaskAccess{}
	if !ta.Update(&task) {
		t.Errorf("update returned false")
	}
}

func TestDeleteTask(t *testing.T) {
	db = &fakeTaskDynamoDB{}

	userID := "test-user-id"
	id := "an-id"

	ta := CockroachTaskAccess{}

	if !ta.Delete(id, userID) {
		t.Errorf("delete returned false")
	}
}
