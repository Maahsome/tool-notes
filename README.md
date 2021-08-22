# tool-notes

The tldr projects lead me to want some extended notes on specific tools

## Features Roadmap

| Complete | Base Cmd   | Target        | Parameters      | Function                                |
| ---------| ---------- | ------------- | --------------- | --------------------------------------- |
|          | new        | repository    | -r,--repository | Create a new repository                 |
|          | new        | tool          | -t,--tool       | Create a new tool entry                 |
|          | new        | section       | -t, --tool      | Create a new section for a tool         |
|          |            |               | -s, --section   |                                         |
|          | new        | example       | -e, --example   | Create a new example for a tool/section |
|          |            |               | -s, --section   |                                         |
|          |            |               | -t, --tool      |                                         |
|          | delete     | tool          | -t,--tool       | Delete a tool entry                     |
|          | delete     | section       | -s, --section   |                                         |
|          |            |               | -t, --tool      |                                         |
|          | delete     | example       | -e, --example   |                                         |
|          |            |               | -s, --section   |                                         |
|          |            |               | -t, --tool      |                                         |
|          | edit       | tool          | -t,--tool       | Edit a tool entry                       |
|          | rename     | tool          | -t,--tool       | Rename a tool                           |
|   mvp    | show|view  | {toolname}    |                 | View information for a tool             |
|          | add        | repository    | -r,--repository | Fetch a remote repository               |
|          | update     |               |                 | Update all remote repositories          |
|          | search     | {search term} |                 | Search for matching notes               |
|          | list       | tool          | -r, --repository| List the tools in the repository        |
|          |            |               |                 |   - default to the primary              |
|          |            | section       | -r, --repository| List the sections for a tool            |
|          |            |               | -t, --tool      |                                         |
|          |            | example       | -r, --repository| List the examples for a tool/section    |
|          |            |               | -t, --tool      |                                         |
|          |            |               | -s, --section   |                                         |
|          |            | repository    |                 | List the repositories                   |

## Example Usage

```bash
tool-notes view {toolname}
tool-notes show {toolname}
  # this will display a prompt of Notes for that specific tool name in the default section
  # if the tool name isn't in the default section, search all other sections
tool-notes new section --section Personal_Projects
  # this will create a top level section named Personal_Projects
tool-notes new tool --section Personal_Projects --tool jq
  # this will create a new tool directory named jq under Personal_Projects
  # if --section is left out, prompt and set the default for the prompt
tool-notes new note --section Personal_Projects --tool jq --note Replace_Values.md
  # this will create a new note file named Replace_Values.md in the proper section/tool
  # if --section is left off prompt and set the default for the prompt.
  # if the --tool is left off, prompt for a selection

tool-notes --section Personal_Projects
  # this will display a list of Tool Names, then display a list of Notes
tool-notes
  # this will display a navigation of Section, Tool Names, Notes -> display
tool-notes {toolname} --section Personal_Projects
  # this will display a list of Notes for <toolname> under Personal_Projects
tool-notes set default-section {sectionname}
  # this will store the default searched section to the config file
```

## External Resources

An external render application will be required.  We will recommend installing
[mdcat](https://github.com/lunaryorn/mdcat) and using that Rust based tool.
Alternatively we will allow `TN_RENDERER=mdcat` as a valid way to identify your
desired rendering engine.

For the `brew install maahsome/tool-notes` we will require the `mdcat` install.

### Adding public repositories

Potentially have a similar github methodology like homebrew does, allowing you
to add a public repository by specifying just github username/repo name.  This
would come with the ability to have RO or RW rights to each repository.  Adding
new objects would be done in the "primary" repository, the one that was initially
created with `tool-notes add repository --repository cmaahs/md-notebooks`.

The multi-repository directory structure may present some annoyances in that the
selections may become quite large and hinder the quick reference tool we are
building.

```bash
tool-notes add repository --repository cmaahs
  # this would presume cmaahs/tool-notebooks
tool-notes add repoistory --repository cmaahs/md-notebooks
  # this would just add cmaahs/md-notebooks directly.
tool-notes update
  # this will perform a git fetch/pull on each of the defined repositories
```
