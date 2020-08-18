package data

var users = []User{
	{
		Name:     "田中 太郎",
		Email:    "tarou@gmail.com",
		Password: "taroupass",
	},
	{
		Name:     "John Smith",
		Email:    "john@gmail.com",
		Password: "johnpass",
	},
}

var characters = []Character{
	{
		Name:  "テストくん",
		Text:  "テストキャラクター１号です",
		Image: "/characters/default.jpg",
	},
	{
		Name:  "",
		Text:  "テストキャラクター２号です",
		Image: "/characters/default.jpg",
	},
}

func setup() {
	CharacterDeleteAll()
	SessionDeleteAll()
	UserDeleteAll()
}
