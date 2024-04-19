# protorefine

### BACKGROUND
Assume we have a big project, it has many sub projects, each sub project may use protobuf types defined in other sub projects,
the easiest way is put all proto files belong to different sub projects into same directory and just name it with different filename,
then compile all proto files into one package.
e,g:
```
proto/
    ...
    order.proto
    product.proto
    ...
```
will generate souce codes in `autogen/pb` directory
```
autogen/
    pb/
        ...
        order.pb.go
        product.pb.go
        ...
```
But the drawback is the final built binary file is too large, it will include all protobuf types defined in all proto files, which is not necessary.

### INTRODUCTION
Fortunately, this tool is developed to find all protobuf types referenced in source codes, then extract the protobuf type definitions from center proto files repository, at last generate new proto source file which only contains protobuf types we needed.

### INSTALL

```bash
go install github.com/rfancn/protorefine
```

### USAGE

1. it find all `<pb-import-path>` protobuf types referenced in source codes in `<project-dir>` 
2. it extract corresponding protobuf type definitions from center proto files repository in `<proto-dir>`
3. it generate a new proto file `<project>.proto` in `<output-dir>` which contains all protobuf types we needed  
4. it match and copy dependent proto files from `<proto-dir>` to `<output-dir>` use import rules, import rules read from config file specified by `--config` option, if no `--config` option specified, it uses default import rules

> if `--output-dir` not set, it will generate a random directory

### EXAMPLE

- find all `project/autogen/pb` protobuf types referenced in `/tmp/project`, extract corresponding protobuf type definitions in `/tmp/proto`, at last generate a new proto file `project.proto` in temporary directory

  ```
  protorefine --project-dir=/tmp/project --proto-dir=/tmp/proto --pb-import-path project/autogen/pb
  ```

- find all `project/autogen/pb` protobuf types referenced in `/tmp/project` exclude `/tmp/project/autogen` directory, extract corresponding protobuf type definitions in `/tmp/proto`, at last generate a new proto file `project.proto` in `/tmp/project/autogen/proto`

  ```
  protorefine --project-dir=/tmp/project --proto-dir=/tmp/proto --pb-import-path project/autogen/pb --output-dir=/tmp/project/autogen/proto --skip-dirs=autogen
  ```

- find all `project/autogen/pb` protobuf types referenced in `/tmp/project`, extract corresponding protobuf type definitions in `/tmp/proto`, generate a new proto file `project.proto` in `/tmp/project/autogen/proto`, copy dependent proto files match the rule defined in `/tmp/config.toml`

  ```
  protorefine --project-dir=/tmp/project --proto-dir=/tmp/proto --pb-import-path project/autogen/pb --output-dir=/tmp/project/autogen/proto --config /tmp/config.toml
  ```

### Config
- match: defines regular expression to match protobuf type definition content
- file:  defines the file name directly imported in generated proto file
- dependents: defines the dependent proto files the new proto file will import indirectly

Default import rules config
```toml
[import]
    [[import.rules]]
        match="gogoproto\\.+"
        file="gogo/protobuf/gogoproto/gogo.proto"
        dependents= ["google/protobuf/descriptor.proto"]
    [[import.rules]]
        match="hdsdk.protobuf\\.+"
        file="hdsdk/protobuf/sdk.proto"
    [[import.rules]]
        match="google.protobuf.Any"
        file="google/protobuf/any.proto"
```

> The `\` character collies with `toml` syntax, you need to escape it with `\\`
