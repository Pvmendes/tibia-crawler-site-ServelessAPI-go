package scraping

import (
	"fmt"

	"strings"

	model "golang-Serveless-characters/pkg/model"
	utils "golang-Serveless-characters/pkg/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// From https://pkg.go.dev/golang.org/x/net/html/atom
// This is an Atom. An Atom is an integer code for a string.
// Instead of importing the whole lib, we thought it would be
// best to just simply use the Br constant value.
const Br = 0x202

var (
	localDivQueryString = ".TableContentContainer tr"
	localTradedString   = " (traded)"

	CharInfo               model.CharacterInfo
	AccountBadgesData      []model.AccountBadges
	AchievementsData       []model.Achievements
	DeathsData             []model.Deaths
	AccountInformationData model.AccountInformation
	OtherCharactersData    []model.OtherCharacters

	// Errors
	characterNotFound bool
	insideError       error
)

func ScrapingInSite(urlSite string) {
	//fmt.Println(urlSite)
	c := colly.NewCollector(colly.AllowedDomains("www.tibia.com", "tibia.com"))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scraping: %s\n", e.Error())
	})

	c.OnHTML("div.TableContainer", func(h *colly.HTMLElement) {
		selection := h.DOM

		caption := selection.Find("div.Text").Text()
		// fmt.Printf(caption)
		switch caption {
		case "Could not find character":
			break
		case "Character Information", "Account Information":
			readAccountInformation(selection, caption)
		case "Account Badges":
			readAccountBadges(selection, caption)
		case "Account Achievements":
			readAccountAchievements(selection, caption)
		case "Character Deaths":
			readCharacterDeaths(selection, caption)
		case "Characters":
			readCharacter(selection, caption)
		}
	})

	c.Visit(urlSite)
}

func readCharacter(selection *goquery.Selection, SectionName string) {
	nodes := selection.Nodes[0]
	characterDivQuery := goquery.NewDocumentFromNode(nodes)

	characterDivQuery.Find(localDivQueryString).EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CharacterListHTML
		CharacterListHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaCharactersCharacterImpl failed at s.Html() inside Characters, err: %s", err)
			return false
		}

		// Removing line breaks
		CharacterListHTML = utils.HTMLRemoveLinebreaks(CharacterListHTML)

		if !strings.Contains(CharacterListHTML, "<td>Name</td><td>World</td><td>Status</td>") {
			const (
				nameIndexer  = `<td style="width: 20%"><nobr>`
				worldIndexer = `<td style="width: 10%"><nobr>`
			)

			nameIdx := strings.Index(
				CharacterListHTML, nameIndexer,
			) + len(nameIndexer)
			nameIdx += strings.Index(
				CharacterListHTML[nameIdx:], " ",
			) + 1
			endNameIdx := strings.Index(
				CharacterListHTML[nameIdx:], `</nobr></td>`,
			) + nameIdx

			tmpCharName := CharacterListHTML[nameIdx:endNameIdx]

			worldIdx := strings.Index(
				CharacterListHTML, worldIndexer,
			) + len(worldIndexer)
			endWorldIdx := strings.Index(
				CharacterListHTML[worldIdx:], `</nobr></td>`,
			) + worldIdx

			world := CharacterListHTML[worldIdx:endWorldIdx]

			var tmpTraded bool
			if strings.Contains(tmpCharName, localTradedString) {
				tmpTraded = true
				tmpCharName = strings.ReplaceAll(tmpCharName, localTradedString, "")
			}

			// If this character is the main character of the account
			var tmpMain bool
			if strings.Contains(tmpCharName, "Main Character") {
				tmpMain = true
				tmp := strings.Split(tmpCharName, "<")
				tmpCharName = strings.TrimSpace(tmp[0])
			}

			// If this character is online or offline
			tmpStatus := "offline"
			if strings.Contains(CharacterListHTML, "<b class=\"green\">online</b>") {
				tmpStatus = "online"
			}

			// Is this character is deleted
			var tmpDeleted bool
			if strings.Contains(CharacterListHTML, "deleted") {
				tmpDeleted = true
			}

			// Is this character having a special position
			var tmpPosition string
			if strings.Contains(CharacterListHTML, "CipSoft Member") {
				tmpPosition = "CipSoft Member"
			}

			// Create the character and append it to the other characters list
			OtherCharactersData = append(OtherCharactersData, model.OtherCharacters{
				Name:     utils.SanitizeStrings(tmpCharName),
				World:    world,
				Status:   tmpStatus,
				Deleted:  tmpDeleted,
				Main:     tmpMain,
				Traded:   tmpTraded,
				Position: tmpPosition,
			})
		}

		return true
	})
}

