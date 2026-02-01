# Conv3D

Conv3D is a Go-based tool for the technical analysis, parsing, and serialization of SCW 3D model and animation data. It provides a bridge between specific binary schemas and structured formats like JSON, enabling inspection and version migration.

### Core Features

* **Binary Parsing:** Custom implementation for reading SCW geometry and animation data.
* **Bidirectional Serialization:** Supports decoding binary files to JSON and encoding JSON back into the original format.
* **Version Management:** Handles logic for various schema iterations (v0, v1, v2) including minor version delta-handling.
* **Optimization:** Implements logic to compute `Node.FrameFlags` to reduce output size by identifying identical properties across animation frames.

### Usage

**Build from source:**

```bash
go build -o conv3d main.go

```

**Commands:**

* **Decode:** `./conv3d --in-file=model.scw` (Outputs `model.scw.json`)
* **Encode:** `./conv3d --in-file=model.scw.json --out-file=model.scw`
* **Migrate:** `./conv3d --in-file=file.scw --scw2scw --out-version=2` (Updates/downgrades format versions)

### Implementation Objectives

* **Data Integrity:** Ensuring lossless transitions during the encoding/decoding process.
* **Schema Mapping:** Defining clear internal representations for cameras, materials, and meshes.
* **Compatibility:** Providing a reliable way to maintain assets across different environment specifications.

### Project Status

* [ ] Implement support for standard intermediate formats (GLB/GLTF).
* [ ] Expand unit test coverage for the parsing engine.
* [ ] Further document internal binary structures.

---

### Community & Contact

* **Discord:** `@s.b`
* **Community:** [discord.peterr.dev](https://www.google.com/search?q=https://discord.peterr.dev)

### Support

If this project served as a helpful reference for your own learning or research, feel free to leave a 🌟!
