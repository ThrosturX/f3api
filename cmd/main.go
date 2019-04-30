package main

import "github.com/ThrosturX/f3api"

// Example API server, for demonstration purposes
// Uses In-Memory storage as opposed to long-term stable storage
func main() {
    store := f3api.NewInMemStore()
    api := f3api.NewGenericApi(store)

    f3api.RunServer(api)
}
