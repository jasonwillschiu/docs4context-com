// Usage: 2 modes (dev and build server) and deployment with git
import { $ } from "bun";
import path from "path";
import fs from "fs/promises";
import { parseArgs } from "util";

// --- Color Helpers ---
const reset = "\x1b[0m";
const colorFormat = "ansi";
const green = (text) => `${Bun.color("#2ecc71", colorFormat)}${text}${reset}`;
const red = (text) => `${Bun.color("#e74c3c", colorFormat)}${text}${reset}`;
const cyan = (text) => `${Bun.color("#3498db", colorFormat)}${text}${reset}`;
const yellow = (text) => `${Bun.color("#f1c40f", colorFormat)}${text}${reset}`;
const blue = cyan;
const bold = (text) => `\x1b[1m${text}${reset}`;

// --- Integrated Bun Spinner ---
// (Keep the createBunSpinner function exactly as it was)
function createBunSpinner(initialText = "", opts = {}) {
  const stream = opts.stream || process.stderr;
  const isTTY = !!stream.isTTY && process.env.TERM !== "dumb" && !process.env.CI;
  const frames = Array.isArray(opts.frames) && opts.frames.length
    ? opts.frames
    : isTTY
      ? ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]
      : ["-"];
  const interval = opts.interval || 80;

  let text = initialText;
  let colorName = opts.color || "yellow";
  let colorFn = mapColor(colorName);
  let idx = 0;
  let timer = null;

  function mapColor(name) {
    switch ((name || "").toLowerCase()) {
      case "red": return red;
      case "green": return green;
      case "yellow": return yellow;
      case "blue": return cyan;
      case "cyan": return cyan;
      default: return (s) => s;
    }
  }

  function render() {
    const frame = frames[idx];
    if (isTTY) {
      // move to start of line, hide cursor, write frame + text
      stream.write("\r\x1b[?25l" + colorFn(frame) + " " + text);
    } else {
      // non-TTY: just print a line once
      stream.write(frame + " " + text + "\n");
    }
  }

  const spinner = {
    start(opts2 = {}) {
      const o = typeof opts2 === "string" ? { text: opts2 } : opts2;
      if (o.text != null) text = o.text;
      if (o.color) { colorName = o.color; colorFn = mapColor(colorName); }
      if (timer) clearInterval(timer);
      render();
      timer = setInterval(() => {
        idx = (idx + 1) % frames.length;
        render();
      }, interval);
      return spinner;
    },

    update(opts2 = {}) {
      const o = typeof opts2 === "string" ? { text: opts2 } : opts2;
      if (o.text != null) text = o.text;
      if (o.color) { colorName = o.color; colorFn = mapColor(colorName); }
      render();
      return spinner;
    },

    stop(opts2 = {}) {
      if (timer) { clearInterval(timer); timer = null; }
      const o = typeof opts2 === "string" ? { text: opts2 } : opts2;
      if (o.text != null) text = o.text;
      if (o.color) colorFn = mapColor(o.color);

      // pick a final mark
      const rawMark = o.mark != null ? o.mark : frames[idx];
      // assume o.mark is already colored if you passed green("✔"), etc.
      const mark = rawMark + " ";
      const out = mark + text + (isTTY ? "\n\x1b[?25h" : "\n");
      // overwrite the spinner line
      stream.write("\r" + out);
      return spinner;
    },

    success(opts2 = {}) {
      const o = typeof opts2 === "string" ? { text: opts2 } : opts2;
      return spinner.stop({
        mark: green("✔"),
        text: o.text != null ? o.text : text,
        color: "green"
      });
    },

    error(opts2 = {}) {
      const o = typeof opts2 === "string" ? { text: opts2 } : opts2;
      return spinner.stop({
        mark: red("✖"),
        text: o.text != null ? o.text : text,
        color: "red"
      });
    },

    warn(opts2 = {}) {
      const o = typeof opts2 === "string" ? { text: opts2 } : opts2;
      return spinner.stop({
        mark: yellow("⚠"),
        text: o.text != null ? o.text : text,
        color: "yellow"
      });
    },

    info(opts2 = {}) {
      const o = typeof opts2 === "string" ? { text: opts2 } : opts2;
      return spinner.stop({
        mark: cyan("ℹ"),
        text: o.text != null ? o.text : text,
        color: "cyan"
      });
    },

    isSpinning() {
      return timer != null;
    }
  };

  return spinner;
}
// --- End Integrated Bun Spinner ---


