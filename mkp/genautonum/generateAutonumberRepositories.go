package genautonum

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
)

type autonumberValue struct {
	Prefix		string	`json:"prefix"`
	SeqValue	int		`json:"seq_value"`
	LeadingZero int		`json:"leading_zero"`
}

type generateAutonumberRepository struct {
	RepoDB Repository
}

func NewGenerateAutonumberRepository(repoDB Repository) generateAutonumberRepository {
	return generateAutonumberRepository{
		RepoDB: repoDB,
	}
}

type GenerateAutonumberRepository interface {
	GenerateAutonumber(p string, v string) (string, error)
	AutonumberValue(prefix string, leadingZero... int) (string, error)
}

// GenerateAutonumber
func (ctx generateAutonumberRepository) GenerateAutonumber(p string, v string) (string, error) {
	var autonumber string

	err := ctx.RepoDB.DB.QueryRow("SELECT fs_gen_autonum($1, $2)", p, v).Scan(&autonumber)
	if err != nil {
		return "", err
	}

	return autonumber, nil
}

func (ctx generateAutonumberRepository) AutonumberValue(prefix string, leadingZero ...int) (string, error) {
	colName := "autonumber_value"
	zeroPadding := 0

	if len(leadingZero) > 0 {
		zeroPadding = leadingZero[0]
	}

	var autonumber autonumberValue
	exists := true
	collect := ctx.RepoDB.MongoDB.Collection(colName).FindOne(ctx.RepoDB.Context, bson.M{"prefix": prefix})
	err := collect.Decode(&autonumber)
	if err != nil {
		exists = false
	}

	if !exists {
		result, err := ctx.RepoDB.MongoDB.Collection(colName).InsertOne(ctx.RepoDB.Context, autonumberValue{
			Prefix: prefix,
			SeqValue: 0,
			LeadingZero: zeroPadding,
		})

		if err != nil {
			return "", err
		}
		fmt.Println("Insert OID:", result.InsertedID)
	}

	filter := bson.M{"prefix": prefix}
	update := bson.M{
		"$inc": bson.M{"seqvalue": 1},
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	out := ctx.RepoDB.MongoDB.Collection(colName).FindOneAndUpdate(ctx.RepoDB.Context, filter, update, &opt)
	if out.Err() != nil {
		return "", err
	}

	err = out.Decode(&autonumber)
	if err != nil {
		return "", err
	}

	autonumberNo := ""
	if zeroPadding != 0 {
		lpad := leftPad(strconv.Itoa(autonumber.SeqValue), "0", autonumber.LeadingZero)
		autonumberNo = fmt.Sprintf("%s%s", prefix, lpad)
	} else {
		autonumberNo = fmt.Sprintf("%s%s", prefix, strconv.Itoa(autonumber.SeqValue))
	}

	return autonumberNo, nil
}

func leftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}