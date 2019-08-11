package client

type BulkDTO struct {
	Id string

	Index string

	Type string

	Data interface{}

	IsDelete bool
}

type MultiGetDTO struct {
	Id string

	Index string

	Type string
}
