# 0.1.7 - Add: README.md
- detailed instructions on how to use docs4context mcp server

# 0.1.6 - Add: cicd.js --full-release flag
- full-release does: --build --commit --tag --push --release
- previously likely releasing without build

# 0.1.5 - Fix: version comparison logic and updater process
- Fixed semantic version comparison in updater.go to properly detect newer versions
- Fixed CICD build process to correctly embed changelog version into binaries
- Previous releases 0.1.2, 0.1.3, and 0.1.4 had incorrect version comparison logic

# 0.1.4 - Fix: cicd.js to use correct version
- previously was broken in 0.1.2 and 0.1.3

# 0.1.3 - Add: install docs are for cursor or claude code
- using the two most common clients
- and to test the updater function

# 0.1.2 - Add: install.sh script to use mcp server
- compiling to 6 architectures, windows/macos/linux in arm/amd64
- install script created
- meant to be used in local directory as local mcp server
- added update ability via github releases
- updated cicd.js to use github releases

# 0.1.1 - Add: mvp for docs4context working
- setup with opencode.ai with opencode.json file
- testing 3x terminal agents: claude code, opencode and gemini cli
- docs4context app downloads up to 100 million token llms.txt from context7
- search and analyze capabilities in docs4context too
- addSaveContextDocumentTool(s)
	addSearchTitlesTool(s)
	addSearchContentTool(s)
	addGetTopicDetailsTool(s)
	addListRepositoriesTool(s)
	addAnalyzeKeywordsTool(s)

# 0.1.0 - Add: initial commit
- sample calculator-mcp go app working
