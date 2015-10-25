package main
import (
    "fmt"
    "net/http"
     "encoding/json"
     "io/ioutil"
     "gopkg.in/mgo.v2/bson"
     "github.com/gorilla/mux"
     "io"
    )

// var RandomID int

type Request struct{
    Name string `json:"name"    bson:"name"`
    Address string `json:"address" bson:"address"`
    City string `json:"city"    bson:"city"`
    State string `json:"state"   bson:"state"`
    Zip int `json:"zip"     bson:"zip"`
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

// type Coordmap struct{
//     lat float32
//     lng float32
// }

// type Coordmap map[string]float32


// func main() {
//     RandomID = 0
//     mux := httprouter.New()
//     mux.GET("/hello/:name", hello)
//     mux.POST("/hello", greeting)
//     session, err := mgo.Dial("localhost")
//     if err != nil {
//     return err
//     // server := http.Server{
//     //         Addr:        "0.0.0.0:8080",
//     //         Handler: mux,
//     // }
//     // server.ListenAndServe()
// }

// func (s *Server) handle(w http.ResponseWriter, r *http.Request){
//     session := s.session.copy()
//     defer seesion.close()

//     var req Request
//     var res Response

//     req := session.DB("mingludb").C("req")
//     if err != nil {
//         panic("error")
//     }
    
//     err = json.Unmarshal(body, &req)
//     if err != nil {
//         panic("error Unmarshal")
//     }
//     res.Id = GenerateID()
//     res.Info = req.Name + "," + req.Address + "," + req.City + "," + req.State + "," + req.Zip
//     b,_:=json.Marshal(res)
//     fmt.Fprintf(w, string(b)+"\n")
// }

// func getCoordinates() ([]byte, error) {
//     coordmap := make(map[string]int)

//     d := Data{coordmap}
//     p
// }

type Service struct{}

func Create(w http.ResponseWriter, r *http.Request) {
    var req Request
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

    Information,err := QueryInfo(fullAddress);

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
    var req Request
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

    Information,err := QueryInfo(fullAddress);

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



// func GenerateID() int{
//     if RandomID == 0{
//         for RandomID == 0{
//             RandomID = rand.Intn(10000)
//         }
//     }else{
//         RandomID = RandomID + 1
//     }
//     return RandomID
// }
