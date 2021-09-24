package mock

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/Flynn81/tkdo/model"

	"go.uber.org/zap"
)

//MockDynamoDB Mock dynamodb
type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

var tasks = []model.Task{}
var users = []model.User{}

//PutItem puts an item
func (d *MockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	if *input.TableName == "task" {
		task := model.Task{}

		err := dynamodbattribute.UnmarshalMap(input.Item, &task)

		if err != nil {
			zap.S().Infof("%e", err)
			return nil, err
		}
		tasks = append(tasks, task)
	} else if *input.TableName == "user" {
		user := model.User{}

		err := dynamodbattribute.UnmarshalMap(input.Item, &user)

		if err != nil {
			zap.S().Infof("%e", err)
			return nil, err
		}
		users = append(users, user)
	}
	return &dynamodb.PutItemOutput{}, nil
}

//GetItem gets an item
func (d *MockDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {

	if *input.TableName == "task" {
		key := input.Key["id"].S
		for _, t := range tasks {
			if t.ID == *key {
        dynamodbattribute.MarshalMap()
				return nil, nil
			}
		}
	} else if *input.TableName == "user" {
		key := input.Key["email"].S
	}
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

//UpdateItem update an item
func (d *MockDynamoDB) UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return &dynamodb.UpdateItemOutput{}, nil
}

//DeleteItem delete an item
func (d *MockDynamoDB) DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return &dynamodb.DeleteItemOutput{}, nil
}
