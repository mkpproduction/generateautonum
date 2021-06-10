package tests

import (
	"context"
	"fmt"
	"github.com/mkpproduction/mkp-sdk-go/mkp/genautonum"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

var (
	ctx = context.Background()
	MONGOHost = "127.0.0.1"
	MONGOPort = "27017"
	MONGODB = "autonumber"
)

func ConnectMongo(ctx context.Context, DBCollection... string) *mongo.Database  {
	connection := fmt.Sprintf("mongodb://%s:%s", MONGOHost, MONGOPort)
	fmt.Println("Connection Mongo:", connection)
	clientOptions := options.Client()
	clientOptions.ApplyURI(connection)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil
	}

	col := MONGODB
	if len(DBCollection) > 0 {
		col = DBCollection[0]
	}

	return client.Database(col)
}

func TestGenerateAutonumber(t *testing.T) {

	mongoConn := ConnectMongo(ctx)
	repo := genautonum.NewRepository(nil , nil, mongoConn)

	genRepository := genautonum.NewGenerateAutonumberRepository(repo)
	outputDocNo, err := genRepository.AutonumberValue("PPL-000000000000001-20210602-", 0)
	if err != nil {
		fmt.Println("Error Output:", err.Error())
	}

	fmt.Println("Output Doc No:", outputDocNo)
}
