package services

import (
	"effective-mobile-test/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func EnrichPerson(name string) models.Person {
	age := getInt(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	gender := getString(fmt.Sprintf("https://api.genderize.io/?name=%s", name), "gender")
	nationality := getString(fmt.Sprintf("https://api.nationalize.io/?name=%s", name), "country", "country_id")

	return models.Person{
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}
}

func getInt(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Request failed:", err)
		return 0
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0
	}
	if age, ok := res["age"].(float64); ok {
		return int(age)
	}
	return 0
}

func getString(url, key string, nested ...string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("Request failed:", err)
	}
	if len(nested) > 0 {
		if arr, ok := res[key].([]interface{}); ok && len(arr) > 0 {
			if item, ok := arr[0].(map[string]interface{}); ok {
				if val, ok := item[nested[0]].(string); ok {
					return val
				}
			}
		}
	} else {
		if val, ok := res[key].(string); ok {
			return val
		}
	}
	return ""
}
