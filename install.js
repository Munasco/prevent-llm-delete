#!/usr/bin/env node

const https = require('https');
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');
const os = require('os');

const VERSION = '1.0.0';
const REPO = 'Munasco/prevent-llm-delete';

// Detect platform and architecture
const platform = os.platform();
const arch = os.arch();

function getPlatformInfo() {
  let osName, archName;

  switch (platform) {
    case 'darwin':
      osName = 'darwin';
      break;
    case 'linux':
      osName = 'linux';
      break;
    case 'win32':
      osName = 'windows';
      break;
    default:
      console.error(`❌ Unsupported platform: ${platform}`);
      process.exit(1);
  }

  switch (arch) {
    case 'x64':
      archName = 'amd64';
      break;
    case 'arm64':
      archName = 'arm64';
      break;
    default:
      console.error(`❌ Unsupported architecture: ${arch}`);
      process.exit(1);
  }

  return { osName, archName };
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    console.log(`⬇️  Downloading from ${url}...`);
    const file = fs.createWriteStream(dest);

    https.get(url, { headers: { 'User-Agent': 'prevent-llm-delete-installer' } }, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        // Follow redirect
        download(response.headers.location, dest).then(resolve).catch(reject);
        return;
      }

      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download: ${response.statusCode}`));
        return;
      }

      response.pipe(file);
      file.on('finish', () => {
        file.close();
        resolve();
      });
    }).on('error', (err) => {
      fs.unlink(dest, () => {});
      reject(err);
    });
  });
}

async function installBinary() {
  console.log('🔒 prevent-llm-delete installer\n');

  const { osName, archName } = getPlatformInfo();
  const isWindows = osName === 'windows';
  const ext = isWindows ? '.zip' : '.tar.gz';
  const filename = `prevent-llm-delete-${osName}-${archName}${ext}`;
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${filename}`;

  const tmpDir = os.tmpdir();
  const downloadPath = path.join(tmpDir, filename);

  try {
    // Download the archive
    await download(url, downloadPath);
    console.log('✅ Download complete');

    // Extract
    console.log('📂 Extracting...');
    const extractDir = path.join(tmpDir, 'prevent-llm-delete-extracted');

    if (!fs.existsSync(extractDir)) {
      fs.mkdirSync(extractDir, { recursive: true });
    }

    if (isWindows) {
      // Windows: use PowerShell to extract
      execSync(`powershell -Command "Expand-Archive -Path '${downloadPath}' -DestinationPath '${extractDir}' -Force"`, { stdio: 'inherit' });
    } else {
      // Unix: use tar
      execSync(`tar -xzf "${downloadPath}" -C "${extractDir}"`, { stdio: 'inherit' });
    }

    // Find the binary
    const binaryName = isWindows ? 'prevent-llm-delete.exe' : `prevent-llm-delete-${osName}-${archName}`;
    const extractedBinary = path.join(extractDir, binaryName);

    // Determine install location
    const installDir = isWindows ? 'C:\\Windows\\System32' : '/usr/local/bin';
    const finalBinary = path.join(installDir, isWindows ? 'prevent-llm-delete.exe' : 'prevent-llm-delete');

    console.log(`📦 Installing to ${installDir}...`);

    // Copy binary
    if (isWindows) {
      // Windows: requires admin
      try {
        fs.copyFileSync(extractedBinary, finalBinary);
      } catch (err) {
        console.log('\n⚠️  Installation requires administrator privileges.');
        console.log('\nPlease run as Administrator:');
        console.log(`  1. Open PowerShell as Administrator`);
        console.log(`  2. Run: npx prevent-llm-delete`);
        process.exit(1);
      }
    } else {
      // Unix: use sudo if needed
      try {
        fs.copyFileSync(extractedBinary, finalBinary);
        fs.chmodSync(finalBinary, '755');
      } catch (err) {
        console.log('   (requires sudo)');
        execSync(`sudo cp "${extractedBinary}" "${finalBinary}"`, { stdio: 'inherit' });
        execSync(`sudo chmod +x "${finalBinary}"`, { stdio: 'inherit' });
      }
    }

    // Cleanup
    fs.unlinkSync(downloadPath);
    fs.rmSync(extractDir, { recursive: true, force: true });

    console.log('\n✅ prevent-llm-delete installed successfully!\n');
    console.log('🚀 Get started:');
    console.log('   prevent-llm-delete install    # Install the protection');
    console.log('   prevent-llm-delete status     # Check status\n');

  } catch (error) {
    console.error('❌ Installation failed:', error.message);
    console.log('\n💡 Try manual installation:');
    console.log(`   curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash`);
    process.exit(1);
  }
}

// If run directly (not via require), install
if (require.main === module) {
  installBinary();
}

module.exports = { installBinary };
