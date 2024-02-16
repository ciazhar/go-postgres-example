package model

type Data struct {
	Service string `json:"service,omitempty"`
	Version string `json:"version,omitempty"`
	Profile string `json:"profile,omitempty"`
	Doc     string `json:"doc"`
}
