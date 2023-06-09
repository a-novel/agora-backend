// Package profile_storage is the storage layer for the profile of a user.
// This profile is non-sensitive, and contains the fake profile used by the person online.
// It must not contain anything sensitive or tied to the user (unless it chooses to). One user may have multiple
// associated profiles, although it is currently not implemented.
package profile_storage
