package res

import (
	"fmt"
	"hash/crc32"
	"strings"

	"gopkg.in/yaml.v3"
)

// RID is a resource identifier. The first 32-bits contain the category, and the last 32-bits contain the ID.
type RID uint64

func (r *RID) UnmarshalYAML(value *yaml.Node) (err error) {
	if value.Kind == yaml.ScalarNode {
		*r, err = RIDFromString(value.Value)
	}
	return err
}

func RIDFromString(s string) (RID, error) {
	var r RID
	err := r.FromString(s)
	return r, err
}

// ErrShortRID is returned when a RID string is too short.
var ErrShortRID = fmt.Errorf("rid string is too short")

// FromString sets the category and ID based on a string in the format of "category:id". If there is no colon, ErrShortRID is returned. If there is more than one colon, all subsequent are treated as part of the ID.
func (r *RID) FromString(s string) error {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return ErrShortRID
	}
	r.Set(parts[0], parts[1])
	return nil
}

// Set sets the category and ID of r.
func (r *RID) Set(cat string, item string) {
	cs := GetCID(cat)
	is := GetIID(item)
	// Assign cs to the first 32 bits of r.
	*r = (*r & 0xffffffff00000000) | RID(cs)
	// Assign is to the last 32 bits of r.
	*r = (*r & 0xffffffff) | RID(is)<<32
}

// SetCat sets the category of r.
func (r *RID) SetCat(cat string) {
	cs := GetCID(cat)
	// Assign cs to the first 32 bits of r.
	*r = (*r & 0xffffffff00000000) | RID(cs)
}

// SetID sets the ID of r.
func (r *RID) SetID(id uint32) {
	// Assign id to the last 32 bits of r.
	*r = (*r & 0xffffffff) | RID(id)<<32
}

// GetCID returns the CRC32 checksum of cat.
func GetCID(cat string) uint32 {
	cs := crc32.ChecksumIEEE([]byte(cat))
	return uint32(cs)
}

// GetIID returns the CRC32 checksum of id.
func GetIID(id string) uint32 {
	i := crc32.ChecksumIEEE([]byte(id))
	return uint32(i)
}
