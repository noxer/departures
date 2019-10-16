# departures
Show departure times for your Berlin public transport station

## install
Currently we are only offering the source code. You need to have a current Go version installed.

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
$ departures -id="900000100003"
  M6 S Hackescher Markt                          10:19 (-1)
  S7 S Potsdam Hauptbahnhof                      10:20
  U2 U Ruhleben                                  10:20
  M4 Falkenberg                                  10:21
  U5 U Kaulsdorf-Nord                            10:21
  U2 S+U Pankow                                  10:21 (+3)
[and so on...]
```
