// Code generated by pg. DO NOT EDIT.
package definition

// Column is a struct that represents a column in a table.
type Column struct {
	Name string
}

// String returns the name of the column.
func (c Column) String() string {
	return c.Name
}