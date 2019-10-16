# departures
Show departure times for your Berlin public transport station

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
