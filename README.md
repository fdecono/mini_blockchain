# Mini Blockchain Implementation

A simple blockchain implementation in Go for educational purposes. Features include automatic block generation, manual block creation, SHA-256 hashing, REST API endpoints, and live browser viewing with real-time updates.

## Attribution

This is an implementation of the educational blockchain tutorial by [nosequeldeebee](https://github.com/nosequeldeebee/blockchain-tutorial).

**Original Resources:**
- [GitHub Repository](https://github.com/nosequeldeebee/blockchain-tutorial)
- [Blog Post](https://mycoralhealth.medium.com/code-your-own-blockchain-in-less-than-200-lines-of-go-e296282bcffc)

## Features

- **Automatic block generation** - Creates new blocks every 2 seconds with random BPM values
- **Manual block creation** - POST endpoint to create blocks with custom BPM values
- **Live browser viewing** - View blocks in real-time at `/view` endpoint
- **Real-time updates** - Server-Sent Events (SSE) stream blockchain updates to the browser
- SHA-256 cryptographic hashing
- REST API (GET/POST endpoints)
- Longest chain consensus rule

## Usage

1. **Start the server**: Run `go run .` (ensure you have a `.env` file with `PORT=8080`)

2. **View blocks live in browser**: Navigate to `http://localhost:8080/view` to see blocks update in real-time

3. **Create blocks manually**: Send a POST request to `http://localhost:8080/` with JSON body:
   ```json
   {"BPM": 72}
   ```

4. **Get blockchain**: GET `http://localhost:8080/` to retrieve the full blockchain as JSON

## Example Output

```json
[
 {
  "Index": 0,
  "Timestamp": "2025-11-20 14:12:24.2924113 -0300 -03 m=+0.002047401",
  "BPM": 0,
  "Hash": "",
  "PrevHash": ""
 },
 {
  "Index": 1,
  "Timestamp": "2025-11-20 14:15:02.8027667 -0300 -03 m=+158.512402801",
  "BPM": 72,
  "Hash": "53e1031854f9984e7bd661d2df025447c9a236f64386a4249238e5d121992b1d",
  "PrevHash": ""
 }
]
```
## Purpose

This project is purely educational. It's designed to help understand blockchain concepts through hands-on implementation rather than focusing on production-ready code or advanced features.


