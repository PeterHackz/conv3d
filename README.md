# conv3d

A Golang tool for loading and parsing .scw model and animation files used in Supercell games.

### Usage
- you need to have golang installed (ofc...)
- clone this repo
- build it with `go build` or just run it with `go run main.go <args>
- ```go run main.go --help``` shows the tool arguments
- ```go run main.go --in-file=sample.scw``` loads your scw model and outputs an `output.scw.json` file
- ```go run main.go --in-file=output.scw.json --out-file=output.scw``` loads the json output and encode it back to SCW

### Extra Info
- in-file is a needed argument, and out-file is optional (by default it uses output.scw.json or output.scw)
- the json file should end with `.scw.json` and not only `.json`

## Disclaimer
This project is solely for educational purposes and must not be used for any malicious intent. Usage of this project is the sole responsibility of the user.

### Features:
- [x] Convert SCW models to JSON.
- [x] Encode JSON output back into SCW.

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
- [ ] Compute Node.FrameFlags to optimize the output (when frames properties, ex: Rotation are identical for the first frame)

### Bugs
Encountered a bug? Please open an issue to help us resolve it.

### Contact
Have a question? (or just want to talk), DM me on discord: @s.b

### Star
Show some love with a ⭐️ because why not? ¯\\_(ツ)_/¯
