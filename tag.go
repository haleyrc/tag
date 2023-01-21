package tag

import (
	"fmt"
	"sort"
)

// A tag is a simple key value pair.
type tag struct {
	Key   string
	Value string
}

// Dict is a map of key/value pairs for creating tags.
type Dict map[string]string

// NewGroup returns a group cf tags. Tags with duplicate keys will be overridden
// in a non-deterministic way, so it's important to ensure you pass unique keys.
func NewGroup(l Dict) Group {
	g := Group{
		tags: map[string]tag{},
	}
	for k, v := range l {
		g.tags[k] = tag{Key: k, Value: v}
	}
	return g
}

// Group is a container for multiple flags.
//
// Groups are additive; a tag can be added but not removed. If your application
// has more complicated state changes than simply an append-only list, the best
// solution for using a Group is to save creation of the group until the end of
// your pipeline when all the state has settled.
type Group struct {
	tags map[string]tag
}

// Add adds a new tag to the group. Passing a key that already exists in the
// group will overwrite the existing value. This may be what you want, but be
// sure you know what you're doing.
func (g *Group) Add(key, value string) {
	g.tags[key] = tag{Key: key, Value: value}
}

// Addf is like Add, but the format string and args are combined to form the
// value using the rules of the fmt package verbs. This can be helpful when you
// have a non-string value like a status code that you want to tag with. For
// example:
//
//	g.Add("status", "%d", statusCode)
func (g *Group) Addf(key string, format string, args ...interface{}) {
	g.tags[key] = tag{
		Key:   key,
		Value: fmt.Sprintf(format, args...),
	}
}

// Get returns the value of the tag with the given key if it exists in the
// group. If no tag with the key exists, Get returns an empty string.
func (g *Group) Get(key string) string {
	t, ok := g.tags[key]
	if !ok {
		return ""
	}
	return t.Value
}

// Lookup returns the value of the tag with the given key if it exists in the
// database as well as a true value indicating that the tag was found. If no
// tag with the given key is found, Lookup returns an empty string and a false
// value.
//
// This can be helpful when you may have tags with empty values and need to
// distinguish that case from a tag's nonexistence.
func (g *Group) Lookup(key string) (string, bool) {
	t, ok := g.tags[key]
	if !ok {
		return "", false
	}
	return t.Value, true
}

// Map returns the tags of the group by flattening them into a hash.
func (g *Group) Map() map[string]string {
	hash := map[string]string{}
	for _, t := range g.tags {
		hash[t.Key] = t.Value
	}
	return hash
}

// Merge adds the tags from the provided group to the existing group's tags.
// Tags with keys that already exist in the group will overwrite the existing
// tag.
func (g *Group) Merge(other Group) {
	for _, t := range other.tags {
		g.tags[t.Key] = t
	}
}

// Slice returns the tags of the group by flattening them into a slice of
// strings where the strings at even indices i are the keys for the string
// values at indices i+1. Visually:
//
//	[key1, value1, key2, value2, ...]
//
// The slice returned is deterministic and is sorted by the keys alphabetically
// in increasing order.
func (g *Group) Slice() []string {
	keys := make([]string, 0, len(g.tags))
	for key := range g.tags {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	slice := make([]string, 0, 2*len(g.tags))
	for _, key := range keys {
		t := g.tags[key]
		slice = append(slice, t.Key)
		slice = append(slice, t.Value)
	}

	return slice
}
