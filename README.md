# google-timeline-extractor
Google Maps Timeline shows an estimate of places that you may have been and routes that you may have taken based on your Location History.

This app will filter timeline records according to the filter.

## how to export information from google time line?
You can use [google takeout](https://takeout.google.com/) to export your information, it will export in json format.

## run
Run main.go to filter your records, these are the valids arguments:
- file: path to the json file
- name: name used to filter records (to use multiple just add | between words)
- output: (file/console): using file the filtered records will be present in a csv file called result, console will just print the records on screen.

Examples
```
go run main.go -file "/2019/2019_SEPTEMBER.json" -name "place1" -output console
```
```
go run main.go -file "/2019/2019_SEPTEMBER.json" -name "place1|plc1" -output console
```
```
go run main.go -file "/2019/2019_SEPTEMBER.json" -name "place1" -output file
```