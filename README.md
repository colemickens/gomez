# Features

## Current
* Not much, everything is broken
* Crappy movie library

## Planned

* Facebook login
* Better DB model

## Haha, right

* uPnp server
* uPnp controller
* MPD? mDNS or custom?

# Whatever

1. Acquire Go

2. Acquire a recent version of Ffmpeg. One that supports json output. (Unfortunately, this means building it yourself in several places. I could be convinced to generate builds and have something auto-grab them, depending on if anyone wants it and if Ffmpeg's [apparently nuanced] license allows it.)

3. Make sure GOPATH/bin is in your PATH.

3. `go get github.com/colemickens/gomez`

3. Generate and then edit a configuration file: `gomez -cmd writeconfig`.

4. `sudo gomez` or `sudo gomez -config ./config.json` will run the server.