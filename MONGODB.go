package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

func MongoCreate(Information *Response) {
//mongodb://<dbuser>:<dbpassword>@ds039484.mongolab.com:39484/<db_name>
	sess, err := mgo.Dial("mongodb://minglu:liu273@ds057862.mongolab.com:57862/mingluliumongodb")
	if err != nil {
		fmt.Printf("MongoDB connection error %v\n", err)
		panic(err)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	db := sess.DB("minglu").C("Location")

	Information.Id = bson.NewObjectId()

	err = db.Insert(&Information)
	if err != nil {
		fmt.Printf("MongoDB create error %v\n", err)
		os.Exit(1)
	}

	var results []Response
	err = db.Find(bson.M{"_id": Information.Id}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}

	err = db.Find(bson.M{}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}
}

func MongoDelete(Id bson.ObjectId) (error) {

	sess, err := mgo.Dial("mongodb://minglu:liu273@ds057862.mongolab.com:57862/mingluliumongodb")
	if err != nil {
		fmt.Printf("MongoDB connection error  %v\n", err)
		panic(err)
		return err;
	}
	defer sess.Close()
	
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("minglu").C("Location")

	err = collection.Delete(bson.M{"_id": Id})
	if err != nil {
		panic(err)
		return err
	}

	return nil
}

func MongoQuery(Id bson.ObjectId) (Response, error) {

    var result []Response
	sess, err := mgo.Dial("mongodb://minglu:liu273@ds057862.mongolab.com:57862/mingluliumongodb")
	if err != nil {
		fmt.Printf("MongoDB connection error %v\n", err)
		panic(err)
		return result[0], err
	}
	defer sess.Close()
	
	sess.SetSafe(&mgo.Safe{})
	collection := sess.("minglu").C("Location")

	err = collection.Find(bson.M{"_id": Id}).All(&result)

	if err != nil {
		panic(err)
		return result[0], err
	}

	return result[0],nil
}

func MongoUpdate(Information Response) (Response, error) {

	var LocationInfo Response

	sess, err := mgo.Dial("mongodb://minglu:liu273@ds057862.mongolab.com:57862/mingluliumongodb")
	if err != nil {
		fmt.Printf("MongoDB connection error %v\n", err)
		panic(err)
		return LocationInfo,err;
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.("minglu").C("Location")

	colQuerier := bson.M{"_id": Information.Id}

	change := bson.M{"$set": bson.M{"address": Information.Address,
									"city": Information.City,
									"state": Information.State,
									"zip": Information.Zip,
									"coordinate": bson.M{"lat":Information.Coordinate.Lat,
											"lng":Information.Coordinate.Lng}}}
	err = collection.Update(colQuerier, change)
	if err != nil {
		panic(err)
		return LocationInfo,err
	}
	Response,error := MongoQuery(Information.Id)
	return Response,error
}


