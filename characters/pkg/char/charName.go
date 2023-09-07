package char

import (
	"fmt"
	model "golang-Serveless-characters/pkg/model"
	scrap "golang-Serveless-characters/pkg/scraping"
	valid "golang-Serveless-characters/pkg/validators"

	"strings"
)

var scrapeUrlBase string = "https://www.tibia.com/community/?name="

func GetCharInfo(name string) (*model.CharacterResponse, error) {

	err := valid.IsCharNameValid(charNameEscapeString(name))

	if checkErr(err) {
		return new(model.CharacterResponse), err
	}

	return scrap.ScrapingInSite(scrapeUrlBase + name)
}

func checkErr(err error) bool {

	if err != nil {
		fmt.Println(err.Error())
		return true
	}
	return false
}

func charNameEscapeString(data string) string {
	data = strings.ReplaceAll(data, "+", " ")

	return data
}