func readCharacterDeaths(selection *goquery.Selection, SectionName string) {
	nodes := selection.Nodes[0]
	characterDivQuery := goquery.NewDocumentFromNode(nodes)

	characterDivQuery.Find(localDivQueryString).EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CharacterListHTML
		CharacterListHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaCharactersCharacterImpl failed at s.Html() inside Character Deaths, err: %s", err)
			return false
		}

		// Removing line breaks
		CharacterListHTML = utils.HTMLRemoveLinebreaks(CharacterListHTML)
		CharacterListHTML = strings.ReplaceAll(CharacterListHTML, ".<br/>Assisted by", ". Assisted by")
		CharacterListHTML = utils.SanitizeStrings(CharacterListHTML)

		dataNoTags := utils.RemoveHtmlTag(CharacterListHTML)

		// defining responses
		DeathKillers := []model.Killers{}
		DeathAssists := []model.Killers{}

		const (
			initIndexer    = `CET`
			levelIndexer   = `at Level `
			killersIndexer = `by `
		)

		initIdx := strings.Index(
			dataNoTags, initIndexer,
		) + len(initIndexer)
		endInitIdx := strings.Index(
			dataNoTags[initIdx:], `by `,
		) + initIdx + len(`by `)

		reasonStart := dataNoTags[initIdx:endInitIdx]
		reasonRest := dataNoTags[endInitIdx:]

		// store for reply later on.. and sanitizing string
		reasonString := reasonStart + reasonRest

		timeIdx := 0
		endTimeIdx := strings.Index(
			dataNoTags[timeIdx:], `CET`,
		) + timeIdx + len(`CET`)

		time := utils.TibiaDataDatetime(dataNoTags[timeIdx:endTimeIdx])

		levelIdx := strings.Index(
			dataNoTags, levelIndexer,
		) + len(levelIndexer)
		endLevelIdx := strings.Index(
			dataNoTags[levelIdx:], " ",
		) + levelIdx

		level := utils.StringToInteger(dataNoTags[levelIdx:endLevelIdx])

		killersIdx := strings.Index(
			CharacterListHTML, killersIndexer,
		) + len(killersIndexer)
		endKillersIdx := strings.Index(
			CharacterListHTML[killersIdx:], `.</td>`,
		) + killersIdx

		rawListofKillers := CharacterListHTML[killersIdx:endKillersIdx]

		// if kill is with assist..
		if strings.Contains(dataNoTags, ". Assisted by ") {
			TmpListOfDeath := strings.Split(CharacterListHTML, ". Assisted by ")
			rawListofKillers = TmpListOfDeath[0][killersIdx:]
			TmpAssist := TmpListOfDeath[1]

			// get a list of killers
			ListOfAssists := strings.Split(TmpAssist, ", ")

			// extract if "and" is in last ss1
			ListOfAssistsTmp := strings.Split(ListOfAssists[len(ListOfAssists)-1], " and ")

			// if there is an "and", then we split it..
			if len(ListOfAssistsTmp) > 1 {
				ListOfAssists[len(ListOfAssists)-1] = ListOfAssistsTmp[0]
				ListOfAssists = append(ListOfAssists, ListOfAssistsTmp[1])
			}

			for i := range ListOfAssists {
				name, isPlayer, isTraded, theSummon := TibiaDataParseKiller(ListOfAssists[i])
				DeathAssists = append(DeathAssists, model.Killers{
					Name:   strings.TrimSuffix(strings.TrimSuffix(name, ".</td>"), "."),
					Player: isPlayer,
					Traded: isTraded,
					Summon: theSummon,
				})
			}
		}

		// get a list of killers
		ListOfKillers := strings.Split(rawListofKillers, ", ")

		// extract if "and" is in last ss1
		ListOfKillersTmp := strings.Split(ListOfKillers[len(ListOfKillers)-1], " and ")

		// if there is an "and", then we split it..
		if len(ListOfKillersTmp) > 1 {
			ListOfKillers[len(ListOfKillers)-1] = ListOfKillersTmp[0]
			ListOfKillers = append(ListOfKillers, ListOfKillersTmp[1])
		}

		// loop through all killers and append to result
		for i := range ListOfKillers {
			name, isPlayer, isTraded, theSummon := TibiaDataParseKiller(ListOfKillers[i])
			DeathKillers = append(DeathKillers, model.Killers{
				Name:   strings.TrimSuffix(strings.TrimSuffix(name, ".</td>"), "."),
				Player: isPlayer,
				Traded: isTraded,
				Summon: theSummon,
			})
		}

		// append deadentry to death list
		DeathsData = append(DeathsData, model.Deaths{
			Time:    time,
			Level:   level,
			Killers: DeathKillers,
			Assists: DeathAssists,
			Reason:  reasonString,
		})

		return true
	})
}

