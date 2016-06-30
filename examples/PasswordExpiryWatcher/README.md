## Password Expiry Watcher

This tool allows you as a system administrator to download a CSV of all of your users with each user's password expiration date. This allows monitoring of users with upcoming password expirations and for you the administrator to prevent these users from abruptly being cut out of their mission critical systems.

### To Install

##### I do not have the Go toolchain installed
> Most users will want to follow these instructions

1. Go to the "Releases" tab of this github project
	- https://github.com/TheJumpCloud/jcapi/releases
2. We have provided `.zip` files for most operating systems
	- Note that macOS users will want to use the `darwin` binaries
3. Download the zip file of your choice 
	- We always recommend downloading from the latest release if possible
	- Most Windows users will want the `386` zip
4. Extract the files from the zip
	- Right click on the zip file and click `Extract All` on Windows
	- Double click on the zip file on macOS/Linux
5. `PasswordExpiryWatcher_myos_myarch` is the relevant binary
	- Some may find this is an unwieldy name to type every time. Feel free to rename to your liking, in fact if this is a tool you plan on running frequently we encourage this for ease of use.

> If you'd like to be able to call this binary from an arbitrary directory you can move it to `/usr/local/bin` on Linux and macOS
>> On macOS the keyboard shortcut `⌘ + Shift + G` with an open Finder window will allow you to directly type in the directory you'd like to access. To move the binary to `/usr/local/bin` with Finder simply use that shortcut and copy/paste `/usr/local/bin` into the dialog box. You can now drag and drop the binary onto the Finder window to easily make it accessible without having to navigate to a specific folder to run it. If you do this we also recommend renaming the binary before moving it to a more memorable name.

##### I have the Go toolchain installed
> If you don't know what "Go" or "Golang" is we have provided pre-made binaries for your convenience. Installation instructions for these binaries is provided in the previous section

1. Clone the `jcapi` repository (it doesn't matter where)
	- `git clone https://github.com/TheJumpCloud/jcapi`
2. Navigate to the `jcapi/examples/PasswordExpiryWatcher` directory
	- `cd jcapi/examples/PasswordExpiryWatcher`
3. Run go install
	- `go install .`

This will install a binary called `passwordExpiryWatcher` to your `$GOBIN`

> If you'd like to be able to call this binary from an arbitrary directory make sure your `$GOBIN` is in your `$PATH` (linux) or `%PATH` (windows)

### To Run

To run this tool you will need to use what is called the "command line". On Windows this is a program called "Power Shell", on macOS it is an app called "Terminal". 

> While running this tool requires no previous experience with either of those programs some might feel wary or nervous working with a tool they don't understand. The following instructions in this section should provide you with all you need to get up and running, but if you would like to learn about the how and why of the command line we highly recommend the excellent (and free!) [Command Line Crash Course by Zed Shaw](http://cli.learncodethehardway.org/). Zed even provides a direct email hotline for users that get stuck. However if you are stuck and would rather just get your users' password expiration dates ASAP please don't hesitate to contact JumpCloud support for assistance running this tool.

This tool only takes two (required) arguments:
- `key` is your JumpCloud API key
- `output` is where you'd like to put the resulting CSV file

For example:

macOS/Linux:

`./PasswordExpiryWatcher -key=82105124f2979e28273d4e8dd32b2355c5012837 -output=password_expirations.csv`

Windows:

`./PasswordExpiryWatcher.exe -key=82105124f2979e28273d4e8dd32b2355c5012837 -output=password_expirations.csv`

> If you have renamed your binary simply replace `PasswordExpiryWatcher` with the new name


> If you installed the binary with the Go toolchain and your `$GOBIN` is in your `$PATH` or `%PATH`, or if you moved the binary to `/usr/local/bin` you can run the above command at any time on your command line excluding the "./"

##### Windows Instructions
1. Open the program `Power Shell`
2. Using the `cd` (stands for "Change Directory") command navigate to where you downloaded your binaries
	- For example, if we downloaded and unzipped the binaries in our `Download` folder we just have to run: `cd Downloads\JumpCloudAPI_Examples_windows_386`. If we downloaded to our desktop the command will probably look something like: `cd Desktop\JumpCloudAPI_Examples_windows_386`
	- If you used the Go install instructions and your `$GOBIN` is in your `%PATH` you can skip step 2 and go right to 3

> To run a command simply type it into the Power Shell window and hit `Enter` or `Return` when finished 

3. Grab your API key from the JumpCloud Admin console
	- Click on your email on the top right hand corner to access the API Settings

4. Run the command
	- `./PasswordExpiryWatcher.exe -key=YOUR_API_KEY_GOES_HERE -output=CSV_FILE_OUTPUT_GOES_HERE`
	- If you would like your CSV file to go somewhere else besides the current directory make sure you include the _full path_ of the file
		- Good Example: `-output=C:\Users\MyUser\CSVFiles\password_expirations.csv` 
		- Bad Example `-output=..\CSVFiles\password_expirations.csv`

##### macOS/Linux Instructions
1. Open the app `Terminal` (this can be found in `Applications/Utilities`)
2. Using the `cd` (stands for "Change Directory") command navigate to where you downloaded your binaries
	- For example, if we downloaded and unzipped the binaries in our `Download` folder we just have to run: `cd Downloads/JumpCloudAPI_Examples_darwin_amd64`. If we downloaded to our desktop the command will probably look something like: `cd Desktop/JumpCloudAPI_Examples_darwin_amd64`
	- If you used the Go install instructions and your `$GOBIN` is in your `$PATH`, or if you manually moved the binary to `/usr/local/bin` you can skip step 2 and go right to 3

> To run a command simply type it into the Terminal window and hit `Enter` or `Return` when finished 

3. Grab your API key from the JumpCloud Admin console
	- Click on your email on the top right hand corner to access the API Settings

4. Run the command
	- `./PasswordExpiryWatcher -key=YOUR_API_KEY_GOES_HERE -output=CSV_FILE_OUTPUT_GOES_HERE`
	- If you would like your CSV file to go somewhere else besides the current directory make sure you include the _full path_ of the file
		- Good Example: `-output=/Users/MyUser/CSVFiles/password_expirations.csv` 
		- Bad Example `-output=../CSVFiles/password_expirations.csv`

### Example Output
![example csv output](https://cloud.githubusercontent.com/assets/712346/16349989/182347a6-3a19-11e6-969a-09a744a0f2ca.png)