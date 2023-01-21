// Package tag provides a basic primitive for standard key/value style tags.
//
// # Conversion
//
// Conversion methods are provided for the two most common use-cases: a hash of
// string key/value pairs and a slice of consective key/value strings. No
// ordering is guaranteed because tags should be considered unordered in nearly
// every case.
//
// # Caveat
//
// A Group is additive only, meaning you can add a tag to an existing group but
// you can't delete one. There are ways to get around this limitation, but in
// general you shouldn't use this package to track changeable state.
//
// The ideal use-case for this functionality is in a web request pipeline where
// various tags may be added at different layers of the stack and then later
// used for monitoring/logging/etc.
//
// # Context
//
// To facilitate the use-case listed above, a number of conveniences are
// provided to use tags in conjunction with a Go context.
//
// The WithTag and WithGroup helpers will return a copy of your original context
// with the new tag(s) attached, but they will also create a new tag group on
// the context if one doesn't exist. This means the first time you add any tag
// you're setup to add them anywhere further down the chain.
//
// Likewise, the FromContext helper will return an empty Group if one doesn't
// exist. The result of all of this is that you can't ever end up in a situation
// where your program panics either because you tried to use an unprepared
// context or you got a nil value back when you expected a usable Group.
package tag
