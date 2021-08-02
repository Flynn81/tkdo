package model

import "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

//Error is a struct used for relaying error messages back to callers of the api
type Error struct {
	Msg string `json:"msg"`
}

var db dynamodbiface.DynamoDBAPI

//Init is use to set the db for the package
func Init(d dynamodbiface.DynamoDBAPI) {
	db = d
}

// func closeRows(r *sql.Rows) {
// 	err := r.Close()
// 	if err != nil {
// 		zap.S().Infof("%e", err)
// 	}
// }
