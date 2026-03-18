#!/usr/bin/env node

import { constants, realpathSync } from 'node:fs';
import { access } from 'node:fs/promises';
import { spawn } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import path from 'node:path';
import process from 'node:process';
import { getVendorBinaryName } from '../install.mjs';

const BIN_DIR = path.dirname(fileURLToPath(import.meta.url));

function isDirectExecution(entryUrl, argv1 = process.argv[1]) {
  if (!argv1) {
    return false;
  }

  try {
    return realpathSync(argv1) === realpathSync(fileURLToPath(entryUrl));
  } catch {
    return false;
  }
}

function resolveBinaryPath(platform = process.platform) {
  return path.resolve(BIN_DIR, '..', 'vendor', getVendorBinaryName(platform));
}

export async function main() {
  const binaryPath = resolveBinaryPath();

  try {
    await access(binaryPath, constants.X_OK);
  } catch {
    console.error(
      `smokepod binary is missing at ${binaryPath}. Re-run npm install or set SMOKEPOD_BINARY=/absolute/path/to/smokepod npm install.`
    );
    process.exit(1);
  }

  await new Promise((resolve, reject) => {
    const child = spawn(binaryPath, process.argv.slice(2), {
      stdio: 'inherit'
    });

    child.on('error', reject);
    child.on('exit', (code, signal) => {
      resolve({ code, signal });
    });
  }).then(({ code, signal }) => {
    process.exit(signal ? 1 : (code ?? 1));
  });
}

if (isDirectExecution(import.meta.url)) {
  try {
    await main();
  } catch (error) {
    console.error(error.message);
    process.exitCode = 1;
  }
}
