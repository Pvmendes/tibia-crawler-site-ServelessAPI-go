package charName

import (
	"fmt"
	model "golang-Serveless-characters/pkg/model"
	valid "golang-Serveless-characters/pkg/validators"
	"strings"
)

var scrapeUrlBase string = "https://www.tibia.com/community/?name="

func GetCharInfo(name string) (*model.Character, error) {
	err := valid.IsCharNameValid(charNameEscapeString(name))

	if checkErr(err) {
		return new(model.Character), err
	}

	item := new(model.Character)
	return item, nil
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
