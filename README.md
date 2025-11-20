# Mini Blockchain Implementation

## Attribution

This is **not my original project**. This repository contains my implementation of the educational blockchain tutorial created by [nosequeldeebee](https://github.com/nosequeldeebee/blockchain-tutorial).

**Original Resources:**
- **GitHub Repository:** [blockchain-tutorial](https://github.com/nosequeldeebee/blockchain-tutorial)
- **Blog Post:** [Code your own blockchain in less than 200 lines of Go](https://mycoralhealth.medium.com/code-your-own-blockchain-in-less-than-200-lines-of-go-e296282bcffc)

## Learning Objectives

This implementation serves as a hands-on learning exercise to understand blockchain fundamentals. The primary goals are:

### Core Learning Goals

1. **Create your own blockchain**
   - Understand the basic structure and components of a blockchain
   - Implement a simple blockchain from scratch in Go

2. **Understand how hashing works in maintaining integrity of the blockchain**
   - Learn how cryptographic hashing (SHA-256) ensures data integrity
   - See how each block's hash depends on its content and the previous block's hash
   - Understand how this creates an immutable chain

3. **See how new blocks get added**
   - Observe the block generation process
   - Understand the relationship between consecutive blocks
   - Learn how blocks are validated before being added to the chain

4. **See how tiebreakers get resolved when multiple nodes generate blocks**
   - Understand the longest chain rule
   - Learn how blockchain networks handle conflicts
   - See the `replaceChain` function in action

5. **View your blockchain in a web browser**
   - Interact with the blockchain through a REST API
   - See the entire chain displayed as JSON
   - Understand how blockchain data can be exposed via HTTP

6. **Write new blocks**
   - Learn how to add new data to the blockchain
   - Understand the POST request flow
   - See how new blocks are created and validated

7. **Get a foundational understanding of the blockchain**
   - Build a solid conceptual foundation for blockchain technology
   - Understand the core principles that apply to all blockchains
   - Prepare for more advanced blockchain concepts

## Example

```bash
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
 },
 {
  "Index": 2,
  "Timestamp": "2025-11-20 14:15:25.8532583 -0300 -03 m=+181.562894401",
  "BPM": 76,
  "Hash": "621308d4de4e90f02a3006297baa45bfcd593fdc89fd15ee3bb2333303a5623c",
  "PrevHash": "53e1031854f9984e7bd661d2df025447c9a236f64386a4249238e5d121992b1d"
 }
]
```

## Purpose

This project is purely educational. It's designed to help understand blockchain concepts through hands-on implementation rather than focusing on production-ready code or advanced features.

