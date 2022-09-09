package model

type Owner struct {
	LastName    string    `json:"lastName"`
	FirstName   string    `json:"firstName"`
	PhoneNumber string    `json:"phoneNumber"`
	CarWashes   []CarWash `json:"carWashes"`
}

type CarWash struct {
	Name    string `json:"carWashName"`
	Address string `json:"carWashAddress"`
}

type Wash struct {
	CarWashID   int    `json:"carWashID,omitempty"`
	CarsEntered int    `json:"carsEntered,omitempty"`
	NumberPlate string `json:"license"`
	DateEntered string `json:"dateEntered,omitempty"`
	TimeEntered string `json:"timeEntered"`
	TimeLeft    string `json:"timeLeft"`
}

type WebSocketResult struct {
	Date    string                  `json:"date"`
	Clients []WebSocketClientResult `json:"clients"`
}

type WebSocketClientResult struct {
	ClientNumber string `json:"client"`
	CarWashName  string `json:"CarWashName"`
	Result       []Wash `json:"result"`
}
