// Code generated by goctl. DO NOT EDIT.
package types

type Request struct {
}

type ResponseCommont struct {
	Res bool `json:"res"`
}

type Response struct {
	List []BookInfo `json:"bookInfo"`
}

type BookInfo struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}
