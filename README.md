# conv3d

A Golang tool for loading and parsing .scw model and animation files used in Supercell games.

### Features:
- Convert SCW models to JSON.
- Encode JSON output back into SCW.
- Update or downgrade scw models

### Usage
- download the tool (check [Installation](#Installation)) or build it yourself (check [Building](#Building))
- `./conv3d --help` shows the tool arguments
- `./conv3d --in-file=sample.scw` loads your scw model and outputs an `output.scw.json` file
- `./conv3d --in-file=output.scw.json --out-file=output.scw` loads the json output and encode it back to SCW
- `./conv3d --in-file=file.scw --scw2scw` updates (or downgrades) a model scw version to another

### Installation
- install the binary for your os from [releases](https://github.com/PeterHackz/conv3d/releases)

**if your os is not there, you can [build](#Building) it yourself**

### Building
- you need to have golang installed (ofc...)
- clone this repo
- build it with `go build` or just run it with `go run main.go <args>`

### Extra Info
- in-file is a needed argument, and out-file is optional (by default it uses output.scw.json or output.scw)
- the json file should end with `.scw.json` and not only `.json`
- the format passed through multiple minor versions during scw v0, which cannot be detected by checking the model itself. for now support for v0.5 (minor version guessed for the tool usage) which is considered from v8 till v12
- you can pass this minor version with `--out-version=5` to specify this, which is the default for now as support for v0.0 (or other if any) is not added yet. major version will always be taken from the scw header
- you can set `out-version` with `scw2scw` to choose if to update it or downgrade it (by default it updates to scw v2)

## Disclaimer
This project is solely for educational purposes and must not be used for any malicious intent. Usage of this project is the sole responsibility of the user.

### Notice
- Encoding JSON back into SCW is still experimental and was tested on 3 models. if you find any bug in the project, open an issue or message me on discord.

### TODO:
- [ ] Test on SCW v1 models and animations
- [x] Test on SCW v2 models/animations
- [ ] Identify unknown field names
- [ ] Better code documentation
- [ ] Support parsing of DAE and GLB/GLTF models
- [ ] Develop a generalized model structure
- [ ] Implement conversion between different model types (e.g., SCW to GLB)
- [x] Compute Node.FrameFlags to optimize the output (when frames properties, ex: Rotation are identical for the first frame)

### Bugs
Encountered a bug? Please open an issue to help us resolve it.

### Contact
Have a question? (or just want to talk), DM me on discord: @s.b

### Star
Show some love with a ⭐️ because why not? ¯\\_(ツ)_/¯
