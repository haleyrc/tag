package tag

import (
	"context"
)

type contextKey int

const groupKey contextKey = 1

// FromContext returns the tag Group stored on the context. If there is no
// Group, FromContext will return a new, empty Group.
func FromContext(ctx context.Context) Group {
	tmp := ctx.Value(groupKey)
	if tmp == nil {
		return NewGroup(nil)
	}
	return tmp.(Group)
}

// WithTag returns a copy of the context with a new tag added to the tag Group.
// If the provided context is empty (it has no tag Group), a new Group will be
// added with the new tag as a member.
func WithTag(ctx context.Context, key, value string) context.Context {
	tmp := ctx.Value(groupKey)

	if tmp == nil {
		return contextWithGroup(ctx, NewGroup(Dict{key: value}))
	}

	g := tmp.(Group)
	g.Add(key, value)
	return contextWithGroup(ctx, g)
}

// WithGroup returns a copy of the context with the tags in the Group added to
// the context. If there is not currently a Group on the context, the provided
// Group will be added as-is. If a Group already exists, the key/value pairs
// from the provided group will be added to the pairs in the existing Group,
// overwriting values in the case where a tag already existed with the same key.
func WithGroup(ctx context.Context, g Group) context.Context {
	tmp := ctx.Value(groupKey)

	if tmp == nil {
		return contextWithGroup(ctx, g)
	}

	curr := tmp.(Group)
	curr.Merge(g)
	return contextWithGroup(ctx, curr)
}

func contextWithGroup(ctx context.Context, g Group) context.Context {
	ctx = context.WithValue(ctx, groupKey, g)
	return ctx
}
