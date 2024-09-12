package OddsApi

import (
	"encoding/json"
	"fmt"
	"go-matcher/config"
	"go-matcher/types"
	"io/ioutil"
	"net/http"
)

// example keys
// soccer_uefa_european_championship
// soccer_epl

func FetchSportsOdds(key string, market string) ([]types.SportsOddsPayload, error) {
	apiUrl := "https://api.the-odds-api.com/v4/sports/" + key + "/odds/?apiKey=" + config.Envs.OddsApiKey + "&regions=au"

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	var jsonPayload []types.SportsOddsPayload
	err = json.Unmarshal(body, &jsonPayload)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		fmt.Println("Raw response body:", string(body))
	}

	return jsonPayload, nil
}
