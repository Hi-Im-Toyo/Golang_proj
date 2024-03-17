package token

import (
	"context"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

var UserData *mongo.Collection = datavase.OpenCollection(database.Client, "Users")

var SECRET_KEY = os.Getenv("SECRET")

func TokenGenerator(email string, first_name string, last_name string, uid string) (signedtoken string, signedrefreshtoken string, err error){
	
	claims := &SignedDetails{
		Email:      email,
		First_Name: first_name,
		Last_Name:  last_name,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
		}
	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims {
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		}
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "",	"", err
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return 
	}
	return token, refreshtoken, err
}


func ValidateToken(signedtoken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token has expired"
		return
	}
	return claims, msg

}


func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateobj primitive.D

	updateobj = append(updateobj, bson.E{Key: "token", Value: signedtoken})
	updateobj = append(updateobj, bson.E{Key: "refreshtoken", Value: signedrefreshtoken})
	updated_at := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateobj = append(updateobj, bson.E{Key: "updatedat", Value: updated_at})

	usert := true

	filter := bson.M{"user_id": userid}	
	opt := options.UpdateOptions{
		Upsert: &usert,
	}
	_, err := UserData.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateobj}}, &opt)
	defer cancel()	

	if err != nil {	
		log.Panic(err)
		return
	}

}