// --- Argument Parsing ---
const { values } = parseArgs({
  args: Bun.argv.slice(2), // Exclude 'bun' and script path
  options: {
    mode: { type: 'string' },
    build: { type: 'boolean', default: false },
    commit: { type: 'boolean', default: false },
    tag: { type: 'boolean', default: false },
    push: { type: 'boolean', default: false },
    release: { type: 'boolean', default: false },
  },
  strict: false,
  allowPositionals: false,
});

const mode = values.mode;
const shouldBuild = values.build;
const shouldCommit = values.commit;
const shouldTag = values.tag;
const shouldPush = values.push;
const shouldRelease = values.release;

// --- Helper Functions ---

/** Builds Go binaries for all supported platforms with version info */
async function buildCrossPlatform(version = null) {
  const platforms = [
    { os: 'darwin', arch: 'amd64', name: 'docs4context-com-darwin-amd64' },
    { os: 'darwin', arch: 'arm64', name: 'docs4context-com-darwin-arm64' },
    { os: 'linux', arch: 'amd64', name: 'docs4context-com-linux-amd64' },
    { os: 'linux', arch: 'arm64', name: 'docs4context-com-linux-arm64' },
    { os: 'windows', arch: 'amd64', name: 'docs4context-com-windows-amd64.exe' },
    { os: 'windows', arch: 'arm64', name: 'docs4context-com-windows-arm64.exe' }
  ];

  const mainSpinner = createBunSpinner(`🚀 Starting cross-platform Go builds...`).start();
  
  try {
    // Clean bin directory
    mainSpinner.update({ text: '🧹 Cleaning bin directory...' });
    await $`rm -rf bin/*`.nothrow();
    await $`mkdir -p bin`.throws(true);

    // Get build info
    const buildDate = new Date().toISOString();
    const gitCommitResult = await $`git rev-parse --short HEAD`.nothrow();
    const gitCommit = gitCommitResult.exitCode === 0 ? gitCommitResult.stdout.toString().trim() : 'unknown';
    const buildVersion = version || 'dev';

    // Build for each platform
    for (const platform of platforms) {
      mainSpinner.update({ text: `🔨 Building ${platform.name}...` });
      
      // Create ldflags for build info
      const ldflags = `-X main.Version=${buildVersion} -X main.BuildDate=${buildDate} -X main.GitCommit=${gitCommit}`;
      
      const buildResult = await $`go build -ldflags ${ldflags} -o bin/${platform.name} .`
        .env({ 
          ...process.env, 
          GOOS: platform.os, 
          GOARCH: platform.arch,
          CGO_ENABLED: '0'
        })
        .nothrow();
        
      if (buildResult.exitCode !== 0) {
        mainSpinner.error({ text: red(`❌ Failed to build ${platform.name}`) });
        console.error(red(buildResult.stderr.toString()));
        throw new Error(`Build failed for ${platform.name}`);
      }
    }
    
    mainSpinner.success({ text: green(`✅ Built ${platforms.length} binaries successfully`) });
    
    // List built files
    console.log(cyan('📦 Built binaries:'));
    const lsResult = await $`ls -la bin/`.nothrow();
    if (lsResult.exitCode === 0) {
      console.log(lsResult.stdout.toString());
    }
    
  } catch (error) {
    if (mainSpinner.isSpinning()) {
      mainSpinner.error({ text: red('❌ Cross-platform build failed') });
    }
    throw error;
  }
}

