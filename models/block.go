package models

type Channel struct {
	Id           int    `json:"id"`
	Channel_name string `json:"channel_name"`
}
type Network struct {
	Id int    `json:"id"`
	Ip string `json:"ip"`
}
type Block struct {
	Block_number int      `json:"block_number"`
	Channel      *Channel `json:"channel"`
	Network      *Network `json:"Network"`
	Prev_hash    string   `json:"Prev_hash"`
	Next_hash    string   `json:"Next_hash"`
	Data         string   `json:"data"`
}
