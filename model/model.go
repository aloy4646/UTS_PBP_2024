package model

type Account struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type Game struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Room struct {
	Id           int                 `json:"id"`
	Room_name    string              `json:"room_name"`
	Participants []ParticipantDetail `json:"participants"`
}

type RoomWithoutGame struct {
	Id        int    `json:"id"`
	Room_name string `json:"room_name"`
}

type Participant struct {
	Id      int     `json:"id"`
	Room    Room    `json:"room"`
	Account Account `json:"account"`
}

type ParticipantDetail struct {
	Id         int    `json:"id"`
	Id_account int    `json:"id_account"`
	Username   string `json:"username"`
}
