# Tests

Here you can find two integration tests which test the normal flow of the system as explained in the main [README](https://github.com/edendoron/bit-framework) file.

In order to run the tests just `cd tests` from the root folder and then run `go test -short`.

In addition, there is a supported simulation function that generates infinite random reports which are ingested into the system and are used by the BIT handler to produce status reports that can be seen using the web UI.

If you wish to only run the simulation (or any other specific test), you can use `go test -run Simulation`.
