package tag_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/haleyrc/tag"
)

func ExampleGroup_Add() {
	g := tag.NewGroup(tag.Dict{
		"env": "prod",
	})
	fmt.Println(g.Slice())
	g.Add("service", "database")
	fmt.Println(g.Slice())
	// Output: [env prod]
	// [env prod service database]
}

func ExampleGroup_Addf() {
	g := tag.NewGroup(tag.Dict{
		"env": "prod",
	})
	fmt.Println(g.Slice())
	g.Addf("status", "%d", 500)
	fmt.Println(g.Slice())
	// Output: [env prod]
	// [env prod status 500]
}

func ExampleGroup_Get() {
	g := tag.NewGroup(tag.Dict{
		"env":     "prod",
		"user_id": "",
	})
	fmt.Printf(
		"%q %q %q",
		g.Get("env"),
		g.Get("user_id"),
		g.Get("environment"),
	)
	// Output: "prod" "" ""
}

func ExampleGroup_Lookup() {
	g := tag.NewGroup(tag.Dict{
		"env":     "prod",
		"user_id": "",
	})

	value, found := g.Lookup("env")
	fmt.Printf("%q %t\n", value, found)

	value, found = g.Lookup("user_id")
	fmt.Printf("%q %t\n", value, found)

	value, found = g.Lookup("environment")
	fmt.Printf("%q %t\n", value, found)

	// Output: "prod" true
	// "" true
	// "" false
}

func ExampleGroup_Map() {
	g := tag.NewGroup(tag.Dict{
		"env":     "prod",
		"service": "database",
		"region":  "us-east1",
	})
	fmt.Println(g.Map())
	// Output: map[env:prod region:us-east1 service:database]
}

func ExampleGroup_Merge() {
	g := tag.NewGroup(tag.Dict{
		"env":    "prod",
		"status": "200",
	})
	g.Merge(tag.NewGroup(tag.Dict{
		"status":  "500",
		"service": "database",
	}))
	fmt.Println(g.Slice())
	// Output: [env prod service database status 500]
}

func ExampleGroup_Slice() {
	g := tag.NewGroup(tag.Dict{
		"env":     "prod",
		"service": "database",
		"region":  "us-east1",
	})
	fmt.Println(g.Slice())
	// Output: [env prod region us-east1 service database]
}

func TestWithTagCreatesANewGroupOnAnEmptyContext(t *testing.T) {
	ctx := tag.WithTag(context.Background(), "env", "prod")
	group := tag.FromContext(ctx)
	expectSliceEqual(t, group.Slice(), []string{"env", "prod"})
}

func TestWithTagAppendsToAGroupOnAContext(t *testing.T) {
	ctx := tag.WithTag(context.Background(), "env", "prod")
	ctx = tag.WithTag(ctx, "service", "database")
	group := tag.FromContext(ctx)
	expectSliceEqual(t, group.Slice(), []string{"env", "prod", "service", "database"})
}

func TestWithTagOverwritesADuplicateKeyedTag(t *testing.T) {
	ctx := tag.WithTag(context.Background(), "status", "200")
	ctx = tag.WithTag(ctx, "status", "500")
	group := tag.FromContext(ctx)
	expectSliceEqual(t, group.Slice(), []string{"status", "500"})
}

func TestWithGroupCreatesANewGroupOnAnEmptyContext(t *testing.T) {
	ctx := tag.WithGroup(context.Background(), tag.NewGroup(tag.Dict{
		"env": "prod",
	}))
	group := tag.FromContext(ctx)
	expectSliceEqual(t, group.Slice(), []string{"env", "prod"})
}

func TestWithGroupMergesTagsWithExistingTags(t *testing.T) {
	ctx := tag.WithGroup(context.Background(), tag.NewGroup(tag.Dict{
		"env":    "prod",
		"status": "200",
	}))
	ctx = tag.WithGroup(ctx, tag.NewGroup(tag.Dict{
		"status":  "500",
		"service": "database",
	}))
	group := tag.FromContext(ctx)
	expectSliceEqual(t, group.Slice(), []string{"env", "prod", "service", "database", "status", "500"})
}

func TestFromContextReturnsAnEmptyGroupFromAnEmptyContext(t *testing.T) {
	group := tag.FromContext(context.Background())
	expectSliceEqual(t, group.Slice(), []string{})
}

func TestFromContextReturnsAGroupFromAContext(t *testing.T) {
	ctx := tag.WithTag(context.Background(), "env", "prod")
	group := tag.FromContext(ctx)
	got := len(group.Map())
	want := 1
	if got != want {
		t.Errorf("Expected group to have %d tag, but got %d.", want, got)
	}
}

func expectSliceEqual(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Errorf(
			"Slices are of different lengths. Wanted %d, but got %d.",
			len(want), len(got),
		)
		return
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Expected to find %q in position %d, but found %q.", want[i], i, got[i])
		}
	}
}
