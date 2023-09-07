package structModel

type CharacterResponse struct {
	Character Character `json:"character"`
}

// Child of CharacterInfo
type Houses struct {
	Name    string `json:"name"`    // The name of the house.
	Town    string `json:"town"`    // The town where the house is located in.
	Paid    string `json:"paid"`    // The date the last paid rent is due.
	HouseID int    `json:"houseid"` // The internal ID of the house.
}

// Child of CharacterInfo
type CharacterGuild struct {
	GuildName string `json:"name,omitempty"` // The name of the guild.
	Rank      string `json:"rank,omitempty"` // The character's rank in the guild.
}

// Child of Character
type CharacterInfo struct {
	Name              string         `json:"name"`                    // The name of the character.
	FormerNames       []string       `json:"former_names,omitempty"`  // List of former names of the character.
	Traded            bool           `json:"traded,omitempty"`        // Whether the character was traded. (last 6 months)
	DeletionDate      string         `json:"deletion_date,omitempty"` // The date when the character will be deleted. (if scheduled for deletion)
	Sex               string         `json:"sex"`                     // The character's sex.
	Title             string         `json:"title"`                   // The character's selected title.
	UnlockedTitles    int            `json:"unlocked_titles"`         // The number of titles the character has unlocked.
	Vocation          string         `json:"vocation"`                // The character's vocation.
	Level             int            `json:"level"`                   // The character's level.
	AchievementPoints int            `json:"achievement_points"`      // The total of achievement points the character has.
	World             string         `json:"world"`                   // The character's current world.
	FormerWorlds      []string       `json:"former_worlds,omitempty"` // List of former worlds the character was in. (last 6 months)
	Residence         string         `json:"residence"`               // The character's current residence.
	MarriedTo         string         `json:"married_to,omitempty"`    // The name of the character's husband/spouse.
	Houses            []Houses       `json:"houses,omitempty"`        // List of houses the character owns currently.
	Guild             CharacterGuild `json:"guild"`                   // The guild that the character is member of.
	LastLogin         string         `json:"last_login,omitempty"`    // The character's last logged in time.
	Position          string         `json:"position,omitempty"`      // The character's special position.
	AccountStatus     string         `json:"account_status"`          // Whether account is Free or Premium.
	Comment           string         `json:"comment,omitempty"`       // The character's comment.
}

// Child of Character
type AccountBadges struct {
	Name        string `json:"name"`        // The name of the badge.
	IconURL     string `json:"icon_url"`    // The URL to the badge's icon.
	Description string `json:"description"` // The description of the badge.
}

// Child of Character
type Achievements struct {
	Name   string `json:"name"`   // The name of the achievement.
	Grade  int    `json:"grade"`  // The grade/stars of the achievement.
	Secret bool   `json:"secret"` // Whether it is a secret achievement or not.
}

// Child of Deaths
type Killers struct {
	Name   string `json:"name"`   // The name of the killer/assist.
	Player bool   `json:"player"` // Whether it is a player or not.
	Traded bool   `json:"traded"` // If the killer/assist was traded after the death.
	Summon string `json:"summon"` // The name of the summoned creature.
}

// Child of Character
type Deaths struct {
	Time    string    `json:"time"`    // The timestamp when the death occurred.
	Level   int       `json:"level"`   // The level when the death occurred.
	Killers []Killers `json:"killers"` // List of killers involved.
	Assists []Killers `json:"assists"` // List of assists involved.
	Reason  string    `json:"reason"`  // The plain text reason of death.
}

// Child of Character
type AccountInformation struct {
	Position     string `json:"position,omitempty"`      // The account's special position.
	Created      string `json:"created,omitempty"`       // The account's date of creation.
	LoyaltyTitle string `json:"loyalty_title,omitempty"` // The account's loyalty title.
}

// Child of Character
type OtherCharacters struct {
	Name     string `json:"name"`               // The name of the character.
	World    string `json:"world"`              // The name of the world.
	Status   string `json:"status"`             // The status of the character being online or offline.
	Deleted  bool   `json:"deleted"`            // Whether the character is scheduled for deletion or not.
	Main     bool   `json:"main"`               // Whether this is the main character or not.
	Traded   bool   `json:"traded"`             // Whether the character has been traded last 6 months or not.
	Position string `json:"position,omitempty"` // // The character's special position.
}

// Child of JSONData
type Character struct {
	CharacterInfo      CharacterInfo      `json:"character"`                     // The character's information.
	AccountBadges      []AccountBadges    `json:"account_badges,omitempty"`      // The account's badges.
	Achievements       []Achievements     `json:"achievements,omitempty"`        // The character's achievements.
	Deaths             []Deaths           `json:"deaths,omitempty"`              // The character's deaths.
	AccountInformation AccountInformation `json:"account_information,omitempty"` // The account information.
	OtherCharacters    []OtherCharacters  `json:"other_characters,omitempty"`    // The account's other characters.
}
