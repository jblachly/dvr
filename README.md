# What is DVR?

dvr is a lightweight digital video recorder. It is web-based, easy to install, has few dependencies, and works with SilconDust's HDHomeRun CONNECT and HDHomeRun PRIME. It could be extended to work with any endpoint that streams video over IP.

DVR is ideal for installation on a VM (virtual machine)

# What does DVR do?

dvr is a digital video recorder. It supports scheduling of recordings, setting recurring recordings, an electronic program guide through (Schedules Direct)[http://www.schedulesdirect.org] *TODO* , all through a web-based interface.

# What does DVR not do?

dvr does NOT play video. As a consequence, it does not "pause live TV". 

# Why would I use a DVR that does not play video?

There are many fine frontends for playing video (Emby, Kodi, MythTV, Plex, to name a few). I anticipate that I will write plugins to export DVR's recorded videos to one or more media center type programs.

# Why did I write DVR?

Because MythTV was too fiddly to install, due to a legacy codebase and too many moving parts. To their credit, the MythTV team are modernising it by adding a web browser interface, cleaning up the GUI, cleaning up the codebase, etc. However, much of MythTV (e.g. capture cards) is irrelevant to me as I am just using the SiliconDust HDHomeRun.


# Installation
## Requirements

couchdb
ffmpeg 3.0

## Installing a binary package

## Installing from source

The go programming language version 1.7 or later is required to build (1.6 or 1.5 may work with GO15VENDOREXPERIMENT)

# Architecture outline

A single process, `dvr` acts as the interactive web site application server, REST interface to the web application, scheduler, and recording engine.

All persistent data are stored in a [CouchDB](https://couchdb.apache.org/) document database.

## Interactive web site

## REST interface

## Scheduler

## Recording engine

Video can be streamed directly to disk, or transcoded by ffmpeg

# Future directions

Post-process videos. This could include a chain of things to run in order. Examples include transcoding with ffmpeg (if for example the dvr is running on a machine too slow to transcode in real time) or commercial detection with mythcommflag.

