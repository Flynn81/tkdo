package model

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type fakeUserDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (d *fakeUserDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, nil
}

func (d *fakeUserDynamoDB) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	name := "alice"
	email := "any@email.com"
	status := "new"

	return &dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
			"status": {
				S: aws.String(status),
			},
			"name": {
				S: aws.String(name),
			},
		}}, nil
}

func TestCreateUser(t *testing.T) {
	db = &fakeUserDynamoDB{}

	name := "alice"
	email := "any@email.com"
	status := "new"

	ua := CockroachUserAccess{}
	user := User{Name: name, Email: email, Status: status}
	returnedUser := ua.Create(&user)
	if returnedUser.Email != email {
		t.Errorf("email not set correctly")
	} else if returnedUser.Name != name {
		t.Errorf("name not set correctly")
	} else if returnedUser.Status != status {
		t.Errorf("status not set correctly")
	}
}

func TestGetUser(t *testing.T) {

	db = &fakeUserDynamoDB{}

	name := "alice"
	email := "any@email.com"
	status := "new"

	ua := CockroachUserAccess{}
	returnedUser, err := ua.Get(email)
	if err != nil {
		t.Errorf("error when getting user, %e", err)
	} else if returnedUser.Email != email {
		t.Errorf("email not set correctly")
	} else if returnedUser.Name != name {
		t.Errorf("name not set correctly")
	} else if returnedUser.Status != status {
		t.Errorf("status not set correctly")
	}
}
