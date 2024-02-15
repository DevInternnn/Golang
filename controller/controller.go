package controller
type Controller struct {
    ctx = context.Background()

    client,err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    // Create a client with appropriate options
    if err != nil {
        fmt.Println("Error connecting to MongoDB:", err)
        return
    }

    // Ping the server to check connectivity
    err = client.Ping(ctx, nil)
    if err != nil {
        fmt.Println("Error pinging MongoDB:", err)
        return
    }

    // Use the client for database operations
    fmt.Println("Successfully connected to MongoDB!")

    // Remember to close the client when done
    defer func() {
        if err := client.Disconnect(ctx); err != nil {
            fmt.Println("Error disconnecting from MongoDB:", err)
        }
    }()
}