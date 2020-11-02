MPRIS Library
=============

This library contains both a client and server for the MPRIS protocol, which the
musicwand library uses.

The client code works similarly to how many clients work. You get a handle, then
are able to search for the player you want and send it commands using methods.

```go
client, err := mpris.NewClient()
if err != nil {
	fmt.Fprintf(os.Stderr, err.Error())
	os.Exit(1)
}

app = client.FindApp("vlc")
if app.Identity() == "" {
	fmt.Fprintf(os.Stderr, "Couldn't find the vlc player interface")
	os.Exit(1)
}
player = app.Player()

player.Play()
player.Next()
```

The server is more complex. You must instantiate a new server, then supply it
with subservers which can answer the necessary functions of the MPRIS API. You
can optionally attach extra interfaces onto the same object to augment the
behavior.

```go
	server, err := mpris.NewServer("musicwand")
	if err != nil {
		log.Fatal(err)
	}

	server.PropertyHandler = &propertyHandler{}
	server.AppServer = &appServer{}
	server.PlayerServer = &playerServer{}

  server.AddInterface("com.github.username.service", &customServer{})

	if err := server.Listen(); err != nil {
		log.Fatal(err)
	}
```

Examples of this can be seen in the musicwand application source code.
