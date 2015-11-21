package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
	"encoding/json"
     "io/ioutil"
     "gopkg.in/mgo.v2/bson"
     "io"
     "net/url"
     "math/rand"
)
var RandomID int

type Arguments struct{
    Name string `json:"name"    bson:"name"`
    Address string `json:"address" bson:"address"`
    City string `json:"city"    bson:"city"`
    State string `json:"state"   bson:"state"`
    Zip string `json:"zip"     bson:"zip"`
}

type Response struct{
    Id bson.ObjectId `json:"id"      bson:"_id,omitempty"`
    Name string `json:"name"    bson:"name"`
    Address string `json:"address" bson:"address"`
    City string `json:"city"    bson:"city"`
    State string `json:"state"   bson:"state"`
    Zip string `json:"zip"     bson:"zip"`
    Coordinate struct {
        Lat float64 `json:"lat" bson:"lat"`
        Lng float64 `json:"lng" bson:"lng"`
    } `json:"coordinate"        bson:"coordinate"`
}


func GenerateID() int{
    if RandomID == 0{
        for RandomID == 0{
            RandomID = rand.Intn(99999)
        }
    }else{
        RandomID = RandomID + 1
    }
    return RandomID
}


func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Create",
		"POST",
		"/locations",
		Create,
	},
	Route{
		"Query",
		"GET",
		"/locations/{location_id}",
		Query,
	},
	Route{
		"Update",
		"PUT",
		"/locations/{location_id}",
		Update,
	},
	Route{
		"Delete",
		"DELETE",
		"/locations/{location_id}",
		Delete,
	},
}

type LocationAllInfo struct {
	Results []struct {
    AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry struct {
			Location struct {
			       Lat float64 `json:"lat"`
			       Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport  struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

func GoogleAPI(Address string) (Response, error) {

	url1 :=  "http://maps.google.com/maps/api/geocode/json?address="
	url2 := url.QueryEscape(Address)
	url3 := "&sensor=false"
	fullUrl := url1 + url2 +url3

	fmt.Println(fullUrl)
	var locationAllInfo Response
	var l LocationAllInfo
	
	res, err := http.Get(fullUrl)
	if err!=nil {
		fmt.Println("GoogleAPI: http.Get",err)
		return locationAllInfo,err
	}
	defer res.Body.Close()

	body,err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Println("GoogleAPI: ioutil.ReadAll",err)
		return locationAllInfo,err
	}

	err = json.Unmarshal(body, &l)

	if err!=nil {
		fmt.Println("GoogleAPI: json.Unmarshal",err)
		return locationAllInfo,err
	}

	locationAllInfo.Coordinate.Lat = l.Results[0].Geometry.Location.Lat;
	locationAllInfo.Coordinate.Lng = l.Results[0].Geometry.Location.Lng;


	return locationAllInfo,nil

}
func Create(w http.ResponseWriter, r *http.Request) {
    var req Arguments
    var res Response

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }

    if err := json.Unmarshal(body, &req); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
        fmt.Println("Unmarshal Json Error.", body)
        return
    }
    fullAddress := req.Address+","+req.City+","+req.State+","+req.Zip

    Information,err := GoogleAPI(fullAddress);

    res.Address = req.Address;
    res.City = req.City;
    res.State = req.State;
    res.Zip = req.Zip;
    res.Name = req.Name;
    res.Coordinate.Lat = Information.Coordinate.Lat
    res.Coordinate.Lng = Information.Coordinate.Lng
    MongoCreate(&res)

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)

    if err := json.NewEncoder(w).Encode(res); err != nil {
        panic(err)
    }
    return
}

func Delete(w http.ResponseWriter, r *http.Request) {

    var res Response

    vars := mux.Vars(r)
    res.Id = bson.ObjectIdHex(vars["location_id"])

    err := MongoDelete(res.Id )
    if err != nil {

        fmt.Printf(err.Error())
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    return
}

func Update(w http.ResponseWriter, r *http.Request) {
    var req Arguments
    var res Response
    vars := mux.Vars(r)
    res.Id = bson.ObjectIdHex(vars["location_id"])

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }

    if err := json.Unmarshal(body, &req); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) 
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
        fmt.Println("Unmarshal Json Error.", body)
        return
    }

    fullAddress := req.Address+","+req.City+","+req.State+","+req.Zip

    Information,err := GoogleAPI(fullAddress);

    res.Address = req.Address;
    res.City = req.City;
    res.State = req.State;
    res.Zip = req.Zip;
    res.Coordinate.Lat = Information.Coordinate.Lat
    res.Coordinate.Lng = Information.Coordinate.Lng

    Information, err = MongoUpdate(res)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    res.Name = Information.Name

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)

    if err := json.NewEncoder(w).Encode(res); err != nil {
        panic(err)
    }

}

func Query(w http.ResponseWriter, r *http.Request) {

    var res Response
    var err error
    vars := mux.Vars(r)
    res.Id = bson.ObjectIdHex(vars["location_id"])


    res, err = MongoQuery(res.Id )
    if err != nil {

        fmt.Printf(err.Error())
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    if err := json.NewEncoder(w).Encode(res); err != nil {
        panic(err)
    }
}

func MongoCreate(Information *Response) {
//mongodb://<dbuser>:<dbpassword>@ds039484.mongolab.com:39484/<db_name>
	sess, err := mgo.Dial("mongodb://minglu:liu273@ds057862.mongolab.com:57862/mingluliumongodb")
	if err != nil {
		fmt.Printf("MongoDB connection error %v\n", err)
		panic(err)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	db := sess.DB("mingluliumongodb").C("Location")

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
	collection := sess.DB("mingluliumongodb").C("Location")

	err = collection.Remove(bson.M{"_id": Id})
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
	collection := sess.DB("mingluliumongodb").C("Location")

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
	collection := sess.DB("mingluliumongodb").C("Location")

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
type Service struct{}



func main() {

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
