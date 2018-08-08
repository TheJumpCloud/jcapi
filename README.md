jcapi
=====

This Repository only supports JumpCloud's V1 API endpoints. For V2 Support please refer to [jcapi-go](https://github.com/TheJumpCloud/jcapi-go).

Please note that this V1 Go SDK is currently out of date with Jumpcloud's v1 API functionality. This is due to a bug in Swagger Code Gen and how we auto-deploy SDK updates. If you need an updated version of the V1 GO SDK please file an issue in this repository and we can evaluate the request on a per need basis.  If you need to only use our V2 set of endpoints, you can refer to our V2 SDK for Go that supports those endpoints as those are currently up to date with V2 API functionality.  

The scripts under the 'examples' folder in this repo are now deprecated.
Please refer to [the support repository](https://github.com/TheJumpCloud/support/tree/master/api-utils/JumpCloud_API_Go_Examples) for the maintained examples scripts.

Binaries for these scripts can be found in the [releases](https://github.com/TheJumpCloud/support/releases) for the support repository.

*JumpCloud's Go (golang) REST API SDK (BETA)*
Copyright (C) 2015 JumpCloud

This JumpCloud SDK is in beta form. The only available documentation comes in the form of the included examples and from jcapi_test.go (best stuff is in there).

Two things to keep in mind as you use this:
 * Use it at your own risk.
 * Please know that this SDK will change over time.

This API exposes several JumpCloud REST APIs:
 * System Users - (see https://github.com/TheJumpCloud/JumpCloudAPI#system-users)
 * Systems - (see https://github.com/TheJumpCloud/JumpCloudAPI#systems)
 * Identity Sources - not yet documented
 * Tags - (see https://github.com/TheJumpCloud/JumpCloudAPI#tags)
 * Authentication and Authorization - (see http://support.jumpcloud.com/knowledgebase/articles/455570)