func readAccountAchievements(selection *goquery.Selection, SectionName string) {
	nodes := selection.Nodes[0]
	characterDivQuery := goquery.NewDocumentFromNode(nodes)

	characterDivQuery.Find(localDivQueryString).EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CharacterListHTML
		CharacterListHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaCharactersCharacterImpl failed at s.Html() inside Account Achievements, err: %s", err)
			return false
		}

		// Removing line breaks
		CharacterListHTML = utils.HTMLRemoveLinebreaks(CharacterListHTML)

		if !strings.Contains(CharacterListHTML, "There are no achievements set to be displayed for this character.") {
			const (
				nameIndexer = `alt="Tibia Achievement"/></td><td>`
			)

			// get the name of the achievement (and ignore the secret image on the right)
			nameIdx := strings.Index(
				CharacterListHTML, nameIndexer,
			) + len(nameIndexer)
			endNameIdx := strings.Index(
				CharacterListHTML[nameIdx:], `<`,
			) + nameIdx

			AchievementsData = append(AchievementsData, model.Achievements{
				Name:   CharacterListHTML[nameIdx:endNameIdx],
				Grade:  strings.Count(CharacterListHTML, "achievement-grade-symbol"),
				Secret: strings.Contains(CharacterListHTML, "achievement-secret-symbol"),
			})
		}

		return true
	})
}

func readAccountBadges(selection *goquery.Selection, SectionName string) {
	nodes := selection.Nodes[0]
	characterDivQuery := goquery.NewDocumentFromNode(nodes)

	characterDivQuery.Find(localDivQueryString + " td span[style]").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CharacterListHTML
		CharacterListHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaCharactersCharacterImpl failed at s.Html() inside Account Badges, err: %s", err)
			return false
		}

		// Removing line breaks
		CharacterListHTML = utils.HTMLRemoveLinebreaks(CharacterListHTML)

		// prevent failure of regex that parses account badges
		if CharacterListHTML != "There are no account badges set to be displayed for this character." {
			const (
				nameIndexer = `alt="`
				iconIndexer = `img src="`
				descIndexer = `&#39;, &#39;`
			)

			nameIdx := strings.Index(
				CharacterListHTML, nameIndexer,
			) + len(nameIndexer)
			endNameIdx := strings.Index(
				CharacterListHTML[nameIdx:], `"`,
			) + nameIdx

			iconIdx := strings.Index(
				CharacterListHTML, iconIndexer,
			) + len(iconIndexer)
			endIconIdx := strings.Index(
				CharacterListHTML[iconIdx:], `"`,
			) + iconIdx

			descIdx := strings.Index(
				CharacterListHTML, descIndexer,
			) + len(descIndexer)
			endDescIdx := strings.Index(
				CharacterListHTML[descIdx:], descIndexer,
			) + descIdx

			AccountBadgesData = append(AccountBadgesData, model.AccountBadges{
				Name:        CharacterListHTML[nameIdx:endNameIdx],
				IconURL:     CharacterListHTML[iconIdx:endIconIdx],
				Description: CharacterListHTML[descIdx:endDescIdx],
			})
		}
		fmt.Print(AccountBadgesData)
		return true
	})
}

