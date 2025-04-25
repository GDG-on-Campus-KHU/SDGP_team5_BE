// db/model/counter.go

package model

type Counter struct {
	ID    string `bson:"_id"`
	Value int    `bson:"value"`
}