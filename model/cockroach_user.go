package model

import (
	"go.uber.org/zap"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	// "github.com/google/uuid"
)

//CockroachUserAccess is a concrete struct implementing UserAccess, backed by CockroachDB
type CockroachUserAccess struct {
}

//Get returns a user given an email address
func (ua CockroachUserAccess) Get(email string) (*User, error) {

	output, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("user"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})

	if err != nil {
		zap.S().Infof("%e", err)
		return nil, err
	}

	if output.Item == nil {
		zap.S().Infof("could not find user for %v", email)
		return nil, nil
	}

	user := User{}

	err = dynamodbattribute.UnmarshalMap(output.Item, &user)

	if err != nil {
		zap.S().Infof("%e", err)
		return nil, err
	}

	return &user, nil

}

//Create takes a user without an id and persists it
func (ua CockroachUserAccess) Create(u *User) *User {
	u.ID = uuid.NewString()
	av, err := dynamodbattribute.MarshalMap(&u)
	if err != nil {
		zap.S().Infof("%e", err)
		return nil
	}
	tableName := "user"
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = db.PutItem(input)
	if err != nil {
		zap.S().Infof("%e", err)
		return nil
	}
	return u

}