func readAccountInformation(selection *goquery.Selection, SectionName string) {

	nodes := selection.Nodes[0]
	characterDivQuery := goquery.NewDocumentFromNode(nodes)

	characterDivQuery.Find(localDivQueryString).Each(func(index int, s *goquery.Selection) {
		rowNameQuery := s.Find("td[class^='Label']")

		fmt.Printf("\n")

		rowName := rowNameQuery.Nodes[0].FirstChild.Data
		rowData := rowNameQuery.Nodes[0].NextSibling.FirstChild.Data

		//fmt.Println(rowName)
		//fmt.Println(sanitizeStrings(rowData))
		switch utils.SanitizeStrings(rowName) {
		case "Name:":
			Tmp := strings.Split(rowData, "<")
			CharInfo.Name = strings.TrimSpace(Tmp[0])
			if strings.Contains(Tmp[0], ", will be deleted at") {
				Tmp2 := strings.Split(Tmp[0], ", will be deleted at ")
				CharInfo.Name = Tmp2[0]
				CharInfo.DeletionDate = utils.TibiaDataDatetime(strings.TrimSpace(Tmp2[1]))
			}
			if strings.Contains(rowData, localTradedString) {
				CharInfo.Traded = true
				CharInfo.Name = strings.Replace(CharInfo.Name, localTradedString, "", -1)
			}
		case "Former Names:":
			CharInfo.FormerNames = strings.Split(rowData, ", ")
		case "Sex:":
			CharInfo.Sex = rowData
		case "Title:":
			leftParenIdx := strings.Index(rowData, "(")
			if leftParenIdx == -1 {
				return
			}

			title := rowData[:leftParenIdx-1]

			spaceIdx := strings.Index(rowData[leftParenIdx:], " ")
			if spaceIdx == -1 {
				return
			}

			unlockedTitles := utils.StringToInteger(
				rowData[leftParenIdx+1 : leftParenIdx+spaceIdx],
			)

			CharInfo.Title = title
			CharInfo.UnlockedTitles = unlockedTitles
		case "Vocation:":
			CharInfo.Vocation = rowData
		case "Level:":
			CharInfo.Level = utils.StringToInteger(rowData)
		case "nobr", "Achievement Points:":
			CharInfo.AchievementPoints = utils.StringToInteger(rowData)
		case "World:":
			CharInfo.World = rowData
		case "Former World:":
			CharInfo.FormerWorlds = strings.Split(rowData, ", ")
		case "Residence:":
			CharInfo.Residence = rowData
		case "Account Status:":
			CharInfo.AccountStatus = rowData
		case "Married To:":
			AnchorQuery := s.Find("a")
			CharInfo.MarriedTo = AnchorQuery.Nodes[0].FirstChild.Data
		case "House:":
			AnchorQuery := s.Find("a")
			HouseName := AnchorQuery.Nodes[0].FirstChild.Data
			HouseHref := AnchorQuery.Nodes[0].Attr[0].Val
			//substring from houseid= to &character in the href for the house
			HouseId := HouseHref[strings.Index(HouseHref, "houseid")+8 : strings.Index(HouseHref, "&character")]
			HouseRawData := rowNameQuery.Nodes[0].NextSibling.LastChild.Data
			HouseTown := HouseRawData[strings.Index(HouseRawData, "(")+1 : strings.Index(HouseRawData, ")")]
			HousePaidUntil := HouseRawData[strings.Index(HouseRawData, "is paid until ")+14:]

			CharInfo.Houses = append(CharInfo.Houses, model.Houses{
				Name:    HouseName,
				Town:    HouseTown,
				Paid:    utils.TibiaDataDate(HousePaidUntil),
				HouseID: utils.StringToInteger(HouseId),
			})
		case "Guild Membership:":
			CharInfo.Guild.Rank = strings.TrimSuffix(rowData, " of the ")

			//TODO: I don't understand why the unicode nbsp is there...
			CharInfo.Guild.GuildName = utils.SanitizeStrings(rowNameQuery.Nodes[0].NextSibling.LastChild.LastChild.Data)
		case "Last Login:":
			if rowData != "never logged in" {
				CharInfo.LastLogin = utils.TibiaDataDatetime(rowData)
			}
		case "Comment:":
			node := rowNameQuery.Nodes[0].NextSibling.FirstChild

			stringBuilder := strings.Builder{}
			for node != nil {
				if node.DataAtom == Br {
					//It appears we can ignore br because either the encoding or goquery adds an \n for us
					//stringBuilder.WriteString("\n")
				} else {
					stringBuilder.WriteString(node.Data)
				}

				node = node.NextSibling
			}

			CharInfo.Comment = stringBuilder.String()
		case "Loyalty Title:":
			if rowData != "(no title)" {
				AccountInformationData.LoyaltyTitle = rowData
			}
		case "Created:":
			AccountInformationData.Created = utils.TibiaDataDatetime(rowData)
		case "Position:":
			TmpPosition := strings.Split(rowData, "<")
			if SectionName == "Character Information" {
				CharInfo.Position = strings.TrimSpace(TmpPosition[0])
			} else if SectionName == "Account Information" {
				AccountInformationData.Position = strings.TrimSpace(TmpPosition[0])
			}

		default:
		}

	})

	fmt.Print(CharInfo)
}

func TibiaDataParseKiller(data string) (string, bool, bool, string) {
	var (
		// local strings used in this function
		localTradedString = " (traded)"

		isPlayer, isTraded bool
		theSummon          string
	)

	// check if killer is a traded player
	if strings.Contains(data, localTradedString) {
		isPlayer = true
		isTraded = true
		data = strings.ReplaceAll(data, localTradedString, "")
	}

	// check if killer is a player
	if strings.Contains(data, "https://www.tibia.com") {
		isPlayer = true
		data = utils.RemoveHtmlTag(data)
	}

	// get summon information
	if strings.HasPrefix(data, "a ") || strings.HasPrefix(data, "an ") {
		if containsCreaturesWithOf(data) {
			// this is not a summon, since it is a creature with a of in the middle
		} else {
			ofIdx := strings.Index(data, "of")
			if ofIdx != -1 {
				theSummon = data[:ofIdx-1]
				data = data[ofIdx+3:]
			}
		}
	}

	// sanitizing string
	data = utils.SanitizeStrings(data)

	return data, isPlayer, isTraded, theSummon
}

