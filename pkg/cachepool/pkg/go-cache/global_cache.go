package gocache

// Global cache use NoSQL database as its container
// Since database runs on other process, it's easy to
// ensure consistency among many cache.
// For NoSQL clients are various, here is no implementation (
// try to implement ICache for your own global_cache
// or see redigo-cache
