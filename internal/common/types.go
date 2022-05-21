package common

type ERC721Attr struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
}

type ERC721Meta struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Image       string       `json:"image"`
	Attributes  []ERC721Attr `json:"attributes"`
}
