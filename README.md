# QSYS CORS Bridge
This is a very simple HTTP server that gets around Q-SYS Cores not setting CORS attributes on their
`/designs/current_design/ExternalControls.xml` document, meaning that webpages or app cannot access it. This server gets
that document on demand, reformats it as JSON, and serves it on `/designData.json` with CORS enabled. Additionally it
will serve any static files in the assets directory.

# Usage
Put any HTML, JS, CSS, etc files you want served in the assets directory and launch the server on the command line
with the `address` flag to specify the target Q-SYS Core address. The internal HTTP server defaults to serving on port
8080, but that can be overridden with the `port` flag.

## Examples
### Default port (8080)
```
./ucibridge.exe -address "192.168.201.165"
2020/11/26 15:26:06 Starting HTTP listener on port 8080
2020/11/26 15:26:06 Proxying requests to Core at 192.168.201.165
2020/11/26 15:26:10 Proxied request for design 'Apartment Core' (RL9OkBTBxlXG)
```

### Custom port
```
./ucibridge.exe -address "192.168.201.165" -port 8082
2020/11/26 15:28:04 Starting HTTP listener on port 8082
2020/11/26 15:28:04 Proxying requests to Core at 192.168.201.165
2020/11/26 15:28:43 Proxied request for design 'Apartment Core' (RL9OkBTBxlXG)
```

## /designData.json Contents
```json5
{
	"design": {
		"Snapshots": [
			// Not sure what these do currently
		],
		// Array containing every External Control in the Design
		"ExternalControls": [
			{
			    // External name of the External Control (the one you set in Designer)
				"Id": "BalconyTrim",
				// Internal name of the External Control
				"ControlId": "fader_2",
				// Pin name
				"ControlName": "Fader 2",
				// Unique internal reference to the component this control is on
				"ComponentId": "cN%DYZ4I=7S5mG;ku8,A",
				// Pretty name for the component this control is on
				"ComponentName": "Zone Controls : Custom Controls",
				// Custom label on the component this control is on
				"ComponentLabel": "Zone Controls",
				// The control name used but the UCI Viewer when using this control in UIs (USING THIS IS NOT DOCUMENTED OR SUPPORTED AS FAR AS I KNOW)
				"MappingName": "cN%DYZ4I=7S5mG;ku8,A_fader_2",
				"Type": "Float",
				"Mode": "RW",
				"MinimumValue": "-100",
				"MaximumValue": "10"
			},
            //... more ...
		],
		// File name of the running Design
		"DesignName": "Apartment Core",
		// Unique ID, changes every time the file is loaded
		"CompileGUID": "RL9OkBTBxlXG"
	},
	// Address of the target core
	"hostname": "192.168.201.165"
}

```