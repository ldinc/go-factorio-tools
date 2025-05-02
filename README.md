# Description

Simple tool to build mod archive.

# Installing

`go install github.com/ldinc/go-factorio-tools`

# Examples

```sh
# build local package (expected mod source located at ./src)
goft -b

# build several mods in mod set (mod_base, mod_extra, mod_resource)
# for example project structure:
#  - modset |
#           | - mod_base |
#           |            | - info.json 
#           |            | - xxx
#           | - mod_resource |
#           |                | - info.json 
#           |                | - xxx
#           | - mod_extra |
#           |             | - info.json 
#           |             | - xxx
#           | - mod_experimental |
#           |                    | - info.json 
#           |                    | - xxx
goft -b mod_base mod_resource mod_extra

# build mod and generate zip file not in current working directory
goft -o /path/to/factorio/mods -b mod_base

```
