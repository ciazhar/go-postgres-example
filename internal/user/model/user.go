package model

type FetchParam struct {
	Page int `query:"page"`
	Size int `query:"size"`
}
