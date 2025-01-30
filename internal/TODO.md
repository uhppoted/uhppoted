# TODO

- [ ] github.released() check should use node-red-contrib version
- [ ] Implement npm.published()
- [ ] Move state to each individual project
      - [ ] Serialise projects list
      - [ ] Make nice json names for fields
- [ ] Sort project list by dependencies
- [ ] --debug that exits after each step/project

## All
   - [ ] Scrub _dist_ folder on release
   - [x] Version check uhppoted-lib-dotnet
   - [x] Version check uhppoted-lib-python
   - [ ] Remove release-notes.md in _bump_
   - [ ] Check CHANGELOG date before release/publish
   - [ ] Check README date before release/publish
         
## _uhppoted-lib-python_
   - [x] Release to github and wait for published
   - [x] Publish to testpy and wait for published
   - [x] Publish to pypi and wait for published
   - [ ] Publish to TestPyPI from github
   - [ ] Publish to PyPI from github
   - [ ] Use twine in venv rather than conda

## _uhppoted-nodejs_
   - [ ] Publish to npm from github

## _node-red-contrib-uhppoted_
   - [ ] Version in README should be node-red-version
   - [ ] Publish to npm from github
   - [ ] Remove old tarballs before publishing

## _uhppoted-app-home-assistant
   - [ ] add to release script

## _uhppoted_
   - [ ] Use sub-project Makefiles to generate uhppoted dist files
