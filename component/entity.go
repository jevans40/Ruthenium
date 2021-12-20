package component

type EntityID int

type Entity struct {
	EntityNum EntityID
	Deleted   bool
}
