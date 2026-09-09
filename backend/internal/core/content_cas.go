package core

// NoActiveContentCommit is an internal compare-and-swap sentinel used when a
// refresh observed that no active snapshot existed. It cannot collide with a
// validated hexadecimal commit object ID.
const NoActiveContentCommit = "<no-active-content>"