/** Parses the latest entry from changelog.md. */
async function parseLatestChangelogEntry() {
  const changelogPath = path.join(import.meta.dir, 'changelog.md');
  console.log(cyan(`ℹ️ Reading changelog: ${changelogPath}`));
  try {
    const content = await fs.readFile(changelogPath, 'utf-8');
    const lines = content.split('\n');
    const headerStartRegex = /^#\s*([0-9]+\.[0-9]+\.[0-9]+[a-z]?)\s*-/i;
    const firstEntryStartIndex = lines.findIndex(line => headerStartRegex.test(line));
    if (firstEntryStartIndex === -1) throw new Error("Could not find any entry starting with '# <version> -' in changelog.md.");
    let nextEntryStartIndex = lines.findIndex((line, index) => index > firstEntryStartIndex && headerStartRegex.test(line));
    if (nextEntryStartIndex === -1) nextEntryStartIndex = lines.length;
    const entryLines = lines.slice(firstEntryStartIndex, nextEntryStartIndex);
    if (entryLines.length === 0) throw new Error("Detected an empty entry block in changelog.md.");
    const headerLine = entryLines[0];
    const headerMatch = headerLine.match(/^#\s*([0-9]+\.[0-9]+\.[0-9]+[a-z]?)\s*-\s*(.*)/i);
    if (!headerMatch || headerMatch.length < 3) throw new Error(`Could not parse header (line ${firstEntryStartIndex + 1}): "${headerLine}". Expected: '# <version> - <summary>'`);
    const version = headerMatch[1].trim();
    const summary = headerMatch[2].trim();
    const descriptionPoints = entryLines.slice(1).map(line => line.trim()).filter(line => line.startsWith('-')).map(line => `* ${line.substring(1).trim()}`).filter(line => line.length > 2);
    const description = descriptionPoints.join('\n');
    console.log(green(`✅ Parsed changelog: v${version} - ${summary}`));
    return { version, summary, description };
  } catch (error) {
    console.error(red(`❌ Error reading/parsing changelog.md: ${error.message}`));
    throw error;
  }
}

/** Checks if there are uncommitted changes. Throws error if clean. */
async function checkGitStatus() {
  console.log(cyan("ℹ️ Checking Git status..."));
  const { stdout, exitCode } = await $`git status --porcelain`.nothrow();
  if (exitCode !== 0) throw new Error(red("Failed to check Git status."));
  if (stdout.toString().trim() === '') throw new Error(red("❌ No changes detected. Working directory clean. Nothing to commit."));
  console.log(green("✅ Git status check passed: Changes detected."));
}

/** Checks if a Git tag already exists. Throws error if it does. */
async function checkGitTagExists(version) {
  const tagToCheck = `v${version}`;
  console.log(cyan(`ℹ️ Checking if Git tag ${tagToCheck} exists...`));
  const { stdout, exitCode } = await $`git tag -l ${tagToCheck}`.nothrow();
  if (exitCode === 0 && stdout.toString().trim() === tagToCheck) {
    throw new Error(red(`❌ Git tag '${tagToCheck}' already exists. Update changelog.md.`));
  }
  console.log(green(`✅ Git tag check passed: Tag '${tagToCheck}' does not exist yet.`));
}

/** Stages all changes. */
async function gitAdd() {
  const spinner = createBunSpinner(`ℹ️ Staging changes (${bold('git add .')})...`).start();
  try {
    await $`git add .`.quiet().throws(true);
    spinner.success({ text: green('✅ Changes staged.') });
  } catch (error) {
    spinner.error({ text: red('❌ Failed to stage changes.') });
    console.error(red(error.stderr?.toString() || error.message));
    throw error;
  }
}

/** Commits staged changes with a formatted message. */
async function gitCommit(summary, description) {
  const spinner = createBunSpinner(`ℹ️ Committing changes...`).start();
  const commitMessage = `${summary}\n\n${description}`;
  try {
    const proc = Bun.spawnSync(['git', 'commit', '--file=-'], { stdin: Buffer.from(commitMessage) });
    if (proc.exitCode !== 0) throw new Error(proc.stderr.toString() || "Git commit failed");
    spinner.success({ text: green('✅ Changes committed.') });
  } catch (error) {
    spinner.error({ text: red('❌ Git commit failed.') });
    console.error(red(error.message));
    throw error;
  }
}

/** Creates an annotated Git tag. */
async function gitTag(version, summary) {
  const tag = `v${version}`;
  const spinner = createBunSpinner(`ℹ️ Creating annotated Git tag ${bold(`'${tag}'`)}...`).start();
  try {
    await $`git tag -a ${tag} -m ${summary}`.quiet().throws(true);
    spinner.success({ text: green(`✅ Tag '${tag}' created.`) });
  } catch (error) {
    spinner.error({ text: red(`❌ Failed to create tag '${tag}'.`) });
    console.error(red(error.stderr?.toString() || error.message));
    throw error;
  }
}

/** Pushes commits and tags to the remote repository. */
async function gitPush() {
  const spinner = createBunSpinner(`ℹ️ Preparing to push...`).start();
  try {
    // Push commits
    spinner.update({ text: `ℹ️ Pushing commits (logs below)...`, color: 'cyan' });
    const pushCommitsResult = await $`git push`.nothrow();
    if (pushCommitsResult.exitCode !== 0) {
      spinner.error({ text: red(`❌ Failed to push commits.`) });
      console.error(red(pushCommitsResult.stderr.toString() || "git push failed"));
      throw new Error("Failed to push commits.");
    }
    spinner.success({ text: green(`✅ Commits pushed.`) });

    // Push tags
    spinner.start({ text: `ℹ️ Pushing tags (logs below)...`, color: 'cyan' });
    const pushTagsResult = await $`git push --tags`.nothrow();
    if (pushTagsResult.exitCode !== 0) {
      spinner.error({ text: red(`❌ Failed to push tags.`) });
      console.error(red(pushTagsResult.stderr.toString() || "git push --tags failed"));
      throw new Error("Failed to push tags.");
    }
    spinner.success({ text: green(`✅ Tags pushed.`) });

  } catch (error) {
    if (spinner.isSpinning()) {
      spinner.error({ text: red(`❌ Push operation failed.`) });
    }
    throw error;
  }
}

/** Builds the main Go binary for current platform */
async function buildLocal() {
  const spinner = createBunSpinner(`🚀 Building docs4context-com for current platform...`).start();
  
  try {
    // Get build info
    const buildDate = new Date().toISOString();
    const gitCommitResult = await $`git rev-parse --short HEAD`.nothrow();
    const gitCommit = gitCommitResult.exitCode === 0 ? gitCommitResult.stdout.toString().trim() : 'unknown';
    
    // Create ldflags for build info
    const ldflags = `-X main.Version=dev -X main.BuildDate=${buildDate} -X main.GitCommit=${gitCommit}`;
    
    const buildResult = await $`go build -ldflags ${ldflags} -o docs4context-com .`
      .env({ ...process.env, CGO_ENABLED: '0' })
      .nothrow();
      
    if (buildResult.exitCode !== 0) {
      spinner.error({ text: red('❌ Local build failed') });
      console.error(red(buildResult.stderr.toString()));
      throw new Error('Local build failed');
    }
    
    spinner.success({ text: green('✅ Local build completed') });
    
  } catch (error) {
    if (spinner.isSpinning()) {
      spinner.error({ text: red('❌ Local build process failed') });
    }
    throw error;
  }
}

/** Creates a GitHub release and uploads binaries */
async function createGitHubRelease(version, summary, description) {
  const spinner = createBunSpinner(`🚀 Creating GitHub release v${version}...`).start();
  
  try {
    // Check if gh CLI is available
    const ghCheckResult = await $`which gh`.nothrow();
    if (ghCheckResult.exitCode !== 0) {
      throw new Error('GitHub CLI (gh) is not installed. Please install it first.');
    }
    
    // Create the release
    spinner.update({ text: `📦 Creating release v${version}...` });
    
    const releaseBody = `${summary}\n\n## Changes\n${description}`;
    const createResult = await $`gh release create v${version} --title "v${version}" --notes ${releaseBody}`.nothrow();
    
    if (createResult.exitCode !== 0) {
      throw new Error(`Failed to create GitHub release: ${createResult.stderr.toString()}`);
    }
    
    // Upload binaries
    spinner.update({ text: `📤 Uploading binaries...` });
    
    const uploadResult = await $`gh release upload v${version} bin/*`.nothrow();
    
    if (uploadResult.exitCode !== 0) {
      throw new Error(`Failed to upload binaries: ${uploadResult.stderr.toString()}`);
    }
    
    spinner.success({ text: green(`✅ GitHub release v${version} created successfully`) });
    
    // Show release URL
    const releaseUrl = `https://github.com/jasonwillschiu/docs4context-com/releases/tag/v${version}`;
    console.log(cyan(`🔗 Release URL: ${releaseUrl}`));
    
  } catch (error) {
    if (spinner.isSpinning()) {
      spinner.error({ text: red('❌ GitHub release creation failed') });
    }
    throw error;
  }
}




// --- Main Logic ---

// Dev Mode
if (mode === 'dev') {
  console.log(cyan("🚀 Building frontend and starting Go backend..."));

  const buildSpinner = createBunSpinner(`🚀 Building Astro frontend...`).start();
  try {
    // Build frontend
    buildSpinner.update({ text: `🚀 Building frontend in ${bold('frontend')}...` });
    const frontendBuildResult = await $`cd frontend && bunx --bun astro build`
      .env({ ...process.env, FORCE_COLOR: '1' })
      .nothrow();

    if (frontendBuildResult.exitCode !== 0) {
      buildSpinner.error({ text: red('❌ Frontend build failed.') });
      console.error(red(frontendBuildResult.stderr.toString()));
      process.exit(1);
    }

    // Copy frontend build to backend
    buildSpinner.update({ text: `📁 Copying frontend build to backend/mpa/...` });
    await $`rm -rf backend/mpa && mkdir -p backend/mpa`.quiet().throws(true);
    await $`cp -r frontend/dist/* backend/mpa/`.quiet().throws(true);

    buildSpinner.success({ text: green('✅ Frontend build and copy completed.') });

    // Start Go backend
    console.log(green("🟢 Starting Go backend with embedded frontend (localhost:8099)..."));
    console.log(yellow("⏳ Press Ctrl+C to stop the server."));

    const backendResult = await $`cd backend && go run main.go`
      .env({ ...process.env, FORCE_COLOR: '1' })
      .nothrow();

    if (backendResult.exitCode !== 0 && backendResult.signal !== 'SIGINT' && backendResult.signal !== 'SIGTERM') {
      console.error(red(`Go server exited unexpectedly: Code ${backendResult.exitCode} Signal ${backendResult.signal}`));
    }

  } catch (error) {
    if (buildSpinner.isSpinning()) {
      buildSpinner.error({ text: red('❌ Build process failed.') });
    }
    console.error(red("🚨 Error during dev build:"), error);
    process.exit(1);
  } finally {
    console.log(yellow("\n🛑 Shutting down development environment..."));
  }
  console.log(green("✅ Development environment stopped."));

  // Build Mode - Local binary only
} else if (mode === 'build') {
  console.log(cyan("🚀 Building local MCP server binary..."));

  try {
    await buildLocal();
    console.log(green("✅ Local build completed successfully."));
    console.log(cyan("📁 Binary available at: ./docs4context-com"));
    console.log(yellow("ℹ️ Run with: ./docs4context-com"));

  } catch (error) {
    console.error(red("🚨 Error during build:"), error);
    process.exit(1);
  }

  // CICD / Build Mode
} else if (shouldBuild || shouldCommit || shouldTag || shouldPush || shouldRelease) {

  console.log(cyan("🚀 Starting CICD actions..."));
  let changelogData;

  try {
    if (shouldCommit || shouldTag || shouldPush || shouldRelease || shouldBuild) {
      changelogData = await parseLatestChangelogEntry();
      const { version, summary, description } = changelogData;
      
      if (shouldBuild) {
        await buildCrossPlatform(version);
      }

      if (shouldCommit) {
        await checkGitStatus();
        await gitAdd();
        await gitCommit(summary, description);
      }

      if (shouldTag) {
        await checkGitTagExists(version);
        await gitTag(version, summary);
      }

      if (shouldPush) {
        await gitPush();
      }

      if (shouldRelease) {
        await createGitHubRelease(version, summary, description);
      }
    }

    console.log(green("\n✅ CICD actions completed successfully."));

  } catch (error) {
    console.error(red(`\n🚨 CICD process failed: ${error.message}`));
    if (error.stack && !error.message.startsWith('❌') && !error.message.includes('failed.')) {
      console.error(error.stack);
    }
    process.exit(1);
  }

} else {
  console.error(red("❌ Invalid or missing mode/action specified."));
  console.log(bold("Usage:"));
  console.log(cyan("  bun run cicd.js --mode <dev|build>"));
  console.log(cyan("  bun run cicd.js --build                           # Cross-platform builds"));
  console.log(cyan("  bun run cicd.js [--commit] [--tag] [--push]       # Git operations"));
  console.log(cyan("  bun run cicd.js --release                         # Create GitHub release"));
  console.log(cyan("  bun run cicd.js --build --commit --tag --release  # Full release flow"));
  process.exit(1);
}
