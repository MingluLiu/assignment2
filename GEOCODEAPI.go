package main
import (
	"net/url"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

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

func GoogleAPI(Address string) (Info, error) {

	url1 :=  "http://maps.google.com/maps/api/geocode/json?address="
	url2 := url.QueryEscape(Address)
	url3 := "&sensor=false"
	fullUrl := url1 + url2 +url3

	fmt.Println(fullUrl)
	var locationAllInfo Info
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
