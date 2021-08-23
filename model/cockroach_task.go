package model

import (
	"github.com/google/uuid"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

//CockroachTaskAccess is a concrete strut implementing TaskAccess, backed by CockroachDB
type CockroachTaskAccess struct {
}

//Get returns an task given an id.
func (ta CockroachTaskAccess) Get(id string, userID string) (*Task, error) {

	output, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("task"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		zap.S().Infof("%e", err)
		return nil, err
	}

	if output.Item == nil {
		zap.S().Infof("could not find task for %v", id)
		return nil, nil
	}

	task := Task{}

	err = dynamodbattribute.UnmarshalMap(output.Item, &task)

	if err != nil {
		zap.S().Infof("%e", err)
		return nil, err
	}

	return &task, nil
}

//Create takes a task
func (ta CockroachTaskAccess) Create(t *Task) *Task {

	//t.ID = t.Name + t.UserID + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	t.ID = uuid.NewString()
	//TODO: set an id in t
	zap.S().Infof("creating a new task for user ID %v", t.UserID)
	av, err := dynamodbattribute.MarshalMap(&t)
	if err != nil {
		zap.S().Infof("%e", err)
		return nil
	}
	tableName := "task"
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = db.PutItem(input)
	if err != nil {
		zap.S().Infof("%e", err)
		return nil
	}
	return t

}

//Update takes a task and attempt to update it
func (ta CockroachTaskAccess) Update(task *Task) bool {

	tableName := "task"

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {
				S: aws.String(task.Name),
			},
			":s": {
				S: aws.String(task.Status),
			},
			":t": {
				S: aws.String(task.TaskType),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(task.ID),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#c": aws.String("name"),
			"#s": aws.String("status"),
		},
		ReturnValues:     aws.String(dynamodb.ReturnValueNone),
		UpdateExpression: aws.String("set #c = :n, #s = :s, taskType = :t"),
	}
	_, err := db.UpdateItem(input)
	if err != nil {
		zap.S().Infof("%e", err)
		return false
	}
	return true

}

//Delete takes an id and attempts to delete the task with the id
func (ta CockroachTaskAccess) Delete(id string, userID string) bool {

	tableName := "task"

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := db.DeleteItem(input)
	if err != nil {
		zap.S().Infof("%e", err)
		return false
	}
	return true
}

//GetMany returns tasks matching the given string and/or task type
func (ta CockroachTaskAccess) GetMany(keyword string, taskType string, userID string) []*Task {

	zap.S().Info("making a list request")

	if keyword == "" && taskType == "" {
		return nil
	}

	//TODO: smarter expression to do paging in the request to dynamodb
	tableName := "task"

	filt := expression.Contains(expression.Name("name"), keyword)

	proj := expression.NamesList(expression.Name("id"), expression.Name("userId"), expression.Name("name"), expression.Name("taskType"), expression.Name("status"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		zap.S().Infof("%e", err)
		return nil
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	result, err := db.Scan(params)
	if err != nil {
		zap.S().Infof("%e", err)
		return nil
	}

	var r = []*Task{}

	for _, i := range result.Items {
		t := Task{}
		err = dynamodbattribute.UnmarshalMap(i, &t)
		if err != nil {
			zap.S().Infof("%e", err)
			return nil
		}
		r = append(r, &t)
	}
	return r
}

func getPageOfTasks(page int, pageSize int, currentPage int, params *dynamodb.ScanInput, r []*Task) []*Task {
	result, err := db.Scan(params)
	if err != nil {
		zap.S().Infof("%e", err)
		return nil
	}

	//zap.S().Info(result.LastEvaluatedKey)

	//can we get to our page's data in the results returned?
	//if not, call again for more
	//else try and get the data
	//	if pagesize is not met, call again
	if len(result.Items) >= pageSize*page {
		//try and get the data
		for index, i := range result.Items {
			if index >= (currentPage-1)*pageSize {
				t := Task{}
				err = dynamodbattribute.UnmarshalMap(i, &t)
				if err != nil {
					zap.S().Infof("%e", err)
					return nil
				}
				r = append(r, &t)
				if len(r) == pageSize {
					return r
				}
			}
		}
		//if pagesize is not met, call again
	}
	return getPageOfTasks(page, pageSize, currentPage+1, params, r)
	//are we on a page that is equal to or greater than the requested page?
	//	yes? collect up to pagesize results
	//		if pagesize results are not met, call method again with currentPage + 1
	//	no?	call again with currentPage + 1
	//
	//
	// if (len(result.Items) <= pageSize*page || len(result.LastEvaluatedKey) == 0) && currentPage >= page {
	// 	var r = []*Task{}
	//
	// 	zap.S().Infof("return size %v, last evaluated key: %v", len(result.Items), result.LastEvaluatedKey)
	//
	// 	for _, i := range result.Items {
	// 		t := Task{}
	// 		err = dynamodbattribute.UnmarshalMap(i, &t)
	// 		if err != nil {
	// 			zap.S().Infof("%e", err)
	// 			return nil
	// 		}
	// 		r = append(r, &t)
	// 		if len(r) == pageSize {
	// 			return r
	// 		}
	// 	}
	// 	return r
	// }
	//
	// params.ExclusiveStartKey = result.LastEvaluatedKey
	//
	// return getPageOfTasks(page, pageSize, currentPage+1, params)
}

//List returns a list of tasks
func (ta CockroachTaskAccess) List(page int, pageSize int, userID string) []*Task {
	zap.S().Info("making a list request")
	zap.S().Infof("page, pageSize, userID: %v, %v, %v", page, pageSize, userID)

	tableName := "task"
	filt := expression.Name("userId").Equal(expression.Value(userID))

	proj := expression.NamesList(expression.Name("id"), expression.Name("userId"), expression.Name("name"), expression.Name("taskType"), expression.Name("status"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		zap.S().Infof("%e", err)
		return nil
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	var r = []*Task{}

	return getPageOfTasks(page, pageSize, 1, params, r)
}
