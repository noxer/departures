# departures
Show departure times for your Berlin public transport station

![screenshot of a departure table](https://raw.githubusercontent.com/noxer/departures/master/images/screen0.png)

## install
Currently we are offering binary (pre-)releases and the source code. The binary releases can be found on the [releases](https://github.com/noxer/departures/releases) page. 
To install from source you need to have a current Go version installed.

```bash
go get -u github.com/noxer/departures
```
Now you should have the `departures` binary installed in your `$GOPATH/bin` directory. You can call it from there or add the directory to your `$PATH`.

## usage
First you need to find out the ID of your station. To do this run the tool with the `-search` parameter.
```bash
~$ departures -search="Alexanderplatz"
Found 5 station(s):
  900000100003 - S+U Alexanderplatz
  900000100024 - S+U Alexanderplatz/Dircksenstr.
  900000100026 - S+U Alexanderplatz/Gontardstr.
  900000100005 - U Alexanderplatz [Tram]
  900000100031 - S+U Alexanderplatz/Memhardstr.
```

This should help you identify the station you want to look at. Now you can request the timetable for Alexanderplatz.

```bash
~$ departures -id="900000100003"
  M6 S Hackescher Markt                          10:19 (-1)
  S7 S Potsdam Hauptbahnhof                      10:20
  U2 U Ruhleben                                  10:20
  M4 Falkenberg                                  10:21
  U5 U Kaulsdorf-Nord                            10:21
  U2 S+U Pankow                                  10:21 (+3)
[and so on...]
```

You can limit the lines and directions shown. Multiple values must be separated by a comma.

```bash
~$ departures -id="900000100003" -filter-line="M4" -filter-destination="S Hackescher Markt"
M4 S Hackescher Markt 10:57 (-1)
M4 S Hackescher Markt 11:02 (-1)
M4 S Hackescher Markt 11:07 (-1)
M4 S Hackescher Markt 11:13
```

You can limit the width of the output to make it fit your terminal or for use in [wtfutil](https://github.com/wtfutil/wtf)*.

(* wtfutil sets the `WTF_WIDGET_WIDTH` environment variable which is automatically recognized by departures)

```bash
~$ departures -id="900000100003" -width=20
  S5 S We 10:58
 100 S+U  10:58
 RE1 Fran 10:59 (+1)
```

You can filter only connections that allow you to take a bike by adding the '-bicycle' argument.

```bash
~$ departures -id 900000029305
M32 S+U Rathaus Spandau    20:02
M32 Staaken, Heidebergplan 20:10 (-1)
M32 S+U Rathaus Spandau    20:20
RE4 Rathenow, Bhf          20:25 (+11)
M32 Staaken, Heidebergplan 20:31
M32 S+U Rathaus Spandau    20:40
RE4 Ludwigsfelde, Bhf      20:44
M32 Staaken, Heidebergplan 20:51
M32 S+U Rathaus Spandau    21:00

~$ departures -id 900000029305 -bicycle
RE4 Rathenow, Bhf     20:25 (+11)
RE4 Ludwigsfelde, Bhf 20:44
```

You can show additional informations like warning bicycle conveyance etc. by adding the '-verbose' argument

```bash
~$ departures -id 900000029305
M32 S+U Rathaus Spandau    20:20
RE4 Rathenow, Bhf          20:23 (+9)

~$ departures -id 900000029305 -verbose
M32 S+U Rathaus Spandau    20:20
    Operator : Berliner Verkehrsbetriebe
    Type     : bus
    Hint     : barrier-free
RE4 Rathenow, Bhf          20:23 (+9)
    Operator : ODEG Ostdeutsche Eisenbahn GmbH
    Type     : regional
    Hint     : barrier-free
    Hint     : Bicycle conveyance
    Hint     : Fahrradmitnahme leicht gemacht: www.vbb.de/radimregio
 ```

## wtfutil
This utility was originally created for use in wtfutil. You can use the following config snippet to get started.

```yml
    departures:
      args: ["-id=900000100003", "-force-color", "-retries=100", "-retry-pause=5s"]
      cmd: "departures"
      enabled: true
      position:
        top: 0
        left: 0
        height: 1
        width: 1
      refreshInterval: 60
      type: cmdrunner
      title: Departures
```

## attribution
I'm using https://2.bvg.transport.rest to request the current timetable data. Thanks to [derhuerst](https://github.com/derhuerst).
