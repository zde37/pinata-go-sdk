# Pinata Go SDK (Unofficial)

**Disclaimer:** This is not the official Pinata Go SDK. This is a community version built by a user interested in Pinata's services(Me).

The Pinata Go SDK is a comprehensive Go library for interacting with the Pinata API, a popular service for pinning content to IPFS (InterPlanetary File System). This SDK provides an easy-to-use interface for developers to integrate Pinata's functionality into their Go applications.

## Purpose

The main purpose of this SDK is to simplify the process of interacting with Pinata's API by providing a set of Go functions and types that abstract away the complexities of making HTTP requests and handling responses. It allows developers to easily pin files and JSON data to IPFS, manage pins, update metadata, and perform other Pinata-related operations.

## Scope

The SDK covers a wide range of Pinata API functionalities, including:

- Authentication
- Pinning files and JSON to IPFS
- Listing pinned files
- Updating file metadata
- Deleting pins
- Querying pins by CID

## Files and Their Purposes

| File | Purpose |
| --- | --- |
| `pinata/auth.go` | Contains the `Auth` struct and related functions for handling authentication with the Pinata API. Supports both API key/secret and JWT token authentication methods. |
| `pinata/client.go` | Defines the main `Client` struct, which is the primary interface for interacting with the Pinata API. Includes the `NewClient` function for creating a new client instance and the `NewRequest` method for initiating API requests. |
| `pinata/pinning.go` | Contains core functionality for pinning operations. Includes structs and methods for pinning files to IPFS, pinning JSON to IPFS, listing pinned files, updating file metadata, deleting pins, and querying pins by CID. |
| `pinata/request_builder.go` | Implements the `requestBuilder` struct and its methods. Handles the construction and execution of HTTP requests to the Pinata API. |
| `pinata/group.go` | Implements functionality for managing Pinata groups, including creating, retrieving, updating, and deleting groups, as well as adding and removing CIDs from groups. |
| `pinata/signature.go` | Provides methods for adding, retrieving, and removing CID signatures in the Pinata API. |
| `pinata/user.go` | Implements user-related functionality, including generating and managing API keys, listing API keys, and revoking API keys. |


## Usage

Follow these steps to use the Pinata Go SDK in your project:

1. Fetch the Pinata Go SDK by running:
```go
go get github.com/zde37/pinata-go-sdk
```

2. Import it as follows:
```go
import "github.com/zde37/pinata-go-sdk/pinata"
```

3. Create a new client with your Pinata API credentials:
```go 
auth := pinata.NewAuthWithJWT("your-jwt-token")
client := pinata.NewClient(auth)
```

4. Use the client to interact with the Pinata API:
```go 
response, err := client.PinFileToIPFS("path/to/file.txt", nil)
if err != nil {
    // handle error
}
fmt.Printf("File pinned successfully. IPFS hash: %s\n", response.IpfsHash)
```

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Future Features

We're always looking to improve and expand the SDK. Here are some features we're considering for future releases:

- Implementation of gateway functionality
- Adding comprehensive test suite
- Support for additional Pinata API endpoints as they become available
- Improved error handling and logging
- Performance optimizations

If you have any suggestions for future features, please open an issue or contribute to the project!



