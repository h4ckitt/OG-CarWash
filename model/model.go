package model

import "mime/multipart"

type Owner struct {
	UUID        string    `json:"uuid"`
	LastName    string    `json:"lastName"`
	FirstName   string    `json:"firstName"`
	PhoneNumber string    `json:"phoneNumber"`
	Email       string    `json:"email"`
	CarWashes   []CarWash `json:"carWashes"`
}

type CarWash struct {
	Name    string `json:"carWashName"`
	Address string `json:"carWashAddress"`
}

type CarWashes struct {
	//ClientNumber string `json:"-"`
	CarWashID   string `json:"carWashID"`
	CarWashName string `json:"carWashName"`
	CarsEntered int    `json:"carsEntered"`
	Cars        []Wash `json:"cars"`
}

type Wash struct {
	ClientNumber string `json:"-" bson:"Phone number"`
	CarWashID    string `json:"carWashID,omitempty" bson:"Car wash id"`
	CarsEntered  int    `json:"carsEntered,omitempty"`
	CarWashName  string `json:"-"`
	ImageName    string
	ImageExt     string
	Image        multipart.File
	NumberPlate  string `json:"license" bson:"Plate"`
	DateEntered  string `json:"dateEntered,omitempty" bson:"Day entered"`
	TimeEntered  string `json:"enteredAt" bson:"Time entered"`
	TimeLeft     string `json:"timeLeft"`
}

type WebSocketResult struct {
	Date    string                `json:"date"`
	Clients WebSocketClientResult `json:"clients"`
}

type WebSocketClientResult struct {
	ClientNumber string      `json:"client"`
	Result       []CarWashes `json:"result"`
}