// containsCreaturesWithOf checks if creature is present in special creatures list
func containsCreaturesWithOf(str string) bool {
	// trim away "an " and "a "
	str = strings.TrimPrefix(strings.TrimPrefix(str, "an "), "a ")

	switch str {
	case "acolyte of darkness",
		"acolyte of the cult",
		"adept of the cult",
		"ancient spawn of morgathla",
		"aspect of power",
		"baby pet of chayenne",
		"bane of light",
		"bloom of doom",
		"bride of night",
		"cloak of terror",
		"energuardian of tales",
		"enlightened of the cult",
		"eruption of destruction",
		"essence of darkness",
		"essence of malice",
		"eye of the seven",
		"flame of omrafir",
		"fury of the emperor",
		"ghastly pet of chayenne",
		"ghost of a planegazer",
		"greater splinter of madness",
		"groupie of skyrr",
		"guardian of tales",
		"gust of wind",
		"hand of cursed fate",
		"harbinger of darkness",
		"herald of gloom",
		"izcandar champion of summer",
		"izcandar champion of winter",
		"lesser splinter of madness",
		"lord of the elements",
		"lost ghost of a planegazer",
		"memory of a banshee",
		"memory of a book",
		"memory of a carnisylvan",
		"memory of a dwarf",
		"memory of a faun",
		"memory of a frazzlemaw",
		"memory of a fungus",
		"memory of a golem",
		"memory of a hero",
		"memory of a hydra",
		"memory of a lizard",
		"memory of a mammoth",
		"memory of a manticore",
		"memory of a pirate",
		"memory of a scarab",
		"memory of a shaper",
		"memory of a vampire",
		"memory of a werelion",
		"memory of a wolf",
		"memory of a yalahari",
		"memory of an amazon",
		"memory of an elf",
		"memory of an insectoid",
		"memory of an ogre",
		"mighty splinter of madness",
		"minion of gaz'haragoth",
		"minion of versperoth",
		"monk of the order",
		"muse of penciljack",
		"nightmare of gaz'haragoth",
		"noble pet of chayenne",
		"novice of the cult",
		"pillar of death",
		"pillar of draining",
		"pillar of healing",
		"pillar of protection",
		"pillar of summoning",
		"priestess of the wild sun",
		"rage of mazoran",
		"reflection of mawhawk",
		"reflection of obujos",
		"reflection of a mage",
		"retainer of baeloc",
		"scorn of the emperor",
		"servant of tentugly",
		"shadow of boreth",
		"shadow of lersatio",
		"shadow of marziel",
		"shard of corruption",
		"shard of magnor",
		"sight of surrender",
		"son of verminor",
		"soul of dragonking zyrtarch",
		"spark of destruction",
		"spawn of despair",
		"spawn of devovorga",
		"spawn of havoc",
		"spawn of the schnitzel",
		"spawn of the welter",
		"sphere of wrath",
		"spirit of earth",
		"spirit of fertility",
		"spirit of fire",
		"spirit of light",
		"spirit of water",
		"spite of the emperor",
		"squire of nictros",
		"stolen knowledge of armor",
		"stolen knowledge of healing",
		"stolen knowledge of lifesteal",
		"stolen knowledge of spells",
		"stolen knowledge of summoning",
		"stolen tome of portals",
		"sword of vengeance",
		"symbol of fear",
		"symbol of hatred",
		"tentacle of the deep terror",
		"the book of death",
		"the book of secrets",
		"the cold of winter",
		"the corruptor of souls",
		"the count of the core",
		"the devourer of secrets",
		"the duke of the depths",
		"the heat of summer",
		"the lily of night",
		"the lord of the lice",
		"the scion of havoc",
		"the scourge of oblivion",
		"the source of corruption",
		"the voice of ruin",
		"tin lizzard of lyxoph",
		"undead pet of chayenne",
		"weak harbinger of darkness",
		"weak spawn of despair",
		"wildness of urmahlullu",
		"wisdom of urmahlullu",
		"wrath of the emperor",
		"zarcorix of yalahar":
		return true
	default:
		return false
	}
}
