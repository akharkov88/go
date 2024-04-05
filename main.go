package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"sort"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Configuration struct {
	Address string
	Port    string
	User    string
	Db      string
}

func main() {

	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	clientOptions := options.Client().ApplyURI("mongodb://" + configuration.Address + ":" + configuration.Port)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Mongo.connect() ERROR: ", err)
		os.Exit(1)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	colletions, err := client.Database(configuration.Db).ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		fmt.Println("ERROR: colletions", err)
		os.Exit(1)
	}
	dict := make(map[string]int)
	for _, value := range colletions {
		count, _ := client.Database(configuration.Db).Collection(value).CountDocuments(ctx, bson.D{{}})
		dict[value] = int(count)
	}
	keys := make([]int, 0, len(dict))
	for _, val := range dict {
		keys = append(keys, val)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	for _, k := range keys[0:3] {
		for key, val := range dict {
			if val == k {
				fmt.Println(key + " : " + strconv.Itoa(val))
				break
			}
		}
	}
}
