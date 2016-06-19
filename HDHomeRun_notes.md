#Documentation for HTTP interface:

[http://www.silicondust.com/hdhomerun/hdhomerun_http_development.pdf]
Last update: 20140407


`http://deviceip:80/discover.json`

`http://deviceip:80/lineup_status.json`

{"ScanInProgress":0,"ScanPossible":1,"Source":"Antenna","SourceList":["Antenna","Cable"]}


http://deviceip:80/lineup.json
Initially, this is an empty array:
[]

However, after running a channel scan, it shows an array of detected channels:
[
{
"GuideNumber": "4.1",
"GuideName": "WCMH-DT",
"HD": 1,
"URL": "http://192.168.1.111:5004/auto/v4.1"
},
{
"GuideNumber": "4.2",
"GuideName": "MeTV",
"URL": "http://192.168.1.111:5004/auto/v4.2"
},
{
"GuideNumber": "4.3",
"GuideName": "ION TV",
"URL": "http://192.168.1.111:5004/auto/v4.3"
},
{
"GuideNumber": "6.1",
"GuideName": "WSYX DT",
"HD": 1,
"URL": "http://192.168.1.111:5004/auto/v6.1"
},
{
"GuideNumber": "6.2",
"GuideName": "MYTV",
"URL": "http://192.168.1.111:5004/auto/v6.2"
},
{
"GuideNumber": "6.3",
"GuideName": "ANTENNA",
"URL": "http://192.168.1.111:5004/auto/v6.3"
}
]

These http URLs can be streamed directly in VLC

There is another URL that can be used to set a channel as a favorite:
method: POST
http://deviceip:80/lineup.post?favorite=value
where value = 
+4.1 (to set favorite)
x4.1 (to disable)
-4.1 (to unset favorite)


Once this is done, a channel lineup JSON can include a "Favorite" field:
[
{
"GuideNumber": "4.1",
"GuideName": "WCMH-DT",
"HD": 1,
"Favorite": 1,
"URL": "http://192.168.1.111:5004/auto/v4.1"
},
{
"GuideNumber": "4.2",
"GuideName": "MeTV",
"URL": "http://192.168.1.111:5004/auto/v4.2"
},
...
]


## Other attributes
Reading javascript source code for the device web page, there may also be attributes Subscribed and DRM



Interestingly, the hdhomerun_gui program instead of passing the http url directly to VLC will stream via RTP , redirect to a local port (5000) and then stream UDP to VLC on localhost

# Transcoding

As of firmware 20150826 adding a ?transcode= will, on the HDHomeRun CONNECT (where transcoding is not supported) stream the full resolution stream, rather than throw an error

http://deviceip:80/transcode.html may exist on higher end devices with transcode 
Native (none)      1920x1080   60fps    ~16Mb/s
Heavy              1920x1080   30fps    ~7Mb/s
Mobile             1280x540    30fps    ~3Mb/s
Internet720        1280x540    30fps    ~3Mb/s
Internet480        848x480     30fps    ~2Mb/s
Internet360        640x360     30fps    ~1.5Mb/s
Internet240        432x240     30fps    ~1Mb/s
(source: https://stuff.purdon.ca/?page_id=393)

# Tuner locking
source: [https://www.silicondust.com/forum/viewtopic.php?f=125&t=6710](https://www.silicondust.com/forum/viewtopic.php?f=125&t=6710)
