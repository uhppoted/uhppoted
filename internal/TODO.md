# TODO

- [ ] github.released() check should use node-red-contrib version
- [ ] Scrub _dist_ folder on release
- [ ] Check CHANGELOG date before release/publish
- [ ] Check README date before release/publish
- [ ] Implement npm.published()
- [ ] Automatically commit all projects post-bump

- [ ] Move state to each individual project
      - [ ] Serialise projects list
      - [ ] Make nice json names for fields
- [ ] Sort project list by dependencies
- [ ] --debug that exits after each step/project

## _uhppoted-httpd_
- [ ] Check _package.json_ version
         
## _uhppoted-lib-python_
   - [ ] Publish to TestPyPI from github
   - [ ] Publish to PyPI from github
   - [ ] Use twine in venv rather than conda

## _uhppoted-nodejs_
   - [ ] Publish to npm from github

## _node-red-contrib-uhppoted_
   - [ ] Version in README should be node-red-version
   - [ ] Publish to npm from github
   - [ ] Remove old tarballs before publishing

## uhppoted-app-home-assistant
   - [ ] Check version in manifest.json
   - [ ] Bump manifest.json

## _uhppoted_
   - [ ] Use sub-project Makefiles to generate uhppoted dist files
