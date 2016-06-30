## Password Expiry Watcher

This tool allows you as a system administrator to download a CSV of all of your users with each user's password expiration date. This allows monitoring of users with upcoming password expirations and for you the administrator to prevent these users from abruptly being cut out of their mission critical systems.

### To Install

##### I have the Go toolchain installed
> If you don't know what "Go" or "Golang" is we have provided pre-made binaries for your convenience. Installation instructions for these binaries is provided in the next section

- Clone the `jcapi` repository (it doesn't matter where)
	- `git clone https://github.com/TheJumpCloud/jcapi`
- Navigate to the `jcapi/examples/PasswordExpiryWatcher` directory
	- On Linux/macOS `cd jcapi/examples/PasswordExpiryWatcher`
- Run go install
	- `go install .`

This will install a binary called `passwordExpiryWatcher` to your `$GOBIN`

> If you'd like to be able to call this binary from an arbitrary directory make sure your `$GOBIN` is in your `$PATH` (linux) or `%PATH` (windows)

##### I do not have the Go toolchain installed
> Most users will want to follow these instructions

- Go to the "Releases" tab of this github project
	- https://github.com/TheJumpCloud/jcapi/releases
- We have provided `.zip` files for most operating systems
	- Note that macOS users will want to use the `darwin` binaries
- Download the zip file of your choice 
	- We always recommend downloading from the latest release if possible
- `PasswordExpiryWatcher_myos_myarch` is the relevant binary
- Some may find this is an unwieldy name to type every time. Feel free to rename to your liking.

> If you'd like to be able to call this binary from an arbitrary directory you can move it to `/usr/local/bin` on Linux

### To Run

This tool is very simple to run and only takes two (required) arguments:
- `key` is your JumpCloud API key
- `output` is where you'd like to put the resulting CSV file

For example:

`passwordExpiryWatcher -key=82105124f2979e28273d4e8dd32b2355c5012837 -output=password_expirations.csv`

### Example Output
![example csv output](https://cloud.githubusercontent.com/assets/712346/16349989/182347a6-3a19-11e6-969a-09a744a0f2ca.png)