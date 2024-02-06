# Itzpapalotl

Itzpapalotl is a PalWorld server wrapper to make it on-demand.

## Motivation

PalWorld is a great game. We can build a PalWorld server ourselves. However, there are some problems with the server.

* Pal health
  * Pal will starve, or get sick if the server is running for a long time.
  * (It has been improved in the recent update probably)
* Electricity cost
  * Running a server 24/7 costs a lot of electricity if the server is running on a physical machine.
* Lagging game
  * PalWorld server consumes much memory if the server is running for a long time.
  * It causes the game to lag.

Itzpapalotl solves these problems.

## What Itzpapalotl does

Itzpapalotl does the following things.

* Start the server when the user wants to play the game.
  * Itzpapalotl does not start the PalWorld server until the user logs in to the server.
* Stop the server when all users log out.
  * Itzpapalotl stops the PalWorld server if no one has played it for 30 minutes.
* Restart the server when the server consumes much memory.
  * It avoids the game lagging.

With Itzpapalotl, the PalWorld server is running only when the user wants to play the game!

## Installation

### Pre-requirements

Itzpapalotl requires the following things.

* Linux
  * It only works on Linux.
* Enable RCON
  * You have to enable RCON in the PalWorld server because Itzpapalotl uses RCON to control the server.
  * Set `RCONEnabled=true` in your `PalWorldSettings.ini`.
* ps(1)

### Install Itzpapalotl

Install the executable file from the [latest release](https://github.com/pocke/itzpapalotl/releases/latest).

```bash
curl -L https://github.com/pocke/itzpapalotl/releases/latest/download/itzpapalotl_Linux_x86_64.tar.gz -o itzpapalotl_Linux_x86_64.tar.gz
tar -xvf itzpapalotl_Linux_x86_64.tar.gz
cp itzpapalotl /path/to/install/
```

### Usage and Configuration

Use `itzpapalotl` command instead of `./PalServer.sh`. The simplest example is:

```bash
itzpapalotl -admin-password ADMIN_PASSWORD -- /path/to/PalServer.sh
```

The `ADMIN_PASSWORD` is the password that is configured as `AdminPassword` in `PalWorldSettings.ini`.

You can specify arguments for `PalServer.sh` after `--`. For example:

```bash
itzpapalotl -admin-password ADMIN_PASSWORD -- /path/to/PalServer.sh -useperfthreads -NoAsyncLoadingThread -UseMultithreadForDS
```

`itzpapalotl` has some options. You can see them by `itzpapalotl --help`.

```bash
$ ./itzpapalotl --help
Usage: itzpapalotl [options] -- [palworld server command]
  -admin-password string
    	Admin password
  -memory-threshold int
    	Memory usage threshold (kb). If the process exceeds this threshold, it will be shut down. (default 10000000)
  -rcon-port int
    	RCON port (default 25575)
  -server-port int
    	PalWorld server port (default 8211)
```

## Why I named it "Itzpapalotl"

1. Open https://ja.wikipedia.org/wiki/%E7%A5%9E%E3%81%AE%E4%B8%80%E8%A6%A7
1. Enter `Ctrl + F` and search `pal`
1. I found it!
