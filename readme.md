With `go-astilectron` create a beautiful GUI for you GO app. It is the official GO bindings of [astilectron](https://github.com/asticode/astilectron) and is powered by [Electron](https://github.com/electron/electron).

# Quick start

WARNING: the code below doesn't handle errors for readibility purposes. However you SHOULD!

### Import astilectron

To import astilectron run:

    $ go get -u github.com/asticode/go-astilectron

### Start astilectron

    // Initialize astilectron
    var a, _ = astilectron.New(astilectron.Options{BaseDirectoryPath: "<where you want the provisioner to install the dependencies>"})
    defer a.Close()
    
    // Start astilectron
    a.Start()

For everything to work properly we need to fetch 2 dependencies : [astilectron](https://github.com/asticode/astilectron) and [Electron](https://github.com/electron/electron). `.Start()` takes care of it by downloading the sources and setting them up properly.

In case you want to embed the sources in the binary to keep a unique binary you can implement your own **Provisioner** and attach it to **Astilectron** with `.SetProvisioner(p Provisioner)`.

If no BaseDirectoryPath is provided, it defaults to the user's home directory path.

The majority of methods are synchrone which means that when executing them **Astilectron** will block until it receives a specific Electron event. This is the case of `.Start()` which will block until it receives the `did-finish-load` Electron WebContents event.

### Create a window

    // Create a new window
    var w, _ = a.NewWindow("http://127.0.0.1:4000", &astilectron.WindowOptions{
        Center: astilectron.PtrBool(true),
        Height: astilectron.PtrInt(600),
        Width:  astilectron.PtrInt(600),
    })
    w.Create()
    
When creating a window you need to indicate a URL as well as options that fit Electron's window options.

This is pretty straightforward except the `astilectron.Ptr*` methods so let me explain: GO doesn't do optional fields when json encoding unless you use pointers whereas Electron does handle optional fields. Therefore I added helper methods to convert int, bool and string into pointers.

### Add listeners

    // Add a listener on Astilectron
    a.On(astilectron.EventNameAppClose, func(e astilectron.Event) (deleteListener bool) {
        a.Stop()
        return
    })
    
    // Add a listener on the window
    w.On(astilectron.EventNameWindowEventResize, func(e astilectron.Event) (deleteListener bool) {
        astilog.Info("Window resized")
        return
    })
    
Nothing much to say here either except that you can add listeners to Astilectron as well.

### Play with the window

    // Play with the window
    w.Resize(200, 200)
    time.Sleep(time.Second)
    w.Maximize()
    
Check out the [Window doc](https://godoc.org/github.com/asticode/go-astilectron#Window) for a list of all exported methods

### Send messages between GO and your webserver

In your webserver add the following script to any of the pages you want to interact with:

    <script>
        // This will wait for the astilectron namespace to be ready
        document.addEventListener('astilectron-ready', function() {
            // This will listen to messages sent by GO
            astilectron.listen(function(message) {
                document.getElementById('message').innerHTML = message
                // This will send a message back to GO
                astilectron.send("I'm good bro")
            });
        })
    </script>
    
In your GO app add the following:
    
    // Listen to messages sent by webserver
    w.On(astilectron.EventNameWindowEventMessage, func(e astilectron.Event) (deleteListener bool) {
        var m string
        e.Message.Unmarshal(&m)
        astilog.Infof("Received message %s", m)
        return
    })
    
    // Send message to webserver
    w.Send("What's up?")
    
And that's it!

NOTE: needless to say that the message can be something other than a string. JSON is used to exchange data between GO and the webserver so any type will do!*

### Final code

    // Set up the logger
    var l <your logger type>
    astilog.SetLogger(l)
    
    // Start an http server
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`<!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <title>Hello world</title>
        </head>
        <body>
            <span id="message">Hello world</span>
            <script>
                // This will wait for the astilectron namespace to be ready
                document.addEventListener('astilectron-ready', function() {
                    // This will listen to messages sent by GO
                    astilectron.listen(function(message) {
                        document.getElementById('message').innerHTML = message
                        // This will send a message back to GO
                        astilectron.send("I'm good bro")
                    });
                })
            </script>
        </body>
        </html>`))
    })
    go http.ListenAndServe("127.0.0.1:4000", nil)
    
    // Initialize astilectron
    var a, _ = astilectron.New(astilectron.Options{BaseDirectoryPath: "<where you want the provisioner to install the dependencies>"})
    defer a.Close()
    
    // Handle quit
    a.HandleSignals()
    a.On(astilectron.EventNameAppClose, func(e astilectron.Event) (deleteListener bool) {
        a.Stop()
        return
    })
    
    // Start astilectron: this will download and set up the dependencies, and start the Electron app
    a.Start()
    
    // Create a new window with a listener on resize
    var w, _ = a.NewWindow("http://127.0.0.1:4000", &astilectron.WindowOptions{
        Center: astilectron.PtrBool(true),
        Height: astilectron.PtrInt(600),
        Width:  astilectron.PtrInt(600),
    })
    w.On(astilectron.EventNameWindowEventResize, func(e astilectron.Event) (deleteListener bool) {
        astilog.Info("Window resized")
        return
    })
	w.On(astilectron.EventNameWindowEventMessage, func(e astilectron.Event) (deleteListener bool) {
		var m string
		e.Message.Unmarshal(&m)
		astilog.Infof("Received message %s", m)
		return
	})
    w.Create()
    
    // Play with the window
    w.Resize(200, 200)
    time.Sleep(time.Second)
    w.Maximize()
    
    // Send a message to the server
    time.Sleep(time.Second)
    w.Send("What's up?")

	// Blocking pattern
	a.Wait()

# I want to see it in actions!

To make things clearer I've tried to split features in different [examples](https://github.com/asticode/go-astilectron/tree/master/examples).

To run any of the examples, run the following command:

    $ go run examples/<name of the example>/main.go -v
    
Here's a list of the examples:

- [1.basic_window](https://github.com/asticode/go-astilectron/tree/master/examples/1.basic_window/main.go) creates a basic window that displays a static .html file
- [2.basic_window_events](https://github.com/asticode/go-astilectron/tree/master/examples/2.basic_window_events/main.go) plays with basic window methods and shows you how to set up your own listeners
- [3.webserver_app](https://github.com/asticode/go-astilectron/tree/master/examples/3.webserver_app/main.go) sets up a basic webserver app
- [4.remote_messaging](https://github.com/asticode/go-astilectron/tree/master/examples/4.remote_messaging/main.go) sends a message to the webserver and listens for any response

# Features and roadmap

- [x] window basic methods (create, show, close, resize, minimize, maximize, ...)
- [x] window basic events (close, blur, focud, unresponsive, crashed, ...)
- [x] remote messaging (messages between GO and the JS in the webserver)
- [ ] menu methods
- [ ] menu events
- [ ] session methods
- [ ] session events
- [ ] window advanced options (add missing ones)
- [ ] window advanced methods (add missing ones)
- [ ] window advanced events (add missing ones)
- [ ] child windows

# Cheers to

[go-thrust](https://github.com/miketheprogrammer/go-thrust) which is awesome but unfortunately not maintained anymore. It inspired this project.